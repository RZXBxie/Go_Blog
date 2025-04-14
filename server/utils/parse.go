package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration 解析持续时间字符串为 time.Duration。
// 持续时间字符串应由数字值和时间单位组成，单位可以是 "d" 表示天，"h" 表示小时，"m" 表示分钟，"s" 表示秒。
// 例如，"1d2h30m" 会被解析为 1 天、2 小时和 30 分钟。
// 如果字符串为空或格式无效，则返回错误。
func ParseDuration(duration string) (time.Duration, error) {
	duration = strings.TrimSpace(duration)
	if len(duration) == 0 {
		return 0, fmt.Errorf("empty duration")
	}

	unitPattern := map[string]time.Duration{
		"d": time.Hour * 24,
		"h": time.Hour,
		"m": time.Minute,
		"s": time.Second,
	}

	var totalDuration time.Duration
	for _, unit := range []string{"d", "h", "m", "s"} {
		for strings.Contains(duration, unit) {
			// 找到单位的位置
			unitIndex := strings.Index(duration, unit)
			part := duration[:unitIndex]
			if part == "" {
				part = "0"
			}
			val, err := strconv.Atoi(part)
			if err != nil {
				return 0, fmt.Errorf("invalid duration part: %v", err)
			}

			totalDuration += time.Duration(val) * unitPattern[unit]
			// 移出已处理的部分
			duration = duration[unitIndex+len(unit):]
		}
	}
	if len(duration) > 0 {
		return 0, fmt.Errorf("unrecognized duration: format: %v", duration)
	}

	return totalDuration, nil
}
