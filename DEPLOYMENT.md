# RHPro-Task 自动化部署文档

## 项目信息

- **代码仓库**: https://codeup.aliyun.com/rhzy/rhzy-task.git
- **编译输出**: `bin/`
- **日志目录**: `logs/`

## 目录结构

```
rhzy-task/
├── bin/                    # 编译输出目录
│   └── rhppro-task         # Linux 可执行文件
├── logs/                   # 日志输出目录
│   └── app.log             # 应用运行日志
├── .env.development        # 开发环境配置
├── .env.production         # 生产环境配置
├── build.sh                # 编译脚本
└── deploy.sh               # 部署脚本
```

## 环境配置

### 开发环境 (.env.development)

```bash
# 服务器配置
SERVER_PORT=8080
GIN_MODE=debug

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=gin_app
DB_SSLMODE=disable

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT配置
JWT_SECRET=your-dev-secret-key
JWT_EXPIRE_HOURS=24

# 任务配置
EXECUTION_PLAN_DEADLINE_HOURS=72

# 微信配置
WECHAT_OPEN_APPID=
WECHAT_OPEN_SECRET=
WECHAT_MP_APPID=
WECHAT_MP_SECRET=
WECHAT_H5_APPID=
WECHAT_H5_SECRET=
```

### 生产环境 (.env.production)

```bash
# 服务器配置
SERVER_PORT=7777
SERVER_MODE=release

# 数据库配置
DB_HOST=10.0.10.114
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=rhpro_task
DB_SSLMODE=disable

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT配置
JWT_SECRET=your_production_secret_key
JWT_EXPIRE_TIME=24

# 任务配置
EXECUTION_PLAN_DEADLINE_HOURS=48

# 微信配置
WECHAT_OPEN_APPID=your_appid
WECHAT_OPEN_SECRET=your_secret
WECHAT_MP_APPID=your_mp_appid
WECHAT_MP_SECRET=your_mp_secret
WECHAT_H5_APPID=your_h5_appid
WECHAT_H5_SECRET=your_h5_secret
```

## 部署流程

### 1. 克隆代码

```bash
git clone https://codeup.aliyun.com/rhzy/rhzy-task.git
cd rhzy-task
```

### 2. 编译应用

```bash
# 赋予执行权限
chmod +x build.sh

# 执行编译（交叉编译为 Linux amd64）
./build.sh
```

编译完成后，二进制文件位于 `bin/rhppro-task`

### 3. 部署到服务器

```bash
# 赋予执行权限
chmod +x deploy.sh

# 执行部署（默认参数）
./deploy.sh

# 或指定参数
./deploy.sh <用户名> <服务器IP> <远程路径>
# 示例：
./deploy.sh root 10.0.10.113 /home/liuyz/rhppro-task
```

### 4. 部署脚本执行流程

1. 检查本地二进制文件是否存在
2. 验证服务器 SSH 连接
3. 创建远程目录
4. 停止旧版本应用
5. 检查并修正服务器时区
6. 上传二进制文件和配置文件
7. 设置执行权限
8. 启动应用

## 日志说明

### 日志位置

- **本地开发**: 控制台输出
- **生产环境**: `<部署路径>/logs/app.log`

### 查看日志

```bash
# 实时查看日志
tail -f /home/liuyz/rhppro-task/logs/app.log

# 查看最近100行
tail -n 100 /home/liuyz/rhppro-task/logs/app.log

# 搜索错误日志
grep -i error /home/liuyz/rhppro-task/logs/app.log
```

### 日志级别

- `INFO` - 常规信息
- `WARN` - 警告信息
- `ERROR` - 错误信息

## 运维命令

### 查看应用状态

```bash
ps aux | grep rhppro-task
```

### 停止应用

```bash
pkill -f rhppro-task
```

### 重启应用

```bash
pkill -f rhppro-task
cd /home/liuyz/rhppro-task
nohup ./rhppro-task > logs/app.log 2>&1 &
```

### 查看端口占用

