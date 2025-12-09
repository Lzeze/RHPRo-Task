-- ============================================
-- 改进后的任务表结构设计
-- ============================================

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
    tags TEXT[],
    
    -- ===== 进度 =====
    progress INTEGER DEFAULT 0,
    
    -- ===== 特殊标识 =====
    is_cross_department BOOLEAN DEFAULT FALSE,
    is_in_pool BOOLEAN DEFAULT FALSE,
    is_template BOOLEAN DEFAULT FALSE,  -- 🆕 是否为模板任务
    
    -- ===== 拆分来源（重要！） =====
    split_from_plan_id INTEGER REFERENCES execution_plans(id),  -- 🆕 从哪个执行计划拆分出来的
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
COMMENT ON COLUMN tasks.tags IS '任务标签数组';
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

CREATE OR REPLACE FUNCTION update_task_hierarchy()
RETURNS TRIGGER AS $$
DECLARE
    parent_level INTEGER;
    parent_path VARCHAR(500);
    parent_root_id INTEGER;
    next_sequence INTEGER;
BEGIN
    -- 如果是顶层任务
    IF NEW.parent_task_id IS NULL THEN
        NEW.root_task_id := NULL;
        NEW.task_level := 0;
        NEW.task_path := NEW.id::VARCHAR;
        NEW.child_sequence := 0;
    ELSE
        -- 获取父任务信息
        SELECT task_level, task_path, root_task_id, COALESCE(total_subtasks, 0) + 1
        INTO parent_level, parent_path, parent_root_id, next_sequence
        FROM tasks
        WHERE id = NEW.parent_task_id;
        
        -- 设置子任务信息
        NEW.task_level := parent_level + 1;
        NEW.task_path := parent_path || '/' || NEW.id::VARCHAR;
        NEW.root_task_id := COALESCE(parent_root_id, NEW.parent_task_id);
        NEW.child_sequence := next_sequence;
        
        -- 更新父任务的子任务统计
        UPDATE tasks 
        SET total_subtasks = total_subtasks + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.parent_task_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_task_hierarchy
    BEFORE INSERT ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_task_hierarchy();

COMMENT ON FUNCTION update_task_hierarchy() IS '自动维护任务层级、路径和序号';

-- ============================================
-- 触发器：更新父任务的完成统计
-- ============================================

CREATE OR REPLACE FUNCTION update_parent_task_completion()
RETURNS TRIGGER AS $$
BEGIN
    -- 如果任务状态变更为已完成
    IF NEW.status_code IN ('req_completed', 'unit_completed') AND 
       OLD.status_code NOT IN ('req_completed', 'unit_completed') AND
       NEW.parent_task_id IS NOT NULL THEN
        
        UPDATE tasks
        SET completed_subtasks = completed_subtasks + 1,
            progress = CASE 
                WHEN total_subtasks > 0 THEN 
                    ROUND((completed_subtasks + 1) * 100.0 / total_subtasks)
                ELSE 0 
            END,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.parent_task_id;
    END IF;
    
    -- 如果任务状态从已完成改为其他状态
    IF OLD.status_code IN ('req_completed', 'unit_completed') AND 
       NEW.status_code NOT IN ('req_completed', 'unit_completed') AND
       NEW.parent_task_id IS NOT NULL THEN
        
        UPDATE tasks
        SET completed_subtasks = GREATEST(completed_subtasks - 1, 0),
            progress = CASE 
                WHEN total_subtasks > 0 THEN 
                    ROUND(GREATEST(completed_subtasks - 1, 0) * 100.0 / total_subtasks)
                ELSE 0 
            END,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.parent_task_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_parent_completion
    AFTER UPDATE OF status_code ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_parent_task_completion();

COMMENT ON FUNCTION update_parent_task_completion() IS '自动更新父任务的完成统计和进度';

-- ============================================
-- 常用查询视图和函数
-- ============================================

-- 视图1：任务详情视图（增强版）
CREATE OR REPLACE VIEW v_task_details AS
SELECT 
    t.id,
    t.task_no,
    t.title,
    t.description,
    t.task_type_code,
    tt.name as task_type_name,
    t.status_code,
    ts.name as status_name,
    t.creator_id,
    u1.username as creator_name,
    t.executor_id,
    u2.username as executor_name,
    t.department_id,
    d.name as department_name,
    t.parent_task_id,
    pt.task_no as parent_task_no,
    pt.title as parent_task_title,
    t.root_task_id,
    rt.task_no as root_task_no,
    rt.title as root_task_title,
    t.task_level,
    t.task_path,
    t.child_sequence,
    t.total_subtasks,
    t.completed_subtasks,
    CASE 
        WHEN t.total_subtasks > 0 THEN 
            ROUND(t.completed_subtasks * 100.0 / t.total_subtasks, 2)
        ELSE 0 
    END as subtask_completion_rate,
    t.priority,
    t.progress,
    t.expected_start_date,
    t.expected_end_date,
    t.actual_start_date,
    t.actual_end_date,
    t.is_cross_department,
    t.is_in_pool,
    t.split_from_plan_id,
    t.split_at,
    t.created_at,
    t.updated_at
