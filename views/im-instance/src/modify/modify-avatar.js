import api from "../api";
import { LoadingOutlined, PlusOutlined } from "@ant-design/icons";
import { message, Upload, Button } from "antd";
import { useState } from "react";

const getBase64 = (img, callback) => {
    const reader = new FileReader();
    reader.addEventListener("load", () => callback(reader.result));
    reader.readAsDataURL(img);
};
const beforeUpload = (file) => {
    const isJpgOrPng = file.type === "image/jpeg" || file.type === "image/png";
    if (!isJpgOrPng) {
        message.error("只能上传 JPG/PNG 文件!");
    }
    const isLt2M = file.size / 1024 / 1024 < 2;
    if (!isLt2M) {
        message.error("图像需要小于 2MB!");
    }
    return isJpgOrPng && isLt2M;
};

const App = (props) => {
    const [loading, setLoading] = useState(false);
    const [imageUrl, setImageUrl] = useState();
    const [showUpload, setShowUpload] = useState({ show: false, message: "点击修改" });
    const [modifyInfo, contextHolder] = message.useMessage();

    const handleChange = (info) => {
        if (info.file.status === "uploading") {
            setLoading(true);
            return;
        }
        if (info.file.status === "done") {
            // Get this url from response in real world.
            getBase64(info.file.originFileObj, (url) => {
                setLoading(false);
                setImageUrl(url);
            });
            modifyInfo.success("修改成功");
        }
    };
    const uploadButton = (
        <div>
            {loading ? <LoadingOutlined /> : <PlusOutlined />}
            <div
                style={{
                    marginTop: 8,
                }}
            >
                Upload
            </div>
        </div>
    );
    return (
        <div
            style={{
                display: "flex",
                alignItems: "center",
                height: "100px",
            }}
        >
            {contextHolder}
            <img
                style={{
                    width: "90px",
                    height: "90px",
                    borderRadius: "100%",
                    borderColor: "black",
                    borderStyle: "solid",
                    borderWidth: "1px",
                    marginLeft: "30%",
                }}
                src={`${api.getAvatar}?id=${props.user.ID}&ram=${parseInt(Math.random() * 1000)}`}
                alt=""
            ></img>
            <Button
                style={{
                    marginLeft: "10px",
                    marginRight: "10px",
                }}
                onClick={() => {
                    if (showUpload.show === true) {
                        setShowUpload({ show: false, message: "点击修改" });
                    } else {
                        setShowUpload({ show: true, message: "点击取消" });
                    }
                }}
            >
                {showUpload.message}
            </Button>
            {showUpload.show === true ? (
                <Upload
                    name="avatar"
                    listType="picture-card"
                    className="avatar-uploader"
                    showUploadList={false}
                    action={api.uploadAvatar}
                    beforeUpload={beforeUpload}
                    onChange={handleChange}
                    headers={{
                        Authorization: localStorage.getItem("token"),
                    }}
                >
                    {imageUrl ? (
                        <img
                            src={imageUrl}
                            alt=""
                            style={{
                                width: "100%",
                            }}
                        />
                    ) : (
                        uploadButton
                    )}
                </Upload>
            ) : null}
        </div>
    );
};

export default App;
