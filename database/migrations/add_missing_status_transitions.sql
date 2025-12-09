-- 补充缺失的 unit_task 状态转换规则

-- unit_pending_assign 相关规则
INSERT INTO task_status_transitions 
(task_type_code, from_status_code, to_status_code, required_role, requires_approval, is_allowed, description) 
VALUES
-- 待指派 → 待接受（创建人指派执行人）
('unit_task', 'unit_pending_assign', 'unit_pending_accept', 'creator', false, true, '指派执行人'),

-- 待接受 → 待指派（执行人拒绝）
('unit_task', 'unit_pending_accept', 'unit_pending_assign', 'executor', false, true, '执行人拒绝，回到待指派')

ON CONFLICT DO NOTHING;
