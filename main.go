package main

// @title 任务管理系统 API
// @version 1.0
// @description 企业级任务管理系统的 RESTful API 接口文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 192.168.12.48:8888
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"RHPRo-Task/config"
	"RHPRo-Task/database"
	"RHPRo-Task/routes"
	"RHPRo-Task/utils"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	// 初始化日志
	utils.InitLogger()
	utils.Logger.Info("Starting application...")

	_ = godotenv.Load()
	// 加载配置
	cfg := config.LoadConfig()

	// 设置JWT密钥
	utils.SetJWTSecret(cfg.JWT.Secret)

	// 初始化数据库
	database.InitPostgres(cfg)
	// 初始化Redis
	// database.InitRedis(cfg)

	// 自动迁移数据库
	// database.AutoMigrate()

	// 初始化路由
	router := routes.SetupRoutes()

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	utils.Logger.Info(fmt.Sprintf("Server starting on %s", addr))

	if err := router.Run(addr); err != nil {
		utils.Logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
