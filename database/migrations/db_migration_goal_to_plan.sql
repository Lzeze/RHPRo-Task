-- ============================================
-- 调整目标关联:从思路方案改为执行计划
-- ============================================

-- 1. 修改 requirement_goals 表,将 task_id 改为 execution_plan_id
ALTER TABLE requirement_goals 
    DROP CONSTRAINT IF EXISTS requirement_goals_task_id_fkey;

ALTER TABLE requirement_goals 
    RENAME COLUMN task_id TO execution_plan_id;

ALTER TABLE requirement_goals 
    ADD CONSTRAINT requirement_goals_execution_plan_id_fkey 
    FOREIGN KEY (execution_plan_id) REFERENCES execution_plans(id) ON DELETE CASCADE;

-- 2. 更新索引
DROP INDEX IF EXISTS idx_requirement_goals_task_id;
CREATE INDEX idx_requirement_goals_execution_plan_id ON requirement_goals(execution_plan_id);

-- 3. 更新注释
COMMENT ON COLUMN requirement_goals.execution_plan_id IS '关联的执行计划ID';

-- 4. 删除 goal_no 和 sort_order 的唯一约束(如果存在)
ALTER TABLE requirement_goals 
    DROP CONSTRAINT IF EXISTS requirement_goals_task_id_goal_no_key;

-- 5. 添加新的唯一约束
ALTER TABLE requirement_goals 
    ADD CONSTRAINT requirement_goals_execution_plan_id_goal_no_key 
    UNIQUE (execution_plan_id, goal_no);

-- 6. 更新状态转换规则:移除 req_pending_goal 和 req_goal_review 相关规则
DELETE FROM task_status_transitions 
WHERE from_status_code IN ('req_pending_goal', 'req_goal_review', 'req_goal_rejected')
   OR to_status_code IN ('req_pending_goal', 'req_goal_review', 'req_goal_rejected');

-- 7. 添加新的状态转换规则:接受任务后直接进入待提交计划
INSERT INTO task_status_transitions (task_type_code, from_status_code, to_status_code, required_role, requires