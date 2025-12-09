package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var LogWriter io.Writer

func InitLogger() {
	Logger = logrus.New()

	// 设置日志格式
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志级别为 Debug，记录所有日志
	Logger.SetLevel(logrus.DebugLevel)

	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Fatal("Failed to create log directory:", err)
	}

	// 配置按天切分的日志文件（保留最近 30 天）
	writer, err := NewDailyFileWriter(logDir, "app", 30)
	if err != nil {
		logrus.Fatal("Failed to initialize daily log writer:", err)
	}

	// 多输出：同时输出到文件和控制台
	multiWriter := io.MultiWriter(writer, os.Stdout)
	Logger.SetOutput(multiWriter)

	// 保存 writer 供 Gin 使用
	LogWriter = multiWriter

	// 设置 Gin 的日志输出也使用这个 writer
	gin.DefaultWriter = multiWriter
	gin.DefaultErrorWriter = multiWriter
}

// DailyFileWriter 每日切分的文件写入器
type DailyFileWriter struct {
	dir        string
	prefix     string
	maxAgeDays int

	mu          sync.Mutex
	currentDate string
	file        *os.File
}

// NewDailyFileWriter 创建每日切分写入器，prefix 为文件名前缀（例如 "app" -> app.2025-12-02.log）
func NewDailyFileWriter(dir, prefix string, maxAgeDays int) (*DailyFileWriter, error) {
	d := &DailyFileWriter{
		dir:        dir,
		prefix:     prefix,
		maxAgeDays: maxAgeDays,
	}
	if err := d.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DailyFileWriter) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.rotateIfNeeded(); err != nil {
		return 0, err
	}
	return d.file.Write(p)
}

func (d *DailyFileWriter) rotateIfNeeded() error {
	today := time.Now().Format("2006-01-02")
	if d.file != nil && d.currentDate == today {
		return nil
	}
	// close existing
	if d.file != nil {
		_ = d.file.Close()
		d.file = nil
	}
	filename := filepath.Join(d.dir, d.prefix+"."+today+".log")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	d.file = f
	d.currentDate = today
	// 清理旧日志
	go d.cleanupOldFiles()
	return nil
}

func (d *DailyFileWriter) cleanupOldFiles() {
	if d.maxAgeDays <= 0 {
		return
	}
	cutoff := time.Now().AddDate(0, 0, -d.maxAgeDays)
	entries, err := os.ReadDir(d.dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, d.prefix+".") || !strings.HasSuffix(name, ".log") {
			continue
		}
		// 格式 prefix.YYYY-MM-DD.log
		parts := strings.Split(name, ".")
		if len(parts) < 3 {
			continue
		}
		dateStr := parts[len(parts)-2]
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if t.Before(cutoff) {
			_ = os.Remove(filepath.Join(d.dir, name))
		}
	}
}
