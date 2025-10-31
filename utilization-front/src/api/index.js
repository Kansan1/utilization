import request from "../utils/request";
import customRequest from "../utils/request";
import axios from "axios";


// 示例API
export const loginApi = {
  // GET请求示例
  // getList: (params) => {
  //   return request({
  //     url: "/api/list",
  //     method: "get",
  //     params,
  //   });
  // },

  // login: (params) => {
  //   return request({
  //     url: "http://localhost:8080/api/user/login",
  //     method: "post",
  //     params,
  //   });
  // },

  // POST请求示例
  login: (data) => {
    const baseUrl = process.env.REACT_APP_API_URL || "http://192.168.150.1:9020";
    // console.log("REACT_APP_API_URL:", process.env.REACT_APP_API_URL); // 添加这行来调
    //
    // console.log("API URL:", `${baseUrl}/api/user/login`); // 添加这行来调试
    return request({
      url: `${baseUrl}/api/user/login`,
      method: "post",
      data,
    });
  },

  // PUT请求示例
  update: (id, data) => {
    return request({
      url: `/api/update/${id}`,
      method: "put",
      data,
    });
  },

  // DELETE请求示例
  delete: (id) => {
    return request({
      url: `/api/delete/${id}`,
      method: "delete",
    });
  },
};

export const homeAPi = {
  // GET请求示例
  getUtilizationList: (params) => {
    return request({
      url: "/api/home/utilization/list",
      method: "get",
      params,
    });
  },
  getUtilizationAllList: () => {
    return request({
      url: "/api/home/utilization/all",
      method: "get",
    });
  },
  getEquipmentRepairCompletionRate: () => {
    return request({
      url: "/api/home/equipmentRepairCompletionRate/list",
      method: "get",
    });
  },



  addOrUpdateEquipmentRepairCompletionRate: (data) => {
    return request({
      url: "/api/home/equipmentRepairCompletionRate/addOrUpdate",
      method: "post",
      data,
    });
  },

  getEquipmentRepairCompletionTimes: () => {
    return request({
      url: "/api/home/equipmentRepairCompletionTimes/list",
      method: "get",
    });

  },

  addOrUpdateEquipmentRepairCompletionTimes: (data) => {
    return request({
      url: "/api/home/equipmentRepairCompletionTimes/addOrUpdate",
      method: "post",
      data,
    });
  },

  downLoadDate: () => {
    return request({
      url: "/api/excel/download",
      method: "get",
      responseType: 'blob'
    });

  },

  getInspection: () => {
    return request({
      baseURL: "http://192.168.0.103:3004",
      url: "/scan/today",
      method: "get",
    });
  }

};
