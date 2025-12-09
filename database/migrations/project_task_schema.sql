-- ============================================
-- RHPRo-Task应用 - PostgreSQL数据库表结构
-- ============================================

-- 1. 部门表 (departments)
CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id INTEGER REFERENCES departments(id),
    status INTEGER DEFAULT 1,  -- 1:正常 0:禁用
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_departments_parent_id ON departments(parent_id);

COMMENT ON TABLE departments IS '部门表';
COMMENT ON COLUMN departments.id IS '主键ID';
COMMENT ON COLUMN departments.name IS '部门名称';
COMMENT ON COLUMN departments.description IS '部门描述';
COMMENT ON COLUMN departments.parent_id IS '父部门ID（支持多级部门）';
COMMENT ON COLUMN departments.status IS '状态：1-正常，0-禁用';
COMMENT ON COLUMN departments.created_at IS '创建时间';
COMMENT ON COLUMN departments.updated_at IS '更新时间';
COMMENT ON COLUMN departments.deleted_at IS '软删除时间';

-- ============================================

-- 2. 用户表扩展 (需要在原有users表基础上添加字段)
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_department_leader BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS job_title VARCHAR(100);


COMMENT ON COLUMN users.is_department_leader IS '是否为部门负责人';
COMMENT ON COLUMN users.job_title IS '职位名称';


-- 3. 任务类型枚举表 (task_types)
CREATE TABLE IF NOT EXISTS task_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,  -- requirement, unit_task
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO task_types (code, name, description) VALUES
    ('requirement', '需求任务', '需要明确需求目标的任务类型'),
    ('unit_task', '最小单元任务', '直接执行的最小单元任务')
ON CONFLICT (code) DO NOTHING;

COMMENT ON TABLE task_types IS '任务类型表';
COMMENT ON COLUMN task_types.id IS '主键ID';
COMMENT ON COLUMN task_types.code IS '任务类型编码（requirement-需求任务, unit_task-最小单元任务）';
COMMENT ON COLUMN task_types.name IS '任务类型名称';
COMMENT ON COLUMN task_types.description IS '任务类型描述';
COMMENT ON COLUMN task_types.created_at IS '创建时间';

-- ============================================

-- 4. 任务状态枚举表 (task_statuses)
CREATE TABLE IF NOT EXISTS task_statuses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    task_type_code VARCHAR(50) REFERENCES task_types(code),
    sort_order INTEGER DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO task_statuses (code, name, task_type_code, sort_order, description) VALUES
    -- 需求任务状态
    ('req_draft', '草稿', 'requirement', 1, '需求任务草稿状态'),
    ('req_pending_assign', '待指派', 'requirement', 2, '发布到待领池，等待执行人领取'),
    ('req_pending_accept', '待接受', 'requirement', 3, '已指派，等待执行人确认接受'),
    ('req_pending_goal', '待提交目标', 'requirement', 4, '执行人已接受，需提交具体目标'),
    ('req_goal_review', '目标审核中', 'requirement', 5, '目标和方案审核中'),
    ('req_goal_rejected', '目标被驳回', 'requirement', 6, '目标和方案被驳回'),
    ('req_pending_plan', '待提交计划', 'requirement', 7, '目标通过，需提交执行计划'),
    ('req_plan_review', '计划审核中', 'requirement', 8, '执行计划审核中'),
    ('req_plan_rejected', '计划被驳回', 'requirement', 9, '执行计划被驳回'),
    ('req_in_progress', '执行中', 'requirement', 10, '子任务执行中'),
    ('req_completed', '已完成', 'requirement', 11, '需求任务已完成'),
    ('req_cancelled', '已取消', 'requirement', 12, '需求任务已取消'),
    
    -- 最小单元任务状态
    ('unit_draft', '草稿', 'unit_task', 1, '单元任务草稿状态'),
    ('unit_pending_accept', '待接受', 'unit_task', 2, '已指派，等待执行人接受'),
    ('unit_in_progress', '进行中', 'unit_task', 3, '任务执行中'),
    ('unit_completed', '已完成', 'unit_task', 4, '任务已完成'),
    ('unit_cancelled', '已取消', 'unit_task', 5, '任务已取消')
ON CONFLICT (code) DO NOTHING;

COMMENT ON TABLE task_statuses IS '任务状态表';
COMMENT ON COLUMN task_statuses.id IS '主键ID';
COMMENT ON COLUMN task_statuses.code IS '状态编码（唯一标识）';
COMMENT ON COLUMN task_statuses.name IS '状态名称';
COMMENT ON COLUMN task_statuses.task_type_code IS '所属任务类型编码';
COMMENT ON COLUMN task_statuses.sort_order IS '排序顺序';
COMMENT ON COLUMN task_statuses.description IS '状态描述';
COMMENT ON COLUMN task_statuses.created_at IS '创建时间';

-- ============================================

-- -- 5. 任务主表 (tasks)
-- CREATE TABLE IF NOT EXISTS tasks (
--     id SERIAL PRIMARY KEY,
--     task_no VARCHAR(50) NOT NULL UNIQUE,  -- 任务编号，如：REQ-2024-001
--     title VARCHAR(255) NOT NULL,
--     description TEXT,
--     task_type_code VARCHAR(50) NOT NULL REFERENCES task_types(code),
--     status_code VARCHAR(50) NOT NULL REFERENCES task_statuses(code),
    
--     -- 关联关系
--     creator_id INTEGER NOT NULL REFERENCES users(id),  -- 创建人
--     executor_id INTEGER REFERENCES users(id),  -- 执行人/负责人
--     department_id INTEGER REFERENCES departments(id),  -- 所属部门
--     parent_task_id INTEGER REFERENCES tasks(id),  -- 父任务ID（子任务关联）
    
--     -- 时间相关
--     expected_start_date DATE,  -- 期望开始日期
--     expected_end_date DATE,    -- 期望完成日期
--     actual_start_date DATE,    -- 实际开始日期
--     actual_end_date DATE,      -- 实际完成日期
    
--     -- 优先级和标签
--     priority INTEGER DEFAULT 2,  -- 1:低 2:中 3:高 4:紧急
--     tags TEXT[],  -- 标签数组
    
--     -- 进度
--     progress INTEGER DEFAULT 0,  -- 进度百分比 0-100
    
--     -- 其他
--     is_cross_department BOOLEAN DEFAULT FALSE,  -- 是否跨部门任务
--     is_in_pool BOOLEAN DEFAULT FALSE,  -- 是否在待领池中
    
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP WITH TIME ZONE
-- );

