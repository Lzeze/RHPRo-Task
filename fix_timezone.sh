#!/bin/bash

# 时区修正脚本 - 修正远程服务器的时区为 UTC+8 (Asia/Shanghai)
# 使用方式：./fix_timezone.sh <username> <server_ip>

set -e

USERNAME=${1:-root}
SERVER_IP=${2:-10.0.10.113}

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}             🕐 远程服务器时区修正脚本${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""

# 检查服务器连接
echo -e "${YELLOW}📡 检查服务器连接...${NC}"
if ! ssh -o ConnectTimeout=5 "${USERNAME}@${SERVER_IP}" "echo 'Connected'" > /dev/null 2>&1; then
    echo -e "${RED}❌ 无法连接到服务器 ${SERVER_IP}${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 服务器连接正常${NC}"
echo ""

# 检查当前时区
echo -e "${YELLOW}🔍 检查当前系统时区...${NC}"
CURRENT_TZ=$(ssh "${USERNAME}@${SERVER_IP}" "cat /etc/timezone 2>/dev/null || echo '未配置'")
CURRENT_TIME=$(ssh "${USERNAME}@${SERVER_IP}" "date '+%Y-%m-%d %H:%M:%S %Z'")
CURRENT_UTC=$(ssh "${USERNAME}@${SERVER_IP}" "date '+%z'")

echo -e "  当前时区文件: ${BLUE}${CURRENT_TZ}${NC}"
echo -e "  当前系统时间: ${BLUE}${CURRENT_TIME}${NC}"
echo -e "  UTC 偏移量:   ${BLUE}${CURRENT_UTC}${NC}"
echo ""

# 检查是否需要修正
if [ "$CURRENT_TZ" = "Asia/Shanghai" ] && [[ "$CURRENT_UTC" == "+0800" ]]; then
    echo -e "${GREEN}✅ 时区已正确配置，无需修正${NC}"
    exit 0
fi

# 提示用户
echo -e "${YELLOW}⚠️  检测到时区不正确，需要修正...${NC}"
echo ""
echo -e "${YELLOW}修正目标：${NC}"
echo -e "  • 时区: Asia/Shanghai (中国标准时间)"
echo -e "  • UTC 偏移: +08:00"
echo ""

# 执行时区修正
echo -e "${YELLOW}🔧 正在修正时区...${NC}"
ssh "${USERNAME}@${SERVER_IP}" << 'EOSSH'
    set -e
    
    echo "检查可用的时区修正工具..."
    
    # 方案1：使用 timedatectl（推荐，systemd 系统）
    if command -v timedatectl &> /dev/null; then
        echo "  找到 timedatectl，使用该工具修正时区..."
        timedatectl set-timezone Asia/Shanghai
        echo "  ✅ 时区已通过 timedatectl 修正"
    # 方案2：直接修改配置文件
    elif [ -f /etc/timezone ]; then
        echo "  使用配置文件方式修正时区..."
        echo "Asia/Shanghai" | tee /etc/timezone > /dev/null
        
        # 更新软链接
        if [ -d /usr/share/zoneinfo ]; then
            ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
            echo "  ✅ 时区已通过配置文件修正"
        else
            echo "  ⚠️  无法找到时区数据库"
        fi
    else
        echo "  ❌ 无法使用任何时区修正工具"
        exit 1
    fi
    
    # 等待系统更新
    sleep 1
    
    # 验证修正结果
    echo ""
    echo "验证修正结果..."
    NEW_TZ=$(cat /etc/timezone 2>/dev/null || echo "未配置")
    NEW_TIME=$(date '+%Y-%m-%d %H:%M:%S %Z')
    NEW_UTC=$(date '+%z')
    
    echo "  新时区文件: $NEW_TZ"
    echo "  新系统时间: $NEW_TIME"
    echo "  新 UTC 偏移: $NEW_UTC"
EOSSH

echo ""

# 验证修正结果
echo -e "${YELLOW}验证修正结果...${NC}"
NEW_TZ=$(ssh "${USERNAME}@${SERVER_IP}" "cat /etc/timezone 2>/dev/null || echo '未配置'")
NEW_TIME=$(ssh "${USERNAME}@${SERVER_IP}" "date '+%Y-%m-%d %H:%M:%S %Z'")
NEW_UTC=$(ssh "${USERNAME}@${SERVER_IP}" "date '+%z'")

echo -e "  时区文件: ${BLUE}${NEW_TZ}${NC}"
echo -e "  系统时间: ${BLUE}${NEW_TIME}${NC}"
echo -e "  UTC 偏移: ${BLUE}${NEW_UTC}${NC}"
echo ""

# 检查修正是否成功
if [ "$NEW_TZ" = "Asia/Shanghai" ] && [[ "$NEW_UTC" == "+0800" ]]; then
    echo -e "${GREEN}✅ 时区修正成功！${NC}"
    echo ""
    echo -e "${YELLOW}📋 修正摘要：${NC}"
    echo -e "  服务器: ${GREEN}${USERNAME}@${SERVER_IP}${NC}"
    echo -e "  新时区: ${GREEN}Asia/Shanghai (UTC+8)${NC}"
    echo -e "  新时间: ${GREEN}${NEW_TIME}${NC}"
else
    echo -e "${RED}❌ 时区修正失败或部分成功${NC}"
    echo -e "${YELLOW}请手动检查和修正：${NC}"
    echo "  ssh ${USERNAME}@${SERVER_IP}"
    echo "  sudo timedatectl set-timezone Asia/Shanghai"
    echo "  # 或者"
    echo "  sudo bash -c 'echo Asia/Shanghai > /etc/timezone && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime'"
    exit 1
fi

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
