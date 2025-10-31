package dao

import (
	"database/sql"
	"fmt"
	"sync"
	"utilization-backend/config"
	"utilization-backend/src/api/models"

	"go.uber.org/zap"
)

// 获取稼动率列表
func GetUtilizationListByMonth(year int, month int) ([]models.UtilizationWorkingHour, []models.UtilizationException) {

	var utilizationWorkingHourList []models.UtilizationWorkingHour
	var utilizationExceptionList []models.UtilizationException

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		utilizationWorkingHourQueryList, err := queryUtilizationWorkingHour(year, month)
		if err != nil {
			errCh <- fmt.Errorf("查询工作时间失败:%w", err)
			return
		}

		utilizationWorkingHourList = append(utilizationWorkingHourList, utilizationWorkingHourQueryList...)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		utilizationExceptionQueryList, err := queryUtilizationException(year, month)

		if err != nil {
			errCh <- fmt.Errorf("查询异常时间失败:%w", err)
		}
		utilizationExceptionList = append(utilizationExceptionList, utilizationExceptionQueryList...)

	}()

	wg.Wait()

	select {
	case err := <-errCh:
		zap.L().Error("查询稼动率失败", zap.Error(err))
		return nil, nil
	default:
		return utilizationWorkingHourList, utilizationExceptionList
	}
}

// 查询异常时间
func queryUtilizationException(year, month int) ([]models.UtilizationException, error) {
	fmt.Printf("查询参数: %d/%d\n", year, month)

	// 使用问号占位符或直接拼接SQL (测试用)
	queryParam := fmt.Sprintf("%d/%d%%", year, month)
	fmt.Printf("查询参数: %s\n", queryParam)
	// 方法1: 使用 ? 参数
	rows, err := config.DB.Query(`
        SELECT time,  abnormal_stop_time, line
        FROM Utilization_exception
        WHERE time like @p1 order by time asc`, sql.Named("p1", queryParam))

	if err != nil {
		fmt.Printf("查询错误: %v\n", err)
		return nil, err
	}
	defer rows.Close() // 正确位置，确保函数返回前关闭

	var utilizationExceptionList []models.UtilizationException
	count := 0

	for rows.Next() {
		count++
		var utilizationException models.UtilizationException

		// 确保字段顺序与SQL查询一致
		err := rows.Scan(
			&utilizationException.Time,
			&utilizationException.AbnormalStopTime,
			&utilizationException.Line, // 补上缺少的字段
		)

		if err != nil {
			fmt.Printf("扫描错误: %v\n", err)
			continue
		}
		utilizationExceptionList = append(utilizationExceptionList, utilizationException)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("遍历错误: %v\n", err)
		return nil, err
	}
	return utilizationExceptionList, nil
}

// 查询稼动率时间
func queryUtilizationWorkingHour(year, month int) ([]models.UtilizationWorkingHour, error) {

	queryParam := fmt.Sprintf("%d/%d%%", year, month)

	rows, err := config.DB.Query(`
	Select time,type,working_hours,line from Utilization_working_hours 
	where time like @yearMonth
	order by time asc`,
		sql.Named("yearMonth", queryParam))
	if err != nil {
		zap.L().Error("查询工作时间失败",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var utilizationWorkingHourList []models.UtilizationWorkingHour
	for rows.Next() {
		var utilizationWorkingHour models.UtilizationWorkingHour
		err := rows.Scan(
			&utilizationWorkingHour.Time,
			&utilizationWorkingHour.Type,
			&utilizationWorkingHour.WorkingHours,
			&utilizationWorkingHour.Line,
		)
		if err != nil {
			zap.L().Error("稼动率解析失败", zap.Error(err))
			continue
		}
		utilizationWorkingHourList = append(utilizationWorkingHourList, utilizationWorkingHour)
	}

	if err = rows.Err(); err != nil {
		zap.L().Error("迭代工作时间数据时出错", zap.Error(err))
		return nil, err
	}

	defer rows.Close()
	return utilizationWorkingHourList, nil
}