-- -- 创建索引
-- CREATE INDEX idx_tasks_task_no ON tasks(task_no);
-- CREATE INDEX idx_tasks_creator_id ON tasks(creator_id);
-- CREATE INDEX idx_tasks_executor_id ON tasks(executor_id);
-- CREATE INDEX idx_tasks_status_code ON tasks(status_code);
-- CREATE INDEX idx_tasks_parent_task_id ON tasks(parent_task_id);
-- CREATE INDEX idx_tasks_department_id ON tasks(department_id);
-- CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);

-- COMMENT ON TABLE tasks IS '任务主表';
-- COMMENT ON COLUMN tasks.id IS '主键ID';
-- COMMENT ON COLUMN tasks.task_no IS '任务编号（唯一，如：REQ-2024-001）';
-- COMMENT ON COLUMN tasks.title IS '任务标题';
-- COMMENT ON COLUMN tasks.description IS '任务描述';
-- COMMENT ON COLUMN tasks.task_type_code IS '任务类型编码（关联task_types表）';
-- COMMENT ON COLUMN tasks.status_code IS '任务状态编码（关联task_statuses表）';
-- COMMENT ON COLUMN tasks.creator_id IS '创建人用户ID';
-- COMMENT ON COLUMN tasks.executor_id IS '执行人/负责人用户ID';
-- COMMENT ON COLUMN tasks.department_id IS '所属部门ID';
-- COMMENT ON COLUMN tasks.parent_task_id IS '父任务ID（用于子任务关联）';
-- COMMENT ON COLUMN tasks.expected_start_date IS '期望开始日期';
-- COMMENT ON COLUMN tasks.expected_end_date IS '期望完成日期';
-- COMMENT ON COLUMN tasks.actual_start_date IS '实际开始日期';
-- COMMENT ON COLUMN tasks.actual_end_date IS '实际完成日期';
-- COMMENT ON COLUMN tasks.priority IS '优先级：1-低，2-中，3-高，4-紧急';
-- COMMENT ON COLUMN tasks.tags IS '任务标签数组';
-- COMMENT ON COLUMN tasks.progress IS '任务进度百分比（0-100）';
-- COMMENT ON COLUMN tasks.is_cross_department IS '是否跨部门任务';
-- COMMENT ON COLUMN tasks.is_in_pool IS '是否在待领池中（未指派执行人）';
-- COMMENT ON COLUMN tasks.created_at IS '创建时间';
-- COMMENT ON COLUMN tasks.updated_at IS '更新时间';
-- COMMENT ON COLUMN tasks.deleted_at IS '软删除时间';

-- 删除原有的 subtasks 表（存在设计冲突）
DROP TABLE IF EXISTS subtasks CASCADE;

-- 改进后的 tasks 表（增强版）
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task_no VARCHAR(50) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    task_type_code VARCHAR(50) NOT NULL REFERENCES task_types(code),
    status_code VARCHAR(50) NOT NULL REFERENCES task_statuses(code),
    
    -- ===== 关联关系（核心改进） =====
    creator_id INTEGER NOT NULL REFERENCES users(id),
    executor_id INTEGER REFERENCES users(id),
    department_id INTEGER REFERENCES departments(id),
    
    -- 父子任务关系（统一使用这个字段）
    parent_task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    root_task_id INTEGER REFERENCES tasks(id),  -- 🆕 根任务ID（快速定位顶层任务）
    task_level INTEGER DEFAULT 0,  -- 🆕 任务层级（0=顶层，1=一级子任务，2=二级子任务...）
    task_path VARCHAR(500),  -- 🆕 任务路径（如：1/5/12，方便查询整个树）
    child_sequence INTEGER DEFAULT 0,  -- 🆕 在父任务中的序号（用于排序）
    
    -- ===== 子任务统计（冗余字段，提升查询性能） =====
    total_subtasks INTEGER DEFAULT 0,  -- 🆕 直接子任务总数
    completed_subtasks INTEGER DEFAULT 0,  -- 🆕 已完成子任务数
    
    -- ===== 时间相关 =====
    expected_start_date DATE,
    expected_end_date DATE,
    actual_start_date DATE,
    actual_end_date DATE,
    
    -- ===== 优先级和标签 =====
    priority INTEGER DEFAULT 2,
    
    -- ===== 进度 =====
    progress INTEGER DEFAULT 0,
    
    -- ===== 特殊标识 =====
    is_cross_department BOOLEAN DEFAULT FALSE,
    is_in_pool BOOLEAN DEFAULT FALSE,
    is_template BOOLEAN DEFAULT FALSE,  -- 🆕 是否为模板任务
    
    -- ===== 拆分来源（重要！） =====
    -- 先以普通整数列存放拆分来源的执行计划 ID，避免在创建表时出现循环外键问题。
    -- 在所有表创建完成后会使用 ALTER 添加外键约束。
    split_from_plan_id INTEGER,  -- 🆕 从哪个执行计划拆分出来的（稍后添加 FK）
    split_at TIMESTAMP WITH TIME ZONE,  -- 🆕 拆分时间
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- 约束：如果有父任务，必须有根任务
    CONSTRAINT check_root_with_parent CHECK (
        (parent_task_id IS NULL AND root_task_id IS NULL) OR 
        (parent_task_id IS NOT NULL AND root_task_id IS NOT NULL)
    )
);

-- 创建索引（优化查询性能）
CREATE INDEX idx_tasks_task_no ON tasks(task_no);
CREATE INDEX idx_tasks_creator_id ON tasks(creator_id);
CREATE INDEX idx_tasks_executor_id ON tasks(executor_id);
CREATE INDEX idx_tasks_status_code ON tasks(status_code);
CREATE INDEX idx_tasks_parent_task_id ON tasks(parent_task_id);
CREATE INDEX idx_tasks_root_task_id ON tasks(root_task_id);  -- 🆕
CREATE INDEX idx_tasks_task_level ON tasks(task_level);  -- 🆕
CREATE INDEX idx_tasks_task_path ON tasks USING gin(string_to_array(task_path, '/'));  -- 🆕 GIN索引，优化路径查询
CREATE INDEX idx_tasks_department_id ON tasks(department_id);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);

