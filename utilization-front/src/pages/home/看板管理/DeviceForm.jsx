import React, { useEffect } from 'react';
import { Form, InputNumber, Button, Card, Row, Col, message } from 'antd';
import {homeAPi as api, homeAPi} from "../../../api";

const fixedDevices = ['A线', 'B线', 'C线', 'D线', 'E线', 'G线', '机械手', 'AGV小车', '空压机', '其他设备'];

export default function DeviceForm() {
    const [repairForm] = Form.useForm();
    const [completionForm] = Form.useForm();

    const handleRepairSubmit = async (values) => {
        const data = fixedDevices.map((key) => ({
            name: key,
            value: values[`repairCount_${key}`] ?? 0
        }));
        try {
            await homeAPi.addOrUpdateEquipmentRepairCompletionTimes(data);
            message.success('维修次数已提交！');
        } catch (err) {
            message.error('提交失败，请重试！');
        }
    };

    const handleCompletionSubmit = async (values) => {
        const data = fixedDevices.map((key) => ({
            name: key,
            value: values[`completionRate_${key}`] ?? 0
        }));
        try {
            await homeAPi.addOrUpdateEquipmentRepairCompletionRate(data);
            message.success('维护完成率已提交！');
        } catch (err) {
            message.error('提交失败，请重试！');
        }
    };

    useEffect(() => {
        async function fetchRepairTimes() {
            try {
                const result = await homeAPi.getEquipmentRepairCompletionTimes();
                const initial = {};
                result.data.forEach(item => {
                    initial[`repairCount_${item.name}`] = item.value;
                });
                repairForm.setFieldsValue(initial);
            } catch (error) {
                console.log(error);
                message.error('加载设备维修次数失败！');
            }
        }

        fetchRepairTimes();
    }, [repairForm]);

    useEffect(() => {
        async function fetchCompletionRate() {
            try {
                const result = await homeAPi.getEquipmentRepairCompletionRate();

                const initial = {};
                result.data.forEach(item => {
                    initial[`completionRate_${item.name}`] = item.value ?? 0;
                });
                completionForm.setFieldsValue(initial);
            } catch (error) {
                console.log(error);
                message.error('加载设备维护完成率失败！');
            }
        }

        fetchCompletionRate();
    }, [completionForm]);


    return (
        <div style={{ minHeight: '100vh', backgroundColor: '#0a192f', padding: '0px' }}>
            <div style={{ maxWidth: 1000, margin: '0 auto', padding: '16px 0', textAlign: 'left' }}>
                <Button
                    type="primary"
                    onClick={handleExport} // 你自己的导出函数
                    style={{
                        backgroundColor: '#1e3a8a',
                        borderColor: '#1e3a8a',
                        color: '#ffffff',
                    }}
                >
                    下载数据
                </Button>
            </div>
            <div style={{ backgroundColor: '#0a192f', minHeight: '100%', width: '100%' }}>
                <Card
                    title={<span style={{ color: '#ccd6f6' }}>本月设备维修数据填写</span>}
                    style={{ maxWidth: 1000, margin: '0 auto', backgroundColor: '#112240' }}
                    headStyle={{ borderBottom: '1px solid #233554' }}
                    bodyStyle={{ backgroundColor: '#112240' }}
                >
                    <Row gutter={40}>
                        <Col span={12}>
                            <Card title={<span style={{ color: '#ccd6f6' }}>维修次数</span>}
                                  style={{ backgroundColor: '#233554' }}
                                  headStyle={{ borderBottom: '1px solid #334567' }}
                                  bodyStyle={{ backgroundColor: '#233554' }}>
                                <Form form={repairForm} layout="vertical" onFinish={handleRepairSubmit} colon={false}>
                                    {fixedDevices.map((name) => (
                                        <Form.Item
                                            key={`repairCount_${name}`}
                                            label={<span style={{ color: '#ccd6f6' }}>{name}</span>}
                                            name={`repairCount_${name}`}
                                            rules={[{ required: true, message: '请输入维修次数' }]}
                                        >
                                            <InputNumber min={0} style={{ width: '100%' }} defaultValue={0} placeholder="请输入维修次数" />
                                        </Form.Item>
                                    ))}
                                    <Form.Item>
                                        <Button type="primary" htmlType="submit" block>提交维修次数</Button>
                                    </Form.Item>
                                </Form>
                            </Card>
                        </Col>

                        <Col span={12}>
                            <Card title={<span style={{ color: '#ccd6f6' }}>维护完成率</span>}
                                  style={{ backgroundColor: '#233554' }}
                                  headStyle={{ borderBottom: '1px solid #334567' }}
                                  bodyStyle={{ backgroundColor: '#233554' }}>
                                <Form form={completionForm} layout="vertical" onFinish={handleCompletionSubmit} colon={false}>
                                    {fixedDevices.map((name) => (
                                        <Form.Item
                                            key={`completionRate_${name}`}
                                            label={<span style={{ color: '#ccd6f6' }}>{name}</span>}
                                            name={`completionRate_${name}`}
                                            rules={[{ required: true, message: '请输入维护完成率' }]}
                                        >
                                            <InputNumber min={0} max={100} style={{ width: '100%' }} defaultValue={0} placeholder="请输入维护完成率" />
                                        </Form.Item>
                                    ))}
                                    <Form.Item>
                                        <Button type="primary" htmlType="submit" block>提交维护完成率</Button>
                                    </Form.Item>
                                </Form>
                            </Card>
                        </Col>
                    </Row>
                </Card>
            </div>
        </div>
    );
}

const handleExport = async () => {
    // 这里可以是调用接口、导出 CSV、Excel 等逻辑

    try {
        const response = await api.downLoadDate()

        // 🟡 打印 headers 看实际返回内容
        console.log('响应头:', response.headers);

        let fileName = '导出数据.xlsx';
        const disposition = response.headers['content-disposition'];
        if (disposition) {
            const match = disposition.match(/filename\*?=UTF-8''(.+)/);
            if (match && match[1]) {
                fileName = decodeURIComponent(match[1]);
            }
        }

        const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = fileName;
        link.click();
        window.URL.revokeObjectURL(url);
    } catch (err) {
        console.error('下载失败:', err);
    }
};
