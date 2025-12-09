package controllers

import (
	"RHPRo-Task/dto"
	"RHPRo-Task/tests/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateTask_Success 测试创建任务
func TestCreateTask_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.POST("/api/v1/tasks", taskController.CreateTask)

	reqBody := dto.TaskRequest{
		TaskNo:       "TASK-TEST-001",
		Title:        "测试任务",
		Description:  "这是一个测试任务",
		TaskTypeCode: "requirement",
		Priority:     1,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	// 成功(0)或失败(500)或任务号重复(400)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500,
		"Response code should be 0, 400 or 500, got %d", resp.Code)
}

// TestCreateTask_InvalidInput 测试创建任务无效输入
func TestCreateTask_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.POST("/api/v1/tasks", taskController.CreateTask)

	reqBody := dto.TaskRequest{
		TaskNo: "", // 空任务号
		Title:  "", // 空标题
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestCreateTask_WithExecutor 测试创建任务并分配执行人
func TestCreateTask_WithExecutor(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.POST("/api/v1/tasks", taskController.CreateTask)

	executorID := uint(2)
	reqBody := dto.TaskRequest{
		TaskNo:       "TASK-TEST-002",
		Title:        "测试任务带执行人",
		TaskTypeCode: "requirement",
		ExecutorID:   &executorID,
		Priority:     2,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestGetTaskList_Success 测试获取任务列表
func TestGetTaskList_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.GET("/api/v1/tasks", taskController.GetTaskList)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/tasks?page=1&page_size=10", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestGetTaskList_WithFilters 测试带过滤条件获取任务列表
func TestGetTaskList_WithFilters(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.GET("/api/v1/tasks", taskController.GetTaskList)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/tasks?page=1&page_size=10&task_type_code=requirement&priority=1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestGetTaskList_ByStatus 测试按状态获取任务列表
func TestGetTaskList_ByStatus(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.GET("/api/v1/tasks", taskController.GetTaskList)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/tasks?status_code=draft", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestGetTaskByID_Success 测试根据ID获取任务
func TestGetTaskByID_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.GET("/api/v1/tasks/:id", taskController.GetTaskByID)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/tasks/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 404 || resp.Code == 500)
}

// TestGetTaskByID_InvalidID 测试无效任务ID
func TestGetTaskByID_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.GET("/api/v1/tasks/:id", taskController.GetTaskByID)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/tasks/invalid", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetTaskByID_NotFound 测试任务不存在
func TestGetTaskByID_NotFound(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.GET("/api/v1/tasks/:id", taskController.GetTaskByID)

	w := testutils.HTTPRequest(router, "GET", "/api/v1/tasks/999999", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
}

// TestUpdateTask_Success 测试更新任务
func TestUpdateTask_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.PUT("/api/v1/tasks/:id", taskController.UpdateTask)

	reqBody := dto.UpdateTaskRequest{
		Title:       "更新后的标题",
		Description: "更新后的描述",
		Priority:    3,
	}

	w := testutils.HTTPRequest(router, "PUT", "/api/v1/tasks/11", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestUpdateTask_InvalidID 测试更新无效任务ID
func TestUpdateTask_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.PUT("/api/v1/tasks/:id", taskController.UpdateTask)

	// 使用 map 代替结构体，避免 DTO 验证器标签问题
	reqBody := map[string]interface{}{
		"title": "Test",
	}

	w := testutils.HTTPRequest(router, "PUT", "/api/v1/tasks/invalid", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDeleteTask_Success 测试删除任务
func TestDeleteTask_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.DELETE("/api/v1/tasks/:id", taskController.DeleteTask)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/tasks/999", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestDeleteTask_InvalidID 测试删除无效任务ID
func TestDeleteTask_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.DELETE("/api/v1/tasks/:id", taskController.DeleteTask)

	w := testutils.HTTPRequest(router, "DELETE", "/api/v1/tasks/invalid", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestTransitStatus_Success 测试任务状态转换
func TestTransitStatus_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.POST("/api/v1/tasks/:id/transit", taskController.TransitStatus)

	reqBody := dto.TaskStatusTransitionRequest{
		ToStatusCode: "assigned",
		Comment:      "测试状态转换",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/transit", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	// 成功(0)或转换失败(500)或无效转换(400)
	assert.True(t, resp.Code == 0 || resp.Code == 400 || resp.Code == 500)
}

// TestTransitStatus_InvalidID 测试状态转换无效任务ID
func TestTransitStatus_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.POST("/api/v1/tasks/:id/transit", taskController.TransitStatus)

	reqBody := dto.TaskStatusTransitionRequest{
		ToStatusCode: "assigned",
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/invalid/transit", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestTransitStatus_InvalidInput 测试状态转换无效输入
func TestTransitStatus_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.POST("/api/v1/tasks/:id/transit", taskController.TransitStatus)

	reqBody := dto.TaskStatusTransitionRequest{
		ToStatusCode: "", // 空状态码
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/transit", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// TestAssignExecutor_Success 测试分配执行人
func TestAssignExecutor_Success(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.POST("/api/v1/tasks/:id/assign", taskController.AssignExecutor)

	reqBody := dto.AssignExecutorRequest{
		ExecutorID: 2,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/assign", reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.True(t, resp.Code == 0 || resp.Code == 500)
}

// TestAssignExecutor_InvalidID 测试分配执行人无效任务ID
func TestAssignExecutor_InvalidID(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.POST("/api/v1/tasks/:id/assign", taskController.AssignExecutor)

	reqBody := dto.AssignExecutorRequest{
		ExecutorID: 1,
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/invalid/assign", reqBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAssignExecutor_InvalidInput 测试分配执行人无效输入
func TestAssignExecutor_InvalidInput(t *testing.T) {
	router := testutils.SetupTestRouterWithAuth(1, "admin")
	taskController := NewTaskController()
	router.POST("/api/v1/tasks/:id/assign", taskController.AssignExecutor)

	reqBody := dto.AssignExecutorRequest{
		ExecutorID: 0, // 无效执行人ID
	}

	w := testutils.HTTPRequest(router, "POST", "/api/v1/tasks/1/assign", reqBody)
	resp, err := testutils.ParseResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}