-- 表和列注释
COMMENT ON TABLE tasks IS '任务主表（统一管理所有类型任务及其层级关系）';
COMMENT ON COLUMN tasks.id IS '主键ID';
COMMENT ON COLUMN tasks.task_no IS '任务编号（唯一，如：REQ-2024-001, UNIT-2024-001）';
COMMENT ON COLUMN tasks.title IS '任务标题';
COMMENT ON COLUMN tasks.description IS '任务描述';
COMMENT ON COLUMN tasks.task_type_code IS '任务类型编码（requirement-需求任务, unit_task-单元任务）';
COMMENT ON COLUMN tasks.status_code IS '任务状态编码';
COMMENT ON COLUMN tasks.creator_id IS '创建人用户ID';
COMMENT ON COLUMN tasks.executor_id IS '执行人/负责人用户ID';
COMMENT ON COLUMN tasks.department_id IS '所属部门ID';
COMMENT ON COLUMN tasks.parent_task_id IS '父任务ID（用于建立父子关系）';
COMMENT ON COLUMN tasks.root_task_id IS '根任务ID（顶层任务的ID，方便追溯到最初的需求）';
COMMENT ON COLUMN tasks.task_level IS '任务层级：0-顶层任务，1-一级子任务，2-二级子任务...';
COMMENT ON COLUMN tasks.task_path IS '任务路径（如：1/5/12，表示任务1的子任务5的子任务12）';
COMMENT ON COLUMN tasks.child_sequence IS '在父任务中的序号（用于子任务排序，从1开始）';
COMMENT ON COLUMN tasks.total_subtasks IS '直接子任务总数（冗余字段，提升查询性能）';
COMMENT ON COLUMN tasks.completed_subtasks IS '已完成的直接子任务数（冗余字段）';
COMMENT ON COLUMN tasks.expected_start_date IS '期望开始日期';
COMMENT ON COLUMN tasks.expected_end_date IS '期望完成日期';
COMMENT ON COLUMN tasks.actual_start_date IS '实际开始日期';
COMMENT ON COLUMN tasks.actual_end_date IS '实际完成日期';
COMMENT ON COLUMN tasks.priority IS '优先级：1-低，2-中，3-高，4-紧急';
COMMENT ON COLUMN tasks.progress IS '任务进度百分比（0-100）';
COMMENT ON COLUMN tasks.is_cross_department IS '是否跨部门任务';
COMMENT ON COLUMN tasks.is_in_pool IS '是否在待领池中（未指派执行人）';
COMMENT ON COLUMN tasks.is_template IS '是否为模板任务（用于快速创建相似任务）';
COMMENT ON COLUMN tasks.split_from_plan_id IS '从哪个执行计划拆分出来的（关联execution_plans表）';
COMMENT ON COLUMN tasks.split_at IS '任务拆分时间';
COMMENT ON COLUMN tasks.created_at IS '创建时间';
COMMENT ON COLUMN tasks.updated_at IS '更新时间';
COMMENT ON COLUMN tasks.deleted_at IS '软删除时间';


-- ============================================
-- 触发器：自动维护任务层级和路径
-- ============================================

-- CREATE OR REPLACE FUNCTION update_task_hierarchy()
-- RETURNS TRIGGER AS $$
-- DECLARE
--     parent_level INTEGER;
--     parent_path VARCHAR(500);
--     parent_root_id INTEGER;
--     next_sequence INTEGER;
-- BEGIN
--     -- 如果是顶层任务
--     IF NEW.parent_task_id IS NULL THEN
--         NEW.root_task_id := NULL;
--         NEW.task_level := 0;
--         NEW.task_path := NEW.id::VARCHAR;
--         NEW.child_sequence := 0;
--     ELSE
--         -- 获取父任务信息
--         SELECT task_level, task_path, root_task_id, COALESCE(total_subtasks, 0) + 1
--         INTO parent_level, parent_path, parent_root_id, next_sequence
--         FROM tasks
--         WHERE id = NEW.parent_task_id;
        
--         -- 设置子任务信息
--         NEW.task_level := parent_level + 1;
--         NEW.task_path := parent_path || '/' || NEW.id::VARCHAR;
--         NEW.root_task_id := COALESCE(parent_root_id, NEW.parent_task_id);
--         NEW.child_sequence := next_sequence;
        
--         -- 更新父任务的子任务统计
--         UPDATE tasks 
--         SET total_subtasks = total_subtasks + 1,
--             updated_at = CURRENT_TIMESTAMP
--         WHERE id = NEW.parent_task_id;
--     END IF;
    
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER trigger_update_task_hierarchy
--     BEFORE INSERT ON tasks
--     FOR EACH ROW
--     EXECUTE FUNCTION update_task_hierarchy();

-- COMMENT ON FUNCTION update_task_hierarchy() IS '自动维护任务层级、路径和序号';

-- ============================================
-- 触发器：更新父任务的完成统计
-- ============================================

-- CREATE OR REPLACE FUNCTION update_parent_task_completion()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     -- 如果任务状态变更为已完成
--     IF NEW.status_code IN ('req_completed', 'unit_completed') AND 
--        OLD.status_code NOT IN ('req_completed', 'unit_completed') AND
--        NEW.parent_task_id IS NOT NULL THEN
        
--         UPDATE tasks
--         SET completed_subtasks = completed_subtasks + 1,
--             progress = CASE 
--                 WHEN total_subtasks > 0 THEN 
--                     ROUND((completed_subtasks + 1) * 100.0 / total_subtasks)
--                 ELSE 0 
--             END,
--             updated_at = CURRENT_TIMESTAMP
--         WHERE id = NEW.parent_task_id;
--     END IF;
    
--     -- 如果任务状态从已完成改为其他状态
--     IF OLD.status_code IN ('req_completed', 'unit_completed') AND 
--        NEW.status_code NOT IN ('req_completed', 'unit_completed') AND
--        NEW.parent_task_id IS NOT NULL THEN
        
