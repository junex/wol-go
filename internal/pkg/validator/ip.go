package validator

import (
	"net"
)

// ValidateIP 验证 IP 地址格式
func ValidateIP(ip string) error {
	if ip == "" {
		return ErrEmptyIP
	}
	if net.ParseIP(ip) == nil {
		return ErrInvalidIP
	}
	return nil
}

// 验证错误
var (
	ErrEmptyIP   = &ValidationError{Field: "IP", Message: "IP 地址不能为空"}
	ErrInvalidIP = &ValidationError{Field: "IP", Message: "IP 地址格式无效"}
)

// ValidationError 验证错误类型
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
