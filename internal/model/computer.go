package model

// Computer 计算机设备模型
type Computer struct {
	Name      string `json:"name"`
	MACAddr   string `json:"mac_address"`
	IPAddress string `json:"ip_address"`
	TestType  string `json:"test_type"` // "icmp", "arp", 或端口号字符串
}

// CronJob 定时任务模型
type CronJob struct {
	Schedule string `json:"schedule"` // Cron 表达式
	MACAddr  string `json:"mac_address"`
	Type     string `json:"type"`     // "wake" or "sleep"
}

// User 用户模型（用于认证）
type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // 注意：生产环境应该存储哈希值
}

// Session 会话模型
type Session struct {
	Token     string `json:"token"`
	Username  string `json:"username"`
	ExpiresAt int64  `json:"expires_at"`
}
