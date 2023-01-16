import axios from "axios";
import api from "../api";

const checkToken = async () => {
    const token = localStorage.getItem("token");
    let tokenValid = false;
    await axios
        .get(api.validateToken, {
            headers: {
                Authorization: token,
            },
        })
        .then((res) => {
            if (res.data.tokenValid === true) {
                tokenValid = true;
            }
        })
        .catch((err) => {
            console.log(err.response.data);
        });
    return tokenValid;
};

const setAuthToken = (token) => {
    if (token) {
        axios.defaults.headers.common["Authorization"] = token;
    } else {
        delete axios.defaults.headers.common["Authorization"];
    }
};

export { checkToken, setAuthToken };
