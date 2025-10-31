import React, { useState, useEffect } from 'react';
import { Table, message } from 'antd';
import {io} from "socket.io-client";

const EquipmentScheduleWithCategory = () => {
    const socket = "http://192.168.0.103:9020";
    const [dataSource, setDataSource] = useState([]);

    // 获取设备列表
    const fetchEquipments = async () => {
        const res = await fetch(socket + '/api/home/equipment/list');
        if (!res.ok) throw new Error('获取设备列表失败');
        return res.json();
    };

    // 获取维护计划列表
    const fetchPlans = async () => {
        const res = await fetch(socket + '/api/home/equipmentPlan/list');
        if (!res.ok) throw new Error('获取维护计划失败');
        return res.json();
    };

    const fetchData = async () => {
        try {
            const equipments = await fetchEquipments();
            const plans = await fetchPlans();

            // 建立设备ID到维护计划的映射
            const planMap = {};

                plans?.data?.forEach(p => {
                    if (!planMap[p.equipment_id]) planMap[p.equipment_id] = {};
                    planMap[p.equipment_id][p.day] = true; // 只记录存在即可
                });

            const sortedEquipments = [...(equipments.data || [])].sort((a, b) => {
                if (a.type !== b.type) return b.type.localeCompare(a.type);
                return b.line.localeCompare(a.line);
            });

            // 组装dataSource格式
            const data = sortedEquipments.map(e => ({
                key: e.id.toString(),
                name: e.name,
                line:e.line,
                category: e.type,
                count: e.qty || 1,
                days: planMap[e.id] || {}
            }));

            setDataSource(data);
        } catch (err) {
            message.error(err.message);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    // 点击格子切换状态，调用接口更新
    const toggleTask = async (recordKey, day) => {
        const device = dataSource.find(d => d.key === recordKey);
        if (!device) return;

        const hasTask = !!device.days?.[day];

        try {
            if (hasTask) {
                // 有计划 → 删除
                const res = await fetch(socket + '/api/home/equipmentPlan/delete', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        equipment_id: Number(recordKey),
                        day: Number(day)
                    })
                });
                if (!res.ok) throw new Error('删除失败');

                // 更新本地状态
                setDataSource(prev =>
                    prev.map(d => {
                        if (d.key === recordKey) {
                            const newDays = { ...d.days };
                            delete newDays[day];
                            return { ...d, days: newDays };
                        }
                        return d;
                    })
                );
            } else {
                // 没有计划 → 添加
                const res = await fetch(socket + '/api/home/equipmentPlan/add', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        equipment_id: Number(recordKey),
                        day: Number(day)
                    })
                });
                if (!res.ok) throw new Error('添加失败');

                // 更新本地状态
                setDataSource(prev =>
                    prev.map(d => {
                        if (d.key === recordKey) {
                            return {
                                ...d,
                                days: {
                                    ...d.days,
                                    [day]: true
                                }
                            };
                        }
                        return d;
                    })
                );
            }
        } catch (err) {
            message.error(err.message);
        }
    };


    // 分类合并计算
    const categoryLineRowSpanMap = (() => {
        const map = new Map();

        dataSource.forEach((item, idx) => {
            if (!map.has(item.category)) {
                map.set(item.category, new Map());
            }
            const lineMap = map.get(item.category);
            if (!lineMap.has(item.line)) {
                lineMap.set(item.line, []);
            }
            lineMap.get(item.line).push(idx);
        });

        return map;
    })();

    const baseColumns = [
        {
            title: '分类',
            dataIndex: 'category',
            width: 60,
            render: (_, record, index) => {
                const lineMap = categoryLineRowSpanMap.get(record.category);
                if (!lineMap) return { children: '', props: {} };

                // 所有 index
                const allIndexes = Array.from(lineMap.values()).flat();
                const firstIndex = allIndexes[0];
                const rowSpan = index === firstIndex ? allIndexes.length : 0;

                return {
                    children: record.category,
                    props: { rowSpan }
                };
            }
        },
        {
            title: '线别',
            dataIndex: 'line',
            width: 60,
            render: (_, record, index) => {
                const lineMap = categoryLineRowSpanMap.get(record.category);
                if (!lineMap) return { children: '', props: {} };

                const indexes = lineMap.get(record.line) || [];
                const firstIndex = indexes[0];
                const rowSpan = index === firstIndex ? indexes.length : 0;

                return {
                    children: record.line,
                    props: { rowSpan }
                };
            }
        },
        { title: '设备名称', dataIndex: 'name', width: 100 },
        { title: '数量', dataIndex: 'count', width: 50 }
    ];

    const dateColumns = Array.from({ length: 31 }, (_, i) => {
        const day = (i + 1).toString();
        return {
            title: day,
            dataIndex: day,
            width: 30,
            align: 'center',
            render: (_, record) => {
                const hasTask = record.days?.[day];
                return (
                    <div
                        onClick={() => toggleTask(record.key, day)}
                        style={{
                            backgroundColor: hasTask ? '#666' : 'transparent',
                            height: 30,
                            cursor: 'pointer',
                            margin: 0,
                            padding: 0,
                            border: 'none',
                            width: '100%',
                        }}
                        title={hasTask ? '有维护计划' : '无维护'}
                    />
                );
            }
        };
    });

    const columns = [...baseColumns, ...dateColumns];

    return (
        <div style={{ padding: 0,marginBottom: 20,marginRight: 10 }}>
            <div style={{ marginBottom: 12, fontSize: 14 }}>
      <span>
        <span
            style={{
                display: 'inline-block',
                backgroundColor: '#1890ff',
                width: 12,
                height: 12,
                marginRight: 4
            }}
        />
        点击格子切换维护完成状态，左边设备分类自动合并。
      </span>
            </div>
            <Table
                bordered
                size="small"
                columns={columns}
                dataSource={dataSource}
                pagination={false}
                scroll={{ x: 1800 }}
                rowClassName={() => 'editable-row'}
            />
        </div>
    );
};

export default EquipmentScheduleWithCategory;
