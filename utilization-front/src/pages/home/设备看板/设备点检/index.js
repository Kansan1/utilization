import React, { useEffect, useState } from "react";
import ColumnChart from "../../../../components/ColumnChart";
import { homeAPi } from "../../../../api";
import { io } from "socket.io-client";

const socket = io("http://192.168.0.103:3004");

const InspectionDay = () => {
    const [lineDataDay, setLineDataDay] = useState([]);

    // 前缀与类型映射
    const prefixMap = {
        PDG: "安全用电",
        SCX: "日常巡查",
        JXS: "日常巡查",
        AGV: "日常巡查",
        BYQ: "温控巡查",
        CQG: "特种设备",
        YSJ: "特种设备",
        GZJ: "特种设备"
    };

    const maxValueMap = {
        "安全用电": 29,
        "日常巡查": 17,
        "温控巡查": 34,
        "特种设备": 4
    };

    // 获取数据方法
    const fetchData = async () => {
        try {
            const result = await homeAPi.getInspection();
            
            // 统计各类型数量
            const countMap = {
                "安全用电": 0,
                "日常巡查": 0,
                "温控巡查": 0,
                "特种设备": 0
            };

            // 根据前缀分类计数
            result.forEach(item => {
                if (item.type && item.type.length >= 3) {
                    const prefix = item.type.substring(0, 3);
                    const typeName = prefixMap[prefix];
                    if (typeName) {
                        countMap[typeName]++;
                    }
                }
            });

            // 计算百分比
            const formattedData = Object.keys(countMap).map(type => {
                const actual = countMap[type] || 0;
                const max = maxValueMap[type] || 1;
                return {
                    type,
                    value: Math.min(Math.round((actual / max) * 100), 100)
                };
            });

            setLineDataDay(formattedData);
        } catch (error) {
            console.error('获取巡检数据失败:', error);
        }
    };

    useEffect(() => {
        fetchData();

        // 监听扫码成功事件
        socket.on("scan-success", () => {
            fetchData();
        });

        return () => {
            socket.off("scan-success");
        };
    }, []);

    const chartConfig = {
        yAxis: {
            min: 0,
            max: 100,
            label: {
                formatter: (v) => `${v}%`,
                style: { fontSize: 14 },
            },
            grid: {
                visible: false,
            },
        },
        label: {
            formatter: (item) => `${item.value}%`,
            position: "middle",
            style: {
                fontSize: 30,
            },
        },
    };

    return (
        <div>
            <ColumnChart data={lineDataDay} title="设备点检（日）" chartConfig={chartConfig} />
        </div>
    );
};

export default InspectionDay;