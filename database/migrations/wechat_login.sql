-- 微信登录支持：扩展用户表
-- 执行前请备份数据

-- 1. 添加微信相关字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS wechat_unionid VARCHAR(64);
ALTER TABLE users ADD COLUMN IF NOT EXISTS wechat_openid VARCHAR(64);
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar VARCHAR(500);

-- 2. 为 unionid 和 openid 创建唯一索引（允许NULL）
CREATE UNIQUE INDEX IF NOT EXISTS users_wechat_unionid_key ON users(wechat_unionid) WHERE wechat_unionid IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS users_wechat_openid_key ON users(wechat_openid) WHERE wechat_openid IS NOT NULL;

-- 3. 添加注释
COMMENT ON COLUMN users.wechat_unionid IS '微信全局唯一ID（跨应用）';
COMMENT ON COLUMN users.wechat_openid IS '微信OpenID（单应用内唯一）';
COMMENT ON COLUMN users.avatar IS '用户头像URL（可存微信头像）';