--         UPDATE tasks
--         SET completed_subtasks = GREATEST(completed_subtasks - 1, 0),
--             progress = CASE 
--                 WHEN total_subtasks > 0 THEN 
--                     ROUND(GREATEST(completed_subtasks - 1, 0) * 100.0 / total_subtasks)
--                 ELSE 0 
--             END,
--             updated_at = CURRENT_TIMESTAMP
--         WHERE id = NEW.parent_task_id;
--     END IF;
    
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER trigger_update_parent_completion
--     AFTER UPDATE OF status_code ON tasks
--     FOR EACH ROW
--     EXECUTE FUNCTION update_parent_task_completion();

-- COMMENT ON FUNCTION update_parent_task_completion() IS '自动更新父任务的完成统计和进度';

-- ============================================
-- 常用查询视图和函数
-- ============================================

-- -- 视图1：任务详情视图（增强版）
-- CREATE OR REPLACE VIEW v_task_details AS
-- SELECT 
--     t.id,
--     t.task_no,
--     t.title,
--     t.description,
--     t.task_type_code,
--     tt.name as task_type_name,
--     t.status_code,
--     ts.name as status_name,
--     t.creator_id,
--     u1.username as creator_name,
--     t.executor_id,
--     u2.username as executor_name,
--     t.department_id,
--     d.name as department_name,
--     t.parent_task_id,
--     pt.task_no as parent_task_no,
--     pt.title as parent_task_title,
--     t.root_task_id,
--     rt.task_no as root_task_no,
--     rt.title as root_task_title,
--     t.task_level,
--     t.task_path,
--     t.child_sequence,
--     t.total_subtasks,
--     t.completed_subtasks,
--     CASE 
--         WHEN t.total_subtasks > 0 THEN 
--             ROUND(t.completed_subtasks * 100.0 / t.total_subtasks, 2)
--         ELSE 0 
--     END as subtask_completion_rate,
--     t.priority,
--     t.progress,
--     t.expected_start_date,
--     t.expected_end_date,
--     t.actual_start_date,
--     t.actual_end_date,
--     t.is_cross_department,
--     t.is_in_pool,
--     t.split_from_plan_id,
--     t.split_at,
--     t.created_at,
--     t.updated_at
-- FROM tasks t
-- LEFT JOIN task_types tt ON t.task_type_code = tt.code
-- LEFT JOIN task_statuses ts ON t.status_code = ts.code
-- LEFT JOIN users u1 ON t.creator_id = u1.id
-- LEFT JOIN users u2 ON t.executor_id = u2.id
-- LEFT JOIN departments d ON t.department_id = d.id
-- LEFT JOIN tasks pt ON t.parent_task_id = pt.id
-- LEFT JOIN tasks rt ON t.root_task_id = rt.id
-- WHERE t.deleted_at IS NULL;

-- COMMENT ON VIEW v_task_details IS '任务详情视图（包含父任务、根任务、子任务统计等信息）';

-- -- 视图2：任务树视图（显示完整层级结构）
-- CREATE OR REPLACE VIEW v_task_tree AS
-- WITH RECURSIVE task_tree AS (
--     -- 顶层任务
--     SELECT 
--         t.id,
--         t.task_no,
--         t.title,
--         t.task_type_code,
--         t.status_code,
--         t.parent_task_id,
--         t.task_level,
--         t.child_sequence,
--         ARRAY[t.id] as path_ids,
--         t.task_no::text as path_display
--     FROM tasks t
--     WHERE t.parent_task_id IS NULL AND t.deleted_at IS NULL
    
--     UNION ALL
    
--     -- 子任务（递归）
--     SELECT 
--         t.id,
--         t.task_no,
--         t.title,
--         t.task_type_code,
--         t.status_code,
--         t.parent_task_id,
--         t.task_level,
--         t.child_sequence,
--         tt.path_ids || t.id,
--         tt.path_display || ' > ' || t.task_no
--     FROM tasks t
--     INNER JOIN task_tree tt ON t.parent_task_id = tt.id
--     WHERE t.deleted_at IS NULL
-- )
-- SELECT * FROM task_tree
-- ORDER BY path_ids;

-- COMMENT ON VIEW v_task_tree IS '任务树形结构视图（递归查询，显示完整层级）';

-- -- ============================================
-- -- 实用函数
-- -- ============================================

-- -- 函数1：获取任务的所有子任务（递归）
-- CREATE OR REPLACE FUNCTION get_all_subtasks(task_id_param INTEGER)
-- RETURNS TABLE (
--     task_id INTEGER,
--     task_no VARCHAR,
--     title VARCHAR,
--     task_level INTEGER,
--     status_code VARCHAR
-- ) AS $$
-- BEGIN
--     RETURN QUERY
--     WITH RECURSIVE subtask_tree AS (
--         SELECT 
--             t.id as task_id,
--             t.task_no,
--             t.title,
--             t.task_level,
--             t.status_code
--         FROM tasks t
--         WHERE t.parent_task_id = task_id_param AND t.deleted_at IS NULL
        
--         UNION ALL
        
--         SELECT 
--             t.id,
--             t.task_no,
--             t.title,
--             t.task_level,
--             t.status_code
--         FROM tasks t
--         INNER JOIN subtask_tree st ON t.parent_task_id = st.task_id
--         WHERE t.deleted_at IS NULL
--     )
--     SELECT * FROM subtask_tree ORDER BY task_level, task_id;
-- END;
-- $$ LANGUAGE plpgsql;

-- COMMENT ON FUNCTION get_all_subtasks(INTEGER) IS '获取指定任务的所有子任务（包括间接子任务）';

