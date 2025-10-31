import React, { useEffect, useRef, useMemo } from "react";
import { Column } from "@antv/g2plot";

// 深度合并函数，确保 chartConfig 可以正确覆盖默认配置
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

const ColumnChart = ({ data, title, chartConfig = {} }) => {
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
      xField: "type",
      yField: "value",
      autoFit: true,
      appendPadding: [60, 10, 0, 0],
      background: { fill: "#0a1f44" }, // 深蓝色背景
      xAxis: {
        label: {
          style: { fontSize: 14, fill: "#b0c4de" }, // x 轴文字颜色
        },
      },
      yAxis: {
        grid: null, // 去掉横向网格线
        label: {
          formatter: (v) => `${v}%`,
          style: { fontSize: 14, fill: "#b0c4de" },
        },
      },
      tooltip: {
        formatter: (datum) => ({
          name: "完成率",
          value: `${datum.value}%`,
        }),
      },
      label: {
        position: "top",
        offsetY: -4,
        style: {
          fontSize: 14,
          fill: "#b0c4de",
        },
      },
      columnWidthRatio: 0.6,
      color: ["#4a90e2"], // 统一柱状颜色
      annotations: [
        {
          type: "text",
          position: ["-4%", "-15%"],
          content: title || "设备维修趋势（月）",
          style: {
            fontSize: 30,
            fontWeight: "bold",
            fill: "#b0c4de",
            textAlign: "left",
          },
        },
      ],
      legend: {
        position: "bottom",
        offsetY: 10,
        itemName: {
          style: { fontSize: 14, fill: "#b0c4de" },
        },
      },
    };

    return deepMerge(defaultConfig, prevChartConfig.current || {});
  }, [validData, title]);

  useEffect(() => {
    if (!chartRef.current || validData.length === 0) return;

    const chart = new Column(chartRef.current, config);
    chart.render();

    return () => chart.destroy();
  }, [validData, title, config]);

  return <div ref={chartRef} />;
};

export default ColumnChart;
