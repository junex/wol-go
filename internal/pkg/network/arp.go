package network

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ARPChecker ARP 状态检测器
type ARPChecker struct {
	Timeout    time.Duration
	Interface  string
}

// NewARPChecker 创建 ARP 检测器
func NewARPChecker(timeoutMs int, iface string) *ARPChecker {
	return &ARPChecker{
		Timeout:   time.Duration(timeoutMs) * time.Millisecond,
		Interface: iface,
	}
}

// Check 检测主机是否在线（使用 arp-scan）
func (c *ARPChecker) Check(ipAddress string) bool {
	// 构造 arp-scan 命令
	timeoutStr := strconv.Itoa(int(c.Timeout.Milliseconds()))
	args := []string{"-qx", "-t", timeoutStr, ipAddress}

	// 如果指定了网络接口
	if c.Interface != "" {
		args = append([]string{"-I", c.Interface}, args...)
	}

	cmd := exec.Command("arp-scan", args...)

	// 运行命令并捕获输出
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// 检查是否有输出
	return strings.TrimSpace(string(output)) != ""
}
