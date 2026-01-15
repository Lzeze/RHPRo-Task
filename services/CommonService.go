package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/models"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CommonService 公共服务，提供编号生成、通用工具方法
type CommonService struct{}

// ========================================
// 编号生成方法
// ========================================

// GenerateProductNo 生成产品主线编号
// 格式：PRD-{年份}-{序号}，如 PRD-2026-001
func (s *CommonService) GenerateProductNo() (string, error) {
	year := time.Now().Format("2006")
	prefix := fmt.Sprintf("PRD-%s-", year)

	// 查询当前年份最大序号
	var maxNo string
	err := database.DB.Model(&models.ProductLine{}).
		Where("product_no LIKE ?", prefix+"%").
		Order("product_no DESC").
		Limit(1).
		Pluck("product_no", &maxNo).Error

	if err != nil {
		return "", fmt.Errorf("查询产品编号失败: %v", err)
	}

	// 计算下一个序号
	nextSeq := 1
	if maxNo != "" {
		// 提取序号部分 PRD-2026-001 -> 001
		parts := strings.Split(maxNo, "-")
		if len(parts) == 3 {
			if seq, err := strconv.Atoi(parts[2]); err == nil {
				nextSeq = seq + 1
			}
		}
	}

	// 生成新编号
	return fmt.Sprintf("%s%03d", prefix, nextSeq), nil
}

// GenerateAnnualPlanNo 生成年度计划编号
// 格式：AP-{年份}-{序号}，如 AP-2026-001
func (s *CommonService) GenerateAnnualPlanNo() (string, error) {
	year := time.Now().Format("2006")
	prefix := fmt.Sprintf("AP-%s-", year)

	// 查询当前年份最大序号
	var maxNo string
	err := database.DB.Model(&models.AnnualPlan{}).
		Where("plan_no LIKE ?", prefix+"%").
		Order("plan_no DESC").
		Limit(1).
		Pluck("plan_no", &maxNo).Error

	if err != nil {
		return "", fmt.Errorf("查询年度计划编号失败: %v", err)
	}

	// 计算下一个序号
	nextSeq := 1
	if maxNo != "" {
		parts := strings.Split(maxNo, "-")
		if len(parts) == 3 {
			if seq, err := strconv.Atoi(parts[2]); err == nil {
				nextSeq = seq + 1
			}
		}
	}

	return fmt.Sprintf("%s%03d", prefix, nextSeq), nil
}

// GeneratePlanNodeNo 生成计划节点编号
// 格式：PN-{年份}-{序号}，如 PN-2026-001
func (s *CommonService) GeneratePlanNodeNo() (string, error) {
	year := time.Now().Format("2006")
	prefix := fmt.Sprintf("PN-%s-", year)

	// 查询当前年份最大序号
	var maxNo string
	err := database.DB.Model(&models.PlanNode{}).
		Where("node_no LIKE ?", prefix+"%").
		Order("node_no DESC").
		Limit(1).
		Pluck("node_no", &maxNo).Error

	if err != nil {
		return "", fmt.Errorf("查询计划节点编号失败: %v", err)
	}

	// 计算下一个序号
	nextSeq := 1
	if maxNo != "" {
		parts := strings.Split(maxNo, "-")
		if len(parts) == 3 {
			if seq, err := strconv.Atoi(parts[2]); err == nil {
				nextSeq = seq + 1
			}
		}
	}

	return fmt.Sprintf("%s%03d", prefix, nextSeq), nil
}

// ========================================
// 阶段相关方法
// ========================================

// PlanStage 计划阶段常量
const (
	StageGermination = "germination" // 萌芽期
	StageExperiment  = "experiment"  // 试验期
	StageMaturity    = "maturity"    // 成熟期
	StagePromotion   = "promotion"   // 推广期
)

// StageOrder 阶段顺序映射
var StageOrder = map[string]int{
	StageGermination: 1,
	StageExperiment:  2,
	StageMaturity:    3,
	StagePromotion:   4,
}

// StageName 阶段中文名称映射
var StageName = map[string]string{
	StageGermination: "萌芽期",
	StageExperiment:  "试验期",
	StageMaturity:    "成熟期",
	StagePromotion:   "推广期",
}

// ValidStages 有效阶段列表
var ValidStages = []string{StageGermination, StageExperiment, StageMaturity, StagePromotion}

// IsValidStage 验证阶段是否有效
func (s *CommonService) IsValidStage(stage string) bool {
	_, exists := StageOrder[stage]
	return exists
}

// GetStageName 获取阶段中文名称
func (s *CommonService) GetStageName(stage string) string {
	if name, exists := StageName[stage]; exists {
		return name
	}
	return stage
}

// GetStageOrder 获取阶段顺序
func (s *CommonService) GetStageOrder(stage string) int {
	if order, exists := StageOrder[stage]; exists {
		return order
	}
	return 0
}

