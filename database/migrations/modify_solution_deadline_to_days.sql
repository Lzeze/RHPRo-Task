-- 修改 tasks 表：将 solution_deadline 从 timestamp 类型改为 INTEGER 类型（天数）

-- 1. 创建临时列用于存储转换后的数据
ALTER TABLE tasks
ADD COLUMN solution_deadline_temp INTEGER DEFAULT NULL;

-- 2. 转换现有数据：计算从现在到 solution_deadline 的天数（如果存在）
UPDATE tasks
SET solution_deadline_temp = EXTRACT(DAY FROM (solution_deadline - NOW()))
WHERE solution_deadline IS NOT NULL;

-- 3. 删除旧列
ALTER TABLE tasks
DROP COLUMN solution_deadline;

-- 4. 重命名临时列
ALTER TABLE tasks
RENAME COLUMN solution_deadline_temp TO solution_deadline;

-- 5. 修改列类型并添加注释
ALTER TABLE tasks
ALTER COLUMN solution_deadline TYPE INTEGER USING solution_deadline;

-- 6. 添加列注释
COMMENT ON COLUMN tasks.solution_deadline IS '思路方案截止天数（需求类任务创建时可设定，表示执行人接受任务后需在N天内提交方案，0表示不限制）';
