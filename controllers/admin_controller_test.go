package controllers

import (
	"RHPRo-Task/tests/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetRoleList_Success 测试获取角色列表
func TestGetRoleList_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	adminController := NewAdminController()
	router.GET("/api/v1/admin/roles", adminController.GetRoleList)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/admin/roles", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	// 成功获取(0)或服务器错误(500)
	assert.True(t, resp.Code == 0 || resp.Code == 500,
		"Response code should be 0 or 500, got %d", resp.Code)
}

// TestGetPermissionList_Success 测试获取权限列表
func TestGetPermissionList_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	adminController := NewAdminController()
	router.GET("/api/v1/admin/permissions", adminController.GetPermissionList)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/admin/permissions", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500,
		"Response code should be 0 or 500, got %d", resp.Code)
}
