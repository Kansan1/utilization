import React, { useEffect, useRef, useState } from 'react';
import { Table, Row, Col } from 'antd';
import './EquipmentMaintenance.css';
import ColumnChart from "../../../components/ColumnChart";
import { io } from "socket.io-client";

const maintainColumns = [
    { title: '线别', dataIndex: 'line', key: 'line', width: 40 },
    { title: '设备名称', dataIndex: 'equipment_name', key: 'equipment_name', width: 90 },
    { title: '维护内容', dataIndex: 'content', key: 'content', width: 150 },
    { title: '状态', dataIndex: 'state', key: 'state', width: 50 },
    { title: '维护人', dataIndex: 'defender', key: 'defender', width: 50 }
];

const repairColumns = [
    { title: '线别', dataIndex: 'line', key: 'line', width: 40 },
    { title: '设备名称', dataIndex: 'device_name', key: 'equipment_name', width: 90 },
    { title: '故障现象', dataIndex: 'fault', key: 'content', width: 150 },
    { title: '状态', dataIndex: 'state', key: 'state', width: 50 },
    { title: '维修人', dataIndex: 'repairer', key: 'repairer', width: 50 }
];

const earlyColumns = [
    {
        title: <span style={{ color: '#ffffff', fontSize: '16px', fontWeight: 'bold' }}>线别</span>,
        dataIndex: 'line',
        key: 'line',
        width: 40,
        render: text => <span style={{ color: '#ff4d4f', fontSize: '16px' }}>{text}</span>,
    },
    {
        title: <span style={{ color: '#ffffff', fontSize: '16px', fontWeight: 'bold' }}>设备名称</span>,
        dataIndex: 'equipment_name',
        key: 'equipment_name',
        width: 90,
        render: text => <span style={{ color: '#ff4d4f', fontSize: '16px' }}>{text}</span>,
    },
    {
        title: <span style={{ color: '#ffffff', fontSize: '16px', fontWeight: 'bold' }}>故障现象</span>,
        dataIndex: 'content',
        key: 'content',
        width: 150,
        render: text => <span style={{ color: '#ff4d4f', fontSize: '16px' }}>{text}</span>,
    },
    {
        title: <span style={{ color: '#ffffff', fontSize: '16px', fontWeight: 'bold' }}>故障次数</span>,
        dataIndex: 'fault_count',
        key: 'fault_count',
        width: 50,
        render: text => <span style={{ color: '#ff4d4f', fontSize: '16px' }}>{text}</span>,
    },
    {
        title: <span style={{ color: '#ffffff', fontSize: '16px', fontWeight: 'bold' }}>状态</span>,
        dataIndex: 'state',
        key: 'state',
        width: 50,
        render: text => <span style={{ color: '#ff4d4f', fontSize: '16px' }}>{text}</span>,
    },
    {
        title: <span style={{ color: '#ffffff', fontSize: '16px', fontWeight: 'bold' }}>维修人</span>,
        dataIndex: 'defender',
        key: 'defender',
        width: 50,
        render: text => <span style={{ color: '#ff4d4f', fontSize: '16px' }}>{text}</span>,
    }
];

