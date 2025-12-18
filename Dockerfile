# 构建阶段
FROM rhzy-harbor.gshbzw.com/rhzy/go:v1.25.5-v1 AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的构建工具
#RUN apk add --no-cache git ca-certificates tzdata
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/rhppro-task main.go

# 运行阶段
FROM alpine:3.23

# 设置时区
 RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/rhppro-task .

# 创建日志目录
RUN mkdir -p /app/logs

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./rhppro-task"]
