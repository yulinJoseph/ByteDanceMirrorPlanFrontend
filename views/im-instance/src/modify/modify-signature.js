import React, { useEffect, useReducer, useState } from "react";
import { Button, Input, message } from "antd";
import { GetWithToken } from "../utils/axios";
import api from "../api";

const { TextArea } = Input;

const App = (props) => {
    const [signatureModified, setSignatureModified] = useState("");
    const [modifyInfo, contextHolder] = message.useMessage();

    useReducer(
        (state, action) => {
            switch (action.type) {
                case "setSignatureModified":
                    return {
                        ...state,
                        signatureModified: action.payload,
                    };
                default:
                    return state;
            }
        },
        {
            signatureModified: "",
        }
    );

    useEffect(() => {
        setSignatureModified(props.user.signature);
    }, [props.user.signature]);

    return (
        <div>
            {contextHolder}
            <div
                style={{
                    marginLeft: "30%",
                }}
            >
                <strong>当前个性签名：</strong>
            </div>
            <div>
                <TextArea
                    rows={4}
                    maxLength={30}
                    style={{
                        width: "300px",
                        marginLeft: "30%",
                    }}
                    type="text"
                    defaultValue={props.user.signature}
                    onChange={(e) => {
                        setSignatureModified(e.target.value);
                    }}
                ></TextArea>

                <span
                    style={{
                        color: "gray",
                        fontSize: "12px",
                    }}
                >
                    还可输入{30 - signatureModified.length}个字符
                </span>
            </div>
            <div>
                <Button
                    style={{
                        marginLeft: "30%",
                        width: "150px",
                    }}
                    onClick={async () => {
                        const res = await GetWithToken(`${api.updateSignature}?signature=${signatureModified}`);
                        if (res.success === true) {
                            modifyInfo.success("修改成功");
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
