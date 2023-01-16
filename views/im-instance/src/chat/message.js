import React from "react";
import "./message.css";
import { Button, message } from "antd";

import api from "../api";
import { GetWithToken } from "../utils/axios";

const Message = (props) => {
    const [errInfo, contextHolder] = message.useMessage();

    if (props.message.type === "agree" || props.message.type === "disagree") {
        return (
            <div>
                <div
                    className="message-text"
                    style={{
                        textAlign: "center",
                        border: "1px solid black",
                    }}
                >
                    {props.message.message}
                </div>
            </div>
        );
    }

    if (props.message.type === "invite" && props.direction === "left") {
        return (
            <div
                className="content"
                style={{
                    float: props.direction,
                }}
            >
                {contextHolder}
                <div className="message-text">{props.message.message}</div>
                <div className="message-button">
                    <Button
                        onClick={async () => {
                            const res = await GetWithToken(`${api.agreeInvitation}?fromID=${props.message.fromID}`);
                            if (res.success === false) {
                                errInfo.error(res.message === "agree" ? "已经同意过了" : "已经拒绝过了");
                            }
                        }}
                    >
                        接受
                    </Button>
                    <Button
                        onClick={async () => {
                            const res = await GetWithToken(`${api.disagreeInvitation}?fromID=${props.message.fromID}`);
                            if (res.success === false) {
                                errInfo.error(res.message === "agree" ? "已经同意过了" : "已经拒绝过了");
                            }
                        }}
                    >
                        拒绝
                    </Button>
                </div>
            </div>
        );
    }
    return (
        <div>
            <div
                className="content"
                style={{
                    float: props.direction,
                }}
            >
                {props.message.message}
                {/* {props.message.type === "single" ? props.message.message : ""} */}
                {/* {props.message.type === "invite" ? props.message.message : ""} */}
            </div>
            <div
                className="have-read"
                style={{
                    float: props.direction,
                }}
            >
                {props.message.haveRead ? "已读" : "未读"}
            </div>
        </div>
    );
};

const App = (props) => {
    const timestamp = new Date(props.message.CreatedAt);
    const now = new Date();
    const year = timestamp.getFullYear() === now.getFullYear() ? "" : timestamp.getFullYear() + "年";
    const month = timestamp.getMonth() === now.getMonth() ? "" : timestamp.getMonth() + 1 + "月";
    const day = timestamp.getDate() === now.getDate() ? "" : timestamp.getDate() + "日";
    const hour = timestamp.getHours() + ":";
    const minute = timestamp.getMinutes() < 10 ? "0" + timestamp.getMinutes() : timestamp.getMinutes();

    return (
        <div className="message-wrapper">
            <div className="timestamp">
                {year}
                {month}
                {day}
                {hour}
                {minute}
            </div>

            <div>
                <Message message={props.message} direction={props.direction}></Message>
            </div>
        </div>
    );
};

export default App;
