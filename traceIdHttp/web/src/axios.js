import axios from 'axios';

// 创建axios实例
const service = axios.create({
    // 服务接口请求
    baseURL: import.meta.env.VITE_APP_BASE_API,
})


export default service;
