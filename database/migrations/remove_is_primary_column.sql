-- ============================================
-- 删除 department_leaders 表的 is_primary 列
-- ============================================

-- 步骤1：删除依赖的视图
DROP VIEW IF EXISTS v_user_departments CASCADE;
DROP VIEW IF EXISTS v_department_leaders_view CASCADE;

-- 步骤2：删除 is_primary 列
ALTER TABLE department_leaders DROP COLUMN IF EXISTS is_primary;

-- 步骤3：重新创建视图（不包含 is_primary）
CREATE OR REPLACE VIEW v_user_departments AS
SELECT 
    dl.user_id,
    u.username,
    d.id as department_id,
    d.name as department_name,
    dl.appointed_at
FROM department_leaders dl
JOIN users u ON dl.user_id = u.id
JOIN departments d ON dl.department_id = d.id
WHERE d.deleted_at IS NULL;

COMMENT ON VIEW v_user_departments IS '用户负责的部门视图';

CREATE OR REPLACE VIEW v_department_leaders_view AS
SELECT 
    dl.department_id,
    d.name as department_name,
    dl.user_id,
    u.username,
    u.email,
    dl.appointed_at
FROM department_leaders dl
JOIN users u ON dl.user_id = u.id
JOIN departments d ON dl.department_id = d.id
WHERE d.deleted_at IS NULL
ORDER BY dl.appointed_at;

COMMENT ON VIEW v_department_leaders_view IS '部门负责人列表视图';
