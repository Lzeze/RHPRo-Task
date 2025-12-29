package controllers

import (
	"RHPRo-Task/upload"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	uploader *upload.Uploader
}

func NewUploadController() *UploadController {
	return &UploadController{
		uploader: upload.NewUploader(),
	}
}

// Upload 通用文件上传
// @Summary 上传文件
// @Description 根据文件类型自动选择存储驱动上传文件
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Param directory formData string false "存储目录"
// @Success 200 {object} map[string]interface{} "上传成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Security BearerAuth
// @Router /upload [post]
func (c *UploadController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择要上传的文件",
		})
		return
	}

	opts := upload.UploadOptions{
		Directory: ctx.PostForm("directory"),
	}

	info, err := c.uploader.UploadFile(ctx.Request.Context(), file, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data":    info,
	})
}

// UploadWithDriver 使用指定驱动上传
// @Summary 使用指定驱动上传文件
// @Description 使用指定的存储驱动上传文件
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Param driver formData string true "驱动类型(local/minio/aliyun)"
// @Param directory formData string false "存储目录"
// @Success 200 {object} map[string]interface{} "上传成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Security BearerAuth
// @Router /upload/driver [post]
func (c *UploadController) UploadWithDriver(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择要上传的文件",
		})
		return
	}

	driverType := ctx.PostForm("driver")
	if driverType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请指定存储驱动",
		})
		return
	}

	opts := upload.UploadOptions{
		Directory: ctx.PostForm("directory"),
	}

	info, err := c.uploader.UploadFileWithDriver(ctx.Request.Context(), upload.DriverType(driverType), file, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data":    info,
	})
}

// UploadAvatar 上传头像（使用本地存储）
// @Summary 上传头像
// @Description 上传用户头像，存储到本地
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "头像文件"
// @Success 200 {object} map[string]interface{} "上传成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Security BearerAuth
// @Router /upload/avatar [post]
func (c *UploadController) UploadAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择要上传的头像",
		})
		return
	}

	// 验证是否为图片
	if !upload.IsImageFile(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请上传图片文件",
		})
		return
	}

	opts := upload.UploadOptions{
		Directory: "avatars",
	}

	// 头像强制使用本地存储
	info, err := c.uploader.UploadFileWithDriver(ctx.Request.Context(), upload.DriverMinIO, file, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data":    info,
	})
}

// UploadMedia 上传媒体文件（视频/音频使用MinIO）
// @Summary 上传媒体文件
// @Description 上传视频或音频文件，存储到MinIO
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "媒体文件"
// @Success 200 {object} map[string]interface{} "上传成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Security BearerAuth
// @Router /upload/media [post]
func (c *UploadController) UploadMedia(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择要上传的媒体文件",
		})
		return
	}

	// 验证是否为媒体文件
	if !upload.IsVideoFile(file.Filename) && !upload.IsAudioFile(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请上传视频或音频文件",
		})
		return
	}

	opts := upload.UploadOptions{
		Directory: "media",
	}

	// 媒体文件使用MinIO存储
	info, err := c.uploader.UploadFileWithDriver(ctx.Request.Context(), upload.DriverMinIO, file, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data":    info,
	})
}

// GetDrivers 获取可用的存储驱动列表
// @Summary 获取存储驱动列表
// @Description 获取当前可用的存储驱动列表
// @Tags 文件上传
// @Produce json
// @Success 200 {object} map[string]interface{} "驱动列表"
// @Security BearerAuth
// @Router /upload/drivers [get]
func (c *UploadController) GetDrivers(ctx *gin.Context) {
	factory := upload.GetFactory()
	drivers := factory.ListDrivers()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"drivers": drivers,
		},
	})
}