-- -- 函数2：获取任务的所有祖先任务
-- CREATE OR REPLACE FUNCTION get_task_ancestors(task_id_param INTEGER)
-- RETURNS TABLE (
--     task_id INTEGER,
--     task_no VARCHAR,
--     title VARCHAR,
--     task_level INTEGER
-- ) AS $$
-- BEGIN
--     RETURN QUERY
--     WITH RECURSIVE ancestor_tree AS (
--         SELECT 
--             t.id as task_id,
--             t.task_no,
--             t.title,
--             t.task_level,
--             t.parent_task_id
--         FROM tasks t
--         WHERE t.id = task_id_param
        
--         UNION ALL
        
--         SELECT 
--             t.id,
--             t.task_no,
--             t.title,
--             t.task_level,
--             t.parent_task_id
--         FROM tasks t
--         INNER JOIN ancestor_tree at ON t.id = at.parent_task_id
--         WHERE t.deleted_at IS NULL
--     )
--     SELECT 
--         ancestor_tree.task_id,
--         ancestor_tree.task_no,
--         ancestor_tree.title,
--         ancestor_tree.task_level
--     FROM ancestor_tree 
--     WHERE ancestor_tree.task_id != task_id_param
--     ORDER BY task_level;
-- END;
-- $$ LANGUAGE plpgsql;

-- COMMENT ON FUNCTION get_task_ancestors(INTEGER) IS '获取指定任务的所有祖先任务（父任务、祖父任务等）';

-- ============================================
-- 查询示例
-- ============================================

-- 示例1：查询顶层任务（没有父任务的）
-- SELECT * FROM v_task_details WHERE parent_task_id IS NULL;

-- 示例2：查询某个任务的直接子任务
-- SELECT * FROM v_task_details WHERE parent_task_id = 1 ORDER BY child_sequence;

-- 示例3：查询某个任务的所有子孙任务（使用函数）
-- SELECT * FROM get_all_subtasks(1);

-- 示例4：查询某个任务的所有祖先任务
-- SELECT * FROM get_task_ancestors(10);

-- 示例5：查询某个根任务下的所有任务
-- SELECT * FROM v_task_details WHERE root_task_id = 1 ORDER BY task_level, child_sequence;

-- 示例6：查询任务树结构
-- SELECT 
--     REPEAT('  ', task_level) || task_no as task_hierarchy,
--     title,
--     status_code
-- FROM v_task_tree
-- WHERE id IN (SELECT task_id FROM get_all_subtasks(1))
-- ORDER BY path_ids;

-- 示例7：统计某个任务的子任务完成情况
-- SELECT 
--     task_no,
--     title,
--     total_subtasks,
--     completed_subtasks,
--     subtask_completion_rate || '%' as completion_rate
-- FROM v_task_details
-- WHERE id = 1;

-- ============================================

-- 6. 需求目标表 (requirement_goals)
CREATE TABLE IF NOT EXISTS requirement_goals (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    goal_no INTEGER NOT NULL,  -- 目标编号（同一任务内）
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    success_criteria TEXT,  -- 成功标准
    priority INTEGER DEFAULT 2,
    status VARCHAR(50) DEFAULT 'pending',  -- pending, approved, rejected
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(task_id, goal_no)
);

CREATE INDEX idx_requirement_goals_task_id ON requirement_goals(task_id);

COMMENT ON TABLE requirement_goals IS '需求目标表（支持多目标）';
COMMENT ON COLUMN requirement_goals.id IS '主键ID';
COMMENT ON COLUMN requirement_goals.task_id IS '关联的任务ID';
COMMENT ON COLUMN requirement_goals.goal_no IS '目标编号（同一任务内的序号）';
COMMENT ON COLUMN requirement_goals.title IS '目标标题';
COMMENT ON COLUMN requirement_goals.description IS '目标描述';
COMMENT ON COLUMN requirement_goals.success_criteria IS '成功标准/验收标准';
COMMENT ON COLUMN requirement_goals.priority IS '目标优先级：1-低，2-中，3-高，4-紧急';
COMMENT ON COLUMN requirement_goals.status IS '目标状态：pending-待审核，approved-已通过，rejected-已驳回';
COMMENT ON COLUMN requirement_goals.sort_order IS '排序顺序';
COMMENT ON COLUMN requirement_goals.created_at IS '创建时间';
COMMENT ON COLUMN requirement_goals.updated_at IS '更新时间';

-- ============================================

