-- ============================================
-- 任务流程调整: 将目标从方案审核移到执行计划审核
-- 执行前请备份数据库！
-- ============================================

-- 1. 新增任务状态
INSERT INTO task_statuses (code, name, task_type_code, sort_order, description) VALUES
    ('req_pending_solution', '待提交方案', 'requirement', 4, '执行人已接受，需提交解决方案'),
    ('req_solution_review', '方案审核中', 'requirement', 5, '解决方案审核中'),
    ('req_solution_rejected', '方案被驳回', 'requirement', 6, '解决方案被驳回，需重新提交')
ON CONFLICT (code) DO NOTHING;

-- 2. 删除旧的目标相关状态(如果不再使用)
-- 注意: 如果现有数据使用了这些状态，请先迁移数据
DELETE FROM task_status_transitions 
WHERE from_status_code IN ('req_pending_goal', 'req_goal_review', 'req_goal_rejected')
   OR to_status_code IN ('req_pending_goal', 'req_goal_review', 'req_goal_rejected');

-- 可选: 删除旧状态定义(谨慎操作)
-- DELETE FROM task_statuses WHERE code IN ('req_pending_goal', 'req_goal_review', 'req_goal_rejected');

-- 3. 修改 requirement_goals 表结构
-- 3.1 删除旧的外键约束
ALTER TABLE requirement_goals 
    DROP CONSTRAINT IF EXISTS requirement_goals_task_id_fkey;

-- 3.2 删除旧的唯一约束
ALTER TABLE requirement_goals 
    DROP CONSTRAINT IF EXISTS requirement_goals_task_id_goal_no_key;

-- 3.3 重命名列
ALTER TABLE requirement_goals 
    RENAME COLUMN task_id TO execution_plan_id;

-- 3.4 添加新的外键约束
ALTER TABLE requirement_goals 
    ADD CONSTRAINT requirement_goals_execution_plan_id_fkey 
    FOREIGN KEY (execution_plan_id) REFERENCES execution_plans(id) ON DELETE CASCADE;

-- 3.5 添加新的唯一约束
ALTER TABLE requirement_goals 
    ADD CONSTRAINT requirement_goals_execution_plan_id_goal_no_key 
    UNIQUE (execution_plan_id, goal_no);

-- 3.6 更新索引
DROP INDEX IF EXISTS idx_requirement_goals_task_id;
CREATE INDEX idx_requirement_goals_execution_plan_id ON requirement_goals(execution_plan_id);

-- 3.7 更新列注释
COMMENT ON COLUMN requirement_goals.execution_plan_id IS '关联的执行计划ID';

-- 4. 更新状态转换规则
-- 4.1 添加新的转换规则
INSERT INTO task_status_transitions (task_type_code, from_status_code, to_status_code, required_role, requires_approval, is_allowed, description) VALUES
    -- 接受任务后进入待提交方案状态
    ('requirement', 'req_pending_accept', 'req_pending_solution', 'executor', false, true, '执行人接受任务，进入待提交方案状态'),
    
    -- 提交方案后进入方案审核
    ('requirement', 'req_pending_solution', 'req_solution_review', 'executor', true, true, '提交解决方案，进入审核'),
    
    -- 方案审核通过，进入待提交计划
    ('requirement', 'req_solution_review', 'req_pending_plan', 'reviewer', true, true, '方案审核通过'),
    
    -- 方案审核驳回
    ('requirement', 'req_solution_review', 'req_solution_rejected', 'reviewer', true, true, '方案审核驳回'),
    
    -- 方案被驳回后重新提交
    ('requirement', 'req_solution_rejected', 'req_pending_solution', 'executor', false, true, '重新提交方案'),
    
    -- 取消任务
    ('requirement', 'req_pending_solution', 'req_cancelled', 'creator', false, true, '取消任务'),
    ('requirement', 'req_solution_review', 'req_cancelled', 'creator', false, true, '取消任务'),
    ('requirement', 'req_solution_rejected', 'req_cancelled', 'creator', false, true, '取消任务'),
    
    -- 标记为受阻
    ('requirement', 'req_pending_solution', 'req_blocked', 'executor', false, true, '标记为受阻')
ON CONFLICT (task_type_code, from_status_code, to_status_code) DO NOTHING;

-- 4.2 删除旧的转换规则(如果不再使用)
-- DELETE FROM task_status_transitions 
-- WHERE (from_status_code = 'req_pending_accept' AND to_status_code = 'req_pending_goal')
--    OR (from_status_code = 'req_pending_goal' AND to_status_code = 'req_goal_review');

-- 5. 数据迁移(如果有现有数据)
-- 注意: 这里需要根据实际情况调整
-- 示例: 将现有目标数据关联到对应的执行计划
-- UPDATE requirement_goals rg
-- SET execution_plan_id = (
--     SELECT ep.id 
--     FROM execution_plans ep 
--     WHERE ep.task_id = rg.execution_plan_id -- 这里 execution_plan_id 原来是 task_id
--     ORDER BY ep.version DESC 
--     LIMIT 1
-- )
-- WHERE EXISTS (
--     SELECT 1 FROM execution_plans ep WHERE ep.task_id = rg.execution_plan_id
-- );

-- 6. 验证迁移结果
-- 检查新状态是否创建成功
SELECT code, name, task_type_code FROM task_statuses 
WHERE code IN ('req_pending_solution', 'req_solution_review', 'req_solution_rejected');

-- 检查转换规则是否创建成功
SELECT task_type_code, from_status_code, to_status_code, description 
FROM task_status_transitions 
WHERE to_status_code IN ('req_pending_solution', 'req_solution_review', 'req_solution_rejected')
   OR from_status_code IN ('req_pending_solution', 'req_solution_review', 'req_solution_rejected')
ORDER BY task_type_code, from_status_code;

-- 检查 requirement_goals 表结构
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'requirement_goals' AND column_name = 'execution_plan_id';

-- ============================================
-- 迁移完成提示
-- ============================================
-- 1. 检查现有任务状态，确保没有任务处于旧状态
-- 2. 如果有旧状态的任务，需要手动迁移到新状态
-- 3. 测试新流程是否正常工作
-- 4. 确认无误后可以删除旧状态定义
