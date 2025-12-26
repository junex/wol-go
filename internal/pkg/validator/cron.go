package validator

import (
	"regexp"
	"strings"
)

// Cron 表达式的 5 个字段正则模式
var (
	minutePattern   = regexp.MustCompile(`^(\*|[1-5]?[0-9](-[1-5]?[0-9])?)(\/[1-9][0-9]*)?(,(\*|[1-5]?[0-9](-[1-5]?[0-9])?)(\/[1-9][0-9]*)?)*$`)
	hourPattern     = regexp.MustCompile(`^(\*|(1?[0-9]|2[0-3])(-(1?[0-9]|2[0-3]))?)(\/[1-9][0-9]*)?(,(\*|(1?[0-9]|2[0-3])(-(1?[0-9]|2[0-3]))?)(\/[1-9][0-9]*)?)*$`)
	dayPattern      = regexp.MustCompile(`^(\*|([1-9]|[1-2][0-9]?|3[0-1])(-([1-9]|[1-2][0-9]?|3[0-1]))?)(\/[1-9][0-9]*)?(,(\*|([1-9]|[1-2][0-9]?|3[0-1])(-([1-9]|[1-2][0-9]?|3[0-1]))?)(\/[1-9][0-9]*)?)*$`)
	monthPattern    = regexp.MustCompile(`^(\*|([1-9]|1[0-2]?)(-([1-9]|1[0-2]?))?)(\/[1-9][0-9]*)?(,(\*|([1-9]|1[0-2]?)(-([1-9]|1[0-2]?))?)(\/[1-9][0-9]*)?)*$`)
	weekdayPattern  = regexp.MustCompile(`^(\*|[0-6](-[0-6])?)(\/[1-9][0-9]*)?(,(\*|[0-6](-[0-6])?)(\/[1-9][0-9]*)?)*$`)
)

// ValidateCron 验证 Cron 表达式格式
func ValidateCron(cronExpr string) error {
	if cronExpr == "" {
		return ErrEmptyCron
	}

	parts := strings.Split(cronExpr, " ")
	if len(parts) != 5 {
		return ErrInvalidCronFormat
	}

	patterns := []*regexp.Regexp{
		minutePattern,
		hourPattern,
		dayPattern,
		monthPattern,
		weekdayPattern,
	}

	fieldNames := []string{"分钟", "小时", "日", "月", "周"}

	for i, part := range parts {
		if !patterns[i].MatchString(part) {
			return &ValidationError{
				Field:   fieldNames[i],
				Message: "字段格式无效: " + part,
			}
		}
	}

	return nil
}

// 验证错误
var (
	ErrEmptyCron       = &ValidationError{Field: "Cron", Message: "Cron 表达式不能为空"}
	ErrInvalidCronFormat = &ValidationError{Field: "Cron", Message: "Cron 表达式必须包含 5 个字段（分 时 日 月 周）"}
)
