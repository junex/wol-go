package validator

import (
	"strconv"
	"strings"
)

// ValidateTestType 验证测试类型（icmp、arp 或端口号）
func ValidateTestType(testType string) error {
	if testType == "" {
		return ErrEmptyTestType
	}

	// 转换为小写进行比较
	testTypeLower := strings.ToLower(testType)

	// 检查是否是 icmp 或 arp
	if testTypeLower == "icmp" || testTypeLower == "arp" {
		return nil
	}

	// 检查是否是有效的端口号
	port, err := strconv.Atoi(testType)
	if err != nil {
		return ErrInvalidTestType
	}

	if port < 1 || port > 65535 {
		return ErrInvalidPortRange
	}

	return nil
}

// 验证错误
var (
	ErrEmptyTestType  = &ValidationError{Field: "TestType", Message: "测试类型不能为空"}
	ErrInvalidTestType = &ValidationError{Field: "TestType", Message: "测试类型必须是 'icmp'、'arp' 或有效的端口号"}
	ErrInvalidPortRange = &ValidationError{Field: "TestType", Message: "端口号必须在 1-65535 范围内"}
)