export default function EquipmentMaintenance({ isFullscreen, currentTime }) {
    const baseURL = "http://192.168.0.103:9020";
    const socket = io(baseURL);

    const [maintainDayData, setMaintainDayData] = useState([]);
    const [maintainWeekData, setMaintainWeekData] = useState([]);
    const [repairDayData, setRepairDayData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [earlyData, setEarlyData] = useState([]);
    const [lineDataDay, setLineDataDay] = useState([]);

    const maintainWeekRef = useRef(null);
    const maintainDayRef = useRef(null);
    const repairRef = useRef(null);
    const earlyRef = useRef(null);

    // 获取数据方法
    const fetchData = async () => {
        setLoading(true);
        try {
            const [dayRes, weekRes, statsRes, repairRes, earlyRes] = await Promise.all([
                fetch(`${baseURL}/api/home/planSpecific/day`).then(res => res.json()),
                fetch(`${baseURL}/api/home/planSpecific/week`).then(res => res.json()),
                fetch(`${baseURL}/api/home/equipment/GetEquipmentTypeMaintainStats`).then(res => res.json()),
                fetch(`${baseURL}/api/home/dailyRepair/todayList`).then(res => res.json()),
                fetch(`${baseURL}/api/home/monthlyFaultFrequency/merged`).then(res => res.json()),
            ]);

            function fillEmptyRows(data, minLength, defaultRow = {}) {
                const filled = [...data];
                const missingCount = minLength - filled.length;
                if (missingCount > 0) {
                    for (let i = 0; i < missingCount; i++) {
                        filled.push({ id: `empty-${Date.now()}-${i}`, ...defaultRow });
                    }
                }
                return filled;
            }

            if (dayRes.data) setMaintainDayData(fillEmptyRows(dayRes.data, 5));
            if (weekRes.data) setMaintainWeekData(fillEmptyRows(weekRes.data, 5));
            if (repairRes.data) setRepairDayData(fillEmptyRows(repairRes.data, 5));
            if (earlyRes.data) setEarlyData(fillEmptyRows(earlyRes.data, 2));

            if (statsRes.data) {
                const formattedData = statsRes.data.flatMap(stat => ([
                    { type: stat.equipment_type, category: '最大值', value: stat.total_count - stat.completed_count, sumValue: stat.total_count },
                    { type: stat.equipment_type, category: '当前值', value: stat.completed_count },
                ]));
                setLineDataDay(formattedData);
            } else {
                setLineDataDay([]);
            }
        } catch (error) {
            console.error('获取巡检数据失败:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();

        // 监听 socket.io 事件
        socket.on("connect", () => {
            console.log("✅ Socket.IO 连接成功:", socket.id);
        });

        socket.on("Change", () => {
            console.log("📡 收到 Change 推送，刷新数据");
            fetchData();
        });

        // 每 10 分钟自动刷新
        const intervalId = setInterval(fetchData, 600000);

        // 自动滚动
        const scrollElements = [
            maintainWeekRef.current?.querySelector('.ant-table-body'),
            maintainDayRef.current?.querySelector('.ant-table-body'),
            repairRef.current?.querySelector('.ant-table-body'),
            earlyRef.current?.querySelector('.ant-table-body'),
        ].filter(Boolean);

        const scrollInterval = setInterval(() => {
            scrollElements.forEach(el => {
                const { scrollTop, scrollHeight, clientHeight } = el;
                if (scrollTop + clientHeight >= scrollHeight) {
                    el.scrollTop = 0;
                } else {
                    el.scrollTop += 30;
                }
            });
        }, 2000);

        return () => {
            clearInterval(intervalId);
            clearInterval(scrollInterval);
            socket.off("Change");
            socket.disconnect();
        };
    }, []);

    const chartConfig = {
        height: 300,
        isStack: true,
        seriesField: 'category',
        columnStyle: (datum) => {
            if (datum.category === '最大值') {
                return { fill: '#fff', stroke: '#e11428', lineWidth: 1.5, lineDash: [6, 2] };
            }
            return { fill: '#1e3974', stroke: '#5B8FF9', lineWidth: 0 };
        },
        yAxis: { min: 0, max: 50, label: { style: { fontSize: 14 } }, grid: { visible: false } },
        label: {
            position: 'middle',
            offsetY: 10,
            formatter: (datum) => {
                if (datum.category === '当前值') return datum.value === 0 ? '' : `${datum.value}`;
                if (datum.category === '最大值') return `${datum.sumValue}\n\n\n\n\n`;
                return '';
            },
            style: { fill: '#5B8FF9', fontSize: 20 },
        },
        color: ({ category }) => category === '当前值' ? '#033b53' : '#ffffff',
    };

    return (
        <div className={`content-body ${isFullscreen ? "fullscreen" : ""}`}>
            <div className="header-container">
                <div className="left-placeholder"></div>
                <div className="content-title">设备维修维护管理看板</div>
                <div className="time-display">{currentTime}</div>
            </div>

            <Row gutter={16}>
                {/* 左列 */}
                <Col span={12}>
                    <Row gutter={[0, 16]}>
                        <Col span={24}>
                            <div ref={maintainWeekRef}>
                                <Table
                                    title={() => <div style={{ fontSize: 24, fontWeight: 'bold', color: '#ffffff' }}>周维护计划</div>}
                                    rowClassName="custom-row2"
                                    dataSource={maintainWeekData}
                                    columns={maintainColumns}
                                    pagination={false}
                                    size="small"
                                    loading={loading}
                                    scroll={{ y: 370 }}
                                    rowKey="id"
                                />
                            </div>
                        </Col>
                        <Col span={24}>
                            <div ref={maintainDayRef}>
                                <Table
                                    title={() => <div style={{ fontSize: 24, fontWeight: 'bold', color: '#ffffff' }}>日维护计划</div>}
                                    rowClassName="custom-row2"
                                    dataSource={maintainDayData}
                                    columns={maintainColumns}
                                    pagination={false}
                                    size="small"
                                    loading={loading}
                                    scroll={{ y: 370 }}
                                    rowKey="id"
                                />
                            </div>
                        </Col>
                    </Row>
                </Col>

                {/* 右列 */}
                <Col span={12}>
                    <Row gutter={[0, 16]}>
                        <Col span={24}>
                            <ColumnChart data={lineDataDay} title="月维护推进柱形图（单位：台）" chartConfig={chartConfig} />
                        </Col>
                        <Col span={24}>
                            <div ref={repairRef}>
                                <Table
                                    title={() => <div style={{ fontSize: 24, fontWeight: 'bold', color: '#ffffff' }}>日维修任务</div>}
                                    rowClassName="custom-row"
                                    dataSource={repairDayData}
                                    columns={repairColumns}
                                    pagination={false}
                                    size="small"
                                    scroll={{ y: 300 }}
                                    rowKey="id"
                                />
                            </div>
                        </Col>
                        <Col span={24}>
                            <div ref={earlyRef}>
                                <Table
                                    title={() => <div style={{ fontSize: 24, fontWeight: 'bold', color: '#ffffff' }}>故障频繁设备预警（每月维修次数≥5次）</div>}
                                    rowClassName="custom-row"
                                    dataSource={earlyData}
                                    columns={earlyColumns}
                                    pagination={false}
                                    bordered
                                    size="small"
                                    scroll={{ y: 120 }}
                                    rowKey="id"
                                />
                            </div>
                        </Col>
                    </Row>
                </Col>
            </Row>
        </div>
    );
}
