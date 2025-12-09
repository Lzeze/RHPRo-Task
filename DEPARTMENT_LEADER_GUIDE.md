# 部门负责人逻辑实现指南

## 📋 概述

已根据需求实现了部门负责人（Department Leader）逻辑，支持一人多部门管理，并在授权信息中包含了负责人状态。

## 🔧 核心变更

### 1. 数据模型

- **DepartmentLeader**: 新增模型，映射数据库中的 `department_leaders` 表。
- **User**: 新增 `ManagedDepartments` 多对多关联，通过 `department_leaders` 表关联。
- **Department**: 新增 `Leaders` 多对多关联，通过 `department_leaders` 表关联。

### 2. JWT 授权信息 (Claims)

Token Claims 中新增了以下字段，用于前端权限控制：

```go
type Claims struct {
    // ... 原有字段
    IsLeader       bool   `json:"is_leader"`        // 是否为部门负责人（只要负责任一部门即为 true）
    ManagedDeptIDs []uint `json:"managed_dept_ids"` // 负责的部门ID列表
}
```

### 3. API 响应

`UserResponse`（用于登录和获取用户信息）新增了 `managed_departments` 字段：

```json
{
  "id": 1,
  "username": "manager",
  "department": {
    "id": 1,
    "name": "研发部"
  },
  "managed_departments": [
    {
      "id": 1,
      "name": "研发部"
    },
    {
      "id": 2,
      "name": "测试部"
    }
  ],
  "roles": [...]
}
```

## 🚀 使用说明

### 前端适配

1. **登录后存储信息**：
   登录接口返回的 `user_info` 中现在包含 `managed_departments` 数组。前端应存储此信息以展示用户负责的部门。

2. **权限判断**：
   解析 JWT Token 获取 `is_leader` 和 `managed_dept_ids`。
   - `is_leader` 可用于显示"部门管理"菜单。
   - `managed_dept_ids` 可用于判断用户是否有权管理特定部门的任务。

### 后端开发

在 Service 层中，可以通过解析 Token 获取当前用户的负责人身份：

```go
claims, _ := c.Get("claims")
userClaims := claims.(*utils.Claims)

if userClaims.IsLeader {
    // 用户是负责人
    for _, deptID := range userClaims.ManagedDeptIDs {
        // 处理负责的部门逻辑
    }
}
```
