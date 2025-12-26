package upload

import (
	"io"
	"sync/atomic"
)

// ProgressReader 带进度追踪的Reader包装器
type ProgressReader struct {
	reader   io.Reader
	total    int64
	uploaded int64
	callback ProgressCallback
}

// NewProgressReader 创建进度Reader
func NewProgressReader(reader io.Reader, total int64, callback ProgressCallback) *ProgressReader {
	return &ProgressReader{
		reader:   reader,
		total:    total,
		callback: callback,
	}
}

// Read 实现io.Reader接口
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n > 0 {
		uploaded := atomic.AddInt64(&pr.uploaded, int64(n))
		if pr.callback != nil {
			pr.callback(uploaded, pr.total)
		}
	}
	return n, err
}

// Progress 获取当前进度
func (pr *ProgressReader) Progress() (uploaded, total int64) {
	return atomic.LoadInt64(&pr.uploaded), pr.total
}

// Percentage 获取进度百分比
func (pr *ProgressReader) Percentage() float64 {
	if pr.total == 0 {
		return 0
	}
	return float64(atomic.LoadInt64(&pr.uploaded)) / float64(pr.total) * 100
}
