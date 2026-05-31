// Package money 在「分」(int64) 与「元」字符串/DECIMAL 之间转换，避免浮点误差。
package money

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lilongjie1137/HelloWorld/common/errcode"
)

// FormatYuan 把分转成两位小数的元字符串（如 1800 -> "18.00"）。
func FormatYuan(cents int64) string {
	neg := cents < 0
	if neg {
		cents = -cents
	}
	s := fmt.Sprintf("%d.%02d", cents/100, cents%100)
	if neg {
		return "-" + s
	}
	return s
}

// ParseYuan 解析元字符串为分，最多两位小数（如 "18.5" -> 1850）。
func ParseYuan(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errcode.ErrInvalidParam.WithMessage("金额为空")
	}
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")

	parts := strings.SplitN(s, ".", 2)
	yuan, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, errcode.ErrInvalidParam.WithMessage("金额格式错误")
	}
	var frac int64
	if len(parts) == 2 {
		f := parts[1]
		if len(f) > 2 {
			return 0, errcode.ErrInvalidParam.WithMessage("金额最多两位小数")
		}
		for len(f) < 2 {
			f += "0"
		}
		frac, err = strconv.ParseInt(f, 10, 64)
		if err != nil {
			return 0, errcode.ErrInvalidParam.WithMessage("金额格式错误")
		}
	}
	cents := yuan*100 + frac
	if neg {
		cents = -cents
	}
	return cents, nil
}
