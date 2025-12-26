package validator

import (
	"regexp"
)

var (
	// macRegex MAC 地址正则表达式 (XX:XX:XX:XX:XX:XX)
	macRegex = regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}([0-9a-fA-F]{2})$`)
)

// ValidateMAC 验证 MAC 地址格式
func ValidateMAC(mac string) error {
	if mac == "" {
		return ErrEmptyMAC
	}
	if !macRegex.MatchString(mac) {
		return ErrInvalidMAC
	}
	return nil
}

// 验证错误
var (
	ErrEmptyMAC   = &ValidationError{Field: "MAC", Message: "MAC 地址不能为空"}
	ErrInvalidMAC = &ValidationError{Field: "MAC", Message: "MAC 地址格式无效，必须是 XX:XX:XX:XX:XX:XX 格式"}
)
