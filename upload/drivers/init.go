package drivers

import (
	"RHPRo-Task/upload"
	"RHPRo-Task/utils"
	"fmt"
)

// InitDrivers 初始化所有启用的驱动
func InitDrivers() error {
	config := upload.GetUploadConfig()
	factory := upload.GetFactory()

	// 初始化本地驱动
	if config.Local.Enabled {
		localDriver, err := NewLocalDriver(config.Local)
		if err != nil {
			return fmt.Errorf("failed to init local driver: %w", err)
		}
		factory.RegisterDriver(upload.DriverLocal, localDriver)
		utils.Logger.Info("Upload local driver initialized")
	}

	// 初始化MinIO驱动
	if config.MinIO.Enabled {
		minioDriver, err := NewMinIODriver(config.MinIO)
		if err != nil {
			// MinIO 初始化失败时记录警告，但不阻止应用启动
			utils.Logger.Warn(fmt.Sprintf("Failed to init MinIO driver (will be skipped): %v", err))
		} else {
			factory.RegisterDriver(upload.DriverMinIO, minioDriver)
			utils.Logger.Info("Upload MinIO driver initialized")
		}
	}

	// 初始化阿里云驱动（预留）
	if config.Aliyun.Enabled {
		utils.Logger.Warn("Aliyun OSS driver not implemented yet")
	}

	// 检查默认驱动是否可用
	if !factory.HasDriver(config.DefaultDriver) {
		return fmt.Errorf("default driver '%s' is not available, please enable it in config", config.DefaultDriver)
	}

	utils.Logger.Info(fmt.Sprintf("Upload module initialized with drivers: %v", factory.ListDrivers()))
	return nil
}
