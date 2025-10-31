package excel

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
	"utilization-backend/src/api/dao"
	"utilization-backend/src/api/v1/home"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// 定义常量
const (
	sheetName = "设备看板数据"
	fileName  = "设备看板数据.xlsx"
)

// 定义设备类型
var equipmentTypes = []string{
	"A线", "B线", "C线", "D线", "E线", "G线",
	"机械手", "AGV小车", "空压机", "其他设备",
}

// DownLoadExcel 导出Excel文件
func DownLoadExcel(ctx *gin.Context) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("关闭Excel文件失败: %v", err)
		}
	}()

	// 设置工作表
	f.SetSheetName("Sheet1", sheetName)

	// 设置表头
	if err := setupHeaders(f); err != nil {
		log.Printf("设置表头失败: %v", err)
		ctx.JSON(500, gin.H{"error": "设置表头失败"})
		return
	}

	// 填充数据
	if err := fillData(f); err != nil {
		log.Printf("填充数据失败: %v", err)
		ctx.JSON(500, gin.H{"error": "填充数据失败"})
		return
	}

	// 自动调整所有列的宽度
	cols, err := f.GetCols(sheetName)
	if err != nil {
		log.Printf("获取列信息失败: %v", err)
	} else {
		for idx, col := range cols {
			largestWidth := 0
			for _, rowCell := range col {
				cellWidth := len(rowCell)
				if cellWidth > largestWidth {
					largestWidth = cellWidth
				}
			}
			// 设置列宽，每个字符大约1.2个单位宽度
			colWidth := float64(largestWidth) * 1.2
			if colWidth < 10 {
				colWidth = 10 // 设置最小列宽
			}
			colName, err := excelize.ColumnNumberToName(idx + 1)
			if err != nil {
				log.Printf("转换列名失败: %v", err)
				continue
			}
			if err := f.SetColWidth(sheetName, colName, colName, colWidth); err != nil {
				log.Printf("设置列宽失败: %v", err)
			}
		}
	}

	// 设置响应头
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(fileName)))
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")

	// 写入响应
	if err := f.Write(ctx.Writer); err != nil {
		log.Printf("写入Excel文件失败: %v", err)
		ctx.JSON(500, gin.H{"error": "写入Excel文件失败"})
		return
	}
}

// setupHeaders 设置Excel表头
func setupHeaders(f *excelize.File) error {
	// 合并单元格
	mergeCells := map[string]string{
		"B1": "D1", // 稼动率
		"E1": "N1", // 设备维护完成率
		"O1": "X1", // 设备维修次数
	}

	for start, end := range mergeCells {
		if err := f.MergeCell(sheetName, start, end); err != nil {
			return fmt.Errorf("合并单元格失败 %s-%s: %v", start, end, err)
		}
	}

	// 设置表头标题
	headers := map[string]string{
		"B1": "稼动率",
		"E1": "设备维护完成率",
		"O1": "设备维修次数",
	}

	for cell, value := range headers {
		if err := f.SetCellValue(sheetName, cell, value); err != nil {
			return fmt.Errorf("设置单元格 %s 失败: %v", cell, err)
		}
	}

	// 设置月份
	months := []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"}
	for i, month := range months {
		if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+3), month); err != nil {
			return fmt.Errorf("设置月份 %s 失败: %v", month, err)
		}
	}

	// 设置列标题
	columnHeaders := map[string]string{
		"B2": "实际工时",
		"C2": "总工时",
		"D2": "实际工时/总工时",
	}

	// 添加设备维护完成率和维修次数的列标题
	for i, equipment := range equipmentTypes {
		columnHeaders[fmt.Sprintf("%c2", 'E'+i)] = equipment
		columnHeaders[fmt.Sprintf("%c2", 'O'+i)] = equipment
	}

	for cell, value := range columnHeaders {
		if err := f.SetCellValue(sheetName, cell, value); err != nil {
			return fmt.Errorf("设置列标题 %s 失败: %v", cell, err)
		}
	}

	return nil
}

// fillData 填充Excel数据
func fillData(f *excelize.File) error {
	currentYear := time.Now().Year()

	for i := 1; i <= 12; i++ {
		currentTime := fmt.Sprintf("%d-%02d", currentYear, i)
		year, month := home.CurrentTimeToYearMonth(currentTime)

		// 获取稼动率数据
		utilizationWorkingHourList, utilizationExceptionList := dao.GetUtilizationListByMonth(year, month)
		valueAll, valueActual, valueRate := home.HandleAllUtilizationList(utilizationWorkingHourList, utilizationExceptionList)
		if err := f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), strconv.FormatFloat(float64(valueRate*100), 'f', 5, 64)+"%"); err != nil {
			return fmt.Errorf("设置稼动率数据失败: %v", err)
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), valueAll); err != nil {
			return fmt.Errorf("设置总工时数据失败: %v", err)
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), valueActual); err != nil {
			return fmt.Errorf("设置实际工时数据失败: %v", err)
		}

		// 获取设备维护完成率数据
		equipmentRepairCompletionRateList, err := dao.GetEquipmentRepairCompletionRateList(currentTime)
		if err != nil {
			return fmt.Errorf("获取设备维护完成率数据失败: %v", err)
		}

		// 获取设备维修次数数据
		equipmentRepairCompletionTimes, err := dao.GetEquipmentRepairCompletionTimes(currentTime)
		if err != nil {
			return fmt.Errorf("获取设备维修次数数据失败: %v", err)
		}

		// 填充设备维护完成率数据
		for _, rate := range equipmentRepairCompletionRateList {
			for j, equipment := range equipmentTypes {
				if rate.Name == equipment {
					if err := f.SetCellValue(sheetName, fmt.Sprintf("%c%d", 'E'+j, i+2), strconv.FormatFloat(rate.Value, 'f', -1, 64)+"%"); err != nil {
						return fmt.Errorf("设置设备维护完成率数据失败: %v", err)
					}
					break
				}
			}
		}

		// 填充设备维修次数数据
		for _, times := range equipmentRepairCompletionTimes {
			for j, equipment := range equipmentTypes {
				if times.Name == equipment {
					if err := f.SetCellValue(sheetName, fmt.Sprintf("%c%d", 'O'+j, i+2), times.Value); err != nil {
						return fmt.Errorf("设置设备维修次数数据失败: %v", err)
					}
					break
				}
			}
		}
	}

	return nil
}
