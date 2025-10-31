import React, { useEffect, useRef, useMemo } from "react";
import { Line } from "@antv/g2plot";

// 深度合并函数，确保 chartConfig 继承默认值
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

const LineChart = ({ data, title, chartConfig = {} }) => {
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
      yField: "rate",
      seriesField: "type",
      smooth: true,
      autoFit: true,
      background: { fill: "#0a1f44" }, // 深蓝色背景
      appendPadding: [70, 20, 20, 20], // 增加左右边距

      xAxis: {
        min: 0,
        max: 30,
        label: {
          style: { fontSize: 14, fill: "#b0c4de" }, // x 轴文字
        },
      },
      yAxis: {
        min: 0,
        max: 100,
        grid: null, // 去掉网格线
        label: {
          formatter: (v) => `${v}%`,
          style: { fontSize: 14, fill: "#b0c4de" }, // y 轴文字
        },
      },
      label: {
        position: "top",
        offsetY: -4,
        style: { fill: "#b0c4de", fontSize: 14 },
        formatter: (item) => `${item.rate}%`,
      },
      point: {
        size: 4,
        shape: "circle",
        style: { fill: "#0a1f44", stroke: "#4a90e2", lineWidth: 2 },
      },
      color: ["#4a90e2", "#ff7f50", "#ffd700"], // 多条线不同颜色
      annotations: [
        {
          type: "text",
          position: ["-8%", "-25%"],
          content: title || "设备维修趋势（月）",
          style: { fontSize: 30, fontWeight: "bold", fill: "#b0c4de" },
        },
      ],
      legend: {
        visible:false,
        position: "bottom",
        offsetY: 10,
        itemName: {
          style: { fontSize: 14, fill: "#b0c4de" },
        },
      },
      tooltip: {
        showTitle: false,
        customContent: (title, items) => {
          const data = items?.[0]?.data;
          if (!data) return '';

          const { value, rate, valueAll, valueActual } = data;

          return `
    <div style="padding: 6px 10px; font-size: 14px; line-height: 1.6;">
  <div><b>${valueAll == null ? '日期' : '月份'}：</b>${value}${valueAll == null ? '日' : '月'}</div>
  <div><b>稼动率：</b>${rate ?? 'N/A'}%</div>
  ${valueAll != null ? `<div><b>计划工时：</b>${valueAll} 小时</div>` : ''}
  ${valueActual != null ? `<div><b>实际工时：</b>${valueActual} 小时</div>` : ''}
</div>
    `;
        },
      }
          };

    return deepMerge(defaultConfig, prevChartConfig.current || {});
  }, [validData, title]);

  useEffect(() => {
    if (!chartRef.current || validData.length === 0) return;
    const chart = new Line(chartRef.current, config);
    chart.render();
    return () => chart.destroy();
  }, [validData, title, config]);

  return <div ref={chartRef} />;
};

export default LineChart;
