# 任务编号自动生成功能说明

## 功能概述

创建任务时，`taskNo`（任务编号）现在支持**自动生成**。如果在创建任务请求中不提供任务编号，系统会自动生成一个全局唯一的编号。

## 生成规则

任务编号由两部分组成：**类型前缀 + 6位随机字符**

### 格式
```
{TaskTypePrefix}-{RandomString6}
```

### 示例
- `REQ-aBc123` （需求类任务）
- `UNIT-Xyz789` （最小单元任务）
- `FEA-KlmNop` （功能任务）

## 前缀映射

| 任务类型编码 | 前缀 | 说明 |
|------------|------|------|
| requirement | REQ | 需求类任务 |
| unit_task | UNIT | 最小单元任务 |
| feature | FEA | 功能任务 |
| bug | BUG | 缺陷任务 |
| task | TASK | 通用任务 |
| 其他 | 自动生成 | 使用编码的前3个字母（大写） |

## 使用方式

### 方式1：让系统自动生成（推荐）

```json
{
  "title": "实现用户登录功能",
  "description": "完成用户登录功能的开发",
  "task_type_code": "requirement",
  "creator_id": 1,
  "executor_id": 2,
  "priority": 3,
  "expected_start_date": "2024-12-15T00:00:00Z",
  "expected_end_date": "2024-12-20T00:00:00Z"
}
```

系统会自动生成编号，如：`REQ-aBc123`

### 方式2：手动指定任务编号

```json
{
  "task_no": "CUSTOM-001",
  "title": "实现用户登录功能",
  "description": "完成用户登录功能的开发",
  "task_type_code": "requirement",
  "creator_id": 1,
  "executor_id": 2,
  "priority": 3
}
```

如果手动指定的编号已存在，系统会返回错误：`任务编号已存在`

## 实现细节

### 核心方法

1. **generateTaskNo(taskTypeCode string)** - 生成任务编号
   - 自动获取任务类型的前缀
   - 生成6位随机字符（大小写字母 + 数字）
   - 检查唯一性（重试最多10次）
   - 如果10次重试都失败，返回错误

2. **getTaskTypePrefix(taskTypeCode string)** - 获取任务类型前缀
   - 首先查询数据库中的 TaskType 表
   - 如果找到，使用其名称的前3个字母
   - 否则使用预定义的映射表
   - 如果都没有，取编码的前3个字母

3. **generateRandomString(length int)** - 生成随机字符串
   - 使用大小写字母（a-z, A-Z）和数字（0-9）
   - 确保生成的字符串长度正确

### 唯一性保证

- 每次生成前都会检查数据库中是否存在相同的编号
- 如果存在，继续生成新的随机串
- 最多重试10次，防止无限循环
- 利用数据库的 uniqueIndex 约束作为最后的保障

## 注意事项

1. **随机性**：由于使用了 `rand.NewSource(time.Now().UnixNano())`，生成的编号具有时间依赖性，请勿在高并发场景中大量生成

2. **并发安全**：如果需要处理高并发的创建任务请求，建议在数据库层面使用唯一索引和异常处理

3. **不支持修改**：任务编号一旦生成就不应该修改（在 UpdateTask 中不允许修改）

4. **客户端选择**：
   - 如果客户端有特殊的编号规则需求，可以在创建时手动指定
   - 否则推荐使用自动生成，减少出错风险

## 相关代码位置

- **TaskService**: `/services/TaskService.go`
  - CreateTask 方法（第17-33行）
  - generateTaskNo 方法（第1160-1190行）
  - getTaskTypePrefix 方法（第1192-1210行）
  - generatePrefixFromCode 方法（第1212-1235行）
  - generateRandomString 方法（第1237-1250行）

- **DTO**: `/dto/TaskDto.go`
  - TaskRequest 结构体（taskNo 字段改为可选）

- **Tests**: `/services/task_no_generator_test.go`
  - 单元测试验证编号生成的正确性

## 后续改进方向

1. 支持自定义前缀规则（通过配置文件）
2. 支持序列号形式的编号（如 REQ-202412-0001）
3. 引入分布式 ID 生成器（如 Snowflake）以支持高并发场景
4. 添加编号格式的验证函数
