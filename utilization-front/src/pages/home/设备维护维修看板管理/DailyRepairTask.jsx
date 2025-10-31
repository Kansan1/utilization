import React, { useEffect, useState } from 'react';
import {Table, Button, Modal, Form, Input, message, Space, Popconfirm} from 'antd';
import dayjs from 'dayjs';
import {io} from "socket.io-client";

const DailyRepairTask = () => {
    const socket = "http://192.168.0.103:9020";
    // --- 日维修任务状态 ---
    const [data, setData] = useState([]);
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingTask, setEditingTask] = useState(null);

    // --- 故障频繁记录状态 ---
    const [freqData, setFreqData] = useState([]);
    const [freqLoading, setFreqLoading] = useState(false);
    const [isFreqModalOpen, setIsFreqModalOpen] = useState(false);
    const [freqEditing, setFreqEditing] = useState(null);
    const [freqForm] = Form.useForm();

    // 获取日维修任务数据
    const fetchData = async () => {
        setLoading(true);
        try {
            const res = await fetch(socket + '/api/home/dailyRepair/todayList');
            const result = await res.json();
            if (result.code === 200) {
                setData(result.data || []);
            } else {
                message.error('获取日维修任务失败');
            }
        } catch (err) {
            message.error('请求失败');
        } finally {
            setLoading(false);
        }
    };

    // 获取故障频繁记录数据
    const fetchFreqData = async () => {
        setFreqLoading(true);
        try {
            const res = await fetch(socket + '/api/home/monthlyFaultFrequency/merged');
            const result = await res.json();
            if (result.code === 200) {
                setFreqData(result.data || []);
            } else {
                message.error('获取故障频繁记录失败');
            }
        } catch (err) {
            message.error('请求失败');
        } finally {
            setFreqLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
        fetchFreqData();
    }, []);

    // --- 日维修任务操作 ---
    const handleDelete = async (id) => {
        try {
            const res = await fetch(socket + `/api/home/dailyRepair/delete/${id}`, {
                method: 'DELETE',
            });
            if (!res.ok) throw new Error('删除失败');
            message.success('删除成功');
            fetchData();
        } catch {
            message.error('删除失败');
        }
    };

    const handleEdit = (record) => {
        setEditingTask(record);
        form.setFieldsValue({
            ...record,
            report_time: record.report_time ? dayjs(record.report_time) : null,
        });
        setIsModalOpen(true);
    };

    const handleAdd = () => {
        setEditingTask(null);
        form.resetFields();
        setIsModalOpen(true);
    };

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            const payload = { ...editingTask, ...values };
            // 格式化报修时间为字符串
            if (values.report_time) {
                payload.report_time = values.report_time.format('YYYY-MM-DD HH:mm:ss');
            }

            const url = editingTask
                ? (socket + '/api/home/dailyRepair/update')
                : (socket + '/api/home/dailyRepair/add');

            const method = editingTask ? 'PUT' : 'POST';

            const res = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });

            const result = await res.json();
            if (result.code === 200) {
                message.success(editingTask ? '更新成功' : '添加成功');
                setIsModalOpen(false);
                fetchData();
            } else {
                message.error(result.message || '操作失败');
            }
        } catch (err) {
            message.error('提交失败');
        }
    };

    // --- 故障频繁记录操作 ---
    const showFreqModal = (record = null) => {
        setFreqEditing(record);
        setIsFreqModalOpen(true);
        if (record) {
            freqForm.setFieldsValue(record);
        } else {
            freqForm.resetFields();
        }
    };

    const handleFreqOk = async () => {
        try {
            const values = await freqForm.validateFields();
            const payload = { ...freqEditing, ...values };

            const url = freqEditing
                ? (socket + '/api/home/monthlyFaultFrequency/update')
                : (socket + '/api/home/monthlyFaultFrequency/add');

            const method = freqEditing ? 'PUT' : 'POST';

            const res = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });

            const result = await res.json();
            if (result.code === 200) {
                message.success(freqEditing ? '更新成功' : '添加成功');
                setIsFreqModalOpen(false);
                fetchFreqData();
            } else {
                message.error(result.message || '操作失败');
            }
        } catch (err) {
            message.error('提交失败');
        }
    };

    const deleteFreq = async (id) => {
        try {
            const res = await fetch(socket + `/api/home/monthlyFaultFrequency/delete/${id}`, {
                method: 'DELETE',
            });
            if (!res.ok) throw new Error('删除失败');
            message.success('删除成功');
            fetchFreqData();
        } catch (err) {
            message.error('删除失败');
        }
    };

    // 日维修任务列
    const columns = [
        { title: '线别', dataIndex: 'line', key: 'line' },
        { title: '设备名称', dataIndex: 'device_name', key: 'device_name' },
        { title: '故障现象', dataIndex: 'fault', key: 'fault' },
        { title: '状态', dataIndex: 'state', key: 'state' },
        { title: '维修人', dataIndex: 'repairer', key: 'repairer' },
        {
            title: '操作',
            key: 'action',
            render: (_, record) => (
                <Space>
                    <Button size="small" onClick={() => handleEdit(record)}>
                        编辑
                    </Button>
                    <Popconfirm
                        title="确认删除这个计划？"
                        onConfirm={() => handleDelete(record.id)}
                    >
                    <Button size="small" danger >
                        删除
                    </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    // 故障频繁记录列
    const freqColumns = [
        { title: '线别', dataIndex: 'line', key: 'line' },
        { title: '设备名称', dataIndex: 'equipment_name', key: 'equipment_name' },
        { title: '故障现象', dataIndex: 'fault', key: 'fault' },
        { title: '故障次数', dataIndex: 'fault_count', key: 'fault_count' },
        { title: '状态', dataIndex: 'state', key: 'state' },
        { title: '负责人', dataIndex: 'repairer', key: 'repairer' },
        {
            title: '操作',
            key: 'action',
            render: (_, record) => (
                <Space>
                    <Button size="small" onClick={() => showFreqModal(record)}>
                        编辑
                    </Button>
                </Space>
            ),
        },
    ];

    return (
        <div style={{ padding: 20 }}>
            {/* 日维修任务 */}
            <h2>日维修任务</h2>
            <Button type="primary" onClick={handleAdd} style={{ marginBottom: 10 }}>
                新增任务
            </Button>
            <Table  pagination={false}  rowKey="id" columns={columns} dataSource={data} loading={loading} bordered size="small" />

            <Modal
                title={editingTask ? '编辑任务' : '新增任务'}
                open={isModalOpen}
                onOk={handleOk}
                onCancel={() => setIsModalOpen(false)}
                destroyOnClose

            >
                <Form layout="vertical" form={form}>
                    <Form.Item name="line" label="线别" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="device_name" label="设备名称" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="fault" label="故障现象">
                        <Input />
                    </Form.Item>
                    {/*<Form.Item name="report_time" label="报修时间" rules={[{ required: true }]}>*/}
                    {/*    <DatePicker showTime format="YYYY-MM-DD HH:mm:ss" style={{ width: '100%' }} />*/}
                    {/*</Form.Item>*/}
                    <Form.Item name="state" label="状态">
                        <Input />
                    </Form.Item>
                    <Form.Item name="repairer" label="维修人">
                        <Input />
                    </Form.Item>
                </Form>
            </Modal>

            {/* 故障频繁记录 */}
            <div style={{ marginTop: 40 }}>
                <h2>故障频繁记录</h2>
                {/*<Button type="primary" onClick={() => showFreqModal()} style={{ marginBottom: 10 }}>*/}
                {/*    新增记录*/}
                {/*</Button>*/}
                <Table
                    rowKey="id"
                    columns={freqColumns}
                    dataSource={freqData}
                    loading={freqLoading}
                    bordered
                    size="small"
                    pagination={false}
                />

                <Modal
                    title={freqEditing ? '编辑故障频繁记录' : '新增故障频繁记录'}
                    open={isFreqModalOpen}
                    onOk={handleFreqOk}
                    onCancel={() => setIsFreqModalOpen(false)}
                    destroyOnClose
                >
                    <Form layout="vertical" form={freqForm}>
                        <Form.Item name="line" label="线别" rules={[{ required: true }]}>
                            <Input disabled/>
                        </Form.Item>
                        <Form.Item name="equipment_name" label="设备名称" rules={[{ required: true }]}>
                            <Input disabled />
                        </Form.Item>
                        <Form.Item name="fault" label="故障现象" rules={[{ required: true }]}>
                            <Input />
                        </Form.Item>
                        {/*<Form.Item*/}
                        {/*    name="fault_count"*/}
                        {/*    label="故障次数"*/}
                        {/*    rules={[{ required: true, type: 'number', min: 0 }]}*/}
                        {/*>*/}
                        {/*    <InputNumber style={{ width: '100%' }} />*/}
                        {/*</Form.Item>*/}
                        <Form.Item name="state" label="状态" rules={[{ required: true }]}>
                            <Input />
                        </Form.Item>
                        <Form.Item name="repairer" label="负责人" rules={[{ required: true }]}>
                            <Input />
                        </Form.Item>
                    </Form>
                </Modal>
            </div>
        </div>
    );
};

export default DailyRepairTask;
