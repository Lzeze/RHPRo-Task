package controllers

import (
	"RHPRo-Task/services"
	"RHPRo-Task/upload"
	"RHPRo-Task/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	uploader      *upload.Uploader
	uploadService *services.UploadService
}

func NewUploadController() *UploadController {
	return &UploadController{
		uploader:      upload.NewUploader(),
		uploadService: services.NewUploadService(),
	}
}

// Upload 通用文件上传
// @Summary 上传文件
// @Description 上传文件并保存附件记录，返回附件ID。task_id可为空，后续创建任务时再绑定
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Param directory formData string false "存储目录"
// @Param task_id formData int false "关联任务ID（可选）"
// @Param attachment_type formData string false "附件类型(requirement/solution/plan/general/task)"
// @Success 200 {object} dto.AttachmentResult "上传成功"
// @Failure 400 {object} utils.Response "请求错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Security BearerAuth
// @Router /upload [post]
func (c *UploadController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.BadRequest(ctx, "请选择要上传的文件")
		return
	}

	opts := upload.UploadOptions{
		Directory: ctx.PostForm("directory"),
	}

	// 上传文件
	info, err := c.uploader.UploadFile(ctx.Request.Context(), file, opts)
	if err != nil {
		utils.Error(ctx, 500, "上传失败: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := ctx.Get("userID")
	var uploadedBy uint
	if userID != nil {
		uploadedBy = userID.(uint)
	}

	// 获取可选参数
	var taskID uint
	if tid := ctx.PostForm("task_id"); tid != "" {
		if tidVal, err := strconv.ParseUint(tid, 10, 32); err == nil {
			taskID = uint(tidVal)
		}
	}
	attachmentType := ctx.PostForm("attachment_type")

	// 保存附件记录
	result, err := c.uploadService.SaveAttachment(info.FileName, info.URL, info.MimeType, info.Size, uploadedBy, taskID, attachmentType)
	if err != nil {
		utils.Error(ctx, 500, "保存附件记录失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "上传成功", result)
}

// UploadWithDriver 使用指定驱动上传
// @Summary 使用指定驱动上传文件
// @Description 使用指定的存储驱动上传文件并保存附件记录
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Param driver formData string true "驱动类型(local/minio/aliyun)"
// @Param directory formData string false "存储目录"
// @Param task_id formData int false "关联任务ID（可选）"
// @Param attachment_type formData string false "附件类型(requirement/solution/plan/general/task)"
// @Success 200 {object} dto.AttachmentResult "上传成功"
// @Failure 400 {object} utils.Response "请求错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Security BearerAuth
// @Router /upload/driver [post]
func (c *UploadController) UploadWithDriver(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.BadRequest(ctx, "请选择要上传的文件")
		return
	}

	driverType := ctx.PostForm("driver")
	if driverType == "" {
		utils.BadRequest(ctx, "请指定存储驱动")
		return
	}

	opts := upload.UploadOptions{
		Directory: ctx.PostForm("directory"),
	}

	info, err := c.uploader.UploadFileWithDriver(ctx.Request.Context(), upload.DriverType(driverType), file, opts)
	if err != nil {
		utils.Error(ctx, 500, "上传失败: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := ctx.Get("userID")
	var uploadedBy uint
	if userID != nil {
		uploadedBy = userID.(uint)
	}

	// 获取可选参数
	var taskID uint
	if tid := ctx.PostForm("task_id"); tid != "" {
		if tidVal, err := strconv.ParseUint(tid, 10, 32); err == nil {
			taskID = uint(tidVal)
		}
	}
	attachmentType := ctx.PostForm("attachment_type")

	// 保存附件记录
	result, err := c.uploadService.SaveAttachment(info.FileName, info.URL, info.MimeType, info.Size, uploadedBy, taskID, attachmentType)
	if err != nil {
		utils.Error(ctx, 500, "保存附件记录失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "上传成功", result)
}

// UploadAvatar 上传头像
// @Summary 上传头像
// @Description 上传用户头像并保存附件记录
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "头像文件"
// @Success 200 {object} dto.AttachmentResult "上传成功"
// @Failure 400 {object} utils.Response "请求错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Security BearerAuth
// @Router /upload/avatar [post]
func (c *UploadController) UploadAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.BadRequest(ctx, "请选择要上传的头像")
		return
	}

	// 验证是否为图片
	if !upload.IsImageFile(file.Filename) {
		utils.BadRequest(ctx, "请上传图片文件")
		return
	}

	opts := upload.UploadOptions{
		Directory: "avatars",
	}

	info, err := c.uploader.UploadFileWithDriver(ctx.Request.Context(), upload.DriverMinIO, file, opts)
	if err != nil {
		utils.Error(ctx, 500, "上传失败: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := ctx.Get("userID")
	var uploadedBy uint
	if userID != nil {
		uploadedBy = userID.(uint)
	}

	// 保存附件记录
	result, err := c.uploadService.SaveAttachment(info.FileName, info.URL, info.MimeType, info.Size, uploadedBy, 0, "avatar")
	if err != nil {
		utils.Error(ctx, 500, "保存附件记录失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "上传成功", result)
}

// UploadMedia 上传媒体文件
// @Summary 上传媒体文件
// @Description 上传视频或音频文件并保存附件记录
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "媒体文件"
// @Param task_id formData int false "关联任务ID（可选）"
// @Success 200 {object} dto.AttachmentResult "上传成功"
// @Failure 400 {object} utils.Response "请求错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Security BearerAuth
// @Router /upload/media [post]
func (c *UploadController) UploadMedia(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.BadRequest(ctx, "请选择要上传的媒体文件")
		return
	}

	// 验证是否为媒体文件
	if !upload.IsVideoFile(file.Filename) && !upload.IsAudioFile(file.Filename) {
		utils.BadRequest(ctx, "请上传视频或音频文件")
		return
	}

	opts := upload.UploadOptions{
		Directory: "media",
	}

	info, err := c.uploader.UploadFileWithDriver(ctx.Request.Context(), upload.DriverMinIO, file, opts)
	if err != nil {
		utils.Error(ctx, 500, "上传失败: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := ctx.Get("userID")
	var uploadedBy uint
	if userID != nil {
		uploadedBy = userID.(uint)
	}

	// 获取可选参数
	var taskID uint
	if tid := ctx.PostForm("task_id"); tid != "" {
		if tidVal, err := strconv.ParseUint(tid, 10, 32); err == nil {
			taskID = uint(tidVal)
		}
	}

	// 保存附件记录
	result, err := c.uploadService.SaveAttachment(info.FileName, info.URL, info.MimeType, info.Size, uploadedBy, taskID, "media")
	if err != nil {
		utils.Error(ctx, 500, "保存附件记录失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "上传成功", result)
}

// GetDrivers 获取可用的存储驱动列表
// @Summary 获取存储驱动列表
// @Description 获取当前可用的存储驱动列表
// @Tags 文件上传
// @Produce json
// @Success 200 {object} utils.Response "驱动列表"
// @Security BearerAuth
// @Router /upload/drivers [get]
func (c *UploadController) GetDrivers(ctx *gin.Context) {
	factory := upload.GetFactory()
	drivers := factory.ListDrivers()

	utils.Success(ctx, gin.H{
		"drivers": drivers,
	})
}
