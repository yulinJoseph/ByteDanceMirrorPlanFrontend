const api = {
    validateToken: "http://192.168.3.50:7898/validateToken",
    register: "http://192.168.3.50:7898/register",
    login: "http://192.168.3.50:7898/login",
    checkUsername: "http://192.168.3.50:7898/user/checkUsername",
    checkEmail: "http://192.168.3.50:7898/user/checkEmail",

    getSelfInfo: "http://192.168.3.50:7898/user/getSelfInfo",
    getContactList: "http://192.168.3.50:7898/user/getContactList",
    getUserList: "http://192.168.3.50:7898/user/getUserList",
    getContactWithMessages: "http://192.168.3.50:7898/user/getContactWithMessages",
    getMessages: "http://192.168.3.50:7898/user/getMessages",
    uploadAvatar: "http://192.168.3.50:7898/uploadAvatar",
    getAvatar: "http://192.168.3.50:7898/getAvatar",
    updateUsername: "http://192.168.3.50:7898/user/updateUsername",
    updateSignature: "http://192.168.3.50:7898/user/updateSignature",
    sendInvitation: "http://192.168.3.50:7898/user/sendInvitation",
    getUserByID: "http://192.168.3.50:7898/user/getUserByID",
    agreeInvitation: "http://192.168.3.50:7898/agreeInvitation",
    disagreeInvitation: "http://192.168.3.50:7898/disagreeInvitation",

    sendMsg: "ws://192.168.3.50:7898/sendMsg",
};

export default api;