-- 7. 思路方案表 (requirement_solutions)
CREATE TABLE IF NOT EXISTS requirement_solutions (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    version INTEGER DEFAULT 1,  -- 方案版本号
    content TEXT,  -- 文字说明
    mindmap_url VARCHAR(500),  -- 脑图文件URL
    file_name VARCHAR(255),
    file_size BIGINT,
    status VARCHAR(50) DEFAULT 'pending',  -- pending, approved, rejected
    submitted_by INTEGER REFERENCES users(id),
    submitted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_requirement_solutions_task_id ON requirement_solutions(task_id);

COMMENT ON TABLE requirement_solutions IS '需求思路方案表';
COMMENT ON COLUMN requirement_solutions.id IS '主键ID';
COMMENT ON COLUMN requirement_solutions.task_id IS '关联的任务ID';
COMMENT ON COLUMN requirement_solutions.version IS '方案版本号（支持多次修改）';
COMMENT ON COLUMN requirement_solutions.content IS '方案文字说明';
COMMENT ON COLUMN requirement_solutions.mindmap_url IS '脑图文件存储URL';
COMMENT ON COLUMN requirement_solutions.file_name IS '脑图文件名';
COMMENT ON COLUMN requirement_solutions.file_size IS '文件大小（字节）';
COMMENT ON COLUMN requirement_solutions.status IS '方案状态：pending-待审核，approved-已通过，rejected-已驳回';
COMMENT ON COLUMN requirement_solutions.submitted_by IS '提交人用户ID';
COMMENT ON COLUMN requirement_solutions.submitted_at IS '提交时间';
COMMENT ON COLUMN requirement_solutions.created_at IS '创建时间';
COMMENT ON COLUMN requirement_solutions.updated_at IS '更新时间';

-- ============================================

-- 8. 执行计划表 (execution_plans)
CREATE TABLE IF NOT EXISTS execution_plans (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    version INTEGER DEFAULT 1,  -- 计划版本号
    tech_stack TEXT NOT NULL,  -- 技术栈选型
    implementation_steps JSONB NOT NULL,  -- 实施步骤（JSON格式）
    resource_requirements TEXT,  -- 资源需求
    risk_assessment TEXT,  -- 风险评估
    status VARCHAR(50) DEFAULT 'pending',  -- pending, approved, rejected
    submitted_by INTEGER REFERENCES users(id),
    submitted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_execution_plans_task_id ON execution_plans(task_id);

COMMENT ON TABLE execution_plans IS '执行计划表';
COMMENT ON COLUMN execution_plans.id IS '主键ID';
COMMENT ON COLUMN execution_plans.task_id IS '关联的任务ID';
COMMENT ON COLUMN execution_plans.version IS '计划版本号（支持多次修改）';
COMMENT ON COLUMN execution_plans.tech_stack IS '技术栈选型说明';
COMMENT ON COLUMN execution_plans.implementation_steps IS '实施步骤JSON：[{step:1, name:"步骤名", description:"描述", duration:3}]';
COMMENT ON COLUMN execution_plans.resource_requirements IS '资源需求说明';
COMMENT ON COLUMN execution_plans.risk_assessment IS '风险评估说明';
COMMENT ON COLUMN execution_plans.status IS '计划状态：pending-待审核，approved-已通过，rejected-已驳回';
COMMENT ON COLUMN execution_plans.submitted_by IS '提交人用户ID';
COMMENT ON COLUMN execution_plans.submitted_at IS '提交时间';
COMMENT ON COLUMN execution_plans.created_at IS '创建时间';
COMMENT ON COLUMN execution_plans.updated_at IS '更新时间';

-- 解决 tasks <-> execution_plans 循环引用：
-- 之前 tasks 中的 split_from_plan_id 暂不声明 REFERENCES，现所有表已创建完毕，补回外键约束。
ALTER TABLE tasks
ADD CONSTRAINT fk_tasks_split_from_plan FOREIGN KEY (split_from_plan_id) REFERENCES execution_plans(id);

-- 为相关列添加索引以优化查询
CREATE INDEX IF NOT EXISTS idx_tasks_split_from_plan_id ON tasks(split_from_plan_id);


-- ============================================

-- -- 9. 子任务表 (subtasks)
-- CREATE TABLE IF NOT EXISTS subtasks (
--     id SERIAL PRIMARY KEY,
--     parent_task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
--     task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,  -- 关联到tasks表
--     subtask_no INTEGER NOT NULL,  -- 子任务编号
--     sort_order INTEGER DEFAULT 0,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     UNIQUE(parent_task_id, subtask_no)
-- );

-- CREATE INDEX idx_subtasks_parent_task_id ON subtasks(parent_task_id);
-- CREATE INDEX idx_subtasks_task_id ON subtasks(task_id);

-- COMMENT ON TABLE subtasks IS '子任务关联表';
-- COMMENT ON COLUMN subtasks.id IS '主键ID';
-- COMMENT ON COLUMN subtasks.parent_task_id IS '父任务ID';
-- COMMENT ON COLUMN subtasks.task_id IS '子任务ID（关联到tasks表）';
-- COMMENT ON COLUMN subtasks.subtask_no IS '子任务编号（父任务内的序号）';
-- COMMENT ON COLUMN subtasks.sort_order IS '排序顺序';
-- COMMENT ON COLUMN subtasks.created_at IS '创建时间';

-- -- ============================================

-- 10. 任务时间节点表 (task_milestones)
CREATE TABLE IF NOT EXISTS task_milestones (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    target_date DATE NOT NULL,
    actual_date DATE,
    status VARCHAR(50) DEFAULT 'pending',  -- pending, completed, delayed
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_task_milestones_task_id ON task_milestones(task_id);

COMMENT ON TABLE task_milestones IS '任务时间节点/里程碑表';
COMMENT ON COLUMN task_milestones.id IS '主键ID';
COMMENT ON COLUMN task_milestones.task_id IS '关联的任务ID';
COMMENT ON COLUMN task_milestones.name IS '节点名称';
COMMENT ON COLUMN task_milestones.description IS '节点描述';
COMMENT ON COLUMN task_milestones.target_date IS '目标完成日期';
COMMENT ON COLUMN task_milestones.actual_date IS '实际完成日期';
COMMENT ON COLUMN task_milestones.status IS '节点状态：pending-待完成，completed-已完成，delayed-延期';
COMMENT ON COLUMN task_milestones.sort_order IS '排序顺序';
COMMENT ON COLUMN task_milestones.created_at IS '创建时间';
COMMENT ON COLUMN task_milestones.updated_at IS '更新时间';

-- ============================================

-- 11. 任务参与人表 (task_participants)
CREATE TABLE IF NOT EXISTS task_participants (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    role VARCHAR(50) NOT NULL,  -- creator, executor, reviewer, jury, observer
    status VARCHAR(50) DEFAULT 'pending',  -- pending, accepted, rejected
    invited_by INTEGER REFERENCES users(id),
    invited_at TIMESTAMP WITH TIME ZONE,
    response_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(task_id, user_id, role)
);

CREATE INDEX idx_task_participants_task_id ON task_participants(task_id);
CREATE INDEX idx_task_participants_user_id ON task_participants(user_id);

COMMENT ON TABLE task_participants IS '任务参与人表';
COMMENT ON COLUMN task_participants.id IS '主键ID';
COMMENT ON COLUMN task_participants.task_id IS '关联的任务ID';
COMMENT ON COLUMN task_participants.user_id IS '参与人用户ID';
COMMENT ON COLUMN task_participants.role IS '参与角色：creator-创建人，executor-执行人，reviewer-审核人，jury-陪审团，observer-观察者';
COMMENT ON COLUMN task_participants.status IS '参与状态：pending-待确认，accepted-已接受，rejected-已拒绝';
COMMENT ON COLUMN task_participants.invited_by IS '邀请人用户ID';
COMMENT ON COLUMN task_participants.invited_at IS '邀请时间';
COMMENT ON COLUMN task_participants.response_at IS '响应时间';
COMMENT ON COLUMN task_participants.created_at IS '创建时间';

-- ============================================

-- -- 12. 审核记录表 (review_records)
-- CREATE TABLE IF NOT EXISTS review_records (
--     id SERIAL PRIMARY KEY,
--     task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
--     review_type VARCHAR(50) NOT NULL,  -- goal_review, solution_review, plan_review
--     target_id INTEGER,  -- 关联的目标/方案/计划ID
--     reviewer_id INTEGER NOT NULL REFERENCES users(id),
--     result VARCHAR(50) NOT NULL,  -- approved, rejected, pending
--     comment TEXT,
--     attachments JSONB,  -- 附件信息
--     review_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- );

-- CREATE INDEX idx_review_records_task_id ON review_records(task_id);
-- CREATE INDEX idx_review_records_reviewer_id ON review_records(reviewer_id);

-- COMMENT ON TABLE review_records IS '审核记录表';
-- COMMENT ON COLUMN review_records.id IS '主键ID';
-- COMMENT ON COLUMN review_records.task_id IS '关联的任务ID';
-- COMMENT ON COLUMN review_records.review_type IS '审核类型：goal_review-目标审核，solution_review-方案审核，plan_review-计划审核';
-- COMMENT ON COLUMN review_records.target_id IS '被审核对象的ID（目标/方案/计划）';
-- COMMENT ON COLUMN review_records.reviewer_id IS '审核人用户ID';
-- COMMENT ON COLUMN review_records.result IS '审核结果：approved-通过，rejected-驳回，pending-审核中';
-- COMMENT ON COLUMN review_records.comment IS '审核意见';
-- COMMENT ON COLUMN review_records.attachments IS '附件信息JSON：[{name:"文件名", url:"地址"}]';
-- COMMENT ON COLUMN review_records.review_at IS '审核时间';
-- COMMENT ON COLUMN review_records.created_at IS '创建时间';

-- -- ============================================

-- 13. 任务变更历史表 (task_change_logs)
CREATE TABLE IF NOT EXISTS task_change_logs (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    change_type VARCHAR(50) NOT NULL,  -- status_change, assign, update, comment
    field_name VARCHAR(100),
    old_value TEXT,
    new_value TEXT,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_task_change_logs_task_id ON task_change_logs(task_id);
CREATE INDEX idx_task_change_logs_created_at ON task_change_logs(created_at);

COMMENT ON TABLE task_change_logs IS '任务变更历史表';
COMMENT ON COLUMN task_change_logs.id IS '主键ID';
COMMENT ON COLUMN task_change_logs.task_id IS '关联的任务ID';
COMMENT ON COLUMN task_change_logs.user_id IS '操作人用户ID';
COMMENT ON COLUMN task_change_logs.change_type IS '变更类型：status_change-状态变更，assign-指派变更，update-信息更新，comment-评论';
COMMENT ON COLUMN task_change_logs.field_name IS '变更字段名称';
COMMENT ON COLUMN task_change_logs.old_value IS '变更前的值';
COMMENT ON COLUMN task_change_logs.new_value IS '变更后的值';
COMMENT ON COLUMN task_change_logs.comment IS '变更说明';
COMMENT ON COLUMN task_change_logs.created_at IS '变更时间';

-- ============================================

-- 14. 任务评论表 (task_comments)
CREATE TABLE IF NOT EXISTS task_comments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    parent_comment_id INTEGER REFERENCES task_comments(id),  -- 支持回复
    attachments JSONB,  -- 附件信息
    is_private BOOLEAN DEFAULT FALSE,  -- 是否私密评论
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_task_comments_task_id ON task_comments(task_id);
CREATE INDEX idx_task_comments_user_id ON task_comments(user_id);
CREATE INDEX idx_task_comments_parent_comment_id ON task_comments(parent_comment_id);

COMMENT ON TABLE task_comments IS '任务评论表';
COMMENT ON COLUMN task_comments.id IS '主键ID';
COMMENT ON COLUMN task_comments.task_id IS '关联的任务ID';
COMMENT ON COLUMN task_comments.user_id IS '评论人用户ID';
COMMENT ON COLUMN task_comments.content IS '评论内容';
COMMENT ON COLUMN task_comments.parent_comment_id IS '父评论ID（用于回复功能）';
COMMENT ON COLUMN task_comments.attachments IS '附件信息JSON：[{name:"文件名", url:"地址", size:123}]';
COMMENT ON COLUMN task_comments.is_private IS '是否为私密评论（仅部分人可见）';
COMMENT ON COLUMN task_comments.created_at IS '创建时间';
COMMENT ON COLUMN task_comments.updated_at IS '更新时间';
COMMENT ON COLUMN task_comments.deleted_at IS '软删除时间';

-- ============================================

-- 15. 任务附件表 (task_attachments)
CREATE TABLE IF NOT EXISTS task_attachments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_type VARCHAR(100),
    file_size BIGINT,
    uploaded_by INTEGER NOT NULL REFERENCES users(id),
    attachment_type VARCHAR(50),  -- requirement, solution, plan, general
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_task_attachments_task_id ON task_attachments(task_id);

COMMENT ON TABLE task_attachments IS '任务附件表';
COMMENT ON COLUMN task_attachments.id IS '主键ID';
COMMENT ON COLUMN task_attachments.task_id IS '关联的任务ID';
COMMENT ON COLUMN task_attachments.file_name IS '文件名';
COMMENT ON COLUMN task_attachments.file_url IS '文件存储URL';
COMMENT ON COLUMN task_attachments.file_type IS '文件类型（MIME类型）';
COMMENT ON COLUMN task_attachments.file_size IS '文件大小（字节）';
COMMENT ON COLUMN task_attachments.uploaded_by IS '上传人用户ID';
COMMENT ON COLUMN task_attachments.attachment_type IS '附件类型：requirement-需求相关，solution-方案相关，plan-计划相关，general-通用附件';
COMMENT ON COLUMN task_attachments.created_at IS '上传时间';

-- ============================================

-- 16. 通知消息表 (notifications)
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,  -- task_assigned, review_request, status_change, comment
    title VARCHAR(255) NOT NULL,
    content TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

COMMENT ON TABLE notifications IS '通知消息表';
COMMENT ON COLUMN notifications.id IS '主键ID';
COMMENT ON COLUMN notifications.user_id IS '接收通知的用户ID';
COMMENT ON COLUMN notifications.task_id IS '关联的任务ID';
COMMENT ON COLUMN notifications.type IS '通知类型：task_assigned-任务指派，review_request-审核请求，status_change-状态变更，comment-评论通知';
COMMENT ON COLUMN notifications.title IS '通知标题';
COMMENT ON COLUMN notifications.content IS '通知内容';
COMMENT ON COLUMN notifications.is_read IS '是否已读';
COMMENT ON COLUMN notifications.read_at IS '阅读时间';
COMMENT ON COLUMN notifications.created_at IS '创建时间';

-- ============================================

-- 17. 任务标签表 (task_tags)
CREATE TABLE IF NOT EXISTS task_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(20),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE task_tags IS '任务标签表';
COMMENT ON COLUMN task_tags.id IS '主键ID';
COMMENT ON COLUMN task_tags.name IS '标签名称（唯一）';
COMMENT ON COLUMN task_tags.color IS '标签颜色（用于前端显示）';
COMMENT ON COLUMN task_tags.description IS '标签描述';
COMMENT ON COLUMN task_tags.created_at IS '创建时间';

-- 初始化默认标签（可在此处扩展或由迁移脚本管理）
INSERT INTO task_tags (name, color, description) VALUES
    ('bug', '#e74c3c', '缺陷/错误'),
    ('feature', '#3498db', '功能需求'),
    ('enhancement', '#2ecc71', '改进/优化'),
    ('documentation', '#9b59b6', '文档'),
    ('urgent', '#e67e22', '紧急'),
    ('low-priority', '#95a5a6', '低优先级'),
    ('research', '#f1c40f', '调研/探索'),
    ('backend', '#34495e', '后端相关'),
    ('frontend', '#1abc9c', '前端相关'),
    ('devops', '#7f8c8d', '运维/部署'),
    ('design', '#d35400', '设计'),
    ('qa', '#8e44ad', '测试'),
    ('security', '#c0392b', '安全'),
    ('performance', '#16a085', '性能'),
    ('refactor', '#27ae60', '重构')
    -- 状态类标签（阻碍/评审/暂停等）
    ,('blocked', '#e74c3c', '阻塞/阻碍')
    ,('on-hold', '#f39c12', '暂停/搁置')
    ,('in-review', '#2980b9', '评审中')
    ,('blocked-by-dependency', '#c0392b', '被依赖阻塞')
ON CONFLICT (name) DO NOTHING;

-- 1) 创建关系表：task_tag_rel（task - tag 多对多）
CREATE TABLE IF NOT EXISTS task_tag_rel (
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    tag_id  INTEGER NOT NULL REFERENCES task_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, tag_id)
);
CREATE INDEX IF NOT EXISTS idx_task_tag_rel_task_id ON task_tag_rel(task_id);
CREATE INDEX IF NOT EXISTS idx_task_tag_rel_tag_id ON task_tag_rel(tag_id);


COMMENT ON TABLE task_tag_rel IS '任务与标签关系表';
COMMENT ON COLUMN task_tag_rel.task_id IS '关联 tasks.id';
COMMENT ON COLUMN task_tag_rel.tag_id IS '关联 task_tags.id';


-- ============================================
-- 创建触发器：自动更新updated_at
-- ============================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为相关表添加触发器
CREATE TRIGGER update_departments_updated_at BEFORE UPDATE ON departments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_requirement_goals_updated_at BEFORE UPDATE ON requirement_goals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_requirement_solutions_updated_at BEFORE UPDATE ON requirement_solutions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_execution_plans_updated_at BEFORE UPDATE ON execution_plans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_task_milestones_updated_at BEFORE UPDATE ON task_milestones
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_task_comments_updated_at BEFORE UPDATE ON task_comments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 初始化基础数据
-- ============================================

-- 插入示例部门
INSERT INTO departments (name, description) VALUES
    ('产品部', '产品设计与规划部门'),
    ('技术部', '技术研发部门'),
    ('运营部', '运营推广部门'),
    ('设计部', 'UI/UX设计部门')
ON CONFLICT DO NOTHING;

-- ============================================
-- 常用查询视图
-- ============================================

-- -- 任务详情视图
-- CREATE OR REPLACE VIEW v_task_details AS
-- SELECT 
--     t.id,
--     t.task_no,
--     t.title,
--     t.description,
--     t.task_type_code,
--     tt.name as task_type_name,
--     t.status_code,
--     ts.name as status_name,
--     t.creator_id,
--     u1.username as creator_name,
--     t.executor_id,
--     u2.username as executor_name,
--     t.department_id,
--     d.name as department_name,
--     t.priority,
--     t.progress,
--     t.expected_start_date,
--     t.expected_end_date,
--     t.actual_start_date,
--     t.actual_end_date,
--     t.is_cross_department,
--     t.is_in_pool,
--     t.parent_task_id,
--     t.created_at,
--     t.updated_at
-- FROM tasks t
-- LEFT JOIN task_types tt ON t.task_type_code = tt.code
-- LEFT JOIN task_statuses ts ON t.status_code = ts.code
-- LEFT JOIN users u1 ON t.creator_id = u1.id
-- LEFT JOIN users u2 ON t.executor_id = u2.id
-- LEFT JOIN departments d ON t.department_id = d.id
-- WHERE t.deleted_at IS NULL;

-- COMMENT ON VIEW v_task_details IS '任务详情视图（包含关联表信息）';

-- ============================================
-- 查询示例
-- ============================================

-- 1. 查询待领池中的任务
-- SELECT * FROM v_task_details WHERE is_in_pool = TRUE AND executor_id IS NULL;

-- 2. 查询某用户的所有任务（作为执行人）
-- SELECT * FROM v_task_details WHERE executor_id = 1;

-- 3. 查询需求任务及其目标
-- SELECT t.*, rg.title as goal_title, rg.description as goal_description
-- FROM tasks t
-- LEFT JOIN requirement_goals rg ON t.id = rg.task_id
-- WHERE t.task_type_code = 'requirement' AND t.id = 1;

-- 4. 查询任务的审核历史
-- SELECT rr.*, u.username as reviewer_name
-- FROM review_records rr
-- LEFT JOIN users u ON rr.reviewer_id = u.id
-- WHERE rr.task_id = 1
-- ORDER BY rr.review_at DESC;

-- 5. 查询跨部门任务
-- SELECT * FROM v_task_details WHERE is_cross_department = TRUE;

-- 6. 查询任务的子任务
-- SELECT t.* FROM tasks t
-- WHERE t.parent_task_id = 1 AND t.deleted_at IS NULL
-- ORDER BY t.created_at;