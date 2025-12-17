-- 用户注册登录调整：手机号必填用于登录，邮箱选填
-- 执行前请备份数据

-- 1. 移除 username 唯一约束（允许同名用户）
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;

-- 2. 先查看有问题的数据（空手机号或重复手机号）
-- SELECT id, username, mobile FROM users WHERE mobile IS NULL OR mobile = '';
-- SELECT mobile, COUNT(*) FROM users GROUP BY mobile HAVING COUNT(*) > 1;

-- 3. 为空手机号的用户生成临时唯一手机号（用用户ID补充）
UPDATE users SET mobile = CONCAT('temp_', id) WHERE mobile IS NULL OR mobile = '';

-- 4. 添加 mobile 唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS users_mobile_key ON users(mobile);

-- 5. 移除 email 唯一约束，修改为可空，再重建唯一索引（允许多个NULL）
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users(email) WHERE email IS NOT NULL AND email != '';

-- 6. 修改 mobile 字段为非空
ALTER TABLE users ALTER COLUMN mobile SET NOT NULL;
