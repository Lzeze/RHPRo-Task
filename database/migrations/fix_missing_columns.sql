-- 修复 review_sessions 表缺少 deleted_at 字段
ALTER TABLE review_sessions ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX IF NOT EXISTS idx_review_sessions_deleted_at ON review_sessions(deleted_at);

-- 修复 department_leaders 表缺少 deleted_at 和 updated_at 字段
ALTER TABLE department_leaders ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE department_leaders ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_department_leaders_deleted_at ON department_leaders(deleted_at);
