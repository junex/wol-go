package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// TCPChecker TCP 端口检测器
type TCPChecker struct {
	Port    int
	Timeout time.Duration
}

// NewTCPChecker 创建 TCP 检测器
func NewTCPChecker(port, timeoutSec int) *TCPChecker {
	return &TCPChecker{
		Port:    port,
		Timeout: time.Duration(timeoutSec) * time.Second,
	}
}

// Check 检测 TCP 端口是否开放
func (c *TCPChecker) Check(ipAddress string) bool {
	address := fmt.Sprintf("%s:%d", ipAddress, c.Port)

	// 尝试连接
	conn, err := net.DialTimeout("tcp", address, c.Timeout)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

// CheckPort 检测指定端口是否开放（静态方法）
func CheckPort(ipAddress string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", ipAddress, port)

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

// ParseTestType 解析测试类型并检测
// testType: "icmp", "arp", 或端口号字符串
func ParseTestTypeAndCheck(ipAddress, testType string, pingTimeout, arpTimeout, tcpTimeout int) bool {
	testTypeLower := strings.ToLower(testType)

	switch testTypeLower {
	case "icmp":
		checker := NewICMPChecker(pingTimeout)
		return checker.Check(ipAddress)

	case "arp":
		checker := NewARPChecker(arpTimeout, "")
		return checker.Check(ipAddress)

	default:
		// 尝试解析为端口号
		port, err := strconv.Atoi(testType)
		if err != nil {
			return false
		}

		checker := NewTCPChecker(port, tcpTimeout)
		return checker.Check(ipAddress)
	}
}
