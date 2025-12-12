#!/bin/bash

# 快速配置脚本 - 为远程服务器配置免密 sudo（用于部署）
# 使用方式：./setup_sudo.sh <username> <server_ip>

set -e

}
SERVER_IP=${2:-10.0.10.113}

echo "🔧 为用户 $USERNAME 配置免密 sudo..."
echo ""

# 连接到远程服务器并配置 sudo
ssh "${USERNAME}@${SERVER_IP}" << 'EOSSH'
    echo "当前用户：$(whoami)"
    echo "用户 ID：$(id)"
    echo ""
    
    # 检查是否已是 root
    if [ "$EUID" -ne 0 ]; then
        echo "⚠️  需要 root 权限来配置 sudoers"
        echo "请使用 sudo 运行此脚本，或手动执行以下命令："
        echo ""
        echo "sudo tee -a /etc/sudoers.d/$(whoami) > /dev/null << 'EOF'"
        echo "$(whoami) ALL=(ALL) NOPASSWD: ALL"
        echo "EOF"
        echo ""
        echo "sudo chmod 0440 /etc/sudoers.d/$(whoami)"
        exit 1
    fi
    
    CURRENT_USER=$(logname)
    echo "配置 sudoers 文件..."
    
    # 添加免密 sudo 配置（使用 sudoers.d 目录，更安全）
    if ! grep -q "^${CURRENT_USER} ALL=(ALL) NOPASSWD: ALL" /etc/sudoers.d/${CURRENT_USER} 2>/dev/null; then
        echo "${CURRENT_USER} ALL=(ALL) NOPASSWD: ALL" | tee /etc/sudoers.d/${CURRENT_USER} > /dev/null
        chmod 0440 /etc/sudoers.d/${CURRENT_USER}
        echo "✅ 已添加免密 sudo 配置"
    else
        echo "✅ 免密 sudo 配置已存在"
    fi
EOSSH

echo ""
echo "✅ 配置完成！"
echo ""
echo "验证配置："
echo "  ssh ${USERNAME}@${SERVER_IP}"
echo "  sudo -l"
echo ""
