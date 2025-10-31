import React from "react";
import "./login.scss";
import { loginApi } from "../../api";
import { Button, Form, Input, Typography, Divider } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import MouseTrail from "../../components/MouseTrail";

const { Title } = Typography;

const Login = () => {
  const [form] = Form.useForm();

  const handleLogin = async (values) => {
    console.log("values=?", values);
    const res = await loginApi.login(values);
    console.log("res=", res);
    if (res.code === 200) {
      alert("登录成功");
      window.location.href = "/home";
    } else {
      alert(res.message);
    }
  };

  return (
    <div className="P-login">
      <MouseTrail />
      <div className="login-container">
        <div className="login-card">
          <div className="login-header">
            <Title level={3}>欢迎登录</Title>
            <p className="login-subtitle">请输入您的账号和密码</p>
          </div>

          <Form
            form={form}
            onFinish={handleLogin}
            layout="vertical"
            className="login-form"
          >
            <Form.Item
              name="username"
              rules={[{ required: true, message: "请输入用户名" }]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="用户名"
                size="large"
              />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[{ required: true, message: "请输入密码" }]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="密码"
                size="large"
              />
            </Form.Item>

            <div className="login-options">
              {/*  <a className="forgot-password">忘记密码?</a> */}
            </div>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                className="login-button"
                size="large"
              >
                登录
              </Button>
            </Form.Item>

            {/* <Divider plain>或者</Divider>

            <div className="register-option">
              <span>还没有账号?</span>
              <a className="register-link">立即注册</a>
            </div> */}
          </Form>
        </div>
      </div>
    </div>
  );
};

export default Login;
