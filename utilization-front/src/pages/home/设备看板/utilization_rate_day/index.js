import React, { useState, useEffect, useMemo } from "react";
import LineChart from "../../../../components/LineChart";
import { homeAPi } from "../../../../api"; // 导入 API

const UtilizationRateDay = () => {
    const [lineDataMonth, setLineDataMonth] = useState([]); // 存储 API 数据
    const currentTime = new Date().toISOString().slice(0, 10); // 获取当前年月，格式化为 YYYY-MM

    useEffect(() => {
        async function fetchData() {
            try {
                const result = await homeAPi.getUtilizationList({ currentTime });
                if (result && result.data) {
                    setLineDataMonth(result.data); // 设置数据
                } else {
                    setLineDataMonth([]); // 若没有数据，置空
                }
            } catch (error) {
                console.error("获取数据失败:", error);
                setLineDataMonth([]); // 错误处理，确保不会崩溃
            }
        }
        fetchData(); // 仅在组件加载时请求一次数据
    }, [currentTime]); // 依赖项为空，意味着只会在 `currentTime` 变化时请求一次

    // 计算每日数据
    const lineDataDay = useMemo(() =>
        lineDataMonth.map((item, index) => ({
            value: (index + 1).toString(), // value 从 1 开始递增
            type: "稼动率",
            rate: parseFloat(parseFloat((parseFloat(item.value || 0).toFixed(4)) * 100).toFixed(4)),
        })), [lineDataMonth]); // 依赖 `lineDataMonth`，只有数据更新时重新计算


 

    const chartConfig = {
        xAxis: {
            min: 0,
            max: 31,
        },
    };

    return (
        <div>
            <LineChart data={lineDataDay} title="稼动率（日）" chartConfig={chartConfig} />
        </div>
    );
};

export default UtilizationRateDay;
