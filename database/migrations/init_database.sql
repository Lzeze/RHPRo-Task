-- ============================================
-- RHPro-Task 数据库初始化脚本
-- 数据库名称: rhpro_task
-- ============================================

-- 1. 创建数据库（需要以 postgres 超级用户执行）
-- 如果数据库已存在，请跳过此步骤
-- CREATE DATABASE rhpro_task
--     WITH OWNER = postgres
--     ENCODING = 'UTF8'
--     LC_COLLATE = 'en_US.UTF-8'
--     LC_CTYPE = 'en_US.UTF-8'
--     TEMPLATE = template0;

-- 2. 连接到 rhpro_task 数据库后执行以下内容
-- \c rhpro_task

-- ============================================
-- 初始化基础数据
-- ============================================

-- 3. 初始化角色数据
INSERT INTO roles (id, name, description, created_at, updated_at) VALUES
(1, 'admin', '系统管理员', NOW(), NOW()),
(2, 'manager', '部门经理', NOW(), NOW()),
(3, 'user', '普通用户', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 重置序列
SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));

-- 4. 初始化权限数据
INSERT INTO permissions (id, name, description, created_at, updated_at) VALUES
-- 用户管理权限
(1, 'user:read', '查看用户', NOW(), NOW()),
(2, 'user:create', '创建用户', NOW(), NOW()),
(3, 'user:update', '更新用户', NOW(), NOW()),
(4, 'user:delete', '删除用户', NOW(), NOW()),
(5, 'user:approve', '审核用户', NOW(), NOW()),
-- 角色管理权限
(6, 'role:read', '查看角色', NOW(), NOW()),
(7, 'role:manage', '管理角色', NOW(), NOW()),
-- 部门管理权限
(8, 'dept:read', '查看部门', NOW(), NOW()),
(9, 'dept:create', '创建部门', NOW(), NOW()),
(10, 'dept:update', '更新部门', NOW(), NOW()),
(11, 'dept:delete', '删除部门', NOW(), NOW()),
(12, 'dept:manage', '管理部门人员', NOW(), NOW()),
-- 任务管理权限
(13, 'task:read', '查看任务', NOW(), NOW()),
(14, 'task:create', '创建任务', NOW(), NOW()),
(15, 'task:update', '更新任务', NOW(), NOW()),
(16, 'task:delete', '删除任务', NOW(), NOW()),
(17, 'task:assign', '分配任务', NOW(), NOW()),
(18, 'task:review', '审核任务', NOW(), NOW()),
-- 系统管理权限
(19, 'permission:manage', '权限管理', NOW(), NOW()),
(20, 'system:config', '系统配置', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 重置序列
SELECT setval('permissions_id_seq', (SELECT MAX(id) FROM permissions));

-- 5. 角色-权限关联（管理员拥有所有权限）
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 部门经理权限
INSERT INTO role_permissions (role_id, permission_id) VALUES
(2, 1),  -- user:read
(2, 8),  -- dept:read
(2, 12), -- dept:manage
(2, 13), -- task:read
(2, 14), -- task:create
(2, 15), -- task:update
(2, 17), -- task:assign
(2, 18)  -- task:review
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 普通用户权限
INSERT INTO role_permissions (role_id, permission_id) VALUES
(3, 1),  -- user:read
(3, 8),  -- dept:read
(3, 13), -- task:read
(3, 14), -- task:create
(3, 15)  -- task:update
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 6. 初始化系统管理员账号
-- 密码: admin123 (bcrypt 加密)
-- 手机号: 13800000000
INSERT INTO users (id, mobile, username, nickname, password, email, status, created_at, updated_at) VALUES
(1, '13800000000', '系统管理员', 'Admin', '$2a$10$3H3S3qGO7CmZSZ0yYlkRPebtUBlbEk9be0J7BCBFo2PfZ1LoHTrDK', 'admin@rhpro.com', 1, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    mobile = EXCLUDED.mobile,
    username = EXCLUDED.username,
    nickname = EXCLUDED.nickname,
    password = EXCLUDED.password,
    email = EXCLUDED.email,
    status = EXCLUDED.status,
    updated_at = NOW();

-- 重置序列
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

-- 7. 管理员角色关联
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1)  -- admin 角色
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 8. 初始化任务类型
INSERT INTO task_types (id, code, name, description, created_at) VALUES
(1, 'requirement', '需求任务', '需要明确需求目标的任务类型', NOW()),
(2, 'unit_task', '最小单元任务', '直接执行的最小单元任务', NOW())
ON CONFLICT (id) DO NOTHING;

-- 重置序列
SELECT setval('task_types_id_seq', (SELECT MAX(id) FROM task_types));

-- 9. 初始化任务状态（按 sort_order 顺序，ID 自增）
INSERT INTO task_statuses (id, code, name, task_type_code, sort_order, description, created_at) VALUES
-- 需求任务状态 (1-14)
(1, 'req_draft', '草稿', 'requirement', 1, '需求任务草稿状态', NOW()),
(2, 'req_pending_assign', '待指派', 'requirement', 2, '发布到待领池，等待执行人领取', NOW()),
(3, 'req_pending_accept', '待接受', 'requirement', 3, '已指派，等待执行人确认接受', NOW()),
(4, 'req_pending_solution', '待提交方案', 'requirement', 4, '执行人已接受，需提交解决方案', NOW()),
(5, 'req_solution_review', '方案审核中', 'requirement', 5, '解决方案审核中', NOW()),
(6, 'req_solution_rejected', '方案被驳回', 'requirement', 6, '解决方案被驳回，需重新提交', NOW()),
(7, 'req_pending_plan', '待提交计划', 'requirement', 7, '方案通过，需提交执行计划', NOW()),
(8, 'req_plan_review', '计划审核中', 'requirement', 8, '执行计划审核中', NOW()),
(9, 'req_plan_rejected', '计划被驳回', 'requirement', 9, '执行计划被驳回', NOW()),
(10, 'req_pending_start', '待开始', 'requirement', 10, '执行计划审核通过，子任务待开始', NOW()),
(11, 'req_in_progress', '执行中', 'requirement', 11, '子任务执行中', NOW()),
(12, 'req_completed', '已完成', 'requirement', 12, '需求任务已完成', NOW()),
(13, 'req_cancelled', '已取消', 'requirement', 13, '需求任务已取消', NOW()),
(14, 'req_blocked', '受阻', 'requirement', 14, '需求任务执行受阻', NOW()),
-- 单元任务状态 (15-22)
(15, 'unit_draft', '草稿', 'unit_task', 1, '单元任务草稿状态', NOW()),
(16, 'unit_pending_assign', '待指派', 'unit_task', 2, '未指派执行人，等待分配', NOW()),
(17, 'unit_pending_accept', '待接受', 'unit_task', 3, '已指派，等待执行人接受', NOW()),
(18, 'unit_pending_start', '待开始', 'unit_task', 4, '执行人已接受任务，待开始执行', NOW()),
(19, 'unit_in_progress', '进行中', 'unit_task', 5, '任务执行中', NOW()),
(20, 'unit_completed', '已完成', 'unit_task', 6, '任务已完成', NOW()),
(21, 'unit_cancelled', '已取消', 'unit_task', 7, '任务已取消', NOW()),
(22, 'unit_blocked', '受阻', 'unit_task', 8, '单元任务执行受阻', NOW())
ON CONFLICT (id) DO NOTHING;

-- 重置序列
SELECT setval('task_statuses_id_seq', (SELECT MAX(id) FROM task_statuses));

-- 10. 初始化状态转换规则（按流程顺序，ID 自增）
INSERT INTO task_status_transitions (id, task_type_code, from_status_code, to_status_code, required_role, requires_approval, is_allowed, description, created_at) VALUES
-- 需求任务状态转换 (1-30)
-- 草稿阶段
(1, 'requirement', 'req_draft', 'req_pending_assign', 'creator', false, true, '发布到待领池', NOW()),
(2, 'requirement', 'req_draft', 'req_pending_accept', 'creator', false, true, '直接指派执行人', NOW()),
(3, 'requirement', 'req_draft', 'req_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 待指派阶段
(4, 'requirement', 'req_pending_assign', 'req_pending_accept', 'executor', false, true, '执行人领取任务', NOW()),
(5, 'requirement', 'req_pending_assign', 'req_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 待接受阶段
(6, 'requirement', 'req_pending_accept', 'req_pending_assign', 'executor', false, true, '执行人拒绝，回到待领池', NOW()),
(7, 'requirement', 'req_pending_accept', 'req_pending_solution', 'executor', false, true, '执行人接受任务，进入待提交方案状态', NOW()),
(8, 'requirement', 'req_pending_accept', 'req_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 待提交方案阶段
(9, 'requirement', 'req_pending_solution', 'req_solution_review', 'executor', true, true, '提交解决方案，进入审核', NOW()),
(10, 'requirement', 'req_pending_solution', 'req_blocked', 'executor', false, true, '标记为受阻', NOW()),
(11, 'requirement', 'req_pending_solution', 'req_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 方案审核阶段
(12, 'requirement', 'req_solution_review', 'req_pending_plan', 'creator', true, true, '方案审核通过', NOW()),
(13, 'requirement', 'req_solution_review', 'req_solution_rejected', 'creator', true, true, '方案审核驳回', NOW()),
(14, 'requirement', 'req_solution_review', 'req_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 方案被驳回阶段
(15, 'requirement', 'req_solution_rejected', 'req_pending_solution', 'executor', false, true, '重新提交方案', NOW()),
(16, 'requirement', 'req_solution_rejected', 'req_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 待提交计划阶段
(17, 'requirement', 'req_pending_plan', 'req_plan_review', 'executor', true, true, '提交执行计划，进入审核', NOW()),
(18, 'requirement', 'req_pending_plan', 'req_blocked', 'executor', false, true, '标记为受阻', NOW()),
-- 计划审核阶段
(19, 'requirement', 'req_plan_review', 'req_in_progress', 'creator', true, true, '计划审核通过，开始执行', NOW()),
(20, 'requirement', 'req_plan_review', 'req_pending_start', NULL, false, true, '执行计划审核通过，拆解子任务后状态变为待开始', NOW()),
(21, 'requirement', 'req_plan_review', 'req_plan_rejected', 'creator', true, true, '计划审核驳回', NOW()),
-- 计划被驳回阶段
(22, 'requirement', 'req_plan_rejected', 'req_pending_plan', 'executor', false, true, '重新提交计划', NOW()),
(23, 'requirement', 'req_plan_rejected', 'req_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 待开始阶段
(24, 'requirement', 'req_pending_start', 'req_in_progress', 'executor', false, true, '子任务开始执行', NOW()),
-- 执行中阶段
(25, 'requirement', 'req_in_progress', 'req_completed', 'executor', false, true, '任务完成', NOW()),
(26, 'requirement', 'req_in_progress', 'req_blocked', 'executor', false, true, '标记为受阻', NOW()),
(27, 'requirement', 'req_in_progress', 'req_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 受阻阶段
(28, 'requirement', 'req_blocked', 'req_in_progress', 'executor', false, true, '解除受阻，继续执行', NOW()),
(29, 'requirement', 'req_blocked', 'req_pending_plan', 'executor', false, true, '回到计划阶段', NOW()),
(30, 'requirement', 'req_blocked', 'req_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 单元任务状态转换 (31-42)
-- 草稿阶段
(31, 'unit_task', 'unit_draft', 'unit_pending_accept', 'creator', false, true, '指派执行人', NOW()),
(32, 'unit_task', 'unit_draft', 'unit_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 待指派阶段
(33, 'unit_task', 'unit_pending_assign', 'unit_pending_accept', 'creator', false, true, '指派执行人', NOW()),
-- 待接受阶段
(34, 'unit_task', 'unit_pending_accept', 'unit_pending_assign', 'executor', false, true, '执行人拒绝，回到待指派', NOW()),
(35, 'unit_task', 'unit_pending_accept', 'unit_pending_start', 'executor', false, true, '执行人接受任务后进入待开始状态', NOW()),
(36, 'unit_task', 'unit_pending_accept', 'unit_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 待开始阶段
(37, 'unit_task', 'unit_pending_start', 'unit_in_progress', 'executor', false, true, '执行人开始执行任务', NOW()),
-- 执行中阶段
(38, 'unit_task', 'unit_in_progress', 'unit_completed', 'executor', false, true, '任务完成', NOW()),
(39, 'unit_task', 'unit_in_progress', 'unit_blocked', 'executor', false, true, '标记为受阻', NOW()),
(40, 'unit_task', 'unit_in_progress', 'unit_cancelled', 'creator', false, true, '取消任务', NOW()),
-- 受阻阶段
(41, 'unit_task', 'unit_blocked', 'unit_in_progress', 'executor', false, true, '解除受阻，继续执行', NOW()),
(42, 'unit_task', 'unit_blocked', 'unit_cancelled', 'creator', false, true, '取消任务', NOW())
ON CONFLICT (id) DO NOTHING;

-- 重置序列
SELECT setval('task_status_transitions_id_seq', (SELECT MAX(id) FROM task_status_transitions));

-- 11. 初始化任务标签
INSERT INTO task_tags (id, name, color, description, created_at) VALUES
(1, 'bug', '#e74c3c', '缺陷/错误', NOW()),
(2, 'feature', '#3498db', '功能需求', NOW()),
(3, 'enhancement', '#2ecc71', '改进/优化', NOW()),
(4, 'documentation', '#9b59b6', '文档', NOW()),
(5, 'urgent', '#e67e22', '紧急', NOW()),
(6, 'low-priority', '#95a5a6', '低优先级', NOW()),
(7, 'research', '#f1c40f', '调研/探索', NOW()),
(8, 'backend', '#34495e', '后端相关', NOW()),
(9, 'frontend', '#1abc9c', '前端相关', NOW()),
(10, 'devops', '#7f8c8d', '运维/部署', NOW()),
(11, 'design', '#d35400', '设计', NOW()),
(12, 'qa', '#8e44ad', '测试', NOW()),
(13, 'security', '#c0392b', '安全', NOW()),
(14, 'performance', '#16a085', '性能', NOW()),
(15, 'refactor', '#27ae60', '重构', NOW()),
(16, 'blocked', '#e74c3c', '阻塞/阻碍', NOW()),
(17, 'on-hold', '#f39c12', '暂停/搁置', NOW()),
(18, 'in-review', '#2980b9', '评审中', NOW()),
(19, 'blocked-by-dependency', '#c0392b', '被依赖阻塞', NOW())
ON CONFLICT (id) DO NOTHING;

-- 重置序列
SELECT setval('task_tags_id_seq', (SELECT MAX(id) FROM task_tags));

-- ============================================
-- 初始化完成
-- ============================================
-- 管理员账号信息：
-- 手机号: 13800000000
-- 密码: admin123
-- ============================================
