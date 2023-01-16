import React, { useEffect } from "react";
import { Button, Form, Input } from "antd";
import { useNavigate } from "react-router-dom";

import { GetWithToken, Get, Post } from "../utils/axios";
import api from "../api";

const App = () => {
    const navigate = useNavigate();

    const onFinish = (values) => {
        console.log(values);
        const res = Post(api.register, values);
        if (res === null) {
            console.log("error");
        } else {
            navigate("/login");
        }
    };
    const onFinishFailed = (errorInfo) => {
        console.log("Failed:", errorInfo);
    };

    useEffect(() => {
        const checkToken = async () => {
            const res = await GetWithToken(api.validateToken);
            if (res.tokenValid === true) {
                navigate("/chat");
            }
        };
        checkToken();
    }, [navigate]);

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
                    label="用户名"
                    name="username"
                    rules={[
                        {
                            required: true,
                            message: "请输入您的用户名",
                        },
                        {
                            pattern: /^[\u4e00-\u9fa5a-zA-Z][\u4e00-\u9fa5a-zA-Z0-9_]{3,15}$/g,
                            message: "用户名需要在4-16个汉字或字母之间",
                        },
                        {
                            validator: async (rule, value) => {
                                if (value.length < 4) {
                                    return Promise.reject;
                                }
                                const res = await Get(api.checkUsername, {
                                    params: {
                                        username: value,
                                    },
                                });

                                console.log(res);
                                if (res.success === true) {
                                    return Promise.reject("用户名已存在");
                                }
                                return Promise.resolve();
                            },
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
                            message: "请输入您的密码",
                        },
                        {
                            min: 4,
                            message: "密码长度不能小于4",
                        },
                    ]}
                >
                    <Input.Password />
                </Form.Item>
                <Form.Item
                    label="确认密码"
                    name="check-passwd"
                    rules={[
                        {
                            required: true,
                            message: "请再次输入您的密码",
                        },
                        ({ getFieldValue }) => ({
                            validator(rule, value) {
                                if (!value || getFieldValue("passwd") === value) {
                                    return Promise.resolve();
                                }
                                return Promise.reject("两次输入的密码不一致");
                            },
                        }),
                    ]}
                >
                    <Input.Password />
                </Form.Item>

                <Form.Item
                    label="邮箱"
                    name="email"
                    rules={[
                        {
                            required: true,
                            message: "请输入您的邮箱",
                        },
                        {
                            type: "email",
                            message: "请输入正确的邮箱格式",
                        },
                        {
                            validator: async (rule, value) => {
                                if (value.length < 6) {
                                    return Promise.resolve();
                                }
                                const res = await Get(api.checkEmail, {
                                    params: {
                                        email: value,
                                    },
                                });
                                console.log(res);
                                if (res.success === true) {
                                    return Promise.reject("邮箱已被注册");
                                }
                                return Promise.resolve();
                            },
                        },
                    ]}
                >
                    <Input />
                </Form.Item>

                <div className="login-register">
                    <span></span>
                    <a href="/login">已有账号？立即登录</a>
                </div>
                <Form.Item
                    wrapperCol={{
                        offset: 8,
                        span: 16,
                    }}
                >
                    <Button type="primary" htmlType="submit">
                        注册
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
};
export default App;
