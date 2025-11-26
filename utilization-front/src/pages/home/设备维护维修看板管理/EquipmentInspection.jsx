import React, { useEffect, useState } from 'react';
import { message, Spin, Table, Tag, Radio } from 'antd';

// 此常量用于HTTP请求的基地址，以保持与其他组件的编码风格一致
// const socket = "http://localhost:9020";
// const socket = "http://192.168.0.103:9020";
const socket=process.env.REACT_APP_API_URL;

const EquipmentInspection = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [filterStatus, setFilterStatus] = useState('全部');
    const [filterType, setFilterType] = useState('全部');

    const fetchData = async () => {
        setLoading(true);
        try {
            // 修正了 API 请求的 URL
            const response = await fetch(`${socket}/api/home/equipment/inspection-status`);
            // 首先获取文本格式的响应，以便调试
            const text = await response.text();
            console.log('Raw response from backend:', text); // 打印原始响应以进行调试

            if (!response.ok) {
                // 如果响应状态码不是 2xx，则抛出错误，并将服务器返回的信息包含在内
                throw new Error(`网络响应错误。服务器返回: ${text}`);
            }

            // 尝试将文本解析为 JSON
            try {
                const result = JSON.parse(text);
                if (result && result.data) {
                    setData(result.data);
                } else {
                    setData([]);
                    console.warn("从后端接收的数据格式不符合预期:", result);
                }
            } catch (jsonError) {
                // 如果 JSON.parse 失败，则会捕获此错误
                console.error("无法将响应解析为 JSON:", jsonError);
                // 原始文本已被打印，这是调试此类问题的关键
                throw new Error("从服务器收到了非 JSON 格式的响应。");
            }

        } catch (error) {
            console.error("查询设备点检状态失败:", error);
            message.error(`查询设备点检状态失败: ${error.message}`);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
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
        return statusFilter && typeFilter;
    });

    const columns = [
        {
            title: '序号',
            dataIndex: 'seq',
            key: 'seq',
            width: 80,
            align: 'center',
            render: (text, record, index) => `${index + 1}`,
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
                    <span style={{ marginRight: 8 }}>设备类型:</span>
                    <Radio.Group onChange={handleTypeFilterChange} value={filterType}>
                        <Radio value="全部">全部</Radio>
                        <Radio value="日常巡查">日常巡查</Radio>
                        <Radio value="温控巡查">温控巡查</Radio>
                        <Radio value="特种设备">特种设备</Radio>
                        <Radio value="安全用电">安全用电</Radio>
                    </Radio.Group>
                </div>
                <div>
                    <span style={{ marginRight: 8 }}>点检状态:</span>
                    <Radio.Group onChange={handleFilterChange} value={filterStatus}>
                        <Radio value="全部">全部</Radio>
                        <Radio value="已点检">已点检</Radio>
                        <Radio value="未点检">未点检</Radio>
                    </Radio.Group>
                </div>
            </div>
            <Spin spinning={loading}>
                <Table
                    columns={columns}
                    dataSource={filteredData}
                    rowKey="equip_code"
                    pagination={false}
                    className="custom-table"
                />
            </Spin>
        </div>
    );
};

export default EquipmentInspection;