```bash
netstat -tlnp | grep 7777
# 或
lsof -i :7777
```

## 数据库部署

### 数据库信息

- **数据库类型**: PostgreSQL 14+
- **数据库名称**: `rhpro_task`
- **字符集**: UTF-8

### 数据库文件说明

```
database/migrations/
├── public.sql              # 完整表结构
├── init_database.sql       # 初始化基础数据
```

### 首次部署（全新安装）

```bash
# 1. 连接到 PostgreSQL
psql -U postgres

# 2. 创建数据库
CREATE DATABASE rhpro_task
    WITH OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE = template0;

# 3. 退出并导入表结构
\q
psql -U postgres -d rhpro_task -f database/migrations/public.sql

# 4. 导入初始化数据
psql -U postgres -d rhpro_task -f database/migrations/init_database.sql
```


### 初始化数据说明

执行 `init_database.sql` 后将创建以下数据：

| 类型 | 内容 |
|------|------|
| 角色 | admin(管理员)、manager(经理)、user(用户) |
| 权限 | 20项系统权限 |
| 管理员 | 手机号: `13800000000`，密码: `admin123` |
| 任务类型 | requirement(需求任务)、unit_task(单元任务) |
| 任务状态 | 22个状态（需求14个 + 单元8个） |
| 状态转换 | 42条转换规则 |
| 任务标签 | 19个预设标签 |

### 生成密码哈希

如需修改管理员密码，可使用工具生成 bcrypt 哈希：

```bash
go run cmd/gen_password/main.go <新密码>
```

### 数据库备份

```bash
# 备份
pg_dump -U postgres -d rhpro_task > backup_$(date +%Y%m%d).sql

# 恢复
psql -U postgres -d rhpro_task < backup_20241217.sql
```

### 远程数据库连接

```bash
# 测试连接
psql -h 10.0.10.114 -U postgres -d rhpro_task

# 执行 SQL 文件
psql -h 10.0.10.114 -U postgres -d rhpro_task -f database/migrations/init_database.sql
```

## 前置要求

### 本地环境

- Go 1.25+
- Git
- SSH 客户端
- PostgreSQL 客户端（psql）

### 服务器环境

- Linux (amd64)
- PostgreSQL 14+
- Redis (可选)
- 开放应用端口 (默认 7777)

### SSH 免密登录配置

```bash
# 生成 SSH 密钥（如果没有）
ssh-keygen -t rsa -b 4096

# 复制公钥到服务器
ssh-copy-id root@10.0.10.113
```

### 免密 sudo 配置（可选）

```bash
# 在服务器上执行
sudo visudo

# 添加以下行（替换 username）
username ALL=(ALL) NOPASSWD: ALL
```

## CI/CD 集成（阿里云 Flow）

### 流水线配置示例

```yaml
# .yunxiao/pipeline.yml
stages:
  - name: 构建
    jobs:
      - name: 编译
        steps:
          - script: |
              go version
              chmod +x build.sh
              ./build.sh
          - upload:
              from: bin/rhppro-task
              to: artifacts/

  - name: 部署
    jobs:
      - name: 部署到生产
        steps:
          - download:
              from: artifacts/rhppro-task
              to: bin/
          - script: |
              chmod +x deploy.sh
              ./deploy.sh root 10.0.10.113 /home/liuyz/rhppro-task
```

## 常见问题

### 1. 编译失败

确保 Go 版本 >= 1.25，并已安装所有依赖：
```bash
go mod tidy
```

### 2. 部署连接失败

检查 SSH 配置和网络连接：
```bash
ssh -v root@10.0.10.113
```

### 3. 应用启动失败

查看日志排查问题：
```bash
cat /home/liuyz/rhppro-task/logs/app.log
```

### 4. 端口被占用

```bash
# 查找占用进程
lsof -i :7777
# 终止进程
kill -9 <PID>
```

### 5. 数据库连接失败

检查 `.env` 配置和数据库服务状态：
```bash
psql -h 10.0.10.114 -U postgres -d rhpro_task
```
