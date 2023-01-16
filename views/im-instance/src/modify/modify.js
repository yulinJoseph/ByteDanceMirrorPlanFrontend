import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

import { GetWithToken } from "../utils/axios";
import api from "../api";
import "./modify.css";

import ModifyAvatar from "./modify-avatar";
import ModifyUsername from "./modify-username";
import ModifySignature from "./modify-signature";

const App = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [selfInfo, setSelfInfo] = useState({});

    useEffect(() => {
        const checkToken = async () => {
            const res = await GetWithToken(api.validateToken);
            if (res.tokenValid === false) {
                navigate("/login");
            }
        };
        checkToken();

        const getSelfInfo = async () => {
            const res = await GetWithToken(api.getSelfInfo);
            setSelfInfo(res.message);
            setLoading(false);
        };
        getSelfInfo();
    }, [navigate]);

    if (loading) {
        return <div>loading</div>;
    }
    return (
        <div>
            <div className="modify-box">
                <a href="/chat">返回</a>
            </div>
            <div className="modify-box">
                <ModifyAvatar user={selfInfo}></ModifyAvatar>
            </div>
            <div className="modify-box">
                <ModifyUsername user={selfInfo}></ModifyUsername>
            </div>
            <div className="modify-box">
                <ModifySignature user={selfInfo}></ModifySignature>
            </div>
        </div>
    );
};

export default App;
