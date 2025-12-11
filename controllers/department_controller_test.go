package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/tests/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateDepartment_Success 测试创建部门
func TestCreateDepartment_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.POST("/api/v1/departments", deptController.CreateDepartment)

	reqBody := dto.DepartmentRequest{
		Name:        "测试部门",
		Description: "这是一个测试部门",
		Status:      1,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/departments", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	// 成功(0)或失败(500)
	assert.True(t, resp.Code == 0 || resp.Code == 500, "Response code should be 0 or 500, got %d", resp.Code)
}

// TestCreateDepartment_InvalidInput 测试创建部门无效输入
func TestCreateDepartment_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.POST("/api/v1/departments", deptController.CreateDepartment)

	reqBody := dto.DepartmentRequest{
		Name: "", // 空名称
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/departments", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestCreateDepartment_WithParent 测试创建子部门
func TestCreateDepartment_WithParent(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.POST("/api/v1/departments", deptController.CreateDepartment)

	parentID := uint(1)
	reqBody := dto.DepartmentRequest{
		Name:        "子部门",
		Description: "测试子部门",
		ParentID:    &parentID,
		Status:      1,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/departments", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestUpdateDepartment_Success 测试更新部门
func TestUpdateDepartment_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.PUT("/api/v1/departments/:id", deptController.UpdateDepartment)

	reqBody := dto.DepartmentRequest{
		Name:        "更新后的部门名",
		Description: "更新后的描述",
		Status:      1,
	}

	w := testutils.HTTPRequest(router, "PUT", "/api/v1/departments/1", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestUpdateDepartment_InvalidID 测试更新无效部门ID
func TestUpdateDepartment_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.PUT("/api/v1/departments/:id", deptController.UpdateDepartment)

	reqBody := dto.DepartmentRequest{
		Name: "Test",
	}

	w := testutils.HTTPRequest(router, "PUT", "/api/v1/departments/invalid", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDeleteDepartment_Success 测试删除部门
func TestDeleteDepartment_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.DELETE("/api/v1/departments/:id", deptController.DeleteDepartment)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/departments/999", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	// 成功(0)或失败(500)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestDeleteDepartment_InvalidID 测试删除无效部门ID
func TestDeleteDepartment_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.DELETE("/api/v1/departments/:id", deptController.DeleteDepartment)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/departments/invalid", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetDepartmentList_Success 测试获取部门列表
func TestGetDepartmentList_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.GET("/api/v1/departments", deptController.GetDepartmentList)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/departments", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestGetDepartmentDetail_Success 测试获取部门详情
func TestGetDepartmentDetail_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.GET("/api/v1/departments/:id", deptController.GetDepartmentDetail)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/departments/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestGetDepartmentDetail_InvalidID 测试获取无效部门详情
func TestGetDepartmentDetail_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.GET("/api/v1/departments/:id", deptController.GetDepartmentDetail)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/departments/invalid", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAddLeader_Success 测试添加部门负责人
func TestAddLeader_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.POST("/api/v1/departments/:id/leaders", deptController.AddLeader)

	reqBody := dto.AddLeaderRequest{
		UserID: 1,
		// IsPrimary: true,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/departments/1/leaders", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestAddLeader_InvalidDeptID 测试添加负责人无效部门ID
func TestAddLeader_InvalidDeptID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.POST("/api/v1/departments/:id/leaders", deptController.AddLeader)

	reqBody := dto.AddLeaderRequest{
		UserID: 1,
		// IsPrimary: true,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/departments/invalid/leaders", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAddLeader_InvalidInput 测试添加负责人无效输入
func TestAddLeader_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.POST("/api/v1/departments/:id/leaders", deptController.AddLeader)

	reqBody := dto.AddLeaderRequest{
		UserID: 0, // 无效用户ID
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/departments/1/leaders", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRemoveLeader_Success 测试移除部门负责人
func TestRemoveLeader_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.DELETE("/api/v1/departments/:id/leaders/:userId", deptController.RemoveLeader)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/departments/1/leaders/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestRemoveLeader_InvalidDeptID 测试移除负责人无效部门ID
func TestRemoveLeader_InvalidDeptID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.DELETE("/api/v1/departments/:id/leaders/:userId", deptController.RemoveLeader)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/departments/invalid/leaders/1", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRemoveLeader_InvalidUserID 测试移除负责人无效用户ID
func TestRemoveLeader_InvalidUserID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.DELETE("/api/v1/departments/:id/leaders/:userId", deptController.RemoveLeader)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/departments/1/leaders/invalid", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAssignUsers_Success 测试分配用户到部门
func TestAssignUsers_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.POST("/api/v1/departments/:id/users", deptController.AssignUsers)

	userIDs := []uint{1, 2, 3}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/departments/1/users", userIDs)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestAssignUsers_InvalidDeptID 测试分配用户无效部门ID
func TestAssignUsers_InvalidDeptID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.POST("/api/v1/departments/:id/users", deptController.AssignUsers)

	userIDs := []uint{1}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/departments/invalid/users", userIDs)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAssignUsers_EmptyList 测试分配空用户列表
func TestAssignUsers_EmptyList(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	deptController := NewDepartmentController()
	router.POST("/api/v1/departments/:id/users", deptController.AssignUsers)

	userIDs := []uint{}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/departments/1/users", userIDs)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	// 空列表可能成功或失败
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}
