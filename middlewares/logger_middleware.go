package middlewares

import (
	"RHPRo-Task/utils"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter 自定义 ResponseWriter 用于捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware 日志中间件 - 记录详细的请求和响应信息
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			// 重新设置请求体，以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 包装 ResponseWriter 以捕获响应体
		blw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 请求信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		// 获取用户ID（如果已登录）
		userID, _ := c.Get("userID")

		// 构建完整路径
		fullPath := path
		if rawQuery != "" {
			fullPath = path + "?" + rawQuery
		}

		// 响应体（限制长度，避免日志过大）
		responseBody := blw.body.String()
		if len(responseBody) > 1000 {
			responseBody = responseBody[:1000] + "...[truncated]"
		}

		// 请求体也限制长度
		if len(requestBody) > 1000 {
			requestBody = requestBody[:1000] + "...[truncated]"
		}

		// 记录详细日志
		logFields := map[string]interface{}{
			"status_code":   statusCode,
			"latency":       latency.String(),
			"latency_ms":    latency.Milliseconds(),
			"client_ip":     clientIP,
			"method":        method,
			"path":          fullPath,
			"user_agent":    userAgent,
			"request_body":  requestBody,
			"response_body": responseBody,
			"response_size": c.Writer.Size(),
		}

		// 如果有用户ID，添加到日志
		if userID != nil {
			logFields["user_id"] = userID
		}

		// 如果有错误，添加错误信息
		if len(c.Errors) > 0 {
			logFields["errors"] = c.Errors.String()
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			utils.Logger.WithFields(logFields).Error("Request failed with server error")
		} else if statusCode >= 400 {
			utils.Logger.WithFields(logFields).Warn("Request failed with client error")
		} else {
			utils.Logger.WithFields(logFields).Info("Request completed")
		}
	}
}