FROM tasks t
LEFT JOIN task_types tt ON t.task_type_code = tt.code
LEFT JOIN task_statuses ts ON t.status_code = ts.code
LEFT JOIN users u1 ON t.creator_id = u1.id
LEFT JOIN users u2 ON t.executor_id = u2.id
LEFT JOIN departments d ON t.department_id = d.id
LEFT JOIN tasks pt ON t.parent_task_id = pt.id
LEFT JOIN tasks rt ON t.root_task_id = rt.id
WHERE t.deleted_at IS NULL;

COMMENT ON VIEW v_task_details IS '任务详情视图（包含父任务、根任务、子任务统计等信息）';

-- 视图2：任务树视图（显示完整层级结构）
CREATE OR REPLACE VIEW v_task_tree AS
WITH RECURSIVE task_tree AS (
    -- 顶层任务
    SELECT 
        t.id,
        t.task_no,
        t.title,
        t.task_type_code,
        t.status_code,
        t.parent_task_id,
        t.task_level,
        t.child_sequence,
        ARRAY[t.id] as path_ids,
        t.task_no::TEXT as path_display  -- 显式转换为TEXT类型
    FROM tasks t
    WHERE t.parent_task_id IS NULL AND t.deleted_at IS NULL
    
    UNION ALL
    
    -- 子任务（递归）
    SELECT 
        t.id,
        t.task_no,
        t.title,
        t.task_type_code,
        t.status_code,
        t.parent_task_id,
        t.task_level,
        t.child_sequence,
        tt.path_ids || t.id,
        tt.path_display || ' > ' || t.task_no::TEXT  -- 显式转换为TEXT类型
    FROM tasks t
    INNER JOIN task_tree tt ON t.parent_task_id = tt.id
    WHERE t.deleted_at IS NULL
)
SELECT * FROM task_tree
ORDER BY path_ids;

COMMENT ON VIEW v_task_tree IS '任务树形结构视图（递归查询，显示完整层级）';

-- ============================================
-- 实用函数
-- ============================================

-- 函数1：获取任务的所有子任务（递归）
CREATE OR REPLACE FUNCTION get_all_subtasks(task_id_param INTEGER)
RETURNS TABLE (
    task_id INTEGER,
    task_no VARCHAR(50),
    title VARCHAR(255),
    task_level INTEGER,
    status_code VARCHAR(50)
) AS $
BEGIN
    RETURN QUERY
    WITH RECURSIVE subtask_tree AS (
        SELECT 
            t.id as task_id,
            t.task_no,
            t.title,
            t.task_level,
            t.status_code
        FROM tasks t
        WHERE t.parent_task_id = task_id_param AND t.deleted_at IS NULL
        
        UNION ALL
        
        SELECT 
            t.id,
            t.task_no,
            t.title,
            t.task_level,
            t.status_code
        FROM tasks t
        INNER JOIN subtask_tree st ON t.parent_task_id = st.task_id
        WHERE t.deleted_at IS NULL
    )
    SELECT * FROM subtask_tree ORDER BY task_level, task_id;
END;
$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_all_subtasks(INTEGER) IS '获取指定任务的所有子任务（包括间接子任务）';

-- 函数2：获取任务的所有祖先任务
CREATE OR REPLACE FUNCTION get_task_ancestors(task_id_param INTEGER)
RETURNS TABLE (
    task_id INTEGER,
    task_no VARCHAR(50),
    title VARCHAR(255),
    task_level INTEGER
) AS $
BEGIN
    RETURN QUERY
    WITH RECURSIVE ancestor_tree AS (
        SELECT 
            t.id as task_id,
            t.task_no,
            t.title,
            t.task_level,
            t.parent_task_id
        FROM tasks t
        WHERE t.id = task_id_param
        
        UNION ALL
        
        SELECT 
            t.id,
            t.task_no,
            t.title,
            t.task_level,
            t.parent_task_id
        FROM tasks t
        INNER JOIN ancestor_tree at ON t.id = at.parent_task_id
        WHERE t.deleted_at IS NULL
    )
    SELECT 
        ancestor_tree.task_id,
        ancestor_tree.task_no,
        ancestor_tree.title,
        ancestor_tree.task_level
    FROM ancestor_tree 
    WHERE ancestor_tree.task_id != task_id_param
    ORDER BY task_level;
END;
$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_task_ancestors(INTEGER) IS '获取指定任务的所有祖先任务（父任务、祖父任务等）';

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