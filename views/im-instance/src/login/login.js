import React, { useEffect, useState } from "react";
import { Button, Form, Input } from "antd";
import { useNavigate } from "react-router-dom";

import { GetWithToken, Post } from "../utils/axios";
import api from "../api";

const App = () => {
    const navigate = useNavigate();
    const [loginFailed, setLoginFailed] = useState(false);

    useEffect(() => {
        const checkToken = async () => {
            const res = await GetWithToken(api.validateToken);
            if (res.tokenValid === true) {
                navigate("/chat");
            }
        };
        checkToken();
    }, [navigate]);

    const onFinish = async (values) => {
        const res = await Post(api.login, values);
        if (res === null) {
            setLoginFailed(true);
        } else {
            localStorage.setItem("token", res.token);
            navigate("/chat");
        }
    };

    const onFinishFailed = (errorInfo) => {
        console.log("Failed:", errorInfo);
    };

    const FailedMsg = () => {
        return <p>用户名或密码不正确</p>;
    };
    return (
        <div className="parent-form">
            <Form
                className="child-form"
                labelCol={{
                    span: 8,
                }}
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
                autoComplete="off"
            >
                <Form.Item
                    label="用户名/邮箱"
                    name="username"
                    rules={[
                        {
                            required: true,
                            message: "请输入您的用户名或邮箱",
                        },
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="密码"
                    name="passwd"
                    rules={[
                        {
                            required: true,
                            message: "请输入密码！",
                        },
                    ]}
                >
                    <Input.Password />
                </Form.Item>

                <div className="login-register">
                    <span>{loginFailed ? <FailedMsg /> : null}</span>
                    <a href="/register">没有账号？立即注册</a>
                </div>
                <Form.Item
                    wrapperCol={{
                        offset: 8,
                        span: 16,
                    }}
                >
                    <Button type="primary" htmlType="submit">
                        登录
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
};
export default App;
