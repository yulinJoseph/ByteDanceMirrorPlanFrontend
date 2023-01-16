import React, { useState } from "react";
import { DownOutlined } from "@ant-design/icons";
import { Dropdown, Space, Modal, Input, message } from "antd";
import { useNavigate } from "react-router-dom";

import { GetWithToken, Get } from "../utils/axios";
import api from "../api";

import "./more.css";

const App = () => {
    const navigate = useNavigate();
    const [someInfo, contextHolder] = message.useMessage();

    const [isModalOpen, setIsModalOpen] = useState(false);
    const showModal = () => {
        setIsModalOpen(true);
    };
    const handleOk = async () => {
        const toUser = document.getElementById("toUser").value;
        const res = await Get(`${api.checkUsername}?username=${toUser}`);
        if (res.success === false) {
            someInfo.error("用户不存在");
            return;
        }
        GetWithToken(`${api.sendInvitation}?toUser=${toUser}`);
        setIsModalOpen(false);
    };
    const handleCancel = () => {
        setIsModalOpen(false);
    };

    const items = [
        {
            key: "1",
            label: (
                <a target="_self" onClick={() => {}} href="/modify">
                    修改个人信息
                </a>
            ),
        },
        {
            key: "2",
            label: <div onClick={showModal}>添加好友</div>,
        },
        {
            key: "3",
            danger: true,
            label: "退出登录",
            onClick: () => {
                localStorage.removeItem("token");
                navigate("/login");
            },
        },
    ];

    return (
        <Dropdown
            menu={{
                items,
            }}
        >
            <button className="more-button" onClick={(e) => e.preventDefault()}>
                <Space className="hover">
                    更多
                    <DownOutlined />
                </Space>
                {contextHolder}

                <Modal
                    style={{
                        top: "30%",
                    }}
                    title="添加好友"
                    open={isModalOpen}
                    onOk={handleOk}
                    onCancel={handleCancel}
                >
                    好友名称/邮箱：<Input id="toUser"></Input>
                </Modal>
            </button>
        </Dropdown>
    );
};
export default App;
