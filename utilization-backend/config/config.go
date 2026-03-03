package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host                   string `yaml:"host"`
		Port                   int    `yaml:"port"`
		Username               string `yaml:"username"`
		Password               string `yaml:"password"`
		DBName                 string `yaml:"dbname"`
		Encrypt                string `yaml:"encrypt"`
		TrustServerCertificate bool   `yaml:"trustServerCertificate"`
	} `yaml:"database"`

	Database2 struct {
		Host                   string `yaml:"host"`
		Port                   int    `yaml:"port"`
		Username               string `yaml:"username"`
		Password               string `yaml:"password"`
		DBName                 string `yaml:"dbname"`
		Encrypt                string `yaml:"encrypt"`
		TrustServerCertificate bool   `yaml:"trustServerCertificate"`
	} `yaml:"database2"`

	JWT struct {
		Secret      string `yaml:"secret"`
		ExpireHours int    `yaml:"expireHours"`
	} `yaml:"jwt"`
}

var AppConfig Config

// LoadConfig 加载配置文件
func LoadConfig() error {
	var data []byte
	var err error

	// 尝试读取打包路径下的配置文件
	data, err = os.ReadFile("config/config.yaml")
	if err != nil {
		// 如果失败，尝试读取开发路径下的配置文件
		data, err = os.ReadFile("../config/config.yaml")
		if err != nil {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 解析YAML
	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func GetDSN() string {
	db := AppConfig.Database
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=%v;TrustServerCertificate=%v",
		db.Host, db.Username, db.Password, db.Port, db.DBName, db.Encrypt, db.TrustServerCertificate)
}

// GetDSN2 获取第二个数据库连接字符串
func GetDSN2() string {
	db := AppConfig.Database2
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=%v;TrustServerCertificate=%v",
		db.Host, db.Username, db.Password, db.Port, db.DBName, db.Encrypt, db.TrustServerCertificate)
}

// GetJWTSecret 获取JWT密钥
func GetJWTSecret() []byte {
	return []byte(AppConfig.JWT.Secret)
}

// GetJWTExpireHours 获取JWT过期时间
func GetJWTExpireHours() int {
	return AppConfig.JWT.ExpireHours
}
