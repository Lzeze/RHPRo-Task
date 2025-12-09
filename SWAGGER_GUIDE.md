# Swagger API 文档集成指南

## ✅ 集成完成

已成功为任务管理系统集成 Swagger API 文档！

## 📦 已安装的包

- `github.com/swaggo/swag` - Swagger 生成工具
- `github.com/swaggo/gin-swagger` - Gin Swagger 中间件
- `github.com/swaggo/files` - Swagger 静态文件服务

## 🎯 完成的工作

### 1. 添加 Swagger 配置 (main.go)

在 `main.go` 中添加了 API 通用配置注释：
- API 标题和版本
- API 描述
- 授权配置（Bearer Token）
- 服务器地址和 Base Path

### 2. 为所有 API 添加注释 (TaskController.go)

为 7 个 API 接口添加了详细的 Swagger 注释：
- ✅ POST /tasks - 创建任务
- ✅ GET /tasks - 获取任务列表
- ✅ GET /tasks/:id - 获取任务详情
- ✅ PUT /tasks/:id - 更新任务
- ✅ DELETE /tasks/:id - 删除任务
- ✅ POST /tasks/:id/transit - 状态转换
- ✅ POST /tasks/:id/assign - 分配执行人

每个接口注释包含：
- Summary（摘要）
- Description（详细描述）
- Tags（分组标签）
- Parameters（请求参数）
- Responses（响应格式）
- Security（认证要求）

### 3. 配置 Swagger UI 路由 (routes.go)

添加 Swagger UI 端点：
```go
router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

### 4. 生成 Swagger 文档

运行 `swag init` 生成了：
- `docs/docs.go` - Go 源码形式的文档
- `docs/swagger.json` - JSON 格式的 API 规范
- `docs/swagger.yaml` - YAML 格式的 API 规范

## 🚀 如何使用

### 1. 启动服务

```bash
cd /Users/liuyize/go/RHPRo-Task
./bin/server

# 或直接运行
go run main.go
```

### 2. 访问 Swagger UI

在浏览器中打开：
```
http://localhost:8888/swagger/index.html
```

### 3. 在 Swagger UI 中测试接口

#### 步骤1：登录获取 Token

1. 找到 "认证相关" 分组
2. 点击 `POST /api/v1/login`
3. 点击 "Try it out"
4. 填写登录信息：
```json
{
  "username": "admin",
  "password": "your_password"
}
```
5. 点击 "Execute"
6. 从响应中复制 token

#### 步骤2：设置认证

1. 点击页面右上角的 "Authorize" 按钮
2. 在弹出框中输入：`Bearer <your_token>`
3. 点击 "Authorize"
4. 点击 "Close"

#### 步骤3：测试任务相关接口

现在你可以测试所有任务接口了！

**创建任务示例**：
1. 找到 "任务管理" 分组
2. 点击 `POST /api/v1/tasks`
3. 点击 "Try it out"
4. 填写请求体：
```json
{
  "task_no": "REQ-2024-001",
  "title": "测试任务",
  "task_type_code": "requirement"
}
```
5. 点击 "Execute"

## 📝 API 文档特性

### 自动生成的内容

- ✅ 完整的 DTO 模型定义
- ✅ 请求和响应示例
- ✅ 参数类型和验证规则
- ✅ HTTP 状态码说明
- ✅ 认证要求标识

### 交互式功能

- ✅ 在线测试所有接口
- ✅ 一键复制请求示例
- ✅ 实时查看响应
- ✅ 支持 Bearer Token 认证
- ✅ 下载 API 规范文件

## 🔧 如何更新文档

当你修改或添加新的 API 接口时：

1. **添加 Swagger 注释**（在 Controller 方法上）：
```go
// @Summary 接口摘要
// @Description 详细描述
// @Tags 分组标签
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param paramName paramType dataType required "参数说明"
// @Success 200 {object} ResponseType "成功响应"
// @Router /path [method]
func (ctrl *Controller) MethodName(c *gin.Context) {
    // ...
}
```

2. **重新生成文档**：
```bash
~/go/bin/swag init
```

3. **重新编译项目**：
```bash
go build -o bin/server .
```

## 📊 Swagger 注释语法参考

### 通用注释

- `@title` - API 标题
- `@version` - API 版本
- `@description` - API 描述
- `@host` - 服务器地址
- `@BasePath` - 基础路径
- `@securityDefinitions.apikey` - 认证配置

### 接口注释

- `@Summary` - 接口摘要（简短）
- `@Description` - 接口详细描述
- `@Tags` - 接口分组
- `@Accept` - 接受的 Content-Type
- `@Produce` - 返回的 Content-Type
- `@Security` - 需要的认证方式
- `@Param` - 参数定义
- `@Success` - 成功响应
- `@Failure` - 失败响应
- `@Router` - 路由路径和方法

### 参数类型

- `query` - URL 查询参数
- `path` - URL 路径参数
- `body` - 请求体
- `header` - 请求头
- `formData` - 表单数据

## 🎨 Swagger UI 功能

### 主要功能

1. **API 浏览器** - 查看所有可用的 API 接口
2. **在线测试** - 直接在浏览器中测试接口
3. **模型查看** - 查看所有 DTO 的结构
4. **认证管理** - 统一管理 Bearer Token
5. **响应查看** - 实时查看接口响应

### 高级功能

- 导出 API 规范（JSON/YAML）
- 生成客户端代码
- 接口分组和搜索
- 请求参数自动填充
- 响应格式化和高亮

## 🌐 外部访问

如果需要从外部网络访问 Swagger 文档，修改 `main.go` 中的 `@host`：

```go
// @host your-domain.com
// @BasePath /api/v1
```

然后重新运行 `swag init` 和重新编译。

## 📚 相关资源

- [Swagger 官方文档](https://swagger.io/docs/)
- [swaggo/swag GitHub](https://github.com/swaggo/swag)
- [Gin Swagger GitHub](https://github.com/swaggo/gin-swagger)
- [OpenAPI 规范](https://swagger.io/specification/)
