import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import "./index.css";

import Home from "./home/home";
import Login from "./login/login";
import About from "./about/about";
import Register from "./register/register";
import Chat from "./chat/dashboard";
import Modify from "./modify/modify";

const router = createBrowserRouter([
    {
        path: "/",
        element: <Login />,
    },
    {
        path: "/home",
        element: <Home />,
    },
    {
        path: "/login",
        element: <Login />,
    },
    {
        path: "/about",
        element: <About />,
    },
    {
        path: "/chat",
        element: <Chat />,
    },
    {
        path: "/register",
        element: <Register />,
    },
    {
        path: "/modify",
        element: <Modify />,
    },
]);

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
);
