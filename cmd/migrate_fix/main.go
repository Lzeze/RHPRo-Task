package main

import (
	"RHPRo-Task/config"
	"RHPRo-Task/database"
	"RHPRo-Task/utils"
	"fmt"
	"log"
	"os"
)

func main() {
	// 初始化日志
	utils.InitLogger()

	// 加载配置
	fmt.Println("Loading config...")
	// 设置环境变量以确保可以连接到测试或开发数据库
	// 使用 .env 中的配置: DB_HOST=10.0.10.114, DB_USER=postgres, DB_PASSWORD=rhzy2025
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "10.0.10.114")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "postgres")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "rhzy2025")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "rhpro_task")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "5432")
	}

	cfg := config.LoadConfig()

	// 初始化数据库连接
	fmt.Println("Initializing database connection...")
	database.InitPostgres(cfg)

	// 读取 SQL 文件
	fmt.Println("Reading SQL file...")
	sqlBytes, err := os.ReadFile("database/migrations/fix_missing_columns.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	// 执行 SQL
	fmt.Println("Executing SQL migration...")
	err = database.DB.Exec(string(sqlBytes)).Error
	if err != nil {
		log.Fatalf("Failed to execute SQL: %v", err)
	}

	fmt.Println("Migration completed successfully!")
}
