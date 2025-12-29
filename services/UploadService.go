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
func (s *UploadService) SaveAttachment(fileName, fileURL, fileType string, fileSize int64, userID, taskID uint, attachmentType string) (*dto.AttachmentResult, error) {
	if attachmentType == "" {
		attachmentType = "general"
	}

	attachment := models.TaskAttachment{
		TaskID:         taskID,
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

// GetAttachmentsByTaskID 获取任务的所有附件（包含上传人信息）
func (s *UploadService) GetAttachmentsByTaskID(taskID uint) ([]dto.AttachmentDetailResult, error) {
	var attachments []models.TaskAttachment
	if err := database.DB.Where("task_id = ?", taskID).Find(&attachments).Error; err != nil {
		return nil, err
	}

	results := make([]dto.AttachmentDetailResult, 0, len(attachments))
	for _, att := range attachments {
		result := dto.AttachmentDetailResult{
			ID:             att.ID,
			TaskID:         att.TaskID,
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

	return results, nil
}

// DeleteAttachment 删除附件记录
func (s *UploadService) DeleteAttachment(attachmentID uint) error {
	return database.DB.Delete(&models.TaskAttachment{}, attachmentID).Error
}
