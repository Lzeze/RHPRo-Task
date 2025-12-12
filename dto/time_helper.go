package dto

import (
	"encoding/json"
	"time"
)

// TimeFormat 统一的时间格式常量
const (
	TimeFormatISO8601  = "2006-01-02T15:04:05Z07:00" // RFC3339
	TimeFormatDate     = "2006-01-02"                // 日期格式
	TimeFormatDatetime = "2006-01-02 15:04:05"       // 日期时间格式
)

// ResponseTime 响应中统一使用的时间类型，自动格式化为 RFC3339 格式
type ResponseTime struct {
	time.Time
}

// MarshalJSON 自定义 JSON 序列化方法，将时间格式化为 RFC3339
func (rt ResponseTime) MarshalJSON() ([]byte, error) {
	if rt.Time.IsZero() {
		return []byte("null"), nil
	}
	// 使用 RFC3339 格式序列化
	return json.Marshal(rt.Time.Format(time.RFC3339))
}

// UnmarshalJSON 自定义 JSON 反序列化方法，支持多种格式
func (rt *ResponseTime) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	if dateStr == "" {
		rt.Time = time.Time{}
		return nil
	}

	// 尝试多种日期格式
	formats := []string{
		time.RFC3339,           // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05Z", // 2006-01-02T15:04:05Z
		"2006-01-02T15:04:05",  // 2006-01-02T15:04:05
		"2006-01-02",           // YYYY-MM-DD
		"2006-01-02 15:04:05",  // YYYY-MM-DD HH:MM:SS
	}

	var t time.Time
	var err error

	for _, format := range formats {
		if t, err = time.Parse(format, dateStr); err == nil {
			rt.Time = t
			return nil
		}
	}

	return err
}

// NullableTime 可空的时间类型，支持 null 值
type NullableTime struct {
	Time  *time.Time
	Valid bool
}

// MarshalJSON 自定义 JSON 序列化方法
func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid || nt.Time == nil || nt.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time.Format(time.RFC3339))
}

// UnmarshalJSON 自定义 JSON 反序列化方法
func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	var dateStr *string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	if dateStr == nil || *dateStr == "" {
		nt.Time = nil
		nt.Valid = false
		return nil
	}

	// 尝试多种日期格式
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, *dateStr); err == nil {
			nt.Time = &t
			nt.Valid = true
			return nil
		}
	}

	nt.Valid = false
	return nil
}

// ToResponseTime 将 time.Time 转换为 ResponseTime
func ToResponseTime(t time.Time) ResponseTime {
	return ResponseTime{t}
}

// PtrToResponseTime 将 *time.Time 转换为 *ResponseTime
func PtrToResponseTime(t *time.Time) *ResponseTime {
	if t == nil {
		return nil
	}
	rt := ToResponseTime(*t)
	return &rt
}
