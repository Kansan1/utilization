package dao

import (
	"fmt"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

func ValidateUser(user *models.LoginUser) (bool, error) {
	var count int

	fmt.Printf("用户名: '%s', 长度: %d\n", user.Username, len(user.Username))

	// 使用参数化查询防止SQL注入
	err := config.DB.QueryRow(
		"SELECT COUNT(*) FROM [tbUsers] WHERE RTRIM(UID) = @p1 AND PWD = @p2",
		user.Username, user.Password).
		Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
