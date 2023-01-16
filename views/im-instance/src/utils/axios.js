import axios from "axios";

async function GetWithToken(url) {
    try {
        const res = await axios.get(url, {
            headers: {
                Authorization: localStorage.getItem("token"),
            },
        });
        return res.data;
    } catch (err) {
        console.log(err);
        return null;
    }
}

async function Get(url, values) {
    try {
        const res = await axios.get(url, values);
        return res.data;
    } catch (err) {
        console.log(err);
        return null;
    }
}

async function PostWithToken(url, values) {
    try {
        const res = await axios.post(url, values, {
            headers: {
                Authorization: localStorage.getItem("token"),
            },
        });
        return res.data;
    } catch (err) {
        console.log(err);
        return null;
    }
}

async function Post(url, values) {
    try {
        const res = await axios.post(url, values);
        return res.data;
    } catch (err) {
        console.log(err);
        return null;
    }
}

export { GetWithToken, Get, PostWithToken, Post };
