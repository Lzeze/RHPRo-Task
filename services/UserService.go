package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/dto"
	"RHPRo-Task/models"
	"RHPRo-Task/utils"
	"errors"
	"math"

	"gorm.io/gorm"
)

type UserService struct{}

// Register 用户注册
func (s *UserService) Register(req *dto.RegisterRequest) (*models.User, error) {
	var existingUser models.User
	// 检查手机号是否存在（手机号为必填且唯一）
	if err := database.DB.Where("mobile = ?", req.Mobile).First(&existingUser).Error; err == nil {
		return nil, errors.New("手机号已被注册")
	}

	// 检查邮箱是否存在（如果填写了邮箱）
	if req.Email != "" {
		if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已被注册")
		}
	}

	// 创建用户
	user := &models.User{
		Mobile:   req.Mobile,
		Username: req.UserName,
		Email:    req.Email,
		Status:   models.UserStatusPending, // 默认为待审核
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// 保存到数据库
	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// 分配默认角色（普通用户）
	var userRole models.Role
	if err := database.DB.Where("name = ?", "user").First(&userRole).Error; err == nil {
		database.DB.Model(user).Association("Roles").Append(&userRole)
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 通过手机号查询用户
	var user models.User
	if err := database.DB.Preload("Roles.Permissions").Preload("Department").Preload("ManagedDepartments").Where("mobile = ?", req.Mobile).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("手机号或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("手机号或密码错误")
	}

	// 检查用户状态
	if user.Status == models.UserStatusDisabled {
		return nil, errors.New("用户已被禁用")
	}
	if user.Status == models.UserStatusPending {
		return nil, errors.New("用户待审核，请联系管理员")
	}

	// // 获取部门信息
	// var deptID uint
	// var deptName string
	// if user.Department != nil {
	// 	deptID = user.Department.ID
	// 	deptName = user.Department.Name
	// }

	// // 获取负责的部门信息
	// var managedDeptIDs []uint
	// if len(user.ManagedDepartments) > 0 {
	// 	for _, dept := range user.ManagedDepartments {
	// 		managedDeptIDs = append(managedDeptIDs, dept.ID)
	// 	}
	// }

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Mobile, 24)
	if err != nil {
		return nil, err
	}

	// 构建用户信息响应
	userInfo := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Status:   user.Status,
	}

	if user.Department != nil {
		userInfo.Department = &dto.DepartmentResponse{
			ID:   user.Department.ID,
			Name: user.Department.Name,
		}
	}

	// 添加管理的部门信息
	if len(user.ManagedDepartments) > 0 {
		for _, dept := range user.ManagedDepartments {
			userInfo.ManagedDepartments = append(userInfo.ManagedDepartments, dto.DepartmentResponse{
				ID:   dept.ID,
				Name: dept.Name,
			})
		}
	}

	// 添加角色信息
	for _, role := range user.Roles {
		roleResp := dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		}
		for _, perm := range role.Permissions {
			roleResp.Permissions = append(roleResp.Permissions, dto.PermissionResponse{
				ID:          perm.ID,
				Name:        perm.Name,
				Description: perm.Description,
			})
		}
		userInfo.Roles = append(userInfo.Roles, roleResp)
	}

	return &dto.LoginResponse{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	var user models.User
	if err := database.DB.Preload("Roles").Preload("Department").Preload("ManagedDepartments").First(&user, id).Error; err != nil {
		return nil, err
	}

	userResp := &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Status:   user.Status,
	}

	if user.Department != nil {
		userResp.Department = &dto.DepartmentResponse{
			ID:   user.Department.ID,
			Name: user.Department.Name,
		}
	}

	// 添加管理的部门信息
	if len(user.ManagedDepartments) > 0 {
		for _, dept := range user.ManagedDepartments {
			userResp.ManagedDepartments = append(userResp.ManagedDepartments, dto.DepartmentResponse{
				ID:   dept.ID,
				Name: dept.Name,
			})
		}
	}

	for _, role := range user.Roles {
		userResp.Roles = append(userResp.Roles, dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	return userResp, nil
}

// GetUserList 获取用户列表（分页）
func (s *UserService) GetUserList(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	var users []models.User
	var total int64

	page := pagination.GetPage()
	pageSize := pagination.GetPageSize()
	offset := (page - 1) * pageSize

	// 查询总数
	database.DB.Model(&models.User{}).Count(&total)

	// 查询列表
	if err := database.DB.Preload("Roles").Preload("Department").Preload("ManagedDepartments").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	var userList []dto.UserResponse
	for _, user := range users {
		userResp := dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Mobile:   user.Mobile,
			Status:   user.Status,
		}
		if user.Department != nil {
			userResp.Department = &dto.DepartmentResponse{
				ID:   user.Department.ID,
				Name: user.Department.Name,
			}
		}
		// 添加管理的部门信息
		if len(user.ManagedDepartments) > 0 {
			for _, dept := range user.ManagedDepartments {
				userResp.ManagedDepartments = append(userResp.ManagedDepartments, dto.DepartmentResponse{
					ID:   dept.ID,
					Name: dept.Name,
				})
			}
		}
		for _, role := range user.Roles {
			userResp.Roles = append(userResp.Roles, dto.RoleResponse{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
			})
		}
		userList = append(userList, userResp)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PaginationResponse{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Data:       userList,
	}, nil
}

// UpdateUser 更新用户信息（只更新请求中有值的字段）
func (s *UserService) UpdateUser(id uint, req *dto.UpdateUserRequest) error {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 构建更新字段map
	updates := make(map[string]interface{})

	// 检查用户名唯一性（如果要更新用户名）
	if req.Username != "" && req.Username != user.Username {
		var existingUser models.User
		if err := database.DB.Where("username = ? AND id != ?", req.Username, id).First(&existingUser).Error; err == nil {
			return errors.New("用户名已被使用")
		}
		updates["username"] = req.Username
	}

	// 检查邮箱唯一性（如果要更新邮箱）
	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return errors.New("邮箱已被使用")
		}
		updates["email"] = req.Email
	}

	// 检查手机号唯一性（如果要更新手机号）
	if req.Mobile != "" && req.Mobile != user.Mobile {
		var existingUser models.User
		if err := database.DB.Where("mobile = ? AND id != ?", req.Mobile, id).First(&existingUser).Error; err == nil {
			return errors.New("手机号已被使用")
		}
		updates["mobile"] = req.Mobile
	}

	// 更新昵称
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}

	// 没有需要更新的字段
	if len(updates) == 0 {
		return nil
	}

	return database.DB.Model(&user).Updates(updates).Error
}

