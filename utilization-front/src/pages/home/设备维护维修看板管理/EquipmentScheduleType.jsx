import React, { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, InputNumber, Space, Popconfirm, message } from 'antd';
import { homeAPi } from '../../../api'; // 导入封装好的API

function EquipmentScheduleType() {
    const [devices, setDevices] = useState([]);
    const [editingDevice, setEditingDevice] = useState(null);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const [form] = Form.useForm();

    const fetchDevices = async () => {
        setLoading(true);
        try {
            // 使用封装好的 homeAPi
            const res = await homeAPi.getEquipmentList();
            setDevices(res.data || []);
        } catch (err) {
            message.error('获取设备列表失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchDevices();
    }, []);

    const showModal = (record = null) => {
        setEditingDevice(record);
        if (record) {
            form.setFieldsValue(record);
        } else {
            form.resetFields();
        }
        setIsModalOpen(true);
    };

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            if (editingDevice) {
                // 使用封装好的 homeAPi
                await homeAPi.updateEquipment({ ...editingDevice, ...values });
                message.success('更新成功');
            } else {
                // 使用封装好的 homeAPi
                await homeAPi.addEquipment(values);
                message.success('添加成功');
            }
            setIsModalOpen(false);
            fetchDevices();
        } catch (err) {
            message.error('保存失败');
        }
    };

    const handleDelete = async (id) => {
        try {
            // 使用封装好的 homeAPi
            await homeAPi.deleteEquipment(id);
            message.success('删除成功');
            fetchDevices();
        } catch (err) {
            message.error('删除失败');
        }
    };

    const columns = [
        { title: '设备种类', dataIndex: 'type', key: 'type' },
        { title: '线别', dataIndex: 'line', key: 'line' },
        { title: '设备名称', dataIndex: 'name', key: 'name' },
        { title: '数量', dataIndex: 'qty', key: 'qty' },
        {
            title: '操作',
            key: 'action',
            render: (_, record) => (
                <Space>
                    <Button size="small" onClick={() => showModal(record)}>编辑</Button>
                    <Popconfirm
                        title="确认删除这个设备？"
                        onConfirm={() => handleDelete(record.id)}
                    >
                        <Button size="small" danger>删除</Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <div>
            <h2>维护设备属性</h2>
            <Button type="primary" onClick={() => showModal()} style={{ marginBottom: 16 }}>
                新增设备
            </Button>
            <Table
                rowKey="id"
                dataSource={devices}
                columns={columns}
                pagination={false}
                loading={loading}
            />

            <Modal
                title={editingDevice ? '编辑设备' : '新增设备'}
                open={isModalOpen}
                onOk={handleOk}
                onCancel={() => setIsModalOpen(false)}
                destroyOnClose
            >
                <Form form={form} layout="vertical">
                    <Form.Item name="type" label="设备种类" rules={[{ required: true, message: '请输入种类' }]}>
                        <Input placeholder="如：生产设备" />
                    </Form.Item>
                    <Form.Item name="line" label="线别">
                        <Input placeholder="如：E线" />
                    </Form.Item>
                    <Form.Item name="name" label="设备名称" rules={[{ required: true, message: '请输入名称' }]}>
                        <Input placeholder="如：流水线A" />
                    </Form.Item>
                    <Form.Item name="qty" label="数量" rules={[{ required: true, message: '请输入数量' }]}>
                        <InputNumber min={1} style={{ width: '100%' }} />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
}

export default EquipmentScheduleType;