package services

import (
	"strings"
	"testing"
)

// TestGenerateTaskNo 测试任务编号生成的正确性
func TestGenerateTaskNo(t *testing.T) {
	s := &TaskService{}

	testCases := []struct {
		taskTypeCode   string
		expectedPrefix string
	}{
		{"requirement", "REQ"},
		{"unit_task", "UNIT"},
		{"feature", "FEA"},
		{"bug", "BUG"},
		{"task", "TASK"},
	}

	for _, tc := range testCases {
		// 跳过实际的数据库测试，因为没有初始化数据库
		prefix := s.generatePrefixFromCode(tc.taskTypeCode)
		if prefix != tc.expectedPrefix {
			t.Errorf("generatePrefixFromCode(%s) = %s, want %s", tc.taskTypeCode, prefix, tc.expectedPrefix)
		}
	}
}

// TestGenerateRandomString 测试随机字符串生成
func TestGenerateRandomString(t *testing.T) {
	s := &TaskService{}

	tests := []int{6, 10, 20}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for _, length := range tests {
		result := s.generateRandomString(length)

		// 检查长度
		if len(result) != length {
			t.Errorf("generateRandomString(%d) length = %d, want %d", length, len(result), length)
		}

		// 检查所有字符是否在允许的字符集中
		for _, ch := range result {
			if !strings.ContainsRune(charset, ch) {
				t.Errorf("generateRandomString(%d) contains invalid character '%c'", length, ch)
			}
		}
	}
}

// TestGetTaskTypePrefix 测试任务类型前缀获取
func TestGetTaskTypePrefix(t *testing.T) {
	s := &TaskService{}

	testCases := []struct {
		taskTypeCode   string
		expectedPrefix string
	}{
		{"requirement", "REQ"},
		{"unit_task", "UNIT"},
		{"unknown_type", "UNK"}, // 未知类型，取前3个字母
	}

	for _, tc := range testCases {
		prefix := s.generatePrefixFromCode(tc.taskTypeCode)
		if prefix != tc.expectedPrefix {
			t.Errorf("generatePrefixFromCode(%s) = %s, want %s", tc.taskTypeCode, prefix, tc.expectedPrefix)
		}
	}
}
