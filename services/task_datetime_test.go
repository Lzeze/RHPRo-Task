package services

import (
	"testing"
)

// TestParseDateTime 测试日期时间解析功能
func TestParseDateTime(t *testing.T) {
	s := &TaskService{}

	testCases := []struct {
		name      string
		input     string
		shouldErr bool
		expected  string // 期望的日期输出格式：2006-01-02
	}{
		{
			name:      "日期格式 YYYY-MM-DD",
			input:     "2025-12-27",
			shouldErr: false,
			expected:  "2025-12-27",
		},
		{
			name:      "RFC3339 格式",
			input:     "2025-12-27T15:30:45Z",
			shouldErr: false,
			expected:  "2025-12-27",
		},
		{
			name:      "RFC3339 with timezone",
			input:     "2025-12-27T15:30:45+08:00",
			shouldErr: false,
			expected:  "2025-12-27",
		},
		{
			name:      "ISO8601 格式",
			input:     "2025-12-27T15:30:45",
			shouldErr: false,
			expected:  "2025-12-27",
		},
		{
			name:      "标准格式 YYYY-MM-DD HH:MM:SS",
			input:     "2025-12-27 15:30:45",
			shouldErr: false,
			expected:  "2025-12-27",
		},
		{
			name:      "空字符串",
			input:     "",
			shouldErr: false,
			expected:  "",
		},
		{
			name:      "无效格式",
			input:     "invalid-date",
			shouldErr: true,
			expected:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := s.parseDateTime(tc.input)

			// 检查错误
			if tc.shouldErr {
				if err == nil {
					t.Errorf("期望出错，但没有出错")
				}
				return
			}

			if !tc.shouldErr && err != nil {
				t.Errorf("不应该出错，但得到错误: %v", err)
				return
			}

			// 检查返回值
			if tc.input == "" {
				if result != nil {
					t.Errorf("期望返回 nil，但得到 %v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("期望返回非 nil 值")
				return
			}

			// 验证日期
			resultDate := result.Format("2006-01-02")
			if resultDate != tc.expected {
				t.Errorf("日期不匹配，期望 %s，得到 %s", tc.expected, resultDate)
			}
		})
	}
}

// TestParseDateTime_EdgeCases 测试日期解析的边界情况
func TestParseDateTime_EdgeCases(t *testing.T) {
	s := &TaskService{}

	testCases := []struct {
		name     string
		input    string
		expected string // YYYY-MM-DD
	}{
		{
			name:     "年初",
			input:    "2025-01-01",
			expected: "2025-01-01",
		},
		{
			name:     "年末",
			input:    "2025-12-31",
			expected: "2025-12-31",
		},
		{
			name:     "闰年二月",
			input:    "2024-02-29",
			expected: "2024-02-29",
		},
		{
			name:     "午夜",
			input:    "2025-12-27T00:00:00Z",
			expected: "2025-12-27",
		},
		{
			name:     "UTC+8 时区",
			input:    "2025-12-27T23:59:59+08:00",
			expected: "2025-12-27",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := s.parseDateTime(tc.input)

			if err != nil {
				t.Errorf("不应该出错: %v", err)
				return
			}

			if result == nil {
				t.Errorf("期望返回非 nil 值")
				return
			}

			resultDate := result.Format("2006-01-02")
			if resultDate != tc.expected {
				t.Errorf("日期不匹配，期望 %s，得到 %s", tc.expected, resultDate)
			}
		})
	}
}
