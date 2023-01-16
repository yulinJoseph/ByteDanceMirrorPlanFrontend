package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"im-instance/utils"
	"net"
	"net/http"
	"sync"
)

type Message struct {
	gorm.Model
	FromID   uint64 `json:"fromID"`
	ToID     uint64 `json:"toID"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	HaveRead bool   `json:"haveRead"`
}

func (table *Message) TableName() string {
	return "message"
}

func SaveMsg(message *Message) error {
	return utils.DB.Create(message).Error
}
func ReadMsg(fromID uint, toID uint) error {
	return utils.DB.Model(&Message{}).
		Where("from_id = ? and to_id = ? and have_read = ?", toID, fromID, false).
		Update("have_read", true).Error
}

func GetMsgList(fromId, toID uint) ([]*Message, error) {
	var data []*Message
	db := utils.DB.
		Order("id").
		Where("from_id = ? and to_id = ? or from_id = ? and to_id = ?", fromId, toID, toID, fromId).
		Find(&data)
	return data, db.Error
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

var ClientMap = map[uint64]*Node{}
var rwmutex = sync.RWMutex{}

func Chat(w http.ResponseWriter, r *http.Request, fromID uint64) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	rwmutex.Lock()
	ClientMap[fromID] = node
	rwmutex.Unlock()
	go sendProc(node)
	go receiveProc(node)
}

// websocket 发送消息
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}

// websocket 收到消息之后发送到 hub 之后转发
func receiveProc(node *Node) {
	for {
		_, msg, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		broadMsg(msg)
		//fmt.Println("[ws]:", msg)
	}
}

// 中转站
var udpChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpChan <- data
}

// 中转站接受和发送消息的动作
func init() {
	go udpSendProc()
	go udpReceiveProc()
}

func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 7899,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(conn)
	for {
		select {
		case data := <-udpChan:
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}

func udpReceiveProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 7899,
	})
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}(conn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		var buffer [1024]byte
		n, err := conn.Read(buffer[0:])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dispatch(buffer[0:n])
	}
}

var Writing sync.RWMutex

// 根据消息类型进行分发
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	switch msg.Type {
	case "single":
		SendSingleMsg(msg.FromID, msg.ToID, data)
		Writing.Lock()
		defer Writing.Unlock()
		if err := SaveMsg(&msg); err != nil {
			fmt.Println(err.Error())
			return
		}
	case "ping", "pong":
		readMsg(msg.FromID, msg.ToID, data)
		Writing.Lock()
		defer Writing.Unlock()
		if err := ReadMsg(uint(msg.FromID), uint(msg.ToID)); err != nil {
			fmt.Println(err.Error())
			return
		}
	case "invite":
	default:
		sendGroupMsg()
	}

}

func SendSingleMsg(fromID uint64, toID uint64, msg []byte) {
	rwmutex.RLock()
	nodeTo, ok1 := ClientMap[toID]
	nodeFrom, ok2 := ClientMap[fromID]
	rwmutex.RUnlock()
	if ok1 && ok2 {
		nodeTo.DataQueue <- msg
	}
	nodeFrom.DataQueue <- msg
}

func readMsg(fromID uint64, toID uint64, msg []byte) {
	rwmutex.RLock()
	nodeTo, ok1 := ClientMap[toID]
	nodeFrom, ok2 := ClientMap[fromID]
	rwmutex.RUnlock()
	if ok1 && ok2 {
		nodeTo.DataQueue <- msg
	}
	nodeFrom.DataQueue <- msg
}

func sendGroupMsg() {

}