// GetNextStage 获取下一阶段
// 返回下一阶段代码，如果已是最后阶段则返回空字符串
func (s *CommonService) GetNextStage(currentStage string) string {
	currentOrder := s.GetStageOrder(currentStage)
	if currentOrder == 0 || currentOrder >= 4 {
		return ""
	}

	for stage, order := range StageOrder {
		if order == currentOrder+1 {
			return stage
		}
	}
	return ""
}

// ValidateStageProgression 验证阶段递进是否有效
// 检查目标阶段是否是源阶段的下一阶段
func (s *CommonService) ValidateStageProgression(sourceStage, targetStage string) bool {
	sourceOrder := s.GetStageOrder(sourceStage)
	targetOrder := s.GetStageOrder(targetStage)

	if sourceOrder == 0 || targetOrder == 0 {
		return false
	}

	return targetOrder == sourceOrder+1
}

// ========================================
// 节点层级计算方法
// ========================================

// CalculateNodeLevel 计算节点层级
// 根节点为0，子节点为父节点层级+1
func (s *CommonService) CalculateNodeLevel(parentNodeID *uint) (int, error) {
	if parentNodeID == nil {
		return 0, nil
	}

	var parentNode models.PlanNode
	if err := database.DB.Select("node_level").First(&parentNode, *parentNodeID).Error; err != nil {
		return 0, fmt.Errorf("父节点不存在: %v", err)
	}

	return parentNode.NodeLevel + 1, nil
}

// CalculateNodePath 计算节点路径
// 格式：父节点路径/父节点ID，根节点路径为空
func (s *CommonService) CalculateNodePath(parentNodeID *uint) (string, error) {
	if parentNodeID == nil {
		return "", nil
	}

	var parentNode models.PlanNode
	if err := database.DB.Select("id, node_path").First(&parentNode, *parentNodeID).Error; err != nil {
		return "", fmt.Errorf("父节点不存在: %v", err)
	}

	if parentNode.NodePath == "" {
		return fmt.Sprintf("%d", parentNode.ID), nil
	}
	return fmt.Sprintf("%s/%d", parentNode.NodePath, parentNode.ID), nil
}

// GetRootNodeID 获取根节点ID
// 如果父节点有根节点ID则继承，否则父节点就是根节点
func (s *CommonService) GetRootNodeID(parentNodeID *uint) (*uint, error) {
	if parentNodeID == nil {
		return nil, nil
	}

	var parentNode models.PlanNode
	if err := database.DB.Select("id, root_node_id").First(&parentNode, *parentNodeID).Error; err != nil {
		return nil, fmt.Errorf("父节点不存在: %v", err)
	}

	if parentNode.RootNodeID != nil {
		return parentNode.RootNodeID, nil
	}
	return parentNodeID, nil
}

// GetNextSortOrder 获取下一个排序序号
// 查询同级节点的最大排序序号+1
func (s *CommonService) GetNextSortOrder(annualPlanID uint, parentNodeID *uint) (int, error) {
	var maxOrder int

	query := database.DB.Model(&models.PlanNode{}).
		Where("annual_plan_id = ? AND deleted_at IS NULL", annualPlanID)

	if parentNodeID != nil {
		query = query.Where("parent_node_id = ?", *parentNodeID)
	} else {
		query = query.Where("parent_node_id IS NULL")
	}

	if err := query.Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder).Error; err != nil {
		return 0, fmt.Errorf("查询排序序号失败: %v", err)
	}

	return maxOrder + 1, nil
}

// ========================================
// 统计计算方法
// ========================================

// CalculateCompletionRate 计算完成率
// 返回百分比，保留两位小数
func (s *CommonService) CalculateCompletionRate(completed, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(completed) / float64(total) * 100
}

// ========================================
// 权限检查辅助方法
// ========================================

// IsDepartmentLeader 检查用户是否是指定部门的负责人
func (s *CommonService) IsDepartmentLeader(userID, departmentID uint) bool {
	var count int64
	database.DB.Model(&models.DepartmentLeader{}).
		Where("user_id = ? AND department_id = ? AND deleted_at IS NULL", userID, departmentID).
		Count(&count)
	return count > 0
}

// IsAnyDepartmentLeader 检查用户是否是任意部门的负责人
func (s *CommonService) IsAnyDepartmentLeader(userID uint) bool {
	var count int64
	database.DB.Model(&models.DepartmentLeader{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count)
	return count > 0
}

// GetUserManagedDepartmentIDs 获取用户管理的所有部门ID
func (s *CommonService) GetUserManagedDepartmentIDs(userID uint) []uint {
	var departmentIDs []uint
	database.DB.Model(&models.DepartmentLeader{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Pluck("department_id", &departmentIDs)
	return departmentIDs
}

// IsSuperAdmin 检查用户是否是超级管理员
func (s *CommonService) IsSuperAdmin(userID uint) bool {
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		return false
	}

	for _, role := range user.Roles {
		if role.Name == "admin" {
			return true
		}
	}
	return false
}
