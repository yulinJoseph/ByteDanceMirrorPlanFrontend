import React from "react";
import { Badge } from "antd";

import "./contact.css";
import api from "../api";

const selectedStyle = {
    backgroundColor: "rgb(209, 228, 253)",
};

const App = (props) => {
    return (
        <div onClick={props.onClick} className="contactItem" style={props.selected ? selectedStyle : {}}>
            <img className="avatar" src={`${api.getAvatar}?id=${props.user.ID}&ram=${parseInt(Math.random() * 1000)}`} alt=""></img>
            <div
                className="name-signature"
                style={{
                    width: "100%",
                }}
            >
                <div className="username">{props.user.username}</div>
                <div className="signature">{props.message.length === 0 ? "" : props.message[props.message.length - 1].message}</div>
            </div>
            <div>
                <Badge
                    count={props.cnt}
                    style={{
                        marginRight: "10px",
                    }}
                ></Badge>
            </div>
        </div>
    );
};

export default App;
