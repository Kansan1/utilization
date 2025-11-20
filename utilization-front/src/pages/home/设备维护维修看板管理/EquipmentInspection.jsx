import React, { useEffect, useState } from 'react';
import { message, Spin, Table, Tag, Radio } from 'antd'; // 导入 Radio 组件
import { homeAPi } from '../../../api'; // 确保你正确导入了 homeAPi
import { io } from "socket.io-client";

// 连接到 WebSocket 服务
const socket = io("http://localhost:9020", {
    transports: ['websocket'],
});

const EquipmentInspection = () => {
    const [data, setData] = useState([]);
    const [filterStatus, setFilterStatus] = useState('全部');
    const [filterType, setFilterType] = useState('全部');

    useEffect(() => {
        const fetchData = async () => {
            try {
                const res = await homeAPi.getEquipmentInspectionStatus(); // 修正 API 调用
                if (res && res.data) {
                    setData(res.data || []);
                }
            } catch (error) {
                console.error("查询设备点检状态失败:", error);
            }
        };

        fetchData();

        const socket = io(process.env.REACT_APP_SOCKET_URL, {
            transports: ['websocket'],
        });

        socket.on('connect', () => {
            console.log('Connected to WebSocket server');
        });

        socket.on('scan-success', (message) => {
            console.log('Received scan-success message, refetching data...');
            fetchData();
        });

        socket.on('disconnect', () => {
            console.log('Disconnected from WebSocket server');
        });

        return () => {
            socket.disconnect();
        };
    }, []);

    const handleFilterChange = (e) => {
        setFilterStatus(e.target.value);
    };

    const handleTypeFilterChange = (e) => {
        setFilterType(e.target.value);
    };

    const filteredData = data.filter(item => {
        const statusFilter = filterStatus === '全部' ||
            (filterStatus === '已点检' &&
                (item.equip_code.startsWith('BYQ') ? item.inspected_am && item.inspected_pm : item.inspected_am)) ||
            (filterStatus === '未点检' &&
                (item.equip_code.startsWith('BYQ') ? !item.inspected_am || !item.inspected_pm : !item.inspected_am));
        const typeFilter = filterType === '全部' || item.equip_type === filterType;
        return statusFilter && typeFilter; // 使用修正后的变量名
    });

    // 定义新的表格列结构
    const columns = [
        {
            title: '序号',
            dataIndex: 'seq',
            key: 'seq',
            width: 80,
            align: 'center',
            render: (text, record, index) => `${index + 1}`, // 自动生成序号
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
            dataIndex: 'inspection_status',
            key: 'inspection_status',
            render: (text, record) => {
                if (record.equip_code.startsWith('BYQ')) {
                    return (
                        <span>
                            <Tag color={record.inspected_am ? 'green' : 'red'}>上午</Tag>
                            <Tag color={record.inspected_pm ? 'green' : 'red'}>下午</Tag>
                        </span>
                    );
                } else {
                    return (
                        <Tag color={record.inspected_am ? 'green' : 'red'}>
                            {record.inspected_am ? '已点检' : '未点检'}
                        </Tag>
                    );
                }
            },
        },
    ];

    return (
        <div className="equipment-inspection-container">
            <h1 className="title-text">设备日常点检看板</h1>
            <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', flexWrap: 'wrap' }}>
                <div style={{ marginBottom: 8 }}>
                    <span style={{ marginBottom: 8 }}>设备类型</span>
                    <Radio.Group onChange={handleTypeFilterChange} value={filterType}>
                        <Radio value="全部">全部</Radio>
                        <Radio value="日常巡查">日常巡查</Radio>
                        <Radio value="温控巡查">温控巡查</Radio>
                        <Radio value="特种设备">特种设备</Radio>
                        <Radio value="安全用电">安全用电</Radio>
                    </Radio.Group>
                </div>
                <div>
                    <span style={{ marginRight: 8 }}>点检状态</span>
                    <Radio.Group onChange={handleFilterChange} value={filterStatus}>
                        <Radio value="全部">全部</Radio>
                        <Radio value="已点检">已点检</Radio>
                        <Radio value="未点检">未点检</Radio>
                    </Radio.Group>
                </div>
            </div>
            <Table
                columns={columns}
                dataSource={filteredData}
                rowKey="equip_code"
                // bordered={false}
                pagination={false}
                className="custom-table"
            />
        </div>
    );
};

export default EquipmentInspection;
