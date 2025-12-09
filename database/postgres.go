package database

import (
	"RHPRo-Task/config"
	"RHPRo-Task/models"
	"RHPRo-Task/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitPostgres 初始化PostgreSQL连接
func InitPostgres(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		utils.Logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	utils.Logger.Info("Database connected successfully")
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)

	if err != nil {
		utils.Logger.Fatal(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	utils.Logger.Info("Database migrated successfully")

	// 初始化默认数据
	initDefaultData()
}

// initDefaultData 初始化默认角色和权限
func initDefaultData() {
	// 创建默认权限
	permissions := []models.Permission{
		{Name: "user:read", Description: "读取用户信息"},
		{Name: "user:create", Description: "创建用户"},
		{Name: "user:update", Description: "更新用户"},
		{Name: "user:delete", Description: "删除用户"},
		{Name: "role:manage", Description: "管理角色"},
		{Name: "permission:manage", Description: "管理权限"},
	}

	for _, perm := range permissions {
		var existing models.Permission
		if err := DB.Where("name = ?", perm.Name).First(&existing).Error; err != nil {
			DB.Create(&perm)
		}
	}

	// 创建默认角色
	var adminRole models.Role
	if err := DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		adminRole = models.Role{
			Name:        "admin",
			Description: "管理员",
		}
		DB.Create(&adminRole)

		// 给管理员分配所有权限
		var allPermissions []models.Permission
		DB.Find(&allPermissions)
		DB.Model(&adminRole).Association("Permissions").Append(&allPermissions)
	}

	var userRole models.Role
	if err := DB.Where("name = ?", "user").First(&userRole).Error; err != nil {
		userRole = models.Role{
			Name:        "user",
			Description: "普通用户",
		}
		DB.Create(&userRole)

		// 给普通用户分配读取权限
		var readPerm models.Permission
		DB.Where("name = ?", "user:read").First(&readPerm)
		DB.Model(&userRole).Association("Permissions").Append(&readPerm)
	}

	utils.Logger.Info("Default data initialized")
}
