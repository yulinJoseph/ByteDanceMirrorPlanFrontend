import { React, useCallback, useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { message, Dropdown, Input } from "antd";
import data from "@emoji-mart/data";
import Picker from "@emoji-mart/react";

import Contact from "./contact";
import api from "../api";
import Message from "./message";
import { useCallbackState } from "../utils/use-state";
import "./dashboard.css";
import More from "./more";

import { GetWithToken } from "../utils/axios";

const { Search } = Input;

const items = [
    {
        key: "1",
        label: (
            <Picker
                data={data}
                onEmojiSelect={(e) => {
                    console.log(e.native);
                    let tc = document.getElementsByClassName("input-box")[0];
                    let len = tc.value.length;
                    tc.focus();
                    let start = tc.selectionStart;
                    tc.value = tc.value.substr(0, tc.selectionStart) + e.native + tc.value.substring(tc.selectionStart, len);
                    tc.selectionStart = tc.selectionEnd = start + e.native.length;
                }}
            />
        ),
    },
];

const App = () => {
    const navigate = useNavigate();
    const ws = useRef();
    const chosen = useRef({});
    const selfID = useRef(0);
    const [loading, setLoading] = useState(true);
    const [contactWithMessages, setContactWithMessages] = useCallbackState([]);
    const [messages, setMessages] = useCallbackState([]);
    const [selfInfo, setSelfInfo] = useState({});
    const [errInfo, contextHolder] = message.useMessage();

    const getContactWithMessages = async () => {
        const res = await GetWithToken(api.getContactWithMessages);
        setContactWithMessages(res.message, () => {
            if (JSON.stringify(chosen.current) === "{}") {
            } else {
                for (let item of res.message) {
                    if (item.contact.ID === chosen.current.contact.ID) {
                        chosen.current = item;
                        setMessages(item.msg, () => updateScroll());
                        break;
                    }
                }
            }
        });
    };

    const webSocketInit = useCallback(() => {
        ws.current = new WebSocket(`${api.sendMsg}?token=${[localStorage.getItem("token")]}`);
        ws.current.onmessage = (e) => {
            e = JSON.parse(e.data);
            if (JSON.stringify(chosen.current) !== "{}" && e.fromID === chosen.current.contact.ID && e.type === "single") {
                readMessage(selfID.current, e.fromID, "ping");
            } else if (e.type === "ping") {
            }
            getContactWithMessages();
        };
    }, [ws]);

    const updateScroll = () => {
        if (JSON.stringify(chosen.current) === "{}") {
            return;
        }
        let messageBox = document.getElementsByClassName("recv-box")[0];
        messageBox.scrollTop = messageBox.scrollHeight;
    };

    const sendMessage = () => {
        const toUser = chosen.current.contact.ID;
        if (chosen.current.status !== "agree") {
            errInfo.error("你还不是对方的好友，无法发送消息");
            return;
        }
        const msg = document.querySelector(".input-box").value;
        if (msg === "") {
            errInfo.error("消息不能为空");
            return;
        }
        ws.current.send(
            JSON.stringify({
                fromID: selfInfo.ID,
                toID: parseInt(toUser),
                type: "single",
                message: msg,
            })
        );
    };

    const getSelfInfo = async () => {
        const res = await GetWithToken(api.getSelfInfo);
        setSelfInfo(res.message);
        setLoading(false);
        selfID.current = res.message.ID;
    };

    const readMessage = async (fromID, toID, type) => {
        ws.current.send(
            JSON.stringify({
                fromID: parseInt(fromID),
                toID: parseInt(toID),
                type: type,
                message: "",
            })
        );
    };

    const GetUnreadCnt = (item) => {
        let cnt = 0;
        item.msg.forEach((msg) => {
            if (msg.haveRead === false && msg.fromID !== selfInfo.ID) {
                cnt++;
            }
        });
        return cnt;
    };

    const onSearch = (value) => {
        if (value === "") {
            getContactWithMessages();
            return;
        }
        const res = contactWithMessages.filter((item) => {
            return item.contact.username.includes(value);
        });
        setContactWithMessages(res);
    };

    useEffect(() => {
        const checkToken = async () => {
            const res = await GetWithToken(api.validateToken);
            if (res.tokenValid === false) {
                navigate("/login");
            }
        };
        checkToken();
        getSelfInfo();
        getContactWithMessages();

        webSocketInit();
    }, [navigate, webSocketInit]);

    const Right = () => {
        return (
            <div className="right">
                <div className="self-info-box to-info-box">
                    <img className="avatar" src={`${api.getAvatar}?id=${chosen.current.contact.ID}&ram=${parseInt(Math.random() * 1000)}`} alt=""></img>
                    <div className="name-signature">
                        <div className="username">
                            <span>{chosen.current.contact.username}</span>
                        </div>
                        <div className="signature">{chosen.current.contact.signature}</div>
                    </div>
                </div>
                <div className="message-box">
                    <div className="recv-box">
                        {messages.map((msg, key) => {
                            return <Message key={key} message={msg} direction={msg.fromID === selfInfo.ID ? "right" : "left"}></Message>;
                        })}
                    </div>
                    <div className="toolbar-box">
                        <Dropdown
                            menu={{
                                items,
                            }}
                            placement="top"
                            trigger={["click"]}
                        >
                            <button className="emoji-button">表情</button>
                        </Dropdown>
                        <button className="send-button" onClick={sendMessage}>
                            发送
                        </button>
                    </div>
                    <textarea
                        className="input-box"
                        type="text"
                        onKeyDown={(e) => {
                            if (e.key === "Enter") {
                                sendMessage();
                            }
                        }}
                    />
                </div>
            </div>
        );
    };

    if (loading) {
        return <div>loading</div>;
    }
    return (
        <div className="parent">
            {contextHolder}
            <div className="left">
                <div className="self-info-box">
                    <img className="avatar" src={`${api.getAvatar}?id=${selfInfo.ID}&ram=${parseInt(Math.random() * 1000)}`} alt="my-avatar"></img>
                    <div className="name-signature">
                        <div className="username">
                            <span>{selfInfo.username}</span>
                            <More />
                        </div>
                        <div className="signature">{selfInfo.signature}</div>
                    </div>
                </div>
                <Search
                    placeholder="搜索好友"
                    onSearch={onSearch}
                    onChange={(e) => {
                        if (e.target.value === "") {
                            onSearch("");
                        }
                    }}
                    style={{
                        margin: "3px",
                    }}
                />
                <div className="contact">
                    {contactWithMessages.map((item, key) => {
                        let cnt = GetUnreadCnt(item);
                        return (
                            <Contact
                                key={key}
                                onClick={() => {
                                    chosen.current = item;
                                    setMessages(item.msg, () => updateScroll());
                                    readMessage(selfInfo.ID, chosen.current.contact.ID, "ping");
                                }}
                                selected={chosen.current === item}
                                user={item.contact}
                                message={item.msg}
                                cnt={cnt}
                            ></Contact>
                        );
                    })}
                </div>
            </div>

            {JSON.stringify(chosen.current) === "{}" ? "" : <Right />}
        </div>
    );
};

export default App;
