# RHPPro-Task 部署指南

## 📋 部署步骤

### 前置条件
- 本地已安装 Go 1.25.4+
- 本地 macOS 环境
- 远程服务器地址：`10.0.10.113`
- 远程服务器可用的 SSH 访问（用户名、密码或密钥）
- 远程服务器已安装 PostgreSQL、Redis（可选）

### 1️⃣ 编译应用

在项目根目录执行：

```bash
# 赋予脚本执行权限
chmod +x build.sh

# 执行编译（编译为 Linux amd64 二进制）
./build.sh
```

**说明**：
- 编译脚本会生成 Linux x86_64 二进制文件（与服务器兼容）
- 二进制文件位置：`bin/rhppro-task`
- 编译后的文件与 Go 版本兼容，可在无 Go 环境的服务器上运行

---

### 2️⃣ 配置生产环境变量

编辑 `.env.production` 文件，根据实际情况配置：

```bash
# 关键配置项：
DB_HOST=<你的PostgreSQL服务器地址>
DB_PORT=5432
DB_USER=<数据库用户>
DB_PASSWORD=<数据库密码>
DB_NAME=rhppro_task
JWT_SECRET=<生成一个安全的JWT密钥>
SERVER_PORT=8888  # 应用监听端口
```

---

### 3️⃣ 部署到远程服务器

```bash
# 赋予脚本执行权限
chmod +x deploy.sh

# 标准部署（默认用户为 root）
./deploy.sh root 10.0.10.113 /opt/rhppro-task

# 或指定自定义用户
./deploy.sh ubuntu 10.0.10.113 /opt/rhppro-task
```

**参数说明**：
- 第一个参数：远程服务器用户名（默认 root）
- 第二个参数：远程服务器 IP（默认 10.0.10.113）
- 第三个参数：远程应用路径（默认 /opt/rhppro-task）

---

## 🔍 部署后验证

### 检查服务状态

```bash
# SSH 连接到服务器
ssh root@10.0.10.113

# 查看应用状态
sudo systemctl status rhppro-task.service

# 查看应用日志
tail -f /opt/rhppro-task/logs/app.log
```

### 测试应用

```bash
# 检查应用是否在监听指定端口
curl http://10.0.10.113:8888/api/v1/swagger/index.html

# 或检查健康状态（如果应用实现了健康检查）
curl http://10.0.10.113:8888/health
```

---

## 🛠️ 应用管理命令

### 启动/停止/重启应用

```bash
# 启动
sudo systemctl start rhppro-task.service

# 停止
sudo systemctl stop rhppro-task.service

# 重启
sudo systemctl restart rhppro-task.service

# 查看状态
sudo systemctl status rhppro-task.service

# 查看完整日志
sudo journalctl -u rhppro-task.service -f
```

### 更新应用

```bash
# 本地重新编译
./build.sh

# 重新部署（会自动停止旧版本，启动新版本）
./deploy.sh root 10.0.10.113 /opt/rhppro-task
```

---

## ⚠️ 常见问题

### 1. 部署失败：Permission denied

**原因**：SSH 脚本权限不足
**解决**：确保部署脚本有执行权限
```bash
chmod +x build.sh deploy.sh
```

### 2. 应用启动后立即停止

**原因**：通常是配置文件错误或数据库连接失败
**解决**：
```bash
# 查看日志
tail -f /opt/rhppro-task/logs/app.log
tail -f /opt/rhppro-task/logs/error.log

# 检查 .env 配置是否正确
cat /opt/rhppro-task/.env
```

### 3. 无法连接数据库

**原因**：PostgreSQL 服务未启动或连接配置错误
**解决**：
```bash
# 检查 PostgreSQL 是否运行
sudo systemctl status postgresql

# 测试数据库连接
psql -h <DB_HOST> -U <DB_USER> -d <DB_NAME>
```

### 4. 端口已被占用

**原因**：其他进程占用了 8888 端口
**解决**：
```bash
# 查看占用端口的进程
lsof -i :8888

# 更改应用监听端口
# 编辑 /opt/rhppro-task/.env，修改 SERVER_PORT
sudo systemctl restart rhppro-task.service
```

---

## 📊 架构说明

```
本地 (macOS)
    ↓
[编译脚本] → Linux x86_64 二进制
    ↓
[部署脚本] → 上传到远程服务器
    ↓
远程服务器 (10.0.10.113)
    ↓
[systemd 服务] → 自动管理应用生命周期
    ↓
[PostgreSQL] + [Redis] + [应用]
```

---

## 📞 获取帮助

如有问题，请检查：
1. 日志文件：`/opt/rhppro-task/logs/`
2. 环境配置：`.env.production`
3. 系统日志：`sudo journalctl -u rhppro-task.service`

---

## 🔐 安全建议

1. **更改 JWT 密钥**：在 `.env.production` 中设置强随机密钥
   ```bash
   openssl rand -base64 32
   ```

2. **配置防火墙**：仅允许必要的端口
   ```bash
   sudo ufw allow 8888/tcp
   ```

3. **定期备份数据库**：制定备份策略
   ```bash
   pg_dump -h localhost -U postgres rhppro_task > backup.sql
   ```

4. **监控日志**：定期检查应用和系统日志，识别异常

5. **SSL/TLS**：考虑使用反向代理（如 Nginx）提供 HTTPS 支持
