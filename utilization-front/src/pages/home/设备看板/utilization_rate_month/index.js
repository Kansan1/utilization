import React, { useState, useEffect, useMemo } from "react";
import LineChart from "../../../../components/LineChart";
import { homeAPi } from "../../../../api"; // 导入 API

const UtilizationRateDay = () => {
    const [lineDataMonth, setLineDataMonth] = useState([]); // 存储 API 数据
    const currentTime = new Date().toISOString().slice(0, 10); // 获取当前年月，格式化为 YYYY-MM

    useEffect(() => {
        async function fetchData() {
            try {
                const result = await homeAPi.getUtilizationAllList({ currentTime });
              //  console.log(result.data)
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
            valueAll: item.valueAll,
            valueActual: item.valueActual
        })), [lineDataMonth]); // 依赖 `lineDataMonth`，只有数据更新时重新计算

    // console.log(lineDataMonth)

    const chartConfig = {
        xAxis: {
            min: 0,
            max: 12,
        },
        // tooltip: {
        //     // 开启 tooltip
        //     visible: true,
        //     fields: ['type', 'value', 'rate', 'valueAll', 'valueActual'],
        //     shared: true, // 当鼠标悬停在某个点时，显示所有系列的数据
        //     showMarkers: false, // 隐藏 tooltip 中的标记点
        //     formatter: (item) => {
        //         // 这里 item 是一个包含当前数据项的对象
        //         // 根据字段 name 和 value 来显示每个数据点的信息
        //         return {
        //             name: item.type, // 显示设备名称或类型
        //             value: `${item.rate}%\n${item.valueAll}/${item.valueActual}`, // 格式化显示月份
        //        //     rate: `${item.rate}%`, // 格式化显示稼动率
        //        //      valueAll: item.valueAll, // 显示 valueAll
        //        //      valueActual: item.valueActual, // 显示 valueActual
        //         };
        //     },
        // },
    };

    return (
        <div>
            <LineChart data={lineDataDay} title="稼动率（月）" chartConfig={chartConfig} />
        </div>
    );
};

export default UtilizationRateDay;
