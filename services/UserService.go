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
	// 检查用户名是否存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.UserName).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否存在
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已被注册")
	}
	// 检查手机号是否存在
	if err := database.DB.Where("mobile = ?", req.Mobile).First(&existingUser).Error; err == nil {
		return nil, errors.New("手机号已被注册")
	}

	// 创建用户
	user := &models.User{
		Username: req.UserName,
		Nickname: req.Nickname,
		Email:    req.Email,
		Mobile:   req.Mobile,
		// DepartmentID: req.DepartmentID,
		Status: models.UserStatusPending, // 默认为待审核
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
	// 查询用户
	var user models.User
	if err := database.DB.Preload("Roles.Permissions").Preload("Department").Preload("ManagedDepartments").Where("username = ?", req.UserName).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status == models.UserStatusDisabled {
		return nil, errors.New("用户已被禁用")
	}
	if user.Status == models.UserStatusPending {
		return nil, errors.New("用户待审核，请联系管理员")
	}

	// 获取部门信息
	var deptID uint
	var deptName string
	if user.Department != nil {
		deptID = user.Department.ID
		deptName = user.Department.Name
	}

	// 获取负责的部门信息
	var managedDeptIDs []uint
	isLeader := false
	if len(user.ManagedDepartments) > 0 {
		isLeader = true
		for _, dept := range user.ManagedDepartments {
			managedDeptIDs = append(managedDeptIDs, dept.ID)
		}
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Mobile, deptID, deptName, isLeader, managedDeptIDs, 24)
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

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, req *dto.UpdateUserRequest) error {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 检查用户名唯一性（如果要更新用户名）
	if req.Username != "" && req.Username != user.Username {
		var existingUser models.User
		if err := database.DB.Where("username = ? AND id != ?", req.Username, id).First(&existingUser).Error; err == nil {
			return errors.New("用户名已被使用")
		}
		user.Username = req.Username
	}

	// 检查邮箱唯一性（如果要更新邮箱）
	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return errors.New("邮箱已被使用")
		}
		user.Email = req.Email
	}

	// 检查手机号唯一性（如果要更新手机号）
	if req.Mobile != "" && req.Mobile != user.Mobile {
		var existingUser models.User
		if err := database.DB.Where("mobile = ? AND id != ?", req.Mobile, id).First(&existingUser).Error; err == nil {
			return errors.New("手机号已被使用")
		}
		user.Mobile = req.Mobile
	}

	// 更新昵称
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	return database.DB.Save(&user).Error
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
	// 复用 Register 的检查逻辑，但状态设置为 Active
	// 检查用户名是否存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.UserName).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已被注册")
	}
	if req.Mobile != "" {
		if err := database.DB.Where("mobile = ?", req.Mobile).First(&existingUser).Error; err == nil {
			return nil, errors.New("手机号已被注册")
		}
	}

	user := &models.User{
		Username: req.UserName,
		Nickname: req.Nickname,
		Email:    req.Email,
		Mobile:   req.Mobile,
		// DepartmentID: req.DepartmentID,
		Status: models.UserStatusActive, // 管理员创建直接激活
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