// AssignRoles 分配角色
func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	var roles []*models.Role
	if err := database.DB.Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	// 清空现有角色并分配新角色
	database.DB.Model(&user).Association("Roles").Clear()
	return database.DB.Model(&user).Association("Roles").Append(roles)
}

// CreateUser 管理员创建用户（直接激活）
func (s *UserService) CreateUser(req *dto.RegisterRequest) (*models.User, error) {
	var existingUser models.User
	// 检查手机号是否存在
	if err := database.DB.Where("mobile = ?", req.Mobile).First(&existingUser).Error; err == nil {
		return nil, errors.New("手机号已被注册")
	}
	// 检查邮箱是否存在（如果填写了邮箱）
	if req.Email != "" {
		if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已被注册")
		}
	}

	user := &models.User{
		Mobile:   req.Mobile,
		Username: req.UserName,
		Email:    req.Email,
		Status:   models.UserStatusActive, // 管理员创建直接激活
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// 分配默认角色
	var userRole models.Role
	if err := database.DB.Where("name = ?", "user").First(&userRole).Error; err == nil {
		database.DB.Model(user).Association("Roles").Append(&userRole)
	}

	return user, nil
}

// ApproveUser 审核用户
func (s *UserService) ApproveUser(id uint) error {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return err
	}

	if user.Status == models.UserStatusActive {
		return errors.New("用户已是正常状态")
	}

	user.Status = models.UserStatusActive
	return database.DB.Save(&user).Error
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(id uint) error {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}
	return database.DB.Delete(&user).Error
}

// DisableUser 禁用用户
func (s *UserService) DisableUser(id uint) error {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	if user.Status == models.UserStatusDisabled {
		return errors.New("用户已是禁用状态")
	}

	user.Status = models.UserStatusDisabled
	return database.DB.Save(&user).Error
}

// GetAssignableUsers 获取可指派的执行人列表
// 包括：1.同部门的所有成员 2.其他部门的负责人（如果负责多部门则产生多条记录）
func (s *UserService) GetAssignableUsers(userID uint, req *dto.GetAssignableUsersRequest) ([]dto.AssignableUserResponse, error) {
	// 获取当前用户信息（用于确定所在部门）
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		return nil, errors.New("当前用户不存在")
	}

	if currentUser.DepartmentID == nil {
		return nil, errors.New("当前用户未指派部门")
	}

	var results []dto.AssignableUserResponse

	// ========== 1. 同部门的所有成员（不包括被禁用的用户） ==========
	var sameDepUsers []models.User
	query := database.DB.Where("department_id = ? AND status != ?", *currentUser.DepartmentID, models.UserStatusDisabled)

	// 如果提供了关键词，则进行模糊搜索
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if err := query.Find(&sameDepUsers).Error; err != nil {
		return nil, err
	}

	// 获取当前部门信息
	var currentDep models.Department
	if err := database.DB.First(&currentDep, *currentUser.DepartmentID).Error; err != nil {
		return nil, errors.New("当前部门不存在")
	}

	// 添加同部门成员到结果
	for _, user := range sameDepUsers {
		results = append(results, dto.AssignableUserResponse{
			ID:                 user.ID,
			Username:           user.Username,
			Nickname:           user.Nickname,
			Email:              user.Email,
			DepartmentID:       currentDep.ID,
			DepartmentName:     currentDep.Name,
			IsDepartmentLeader: user.IsDepartmentLeader,
		})
	}

	// ========== 2. 其他部门的负责人（使用GORM关联查询） ==========
	// 获取所有部门领导者的部门ID及对应用户信息
	var leaders []struct {
		UserID       uint
		Nickname     string
		Username     string
		Email        string
		DepartmentID uint
		DepName      string
	}

	// 使用GORM的Joins进行关联查询
	query = database.DB.
		Table("users u").
		Select("DISTINCT u.id as user_id, u.nickname,u.username, u.email, d.id as department_id, d.name as dep_name").
		Joins("INNER JOIN department_leaders dl ON u.id = dl.user_id").
		Joins("INNER JOIN departments d ON dl.department_id = d.id").
		Where("u.status != ? AND u.id != ? AND dl.department_id != ?",
			models.UserStatusDisabled, userID, *currentUser.DepartmentID)

	// 如果提供了关键词，则进行模糊搜索
	if req.Keyword != "" {
		query = query.Where("u.username LIKE ? OR u.nickname LIKE ? OR u.email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if err := query.Scan(&leaders).Error; err != nil {
		return nil, err
	}

	// 添加其他部门的负责人到结果
	for _, leader := range leaders {
		results = append(results, dto.AssignableUserResponse{
			ID:                 leader.UserID,
			Username:           leader.Username,
			Nickname:           leader.Nickname,
			Email:              leader.Email,
			DepartmentID:       leader.DepartmentID,
			DepartmentName:     leader.DepName,
			IsDepartmentLeader: true,
		})
	}

	return results, nil
}
