-- ============================================
-- 添加"待开始"状态
-- ============================================

-- 需求类任务待开始状态（执行计划审核通过后，拆解出的子任务默认状态）
INSERT INTO task_statuses (code, name, task_type_code, sort_order, description, created_at)
VALUES ('req_pending_start', '待开始', 'requirement', 13, '执行计划审核通过，子任务待开始', CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;

-- 最小单元任务待开始状态（执行人接受后的状态）
INSERT INTO task_statuses (code, name, task_type_code, sort_order, description, created_at)
VALUES ('unit_pending_start', '待开始', 'unit_task', 4, '执行人已接受任务，待开始执行', CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- 添加状态转换规则
-- ============================================

-- 需求类任务：执行计划审核通过 -> 待开始（创建子任务时使用）
INSERT INTO task_status_transitions (task_type_code, from_status_code, to_status_code, required_role, requires_approval, is_allowed, description, created_at)
VALUES ('requirement', 'req_plan_review', 'req_pending_start', NULL, false, true, '执行计划审核通过，拆解子任务后状态变为待开始', CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- 需求类任务：待开始 -> 执行中
INSERT INTO task_status_transitions (task_type_code, from_status_code, to_status_code, required_role, requires_approval, is_allowed, description, created_at)
VALUES ('requirement', 'req_pending_start', 'req_in_progress', 'executor', false, true, '子任务开始执行', CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- 最小单元任务：待接受 -> 待开始（原来是直接到进行中，现在改为待开始）
-- 先删除旧规则
DELETE FROM task_status_transitions 
WHERE task_type_code = 'unit_task' 
  AND from_status_code = 'unit_pending_accept' 
  AND to_status_code = 'unit_in_progress';

-- 添加新规则：待接受 -> 待开始
INSERT INTO task_status_transitions (task_type_code, from_status_code, to_status_code, required_role, requires_approval, is_allowed, description, created_at)
VALUES ('unit_task', 'unit_pending_accept', 'unit_pending_start', 'executor', false, true, '执行人接受任务后进入待开始状态', CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- 最小单元任务：待开始 -> 进行中
INSERT INTO task_status_transitions (task_type_code, from_status_code, to_status_code, required_role, requires_approval, is_allowed, description, created_at)
VALUES ('unit_task', 'unit_pending_start', 'unit_in_progress', 'executor', false, true, '执行人开始执行任务', CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- ============================================
-- 验证新增的状态
-- ============================================
-- SELECT * FROM task_statuses WHERE code IN ('req_pending_start', 'unit_pending_start');
-- SELECT * FROM task_status_transitions WHERE to_status_code IN ('req_pending_start', 'unit_pending_start', 'unit_in_progress');
