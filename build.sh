#!/bin/bash

# 编译脚本 - 编译 Go 应用为 Linux 二进制
set -e

echo "🔨 开始编译应用..."

# 编译应用（仅在子进程中设置环境变量）
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/rhppro-task main.go

echo "✅ 编译完成！二进制文件位置：bin/rhppro-task"
echo "📦 二进制文件大小："
ls -lh bin/rhppro-task
