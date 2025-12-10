-- 为 requirement_solutions 表添加 title 字段
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='requirement_solutions' AND column_name='title') THEN
        ALTER TABLE requirement_solutions
        ADD COLUMN title VARCHAR(500) DEFAULT '';
    END IF;
END $$;

-- 为 execution_plans 表添加 title 字段
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='execution_plans' AND column_name='title') THEN
        ALTER TABLE execution_plans
        ADD COLUMN title VARCHAR(500) DEFAULT '';
    END IF;
END $$;

-- 为现有的记录设置默认标题（基于版本号）
UPDATE requirement_solutions
SET title = '思路方案 v' || version::text
WHERE title = '' OR title IS NULL;

UPDATE execution_plans
SET title = '执行计划 v' || version::text
WHERE title = '' OR title IS NULL;

-- 修改字段为非空
ALTER TABLE requirement_solutions
ALTER COLUMN title SET NOT NULL;

ALTER TABLE execution_plans
ALTER COLUMN title SET NOT NULL;

-- 添加列注释
COMMENT ON COLUMN requirement_solutions.title IS '方案标题（用于在列表中快速识别）';
COMMENT ON COLUMN execution_plans.title IS '执行计划标题（用于在列表中快速识别）';
