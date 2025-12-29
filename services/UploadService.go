package services

import (
	"RHPRo-Task/database"
	"RHPRo-Task/dto"
	"RHPRo-Task/models"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

// SaveAttachment 保存附件记录到数据库
// 如果提供了 taskID，会记录变更日志
// 如果 attachmentType 为 solution，需要提供 solutionID
// 如果 attachmentType 为 plan，需要提供 planID
func (s *UploadService) SaveAttachment(fileName, fileURL, fileType string, fileSize int64, userID, taskID, solutionID, planID uint, attachmentType string) (*dto.AttachmentResult, error) {
	if attachmentType == "" {
		attachmentType = "general"
	}

	attachment := models.TaskAttachment{
		TaskID:         taskID,
		SolutionID:     solutionID,
		PlanID:         planID,
		FileName:       fileName,
		FileURL:        fileURL,
		FileType:       fileType,
		FileSize:       fileSize,
		UploadedBy:     userID,
		AttachmentType: attachmentType,
	}

	if err := database.DB.Create(&attachment).Error; err != nil {
		return nil, err
	}

	// 如果有 taskID，记录变更日志
	if taskID > 0 && userID > 0 {
		comment := "添加附件: " + fileName
		if attachmentType == "solution" && solutionID > 0 {
			comment = "添加方案附件: " + fileName
		} else if attachmentType == "plan" && planID > 0 {
			comment = "添加执行计划附件: " + fileName
		}

		changeLog := models.TaskChangeLog{
			TaskID:     taskID,
			UserID:     userID,
			ChangeType: "attachment_add",
			FieldName:  "attachment",
			OldValue:   "",
			NewValue:   fileName,
			Comment:    comment,
		}
		database.DB.Create(&changeLog)
	}

	return &dto.AttachmentResult{
		ID:       attachment.ID,
		FileName: attachment.FileName,
		FileURL:  attachment.FileURL,
		FileType: attachment.FileType,
		FileSize: attachment.FileSize,
	}, nil
}

// BindAttachmentsToTask 将附件绑定到任务
func (s *UploadService) BindAttachmentsToTask(attachmentIDs []uint, taskID uint) error {
	if len(attachmentIDs) == 0 {
		return nil
	}
	return database.DB.Model(&models.TaskAttachment{}).
		Where("id IN ?", attachmentIDs).
		Update("task_id", taskID).Error
}

// BindAttachmentsToSolution 将附件绑定到方案
func (s *UploadService) BindAttachmentsToSolution(attachmentIDs []uint, taskID, solutionID uint) error {
	if len(attachmentIDs) == 0 {
		return nil
	}
	return database.DB.Model(&models.TaskAttachment{}).
		Where("id IN ?", attachmentIDs).
		Updates(map[string]interface{}{
			"task_id":     taskID,
			"solution_id": solutionID,
		}).Error
}

// BindAttachmentsToPlan 将附件绑定到执行计划
func (s *UploadService) BindAttachmentsToPlan(attachmentIDs []uint, taskID, planID uint) error {
	if len(attachmentIDs) == 0 {
		return nil
	}
	return database.DB.Model(&models.TaskAttachment{}).
		Where("id IN ?", attachmentIDs).
		Updates(map[string]interface{}{
			"task_id": taskID,
			"plan_id": planID,
		}).Error
}

// GetAttachmentsByTaskID 获取任务的所有附件（包含上传人信息）
func (s *UploadService) GetAttachmentsByTaskID(taskID uint) ([]dto.AttachmentDetailResult, error) {
	var attachments []models.TaskAttachment
	if err := database.DB.Where("task_id = ?", taskID).Find(&attachments).Error; err != nil {
		return nil, err
	}

	return s.toAttachmentDetailResults(attachments), nil
}

// GetTaskOwnAttachments 获取任务本身的附件（不含方案和计划附件）
func (s *UploadService) GetTaskOwnAttachments(taskID uint) []dto.AttachmentDetailResult {
	var attachments []models.TaskAttachment
	if err := database.DB.Where("task_id = ? AND solution_id = 0 AND plan_id = 0", taskID).Find(&attachments).Error; err != nil {
		return nil
	}

	return s.toAttachmentDetailResults(attachments)
}

// GetSolutionAttachments 获取方案的附件
func (s *UploadService) GetSolutionAttachments(solutionID uint) []dto.AttachmentDetailResult {
	var attachments []models.TaskAttachment
	if err := database.DB.Where("solution_id = ?", solutionID).Find(&attachments).Error; err != nil {
		return nil
	}

	return s.toAttachmentDetailResults(attachments)
}

// GetPlanAttachments 获取执行计划的附件
func (s *UploadService) GetPlanAttachments(planID uint) []dto.AttachmentDetailResult {
	var attachments []models.TaskAttachment
	if err := database.DB.Where("plan_id = ?", planID).Find(&attachments).Error; err != nil {
		return nil
	}

	return s.toAttachmentDetailResults(attachments)
}

// toAttachmentDetailResults 转换附件列表为响应结果
func (s *UploadService) toAttachmentDetailResults(attachments []models.TaskAttachment) []dto.AttachmentDetailResult {
	results := make([]dto.AttachmentDetailResult, 0, len(attachments))
	for _, att := range attachments {
		result := dto.AttachmentDetailResult{
			ID:             att.ID,
			TaskID:         att.TaskID,
			SolutionID:     att.SolutionID,
			PlanID:         att.PlanID,
			FileName:       att.FileName,
			FileURL:        att.FileURL,
			FileType:       att.FileType,
			FileSize:       att.FileSize,
			AttachmentType: att.AttachmentType,
			UploadedBy:     att.UploadedBy,
			CreatedAt:      dto.ToResponseTime(att.CreatedAt),
		}

		// 查询上传人信息
		if att.UploadedBy > 0 {
			var user models.User
			if err := database.DB.Select("username, nickname").First(&user, att.UploadedBy).Error; err == nil {
				result.UploaderUsername = user.Username
				result.UploaderNickname = user.Nickname
			}
		}

		results = append(results, result)
	}

	return results
}

// DeleteAttachment 删除附件记录
// 如果提供了 taskID 和 userID，会记录变更日志
func (s *UploadService) DeleteAttachment(attachmentID uint, taskID uint, userID uint) error {
	// 先查询附件信息用于记录日志
	var attachment models.TaskAttachment
	if err := database.DB.First(&attachment, attachmentID).Error; err != nil {
		return err
	}

	// 删除附件
	if err := database.DB.Delete(&models.TaskAttachment{}, attachmentID).Error; err != nil {
		return err
	}

	// 如果有 taskID，记录变更日志
	if taskID > 0 && userID > 0 {
		changeLog := models.TaskChangeLog{
			TaskID:     taskID,
			UserID:     userID,
			ChangeType: "attachment_delete",
			FieldName:  "attachment",
			OldValue:   attachment.FileName,
			NewValue:   "",
			Comment:    "删除附件: " + attachment.FileName,
		}
		database.DB.Create(&changeLog)
	}

	return nil
}
