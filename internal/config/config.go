package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config 应用配置结构体
type Config struct {
	// 服务器配置
	Port int
	TZ   string

	// 认证配置
	EnableLogin bool
	Username    string
	Password    string

	// 功能开关
	EnableAddDel  bool
	EnableRefresh bool

	// 状态检测配置
	RefreshInterval int // 秒
	PingTimeout     int // 毫秒
	ARPTimeout      int // 毫秒
	TCPTimeout      int // 秒
	ARPInterface    string

	// 数据路径
	DataPath string
	CronPath string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	cfg := &Config{
		Port:            getEnvInt("PORT", 5000),
		TZ:              getEnv("TZ", "UTC"),
		EnableLogin:     getEnvBool("ENABLE_LOGIN", false),
		Username:        getEnv("USERNAME", "admin"),
		Password:        getEnv("PASSWORD", "admin"),
		EnableAddDel:    getEnvBool("ENABLE_ADD_DEL", true),
		EnableRefresh:   getEnvBool("ENABLE_REFRESH", true),
		RefreshInterval: getEnvInt("REFRESH_INTERVAL", 30),
		PingTimeout:     getEnvInt("PING_TIMEOUT", 300),
		ARPTimeout:      getEnvInt("ARP_TIMEOUT", 300),
		TCPTimeout:      getEnvInt("TCP_TIMEOUT", 1),
		ARPInterface:    getEnv("ARP_INTERFACE", ""),
		DataPath:        getEnv("DATA_PATH", "/app/db"),
		CronPath:        getEnv("CRON_PATH", "/etc/cron.d"),
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("配置验证失败: %v", err))
	}

	return cfg
}

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	// 验证端口
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", c.Port)
	}

	// 验证刷新间隔
	validIntervals := map[int]bool{15: true, 30: true, 60: true}
	if !validIntervals[c.RefreshInterval] {
		return fmt.Errorf("无效的刷新间隔: %d, 必须是 15, 30 或 60", c.RefreshInterval)
	}

	// 验证超时值
	if c.PingTimeout < 100 || c.PingTimeout > 5000 {
		return fmt.Errorf("无效的 Ping 超时值: %d, 必须在 100-5000 毫秒之间", c.PingTimeout)
	}

	if c.ARPTimeout < 100 || c.ARPTimeout > 5000 {
		return fmt.Errorf("无效的 ARP 超时值: %d, 必须在 100-5000 毫秒之间", c.ARPTimeout)
	}

	if c.TCPTimeout < 1 || c.TCPTimeout > 10 {
		return fmt.Errorf("无效的 TCP 超时值: %d, 必须在 1-10 秒之间", c.TCPTimeout)
	}

	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数类型的环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvBool 获取布尔类型的环境变量
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}
