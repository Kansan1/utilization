import React, { useEffect, useState } from 'react';
import {Table, Button, Space, Modal, Form, Input, message, Select, InputNumber} from 'antd';
import dayjs from 'dayjs';
import {io} from "socket.io-client";

const { Option } = Select;

const EquipmentPlanSpecific = () => {

    const socket = "http://192.168.0.103:9020";

    const [weeklyData, setWeeklyData] = useState([]);
    const [dailyData, setDailyData] = useState([]);
    const [loading, setLoading] = useState(false);

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [form] = Form.useForm();
    const [editingItem, setEditingItem] = useState(null);
    const [isWeekly, setIsWeekly] = useState(true);

    const columns = [
        { title: '产线', dataIndex: 'line', key: 'line' },
        { title: '设备名称', dataIndex: 'equipment_name', key: 'equipment_name' },
        { title: '内容', dataIndex: 'content', key: 'content' },
        { title: '数量', dataIndex: 'shu', key: 'shu' },
        { title: '状态', dataIndex: 'state', key: 'state' },
        { title: '责任人', dataIndex: 'defender', key: 'defender' },
        {
            title: '操作',
            key: 'action',
            render: (_, record) => (
                <Space>
                    <Button size="small" onClick={() => showModal(record, isWeekly)}>编辑</Button>

                </Space>
            )
        }
    ];

    const fetchData = async () => {
        setLoading(true);
        try {
            const [weekRes, dayRes] = await Promise.all([
                fetch(socket + '/api/home/planSpecific/week').then(res => res.json()),
                fetch(socket + '/api/home/planSpecific/day').then(res => res.json())
            ]);
            setWeeklyData(weekRes.data || []);
            setDailyData(dayRes.data || []);
        } catch (err) {
            message.error('获取数据失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    const showModal = (record = null, isWeek = true) => {
        setIsWeekly(isWeek);
        setIsModalOpen(true);
        setEditingItem(record);
        if (record) {
            form.setFieldsValue(record);
        } else {
            form.resetFields();
        }
    };

    const handleOk = async () => {
        try {
            const values = await form.validateFields();

            const payload = {
                ...editingItem,
                ...values
            };

            const res = await fetch(socket + '/api/home/planSpecific/update', {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (!res.ok) throw new Error('保存失败');
            message.success('保存成功');
            setIsModalOpen(false);
            setWeeklyData([]);
            setDailyData([]);
            fetchData();
        } catch (err) {
            message.error(err.message);
        }
    };

    return (
        <div style={{ padding: 20 }}>
            <h2>周维护计划</h2>
            <Table
                rowKey="id"
                columns={columns}
                dataSource={weeklyData}
                loading={loading}
                pagination={false}
                bordered
                size="small"
                style={{ marginBottom: 40 }}
            />

            <h2>日维护计划</h2>
            <Table
                rowKey="id"
                columns={columns}
                dataSource={dailyData}
                loading={loading}
                pagination={false}
                bordered
                size="small"
            />

            <Modal
                title={editingItem ? '编辑计划' : '新增计划'}
                open={isModalOpen}
                onOk={handleOk}
                onCancel={() => setIsModalOpen(false)}
                destroyOnClose
            >
                <Form layout="vertical" form={form}>
                    <Form.Item name="equipment_name" label="设备名称">
                        <Input disabled />
                    </Form.Item>
                    <Form.Item name="line" label="产线" rules={[{ required: true }]}>
                        <Input disabled/>
                    </Form.Item>
                    <Form.Item name="content" label="内容" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="shu" label="数量" rules={[{ required: true }]}>
                        <InputNumber />
                    </Form.Item>
                    <Form.Item
                        name="state"
                        label="状态"
                        rules={[{ required: true }]}
                    >
                        <Select>
                            <Option value="待维护">待维护</Option>
                            <Option value="完成">完成</Option>
                        </Select>
                    </Form.Item>
                    <Form.Item name="defender" label="责任人" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default EquipmentPlanSpecific;
