#!/bin/bash

# 编译脚本 - 编译 Go 应用为 Linux 二进制
set -e

# 默认环境
ENV=${1:-development}

# 验证环境参数
if [[ "$ENV" != "development" && "$ENV" != "production" ]]; then
    echo "❌ 无效的环境参数: $ENV"
    echo "用法: ./build.sh [development|production]"
    exit 1
fi

echo "🔨 开始编译应用 (环境: $ENV)..."

# 根据环境参数复制对应的配置文件
if [[ "$ENV" == "production" ]]; then
    echo "📋 使用生产环境配置..."
    cp .env.production .env
else
    echo "📋 使用测试环境配置..."
    cp .env.development .env
fi

# 编译应用（仅在子进程中设置环境变量）
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/rhppro-task main.go

echo "✅ 编译完成！二进制文件位置：bin/rhppro-task"
echo "📦 二进制文件大小："
ls -lh bin/rhppro-task
echo "🌍 当前环境: $ENV"
