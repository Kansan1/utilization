import React, {useEffect, useMemo, useState} from "react";
import ColumnChart from "../../../../components/ColumnChart";
import BarChart from "../../../../components/BarChart";
import {homeAPi} from "../../../../api"; // 导入 LineChart 组件

const MaintenanceMonth = () => {

    const [lineDataMonth, setLineDataMonth] = useState([]); // 存储 API 数据
    const currentTime = new Date().toISOString().slice(0, 10); // 获取当前年月，格式化为 YYYY-MM

    useEffect(() => {
        async function fetchData() {
            try {
                const result = await homeAPi.getEquipmentRepairCompletionRate();
            //    console.log(result)
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

    const lineDataDay = useMemo(() =>
        lineDataMonth.map((item, index) => ({
            value: item.value, // value 从 1 开始递增
            type: item.name,
        })), [lineDataMonth]); // 依赖 `lineDataMonth`，只有数据更新时重新计算


    const chartConfig = {
        xAxis: {
            grid: {
                visible: false,
            },
            min: 0,
            max: 100,
        },
    };
    return (
        <div>
            {/* 将数据和标题传递给 LineChart 组件 */}
            <BarChart data={lineDataDay} title="设备维护完成率（月）" chartConfig={chartConfig} />
        </div>
    );
};

export default MaintenanceMonth;
