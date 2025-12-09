-- ============================================
-- RHPRo-Task应用 - PostgreSQL数据库表结构
-- ============================================

-- 1. 用户表 (users)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    mobile VARCHAR(20),
    status INTEGER DEFAULT 1,  -- 1:正常 0:禁用
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 创建索引
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- 添加注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.id IS '用户ID';
COMMENT ON COLUMN users.username IS '用户名';
COMMENT ON COLUMN users.email IS '邮箱';
COMMENT ON COLUMN users.password IS '密码（加密）';
COMMENT ON COLUMN users.mobile IS '手机号';
COMMENT ON COLUMN users.status IS '状态：1-正常，0-禁用';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';
COMMENT ON COLUMN users.deleted_at IS '软删除时间';

-- ============================================

-- 2. 角色表 (roles)
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 创建索引
CREATE INDEX idx_roles_deleted_at ON roles(deleted_at);
CREATE INDEX idx_roles_name ON roles(name);

-- 添加注释
COMMENT ON TABLE roles IS '角色表';
COMMENT ON COLUMN roles.id IS '角色ID';
COMMENT ON COLUMN roles.name IS '角色名称';
COMMENT ON COLUMN roles.description IS '角色描述';
COMMENT ON COLUMN roles.created_at IS '创建时间';
COMMENT ON COLUMN roles.updated_at IS '更新时间';
COMMENT ON COLUMN roles.deleted_at IS '软删除时间';

-- ============================================

-- 3. 权限表 (permissions)
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 创建索引
CREATE INDEX idx_permissions_deleted_at ON permissions(deleted_at);
CREATE INDEX idx_permissions_name ON permissions(name);

-- 添加注释
COMMENT ON TABLE permissions IS '权限表';
COMMENT ON COLUMN permissions.id IS '权限ID';
COMMENT ON COLUMN permissions.name IS '权限名称（如：user:read）';
COMMENT ON COLUMN permissions.description IS '权限描述';
COMMENT ON COLUMN permissions.created_at IS '创建时间';
COMMENT ON COLUMN permissions.updated_at IS '更新时间';
COMMENT ON COLUMN permissions.deleted_at IS '软删除时间';

-- ============================================

-- 4. 用户角色关联表 (user_roles)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- 添加注释
COMMENT ON TABLE user_roles IS '用户角色关联表（多对多）';
COMMENT ON COLUMN user_roles.user_id IS '用户ID';
COMMENT ON COLUMN user_roles.role_id IS '角色ID';

-- ============================================

-- 5. 角色权限关联表 (role_permissions)
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- 添加注释
COMMENT ON TABLE role_permissions IS '角色权限关联表（多对多）';
COMMENT ON COLUMN role_permissions.role_id IS '角色ID';
COMMENT ON COLUMN role_permissions.permission_id IS '权限ID';

-- ============================================
-- 初始化默认数据
-- ============================================

-- 插入默认权限
INSERT INTO permissions (name, description) VALUES
    ('user:read', '读取用户信息'),
    ('user:create', '创建用户'),
    ('user:update', '更新用户'),
    ('user:delete', '删除用户'),
    ('role:manage', '管理角色'),
    ('permission:manage', '管理权限')
ON CONFLICT (name) DO NOTHING;

-- 插入默认角色
INSERT INTO roles (name, description) VALUES
    ('admin', '管理员'),
    ('user', '普通用户')
ON CONFLICT (name) DO NOTHING;

-- 为管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'admin'),
    id
FROM permissions
ON CONFLICT DO NOTHING;

-- 为普通用户角色分配读取权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'user'),
    id
FROM permissions
WHERE name = 'user:read'
ON CONFLICT DO NOTHING;

-- ============================================
-- 创建更新时间自动触发器
-- ============================================

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有表添加更新时间触发器
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_permissions_updated_at BEFORE UPDATE ON permissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 查询示例
-- ============================================

-- 查看所有表
-- SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';

-- 查看用户及其角色
-- SELECT u.id, u.username, u.email, r.name as role_name
-- FROM users u
-- LEFT JOIN user_roles ur ON u.id = ur.user_id
-- LEFT JOIN roles r ON ur.role_id = r.id;

-- 查看角色及其权限
-- SELECT r.name as role_name, p.name as permission_name, p.description
-- FROM roles r
-- LEFT JOIN role_permissions rp ON r.id = rp.role_id
-- LEFT JOIN permissions p ON rp.permission_id = p.id
-- ORDER BY r.name, p.name;

-- 查看用户的所有权限（通过角色）
-- SELECT DISTINCT u.username, p.name as permission_name
-- FROM users u
-- JOIN user_roles ur ON u.id = ur.user_id
-- JOIN roles r ON ur.role_id = r.id
-- JOIN role_permissions rp ON r.id = rp.role_id
-- JOIN permissions p ON rp.permission_id = p.id
-- WHERE u.username = 'your_username';

-- ============================================
-- 数据库统计
-- ============================================

-- 统计各表记录数
-- SELECT 
--     'users' as table_name, COUNT(*) as count FROM users
-- UNION ALL
-- SELECT 'roles', COUNT(*) FROM roles
-- UNION ALL
-- SELECT 'permissions', COUNT(*) FROM permissions
-- UNION ALL
-- SELECT 'user_roles', COUNT(*) FROM user_roles
-- UNION ALL
-- SELECT 'role_permissions', COUNT(*) FROM role_permissions;