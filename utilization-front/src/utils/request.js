import axios from "axios";
import { message } from "antd";

// 创建axios实例
const request = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  timeout: 15000, // 请求超时时间
  headers: {
    "Content-Type": "application/json",
  },
  // 自定义状态码验证函数，返回true表示不触发错误
  validateStatus: function (status) {
    return status >= 200 && status < 600; // 接受所有状态码，自己处理错误
  },
});

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 在这里可以添加token等认证信息
    const token = localStorage.getItem("token");
    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    console.log("请求配置：", config);

    return config;
  },
  (error) => {
    console.error("请求错误：", error);
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
    (response) => {
        // 如果是 blob 类型（文件下载），就返回完整 response
        const contentType = response.headers['content-type'];
        if (response.config.responseType === 'blob' || contentType?.includes('application/vnd.openxmlformats-officedocument.spreadsheetml.sheet')) {
            return response; // ⬅️ 关键：返回完整 response 包含 headers
        }

        const res = response.data;

        // 可选：根据后端结构做通用判断
        // if (res.code !== 200) {
        //     message.error(res.message || "请求失败");
        //     return Promise.reject(new Error(res.message || "请求失败"));
        // }

        return res; // 默认返回 data
    },
    (error) => {
        console.error("响应错误：", error);
        if (error.response && error.response.status === 401) {
            message.error("用户未登录或登录已过期，请重新登录");
            window.location.href = "/login";
        } else {
            message.error(error.message || "网络错误");
        }
        return Promise.reject(error);
    }
);

export default request;
