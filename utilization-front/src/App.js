import React, { Fragment } from "react";
import Login from "./pages/login";
import Home from "./pages/home";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ConfigProvider, theme } from "antd";

function App() {
    return (
        <ConfigProvider
            theme={{
                algorithm: theme.darkAlgorithm, // 开启深色主题
            }}
        >
            <Fragment>
                <BrowserRouter>
                    <Routes>
                        <Route path="/login" element={<Login />} />
                        <Route path="/home" element={<Home />} />
                        <Route path="*" element={<Navigate to="/login" />} />
                    </Routes>
                </BrowserRouter>
            </Fragment>
        </ConfigProvider>
    );
}

export default App;
