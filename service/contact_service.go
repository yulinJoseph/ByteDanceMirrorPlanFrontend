package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"im-instance/models"
	"im-instance/utils"
	"net/http"
	"strconv"
)

func GetContactList(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if utils.CheckToken(token) {
		list, err := models.GetContactList(utils.GetIdFromToken(token))
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": list,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token invalid",
		})
	}
}

type ContactWithMessages struct {
	Contact *models.UserBasic `json:"contact"`
	Msg     []*models.Message `json:"msg"`
	Status  string            `json:"status"`
}

// 123

func GetContactWithMessages(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if utils.CheckToken(token) {
		fromID := utils.GetIdFromToken(token)
		users, err := models.GetContactList(fromID)
		var list = make([]ContactWithMessages, len(users))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "get contact list failed",
			})
			return
		}
		models.Writing.RLock()
		for i, v := range users {
			list[i].Contact = v
			messages, err := models.GetMsgList(fromID, v.ID)
			list[i].Msg = messages
			status, _ := models.GetContactStatus(fromID, v.ID)
			list[i].Status = status
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "get messages failed",
				})
			}
		}
		models.Writing.RUnlock()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": list,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token invalid",
		})
	}
}

func SendInvitation(c *gin.Context) {
	fromID := utils.GetIdFromToken(c.GetHeader("Authorization"))
	toUser := c.Query("toUser")
	user, err := models.FindUserByUsername(toUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "user not exists",
		})
	}
	toID := user.ID
	if models.CheckContact(fromID, toID) && models.CheckContact(toID, fromID) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "already friends",
		})
	}
	fromUser, _ := models.GetUserByID(fromID)
	toUser_, _ := models.GetUserByID(toID)
	message := models.Message{
		Model:    gorm.Model{},
		FromID:   uint64(fromID),
		ToID:     uint64(toID),
		Type:     "invite",
		Message:  fromUser.Username + " 邀请 " + toUser_.Username + " 成为好友",
		HaveRead: false,
	}
	err = models.SaveMsg(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "save invitation failed",
		})
		return
	}
	err = models.Invite(fromID, toID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invite failed",
		})
		return
	}
	if models.ClientMap[uint64(toID)] != nil {
		data, err := json.Marshal(message)
		if err != nil {
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "json unmarshal failed",
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "invite failed",
			})
			return
		}
		models.SendSingleMsg(uint64(fromID), uint64(toID), data)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "send invitation success",
	})

}

func AgreeInvitation(c *gin.Context) {
	toID := utils.GetIdFromToken(c.GetHeader("Authorization"))
	fromID, _ := strconv.ParseUint(c.Query("fromID"), 10, 64)

	status, err := models.GetContactStatus(uint(fromID), toID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "get status failed",
		})
		return
	}
	if status == "agree" || status == "disagree" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": status,
		})
		return
	}

	err = models.Agree(uint(fromID), toID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "db failed",
		})
		return
	}
	message := models.Message{
		Model:    gorm.Model{},
		FromID:   fromID,
		ToID:     uint64(toID),
		Type:     "agree",
		Message:  "已同意好友申请，可以开始聊天了",
		HaveRead: true,
	}
	err = models.SaveMsg(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "save invitation failed",
		})
		return
	}
	if models.ClientMap[fromID] != nil {
		data, err := json.Marshal(message)
		if err != nil {
			return
		}
		models.SendSingleMsg(fromID, uint64(toID), data)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "agree invitation success",
	})
}

func DisagreeInvitation(c *gin.Context) {
	toID := utils.GetIdFromToken(c.GetHeader("Authorization"))
	fromID, _ := strconv.ParseUint(c.Query("fromID"), 10, 64)

	status, err := models.GetContactStatus(uint(fromID), toID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "get status failed",
		})
		return
	}
	if status == "agree" || status == "disagree" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": status,
		})
		return
	}
	err = models.Disagree(uint(fromID), toID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "db failed",
		})
		return
	}
	message := models.Message{
		Model:    gorm.Model{},
		FromID:   fromID,
		ToID:     uint64(toID),
		Type:     "disagree",
		Message:  "已拒绝好友申请",
		HaveRead: true,
	}
	err = models.SaveMsg(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "save invitation failed",
		})
		return
	}
	if models.ClientMap[fromID] != nil {
		data, err := json.Marshal(message)
		if err != nil {
			return
		}
		models.SendSingleMsg(fromID, uint64(toID), data)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "disagree invitation success",
	})
}
