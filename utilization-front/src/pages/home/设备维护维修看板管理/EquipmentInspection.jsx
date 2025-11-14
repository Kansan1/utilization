import React, { useEffect, useState } from 'react';
import { Table, message, Tag, Spin } from 'antd';
import { homeAPi } from '../../../api';
import { io } from "socket.io-client";

// 连接到您的 Go 后端 WebSocket 服务
// 请确保这里的 IP 和端口是您后端服务的正确地址
const socket = io("http://localhost:9020", {
    // transports: ['websocket'], // 强制使用 WebSocket 协议
});

const EquipmentInspection = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);

    // 定义新的表格列结构
    const columns = [
        {
            title: '序号',
            dataIndex: 'seq',
            key: 'seq',
            width: 80,
            align: 'center',
        },
        {
            title: '设备名称',
            dataIndex: 'equip_name',
            key: 'equip_name',
        },
        {
            title: '设备编号',
            dataIndex: 'equip_code',
            key: 'equip_code',
        },
        {
            title: '设备位置',
            dataIndex: 'equip_location',
            key: 'equip_location',
        },
        {
            title: '设备类型',
            dataIndex: 'equip_type',
            key: 'equip_type',
        },
        {
            title: '当日点检状态',
            dataIndex: 'inspected',
            key: 'inspected',
            align: 'center',
            width: 150,
            render: (inspected) => (
                inspected
                    ? <Tag color="green">已点检</Tag>
                    : <Tag color="red">未点检</Tag>
            ),
        },
    ];

    // 从新的 API 获取数据
    const fetchData = async () => {
        if (!loading) setLoading(true);
        try {
            const res = await homeAPi.getEquipmentInspectionStatus();
            if (res.code === 200) {
                setData(res.data || []);
            } else {
                message.error(res.message || '获取设备点检数据失败');
            }
        } catch (error) {
            console.error('获取设备点检数据失败：', error);
            message.error('网络错误或服务器无响应');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        // 初始加载数据
        fetchData();

        // 监听 WebSocket 连接事件
        socket.on('connect', () => {
            console.log('WebSocket 连接成功');
        });

        // 监听后端广播的 'inspection-update' 事件
        socket.on('inspection-update', (msg) => {
            console.log('收到点检更新通知:', msg);
            message.success('点检状态已更新，正在刷新列表...');
            fetchData(); // 重新获取数据
        });

        // 监听断开连接事件
        socket.on('disconnect', () => {
            console.log('WebSocket 断开连接');
        });

        // 组件卸载时清理监听器
        return () => {
            socket.off('connect');
            socket.off('inspection-update');
            socket.off('disconnect');
        };
    }, []); // 空依赖数组确保只在组件挂载时执行一次

    return (
        <Spin spinning={loading} tip="正在加载数据...">
            <Table
                columns={columns}
                dataSource={data}
                rowKey="equip_code" // 使用唯一的 equip_code 作为 key
                bordered
                pagination={{ pageSize: 20 }} // 调整分页大小
                title={() => <h2 style={{ textAlign: 'center', margin: 0 }}>设备每日点检状态看板</h2>}
            />
        </Spin>
    );
};

export default EquipmentInspection;
