-- 添加思路方案截止时间字段
-- 用于需求类任务，创建人可设定执行人需在此时间前提交思路方案

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS solution_deadline TIMESTAMP WITH TIME ZONE;

COMMENT ON COLUMN tasks.solution_deadline IS '思路方案截止时间（执行人需在此时间前提交思路方案）';

-- 创建索引，用于查询即将到期的任务
CREATE INDEX IF NOT EXISTS idx_tasks_solution_deadline ON tasks(solution_deadline) WHERE solution_deadline IS NOT NULL;
