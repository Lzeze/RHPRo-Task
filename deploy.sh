#!/bin/bash

# 部署脚本 - 部署应用到远程服务器
# 使用方式：./deploy.sh <username> <server_ip> <remote_path>
#
# 前置要求：
# 1. 本地已通过 SSH 密钥认证到远程服务器（支持免密登录）
# 2. 远程用户已配置免密 sudo（或在脚本中提供密码）

set -e

# 配置参数
USERNAME=${1:-root}
SERVER_IP=${2:-10.0.10.113}
REMOTE_PATH=${3:-/home/liuyz/rhppro-task}
APP_NAME="rhppro-task"
BINARY_PATH="bin/rhppro-task"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 开始部署应用到 ${SERVER_IP}...${NC}"
echo ""
echo -e "${YELLOW}⚠️  前置检查：${NC}"
echo "  • 确保本地 SSH 密钥已配置，支持免密登录"
echo "  • 确保远程用户已配置免密 sudo，或手动输入密码"
echo ""

# 1. 检查本地二进制文件
if [ ! -f "$BINARY_PATH" ]; then
    echo -e "${RED}❌ 错误：未找到二进制文件 $BINARY_PATH${NC}"
    echo "请先运行 ./build.sh 编译应用"
    exit 1
fi

echo -e "${GREEN}✅ 二进制文件已准备${NC}"

# 2. 检查远程服务器连接
echo -e "${YELLOW}📡 检查服务器连接...${NC}"
if ! ssh -o ConnectTimeout=5 "${USERNAME}@${SERVER_IP}" "echo 'Connected'" > /dev/null 2>&1; then
    echo -e "${RED}❌ 无法连接到服务器 ${SERVER_IP}${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 服务器连接正常${NC}"

# 3. 远程服务器准备（创建目录、停止旧服务）
echo -e "${YELLOW}📁 准备远程环境...${NC}"
ssh "${USERNAME}@${SERVER_IP}" << EOSSH
    set -e
    
    # 创建应用目录
    if [ ! -d "$REMOTE_PATH" ]; then
        mkdir -p "$REMOTE_PATH"
        echo "已创建目录：$REMOTE_PATH"
    fi
    
    # 停止运行的应用（如果存在）
    if pgrep -f "$APP_NAME" > /dev/null 2>&1; then
        echo "停止运行中的 $APP_NAME..."
        pkill -f "$APP_NAME" || true
        sleep 2
    fi
    
    # 检查并修正系统时区（UTC+8 中国标准时间）
    echo "检查系统时区..."
    CURRENT_TZ=\$(cat /etc/timezone 2>/dev/null || echo "未配置")
    if [ "\$CURRENT_TZ" != "Asia/Shanghai" ]; then
        echo "⚠️  当前时区: \$CURRENT_TZ"
        echo "正在修正时区为 Asia/Shanghai (UTC+8)..."
        # 尝试使用 timedatectl（systemd 系统）
        if command -v timedatectl &> /dev/null; then
            timedatectl set-timezone Asia/Shanghai || echo "⚠️  timedatectl 修改失败，请手动配置"
        # 备选方案：直接修改 /etc/timezone
        elif [ -f /etc/timezone ]; then
            echo "Asia/Shanghai" | tee /etc/timezone > /dev/null
            # 更新时区文件软链接
            ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
        fi
        echo "✅ 时区已修正为 Asia/Shanghai"
    else
        echo "✅ 时区已正确配置: \$CURRENT_TZ"
    fi
    
    echo "远程环境准备完毕"
EOSSH

echo -e "${GREEN}✅ 远程环境已准备${NC}"

# 4. 上传文件
echo -e "${YELLOW}📤 上传应用文件...${NC}"
scp -q "$BINARY_PATH" "${USERNAME}@${SERVER_IP}:${REMOTE_PATH}/"
if [ -f ".env.production" ]; then
    scp -q ".env.production" "${USERNAME}@${SERVER_IP}:${REMOTE_PATH}/.env"
else
    echo -e "${YELLOW}⚠️  警告：未找到 .env.production 文件，跳过上传${NC}"
fi

echo -e "${GREEN}✅ 文件上传完成${NC}"

# 5. 设置权限并启动应用
echo -e "${YELLOW}🔧 配置权限并启动应用...${NC}"
ssh "${USERNAME}@${SERVER_IP}" << EOSSH
    set -e
    
    # 显示当前系统时区（验证时区修正）
    CURRENT_TZ=\$(cat /etc/timezone 2>/dev/null || echo "未配置")
    SYSTEM_TIME=\$(date '+%Y-%m-%d %H:%M:%S %Z')
    echo "系统时区: \$CURRENT_TZ"
    echo "系统时间: \$SYSTEM_TIME"
    
    # 设置二进制文件执行权限
    chmod +x "$REMOTE_PATH/$APP_NAME"
    
    # 创建日志目录
    mkdir -p "$REMOTE_PATH/logs"
    
    # 启动应用（使用 nohup）
    cd "$REMOTE_PATH"
    nohup ./$APP_NAME > logs/app.log 2>&1 &
    echo "应用已启动，PID: \$!"
    
    # 等待应用启动
    sleep 2
    
    # 检查应用是否成功启动
    if pgrep -f "$APP_NAME" > /dev/null 2>&1; then
        echo "✅ 应用启动成功"
        echo "日志文件: $REMOTE_PATH/logs/app.log"
    else
        echo "⚠️  应用启动检查失败，请查看日志: $REMOTE_PATH/logs/app.log"
    fi
EOSSH

echo -e "${GREEN}✅ 部署完成！${NC}"
echo -e "${YELLOW}📋 部署摘要：${NC}"
echo "  服务器：${USERNAME}@${SERVER_IP}"
echo "  应用路径：${REMOTE_PATH}"
echo "  应用名称：${APP_NAME}"
echo "  配置文件：${REMOTE_PATH}/.env"
echo "  启动方式：nohup"
echo ""
echo -e "${YELLOW}🔍 查看应用进程：${NC}"
echo "  ssh ${USERNAME}@${SERVER_IP}"
echo "  ps aux | grep ${APP_NAME}"
echo ""
echo -e "${YELLOW}📝 查看应用日志：${NC}"
echo "  ssh ${USERNAME}@${SERVER_IP}"
echo "  tail -f ${REMOTE_PATH}/logs/app.log"
echo ""
echo -e "${YELLOW}🛑 停止应用：${NC}"
echo "  ssh ${USERNAME}@${SERVER_IP}"
echo "  pkill -f ${APP_NAME}"
echo ""
echo -e "${YELLOW}🔄 重启应用：${NC}"
echo "  ssh ${USERNAME}@${SERVER_IP}"
echo "  pkill -f ${APP_NAME}"
echo "  cd ${REMOTE_PATH} && nohup ./${APP_NAME} > logs/app.log 2>&1 &"
echo ""
echo -e "${YELLOW}💡 如果部署时遇到 sudo 密码问题，请在远程服务器配置免密 sudo：${NC}"
echo "  ssh ${USERNAME}@${SERVER_IP}"
echo "  sudo visudo"
echo "  # 在文件末尾添加以下行（将 <username> 替换为实际用户名）："
echo "  <username> ALL=(ALL) NOPASSWD: ALL"
echo ""
