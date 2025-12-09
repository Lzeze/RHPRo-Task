package models

import "gorm.io/gorm"

// PreloadUserWithRolesPermissions 示例：预加载用户及其角色和角色的权限
func PreloadUserWithRolesPermissions(db *gorm.DB, userID uint) (*User, error) {
	var u User
	if err := db.Preload("Roles.Permissions").First(&u, userID).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// PreloadRoleWithPermissions 示例：预加载角色及其权限
func PreloadRoleWithPermissions(db *gorm.DB, roleID uint) (*Role, error) {
	var r Role
	if err := db.Preload("Permissions").First(&r, roleID).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// PreloadTaskWithTags 示例：预加载任务及其标签
func PreloadTaskWithTags(db *gorm.DB, taskID uint) (*Task, error) {
	var t Task
	if err := db.Preload("Tags").First(&t, taskID).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
