-- 补充 unit_pending_assign 状态
INSERT INTO task_statuses (code, name, task_type_code, sort_order, description) VALUES
    ('unit_pending_assign', '待指派', 'unit_task', 2, '未指派执行人，等待分配')
ON CONFLICT (code) DO NOTHING;

-- 调整 unit_pending_accept 的 sort_order
UPDATE task_statuses SET sort_order = 3 WHERE code = 'unit_pending_accept';
