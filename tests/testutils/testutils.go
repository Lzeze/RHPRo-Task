package testutils

import (
	"RHPRo-Task/config"
	"RHPRo-Task/database"
	"RHPRo-Task/middlewares"
	"RHPRo-Task/models"
	"RHPRo-Task/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var initOnce sync.Once

// InitTestEnv 初始化测试环境
func InitTestEnv() {
	initOnce.Do(func() {
		// 初始化日志
		utils.InitLogger()

		// 尝试加载环境变量（从多个可能的位置）
		_ = godotenv.Load("../../.env")
		_ = godotenv.Load("../.env")
		_ = godotenv.Load(".env")

		// 加载配置
		cfg := config.LoadConfig()

		// 设置JWT密钥
		utils.SetJWTSecret(cfg.JWT.Secret)

		// 初始化数据库
		database.InitPostgres(cfg)
	})
}

// TestConfig 测试配置
type TestConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

// DefaultTestConfig 默认测试配置
var DefaultTestConfig = TestConfig{
	DBHost:     "localhost",
	DBPort:     5432,
	DBUser:     "postgres",
	DBPassword: "password",
	DBName:     "gin_app_test", // 使用测试数据库
}

// SetupTestDB 设置测试数据库
func SetupTestDB(cfg *TestConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// 设置全局 DB
	database.DB = db

	return db, nil
}

// SetupTestDBWithDefault 使用默认配置设置测试数据库
func SetupTestDBWithDefault() (*gorm.DB, error) {
	return SetupTestDB(&DefaultTestConfig)
}

// CleanupTestDB 清理测试数据库连接
func CleanupTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// SetupTestRouter 设置测试路由器（不需要认证）
func SetupTestRouter() *gin.Engine {
	// 确保测试环境已初始化
	InitTestEnv()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.RecoveryMiddleware())
	router.Use(middlewares.ValidatorMiddleware())
	return router
}

// SetupTestRouterWithAuth 设置带模拟认证的测试路由器
func SetupTestRouterWithAuth(userID uint, username string) *gin.Engine {
	// 确保测试环境已初始化
	InitTestEnv()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.RecoveryMiddleware())
	router.Use(middlewares.ValidatorMiddleware())
	router.Use(MockAuthMiddleware(userID, username))
	return router
}

// MockAuthMiddleware 模拟认证中间件
func MockAuthMiddleware(userID uint, username string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 同时设置两种命名方式，兼容不同控制器
		c.Set("user_id", userID)
		c.Set("userID", userID)
		c.Set("username", username)
		c.Next()
	}
}

// GenerateTestToken 生成测试用 JWT Token
func GenerateTestToken(userID uint, username string) (string, error) {
	return utils.GenerateToken(userID, username, "13800138000", 24)
}

// HTTPRequest 发送 HTTP 请求
func HTTPRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request

	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// HTTPRequestWithAuth 发送带认证的 HTTP 请求
func HTTPRequestWithAuth(router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var req *http.Request

	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ParseResponse 解析响应
func ParseResponse(w *httptest.ResponseRecorder) (utils.Response, error) {
	var resp utils.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	return resp, err
}

// CreateTestUser 创建测试用户
func CreateTestUser(db *gorm.DB, username string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Email:    username + "@test.com",
		Password: "hashed_password",
		Mobile:   "13800138000",
		Status:   1,
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// CreateTestDepartment 创建测试部门
func CreateTestDepartment(db *gorm.DB, name string) (*models.Department, error) {
	dept := &models.Department{
		Name:   name,
		Status: 1,
	}
	if err := db.Create(dept).Error; err != nil {
		return nil, err
	}
	return dept, nil
}

// CreateTestTask 创建测试任务
func CreateTestTask(db *gorm.DB, title string, creatorID uint) (*models.Task, error) {
	task := &models.Task{
		TaskNo:       fmt.Sprintf("TASK-%d", creatorID),
		Title:        title,
		TaskTypeCode: "requirement",
		StatusCode:   "draft",
		CreatorID:    creatorID,
		Priority:     1,
	}
	if err := db.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

// AutoMigrateTest 自动迁移测试数据库
func AutoMigrateTest(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Department{},
		&models.DepartmentLeader{},
		&models.Task{},
		&models.TaskType{},
		&models.TaskStatus{},
		&models.TaskStatusTransition{},
		&models.ReviewSession{},
		&models.ReviewRecord{},
		&models.RequirementGoal{},
		&models.RequirementSolution{},
		&models.ExecutionPlan{},
	)
}

// CleanupTestData 清理测试数据,谨慎使用
func CleanupTestData(db *gorm.DB) {
	// 按照外键依赖顺序删除
	db.Exec("TRUNCATE TABLE review_records CASCADE")
	db.Exec("TRUNCATE TABLE review_sessions CASCADE")
	db.Exec("TRUNCATE TABLE execution_plans CASCADE")
	db.Exec("TRUNCATE TABLE requirement_solutions CASCADE")
	db.Exec("TRUNCATE TABLE requirement_goals CASCADE")
	db.Exec("TRUNCATE TABLE tasks CASCADE")
	db.Exec("TRUNCATE TABLE department_leaders CASCADE")
	db.Exec("TRUNCATE TABLE departments CASCADE")
	db.Exec("TRUNCATE TABLE user_roles CASCADE")
	db.Exec("TRUNCATE TABLE users CASCADE")
}
