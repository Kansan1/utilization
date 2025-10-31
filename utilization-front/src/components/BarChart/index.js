import React, { useEffect, useRef, useMemo } from "react";
import { Bar } from "@antv/g2plot";

const deepMerge = (target, source) => {
    if (typeof target !== "object" || typeof source !== "object") return source;
    for (const key in source) {
        if (source[key] && typeof source[key] === "object") {
            target[key] = deepMerge(target[key] || {}, source[key]);
        } else {
            target[key] = source[key];
        }
    }
    return target;
};

const BarChart = ({ data, title, chartConfig = {} }) => {
    const validData = useMemo(() => {
        return Array.isArray(data)
            ? data.map((item) => ({
                ...item,
                rate: typeof item.rate === "number" ? item.rate : 0,
            }))
            : [];
    }, [data]);

    const chartRef = useRef(null);
    const prevChartConfig = useRef(chartConfig);

    useEffect(() => {
        prevChartConfig.current = chartConfig;
    }, [chartConfig]);

    const config = useMemo(() => {
        const defaultConfig = {
            data: validData,
            xField: "value",
            yField: "type",
            color: ["#4a90e2"], // 柱状颜色
            background: { fill: "#0a1f44" }, // 深蓝色背景
            colorField: "type",
            autoFit: true,
            appendPadding: [60, 0, 0, 0],
            xAxis: {
                grid: null, // 去掉横线
                label: {
                    formatter: (v) => `${v}%`,
                    style: { fontSize: 14, fill: "#b0c4de" }, // 文字颜色
                },
            },
            yAxis: {
                label: { style: { fontSize: 14, fill: "#b0c4de" } },
            },
            legend: {
                position: "top",
                itemName: { style: { fill: "#b0c4de" } }, // 图例颜色
            },
            label: {
                position: "middle",
                style: { fontSize: 14, fill: "#b0c4de" },
                formatter: (item) => `${item.value}%`,
            },
            columnWidthRatio: 0.6,
            annotations: [
                {
                    type: "text",
                    position: ["0%", "0%"],
                    content: title || "设备维护（月）",
                    offsetX: -60,
                    offsetY: -40,
                    style: {
                        fontSize: 30,
                        fontWeight: "bold",
                        fill: "#b0c4de",
                        textAlign: "left",
                    },
                },
            ],
            tooltip: {
                formatter: (datum) => ({
                    name: "完成率",
                    value: `${datum.value}%`,
                }),
            },
        };

        // **深度合并外部 `chartConfig` 和默认 `defaultConfig`**
        return deepMerge(defaultConfig, prevChartConfig.current || {});
    }, [validData, title]);

    useEffect(() => {
        if (!chartRef.current || validData.length === 0) return;

        const chart = new Bar(chartRef.current, config);
        chart.render();

        return () => chart.destroy();
    }, [validData, title, config]);

    return <div ref={chartRef} />;
};

export default BarChart;
