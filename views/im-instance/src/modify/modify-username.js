import React, { useState } from "react";
import { Button, Input, message } from "antd";
import { Get, GetWithToken } from "../utils/axios";
import api from "../api";

const App = (props) => {
    const [usernameExits, setUsernameExits] = useState(false);
    const [usernameValid, setUsernameValid] = useState(true);
    const [usernameModified, setUsernameModified] = useState("");
    const [usernameSame, setUsernameSame] = useState(false);
    const [modifyInfo, contextHolder] = message.useMessage();

    const checkUsername = async (e) => {
        let pattern = /^[\u4e00-\u9fa5a-zA-Z][\u4e00-\u9fa5a-zA-Z0-9_]{3,15}$/;
        if (!pattern.test(e.target.value)) {
            setUsernameValid(false);
        } else {
            setUsernameValid(true);
            if (e.target.value === props.user.username) {
                setUsernameExits(false);
                setUsernameSame(true);
                return;
            }
            const res = await Get(api.checkUsername, {
                params: {
                    username: e.target.value,
                },
            });
            if (res.success === true) {
                setUsernameExits(true);
                setUsernameSame(false);
            } else {
                setUsernameSame(false);
                setUsernameExits(false);
                setUsernameModified(e.target.value);
            }
        }
    };

    return (
        <div>
            {contextHolder}
            <div
                style={{
                    marginLeft: "30%",
                }}
            >
                <strong>当前用户名：</strong>
            </div>
            <div>
                <Input
                    style={{
                        width: "150px",
                        marginLeft: "30%",
                    }}
                    type="text"
                    defaultValue={props.user.username}
                    onBlur={checkUsername}
                    onChange={checkUsername}
                ></Input>
                <span
                    style={{
                        color: "red",
                        fontSize: "12px",
                    }}
                >
                    {usernameExits ? "用户名已存在" : ""}
                    {usernameValid ? "" : "用户名需要在4-16个汉字或字母之间"}
                    {usernameSame ? "用户名未修改" : ""}
                </span>
            </div>
            <div>
                <Button
                    style={{
                        marginLeft: "30%",
                        width: "150px",
                    }}
                    onClick={async () => {
                        if (usernameExits || !usernameValid || usernameSame) {
                            return;
                        }
                        const res = await GetWithToken(`${api.updateUsername}?username=${usernameModified}`);
                        if (res.success === true) {
                            modifyInfo.success("用户名已存在");
                        }
                    }}
                >
                    点击提交修改
                </Button>
            </div>
        </div>
    );
};

export default App;
