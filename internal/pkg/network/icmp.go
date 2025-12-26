package network

import (
	"os/exec"
	"strconv"
	"time"
)

// ICMPChecker ICMP 状态检测器
type ICMPChecker struct {
	Timeout time.Duration
}

// NewICMPChecker 创建 ICMP 检测器
func NewICMPChecker(timeoutMs int) *ICMPChecker {
	return &ICMPChecker{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}
}

// Check 检测主机是否在线（使用 fping）
func (c *ICMPChecker) Check(ipAddress string) bool {
	// 使用 fping 命令
	timeoutStr := strconv.Itoa(int(c.Timeout.Milliseconds()))
	cmd := exec.Command("fping", "-t", timeoutStr, "-c", "1", ipAddress)

	// 运行命令，不输出任何内容
	err := cmd.Run()
	return err == nil
}

// CheckWithPing 使用系统 ping 命令检测（备选方案）
func (c *ICMPChecker) CheckWithPing(ipAddress string) bool {
	timeoutStr := strconv.Itoa(int(c.Timeout.Seconds()))
	cmd := exec.Command("ping", "-c", "1", "-W", timeoutStr, ipAddress)

	err := cmd.Run()
	return err == nil
}
