-- ============================================
-- 改进后的任务管理系统 - 完整版
-- ============================================

-- ============================================
-- 改进1：支持多部门负责人
-- ============================================

-- 1.1 部门负责人关联表（多对多）
CREATE TABLE IF NOT EXISTS department_leaders (
    id SERIAL PRIMARY KEY,
    department_id INTEGER NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_primary BOOLEAN DEFAULT FALSE,  -- 是否为主要负责人
    appointed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    appointed_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(department_id, user_id)
);

CREATE INDEX idx_department_leaders_department_id ON department_leaders(department_id);
CREATE INDEX idx_department_leaders_user_id ON department_leaders(user_id);

COMMENT ON TABLE department_leaders IS '部门负责人关联表（支持一人多部门、一部门多负责人）';
COMMENT ON COLUMN department_leaders.id IS '主键ID';
COMMENT ON COLUMN department_leaders.department_id IS '部门ID';
COMMENT ON COLUMN department_leaders.user_id IS '负责人用户ID';
COMMENT ON COLUMN department_leaders.is_primary IS '是否为主要负责人（用于区分正副职）';
COMMENT ON COLUMN department_leaders.appointed_at IS '任命时间';
COMMENT ON COLUMN department_leaders.appointed_by IS '任命人用户ID';
COMMENT ON COLUMN department_leaders.created_at IS '创建时间';

-- 1.2 修改 departments 表（移除单一 leader_id，改为关联表）
ALTER TABLE departments DROP COLUMN IF EXISTS leader_id;

-- 1.3 修改 users 表（保留标识字段）
-- is_department_leader 改为查询 department_leaders 表判断

-- 1.4 查询某用户负责的所有部门
CREATE OR REPLACE VIEW v_user_departments AS
SELECT 
    dl.user_id,
    u.username,
    d.id as department_id,
    d.name as department_name,
    dl.is_primary,
    dl.appointed_at
FROM department_leaders dl
JOIN users u ON dl.user_id = u.id
JOIN departments d ON dl.department_id = d.id
WHERE d.deleted_at IS NULL;

COMMENT ON VIEW v_user_departments IS '用户负责的部门视图';

-- 1.5 查询某部门的所有负责人
CREATE OR REPLACE VIEW v_department_leaders_view AS
SELECT 
    dl.department_id,
    d.name as department_name,
    dl.user_id,
    u.username,
    u.email,
    dl.is_primary,
    dl.appointed_at
FROM department_leaders dl
JOIN users u ON dl.user_id = u.id
JOIN departments d ON dl.department_id = d.id
WHERE d.deleted_at IS NULL
ORDER BY dl.is_primary DESC, dl.appointed_at;

COMMENT ON VIEW v_department_leaders_view IS '部门负责人列表视图';

-- ============================================
-- 改进2：添加受阻状态和受阻任务管理
-- ============================================

-- 2.1 在 task_statuses 表中添加受阻状态
INSERT INTO task_statuses (code, name, task_type_code, sort_order, description) VALUES
    ('req_blocked', '受阻', 'requirement', 15, '需求任务执行受阻'),
    ('unit_blocked', '受阻', 'unit_task', 6, '单元任务执行受阻')
ON CONFLICT (code) DO NOTHING;

-- 2.2 受阻任务表
CREATE TABLE IF NOT EXISTS blocked_tasks (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    blocked_reason TEXT NOT NULL,  -- 受阻原因
    blocker_type VARCHAR(50) NOT NULL,  -- 受阻类型：dependency, resource, technical, external
    blocking_task_id INTEGER REFERENCES tasks(id),  -- 阻塞任务ID（如果是依赖其他任务）
    
    -- 解决方案
    solution_description TEXT,  -- 解决方案描述
    resolution_task_id INTEGER REFERENCES tasks(id),  -- 创建的解决任务ID
    
    -- 状态
    status VARCHAR(50) DEFAULT 'open',  -- open, in_progress, resolved
    
    -- 时间记录
    blocked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,
    
    -- 关联人员
    reported_by INTEGER NOT NULL REFERENCES users(id),
    assigned_to INTEGER REFERENCES users(id),  -- 指派解决人
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_blocked_tasks_task_id ON blocked_tasks(task_id);
CREATE INDEX idx_blocked_tasks_blocking_task_id ON blocked_tasks(blocking_task_id);
CREATE INDEX idx_blocked_tasks_resolution_task_id ON blocked_tasks(resolution_task_id);
CREATE INDEX idx_blocked_tasks_status ON blocked_tasks(status);

COMMENT ON TABLE blocked_tasks IS '受阻任务表';
COMMENT ON COLUMN blocked_tasks.id IS '主键ID';
COMMENT ON COLUMN blocked_tasks.task_id IS '被阻塞的任务ID';
COMMENT ON COLUMN blocked_tasks.blocked_reason IS '受阻原因详细描述';
COMMENT ON COLUMN blocked_tasks.blocker_type IS '受阻类型：dependency-依赖阻塞, resource-资源不足, technical-技术难题, external-外部因素';
COMMENT ON COLUMN blocked_tasks.blocking_task_id IS '阻塞任务ID（如果是依赖其他任务导致的阻塞）';
COMMENT ON COLUMN blocked_tasks.solution_description IS '解决方案描述';
COMMENT ON COLUMN blocked_tasks.resolution_task_id IS '为解决阻塞而创建的任务ID';
COMMENT ON COLUMN blocked_tasks.status IS '受阻状态：open-未解决, in_progress-解决中, resolved-已解决';
COMMENT ON COLUMN blocked_tasks.blocked_at IS '受阻时间';
COMMENT ON COLUMN blocked_tasks.resolved_at IS '解决时间';
COMMENT ON COLUMN blocked_tasks.reported_by IS '报告人用户ID';
COMMENT ON COLUMN blocked_tasks.assigned_to IS '指派解决人用户ID';
COMMENT ON COLUMN blocked_tasks.created_at IS '创建时间';
COMMENT ON COLUMN blocked_tasks.updated_at IS '更新时间';

-- 2.3 触发器：任务状态变为受阻时，记录到 blocked_tasks
CREATE OR REPLACE FUNCTION handle_task_blocked()
RETURNS TRIGGER AS $$
BEGIN
    -- 如果任务状态变为受阻
    IF NEW.status_code IN ('req_blocked', 'unit_blocked') AND 
       OLD.status_code NOT IN ('req_blocked', 'unit_blocked') THEN
        
        -- 记录到变更日志
        INSERT INTO task_change_logs (task_id, user_id, change_type, field_name, old_value, new_value, comment)
        VALUES (NEW.id, NEW.executor_id, 'status_change', 'status_code', OLD.status_code, NEW.status_code, '任务受阻');
        
        -- 发送通知给创建人
        INSERT INTO notifications (user_id, task_id, type, title, content)
        VALUES (NEW.creator_id, NEW.id, 'status_change', '任务受阻', 
                '任务 ' || NEW.task_no || ' 已标记为受阻状态');
    END IF;
    
    -- 如果任务从受阻状态恢复
    IF OLD.status_code IN ('req_blocked', 'unit_blocked') AND 
       NEW.status_code NOT IN ('req_blocked', 'unit_blocked') THEN
        
        -- 更新 blocked_tasks 表
        UPDATE blocked_tasks 
        SET status = 'resolved',
            resolved_at = CURRENT_TIMESTAMP,
            updated_at = CURRENT_TIMESTAMP
        WHERE task_id = NEW.id AND status != 'resolved';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_handle_task_blocked
    AFTER UPDATE OF status_code ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION handle_task_blocked();

COMMENT ON FUNCTION handle_task_blocked() IS '处理任务受阻状态变更';

-- ============================================
-- 改进3：增强状态机约束
-- ============================================

-- 3.1 任务状态转换规则表
CREATE TABLE IF NOT EXISTS task_status_transitions (
    id SERIAL PRIMARY KEY,
    task_type_code VARCHAR(50) NOT NULL REFERENCES task_types(code),
    from_status_code VARCHAR(50) NOT NULL REFERENCES task_statuses(code),
    to_status_code VARCHAR(50) NOT NULL REFERENCES task_statuses(code),
    required_role VARCHAR(50),  -- 需要的角色：creator, executor, reviewer
    requires_approval BOOLEAN DEFAULT FALSE,  -- 是否需要审核
    is_allowed BOOLEAN DEFAULT TRUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(task_type_code, from_status_code, to_status_code)
);

CREATE INDEX idx_status_transitions_task_type ON task_status_transitions(task_type_code);
CREATE INDEX idx_status_transitions_from_status ON task_status_transitions(from_status_code);

COMMENT ON TABLE task_status_transitions IS '任务状态转换规则表（状态机配置）';
COMMENT ON COLUMN task_status_transitions.id IS '主键ID';
COMMENT ON COLUMN task_status_transitions.task_type_code IS '任务类型编码';
COMMENT ON COLUMN task_status_transitions.from_status_code IS '源状态编码';
COMMENT ON COLUMN task_status_transitions.to_status_code IS '目标状态编码';
COMMENT ON COLUMN task_status_transitions.required_role IS '需要的角色：creator-创建人, executor-执行人, reviewer-审核人';
COMMENT ON COLUMN task_status_transitions.requires_approval IS '是否需要审核批准';
COMMENT ON COLUMN task_status_transitions.is_allowed IS '是否允许此转换';
COMMENT ON COLUMN task_status_transitions.description IS '转换说明';
COMMENT ON COLUMN task_status_transitions.created_at IS '创建时间';

-- 3.2 初始化需求任务的状态转换规则
INSERT INTO task_status_transitions (task_type_code, from_status_code, to_status_code, required_role, requires_approval, description) VALUES
    -- 草稿阶段
    ('requirement', 'req_draft', 'req_pending_assign', 'creator', FALSE, '发布到待领池'),
    ('requirement', 'req_draft', 'req_pending_accept', 'creator', FALSE, '直接指派执行人'),
    ('requirement', 'req_draft', 'req_cancelled', 'creator', FALSE, '取消任务'),
    
    -- 待领池
    ('requirement', 'req_pending_assign', 'req_pending_accept', 'executor', FALSE, '执行人领取任务'),
    ('requirement', 'req_pending_assign', 'req_cancelled', 'creator', FALSE, '取消任务'),
    
    -- 待接受
    ('requirement', 'req_pending_accept', 'req_pending_goal', 'executor', FALSE, '执行人接受任务'),
    ('requirement', 'req_pending_accept', 'req_pending_assign', 'executor', FALSE, '执行人拒绝，回到待领池'),
    ('requirement', 'req_pending_accept', 'req_cancelled', 'creator', FALSE, '取消任务'),
    
    -- 待提交目标
    ('requirement', 'req_pending_goal', 'req_goal_review', 'executor', TRUE, '提交目标和方案，进入审核'),
    ('requirement', 'req_pending_goal', 'req_blocked', 'executor', FALSE, '标记为受阻'),
    
    -- 目标审核中
    ('requirement', 'req_goal_review', 'req_pending_plan', 'reviewer', TRUE, '目标审核通过'),
    ('requirement', 'req_goal_review', 'req_goal_rejected', 'reviewer', TRUE, '目标审核驳回'),
    
    -- 目标被驳回
    ('requirement', 'req_goal_rejected', 'req_pending_goal', 'executor', FALSE, '重新提交目标'),
    ('requirement', 'req_goal_rejected', 'req_cancelled', 'creator', FALSE, '取消任务'),
    
    -- 待提交计划
    ('requirement', 'req_pending_plan', 'req_plan_review', 'executor', TRUE, '提交执行计划，进入审核'),
    ('requirement', 'req_pending_plan', 'req_blocked', 'executor', FALSE, '标记为受阻'),
    
    -- 计划审核中
    ('requirement', 'req_plan_review', 'req_in_progress', 'reviewer', TRUE, '计划审核通过，开始执行'),
    ('requirement', 'req_plan_review', 'req_plan_rejected', 'reviewer', TRUE, '计划审核驳回'),
    
    -- 计划被驳回
    ('requirement', 'req_plan_rejected', 'req_pending_plan', 'executor', FALSE, '重新提交计划'),
    ('requirement', 'req_plan_rejected', 'req_cancelled', 'creator', FALSE, '取消任务'),
    
    -- 执行中
    ('requirement', 'req_in_progress', 'req_completed', 'executor', FALSE, '任务完成'),
    ('requirement', 'req_in_progress', 'req_blocked', 'executor', FALSE, '标记为受阻'),
    ('requirement', 'req_in_progress', 'req_cancelled', 'creator', FALSE, '取消任务'),
    
    -- 受阻状态
    ('requirement', 'req_blocked', 'req_in_progress', 'executor', FALSE, '解除受阻，继续执行'),
    ('requirement', 'req_blocked', 'req_pending_plan', 'executor', FALSE, '回到计划阶段'),
    ('requirement', 'req_blocked', 'req_cancelled', 'creator', FALSE, '取消任务')
ON CONFLICT (task_type_code, from_status_code, to_status_code) DO NOTHING;

-- 3.3 初始化单元任务的状态转换规则
INSERT INTO task_status_transitions (task_type_code, from_status_code, to_status_code, required_role, requires_approval, description) VALUES
    ('unit_task', 'unit_draft', 'unit_pending_accept', 'creator', FALSE, '指派执行人'),
    ('unit_task', 'unit_draft', 'unit_cancelled', 'creator', FALSE, '取消任务'),
    
    ('unit_task', 'unit_pending_accept', 'unit_in_progress', 'executor', FALSE, '执行人接受并开始'),
    ('unit_task', 'unit_pending_accept', 'unit_cancelled', 'creator', FALSE, '取消任务'),
    
    ('unit_task', 'unit_in_progress', 'unit_completed', 'executor', FALSE, '任务完成'),
    ('unit_task', 'unit_in_progress', 'unit_blocked', 'executor', FALSE, '标记为受阻'),
    ('unit_task', 'unit_in_progress', 'unit_cancelled', 'creator', FALSE, '取消任务'),
    
    ('unit_task', 'unit_blocked', 'unit_in_progress', 'executor', FALSE, '解除受阻，继续执行'),
    ('unit_task', 'unit_blocked', 'unit_cancelled', 'creator', FALSE, '取消任务')
ON CONFLICT (task_type_code, from_status_code, to_status_code) DO NOTHING;

-- 3.4 状态转换验证函数
CREATE OR REPLACE FUNCTION validate_status_transition()
RETURNS TRIGGER AS $$
DECLARE
    is_valid BOOLEAN;
    transition_rule RECORD;
BEGIN
    -- 如果状态没有变化，直接通过
    IF NEW.status_code = OLD.status_code THEN
        RETURN NEW;
    END IF;
    
    -- 查询状态转换规则
    SELECT * INTO transition_rule
    FROM task_status_transitions
    WHERE task_type_code = NEW.task_type_code
      AND from_status_code = OLD.status_code
      AND to_status_code = NEW.status_code
      AND is_allowed = TRUE;
    
    -- 如果没有找到允许的转换规则
    IF NOT FOUND THEN
        RAISE EXCEPTION '不允许的状态转换: % 从 % 到 %', 
            NEW.task_type_code, OLD.status_code, NEW.status_code;
    END IF;
    
    -- 记录状态变更日志
    INSERT INTO task_change_logs (task_id, user_id, change_type, field_name, old_value, new_value, comment)
    VALUES (NEW.id, COALESCE(NEW.executor_id, NEW.creator_id), 'status_change', 
            'status_code', OLD.status_code, NEW.status_code, 
            transition_rule.description);
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_validate_status_transition
    BEFORE UPDATE OF status_code ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION validate_status_transition();

COMMENT ON FUNCTION validate_status_transition() IS '验证任务状态转换是否符合状态机规则';

-- 3.5 查询允许的状态转换
CREATE OR REPLACE FUNCTION get_allowed_transitions(
    p_task_id INTEGER,
    p_user_id INTEGER
)
RETURNS TABLE (
    to_status_code VARCHAR,
    to_status_name VARCHAR,
    required_role VARCHAR,
    requires_approval BOOLEAN,
    description TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        tst.to_status_code,
        ts.name as to_status_name,
        tst.required_role,
        tst.requires_approval,
        tst.description
    FROM tasks t
    JOIN task_status_transitions tst ON 
        t.task_type_code = tst.task_type_code AND
        t.status_code = tst.from_status_code
    JOIN task_statuses ts ON tst.to_status_code = ts.code
    WHERE t.id = p_task_id
      AND tst.is_allowed = TRUE
      AND (
          tst.required_role IS NULL OR
          (tst.required_role = 'creator' AND t.creator_id = p_user_id) OR
          (tst.required_role = 'executor' AND t.executor_id = p_user_id) OR
          (tst.required_role = 'reviewer' AND EXISTS (
              SELECT 1 FROM task_participants tp 
              WHERE tp.task_id = t.id 
                AND tp.user_id = p_user_id 
                AND tp.role = 'reviewer'
          ))
      );
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_allowed_transitions(INTEGER, INTEGER) IS '获取当前用户对指定任务允许的状态转换';

-- ============================================
-- 改进4：增强审核记录与审核对象关联
-- ============================================

-- 4.1 审核会话表（审核流程管理）
CREATE TABLE IF NOT EXISTS review_sessions (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    review_type VARCHAR(50) NOT NULL,  -- goal_review, solution_review, plan_review
    target_type VARCHAR(50) NOT NULL,  -- requirement_goals, requirement_solutions, execution_plans
    target_id INTEGER NOT NULL,  -- 被审核对象的ID
    
    -- 审核发起
    initiated_by INTEGER NOT NULL REFERENCES users(id),
    initiated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- 审核状态
    status VARCHAR(50) DEFAULT 'pending',  -- pending, in_review, approved, rejected, cancelled
    
    -- 审核模式
    review_mode VARCHAR(50) NOT NULL,  -- single-单人审核, jury-陪审团审核
    required_approvals INTEGER DEFAULT 1,  -- 需要的通过票数
    
    -- 最终决策
    final_decision VARCHAR(50),  -- approved, rejected
    final_decision_by INTEGER REFERENCES users(id),  -- 最终决策人（创建人）
    final_decision_at TIMESTAMP WITH TIME ZONE,
    final_decision_comment TEXT,
    
    -- 时间记录
    completed_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_review_sessions_task_id ON review_sessions(task_id);
CREATE INDEX idx_review_sessions_status ON review_sessions(status);
CREATE INDEX idx_review_sessions_target ON review_sessions(target_type, target_id);

COMMENT ON TABLE review_sessions IS '审核会话表（管理整个审核流程）';
COMMENT ON COLUMN review_sessions.id IS '主键ID';
COMMENT ON COLUMN review_sessions.task_id IS '关联的任务ID';
COMMENT ON COLUMN review_sessions.review_type IS '审核类型：goal_review-目标审核, solution_review-方案审核, plan_review-计划审核';
COMMENT ON COLUMN review_sessions.target_type IS '被审核对象的表名：requirement_goals, requirement_solutions, execution_plans';
COMMENT ON COLUMN review_sessions.target_id IS '被审核对象的ID';
COMMENT ON COLUMN review_sessions.initiated_by IS '发起审核的用户ID（通常是执行人）';
COMMENT ON COLUMN review_sessions.initiated_at IS '发起审核时间';
COMMENT ON COLUMN review_sessions.status IS '审核状态：pending-待审核, in_review-审核中, approved-已通过, rejected-已驳回, cancelled-已取消';
COMMENT ON COLUMN review_sessions.review_mode IS '审核模式：single-单人审核（创建人）, jury-陪审团审核（多人投票）';
COMMENT ON COLUMN review_sessions.required_approvals IS '需要的通过票数（陪审团模式）';
COMMENT ON COLUMN review_sessions.final_decision IS '最终决策：approved-通过, rejected-驳回';
COMMENT ON COLUMN review_sessions.final_decision_by IS '最终决策人用户ID（通常是任务创建人）';
COMMENT ON COLUMN review_sessions.final_decision_at IS '最终决策时间';
COMMENT ON COLUMN review_sessions.final_decision_comment IS '最终决策说明';
COMMENT ON COLUMN review_sessions.completed_at IS '审核完成时间';
COMMENT ON COLUMN review_sessions.created_at IS '创建时间';
COMMENT ON COLUMN review_sessions.updated_at IS '更新时间';

-- 4.2 重构审核记录表（关联到审核会话）
DROP TABLE IF EXISTS review_records CASCADE;

CREATE TABLE IF NOT EXISTS review_records (
    id SERIAL PRIMARY KEY,
    review_session_id INTEGER NOT NULL REFERENCES review_sessions(id) ON DELETE CASCADE,
    reviewer_id INTEGER NOT NULL REFERENCES users(id),
    reviewer_role VARCHAR(50),  -- creator, jury, expert
    
    -- 审核意见
    opinion VARCHAR(50) NOT NULL,  -- approve, reject, abstain
    comment TEXT,
    score INTEGER,  -- 评分（可选，1-5分）
    attachments JSONB,
    
    -- 审核权重（陪审团模式）
    vote_weight DECIMAL(3,2) DEFAULT 1.0,  -- 投票权重（0.5-1.5）
    
    -- 时间记录
    reviewed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(review_session_id, reviewer_id)
);

CREATE INDEX idx_review_records_session_id ON review_records(review_session_id);
CREATE INDEX idx_review_records_reviewer_id ON review_records(reviewer_id);

COMMENT ON TABLE review_records IS '审核记录表（存储每个审核人的意见）';
COMMENT ON COLUMN review_records.id IS '主键ID';
COMMENT ON COLUMN review_records.review_session_id IS '关联的审核会话ID';
COMMENT ON COLUMN review_records.reviewer_id IS '审核人用户ID';
COMMENT ON COLUMN review_records.reviewer_role IS '审核人角色：creator-创建人, jury-陪审团成员, expert-专家';
COMMENT ON COLUMN review_records.opinion IS '审核意见：approve-同意, reject-拒绝, abstain-弃权';
COMMENT ON COLUMN review_records.comment IS '审核意见详细说明';
COMMENT ON COLUMN review_records.score IS '评分（1-5分，可选）';
COMMENT ON COLUMN review_records.attachments IS '附件信息JSON';
COMMENT ON COLUMN review_records.vote_weight IS '投票权重（陪审团模式，默认1.0）';
COMMENT ON COLUMN review_records.reviewed_at IS '审核时间';
COMMENT ON COLUMN review_records.created_at IS '创建时间';
COMMENT ON COLUMN review_records.updated_at IS '更新时间';

-- 4.3 审核会话统计视图

CREATE TRIGGER update_review_records_updated_at BEFORE UPDATE ON review_records
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_review_sessions_updated_at BEFORE UPDATE ON review_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE OR REPLACE VIEW v_review_session_stats AS
SELECT 
    rs.id as session_id,
    rs.task_id,
    rs.review_type,
    rs.status,
    rs.review_mode,
    rs.required_approvals,
    COUNT(rr.id) as total_reviewers,
    SUM(CASE WHEN rr.opinion = 'approve' THEN 1 ELSE 0 END) as approve_count,
    SUM(CASE WHEN rr.opinion = 'reject' THEN 1 ELSE 0 END) as reject_count,
    SUM(CASE WHEN rr.opinion = 'abstain' THEN 1 ELSE 0 END) as abstain_count,
    SUM(CASE WHEN rr.opinion = 'approve' THEN rr.vote_weight ELSE 0 END) as approve_weight,
    AVG(rr.score) as average_score,
    rs.final_decision,
    rs.final_decision_by,
    rs.final_decision_at
FROM review_sessions rs
LEFT JOIN review_records rr ON rs.id = rr.review_session_id
GROUP BY rs.id;

COMMENT ON VIEW v_review_session_stats IS '审核会话统计视图（汇总审核意见）';

-- 4.4 触发器：更新审核会话状态
CREATE OR REPLACE FUNCTION update_review_session_status()
RETURNS TRIGGER AS $$
DECLARE
    session_record RECORD;
    stats_record RECORD;
BEGIN
    -- 获取审核会话信息
    SELECT * INTO session_record
    FROM review_sessions
    WHERE id = NEW.review_session_id;
    
    -- 获取统计信息
    SELECT * INTO stats_record
    FROM v_review_session_stats
    WHERE session_id = NEW.review_session_id;
    
    -- 如果是单人审核模式（创建人审核）
    IF session_record.review_mode = 'single' THEN
        -- 创建人提交意见后，直接更新状态
        IF NEW.reviewer_role = 'creator' THEN
            UPDATE review_sessions
            SET status = CASE 
                    WHEN NEW.opinion = 'approve' THEN 'approved'
                    WHEN NEW.opinion = 'reject' THEN 'rejected'
                    ELSE status
                END,
                final_decision = NEW.opinion,
                final_decision_by = NEW.reviewer_id,
                final_decision_at = CURRENT_TIMESTAMP,
                completed_at = CURRENT_TIMESTAMP,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.review_session_id;
        END IF;
    
    -- 如果是陪审团模式
    ELSIF session_record.review_mode = 'jury' THEN
        -- 更新会话状态为审核中
        IF session_record.status = 'pending' THEN
            UPDATE review_sessions
            SET status = 'in_review',
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.review_session_id;
        END IF;
        
        -- 这里不自动决策，等待创建人最终确认
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_review_session_status
    AFTER INSERT OR UPDATE ON review_records
    FOR EACH ROW
    EXECUTE FUNCTION update_review_session_status();

COMMENT ON FUNCTION update_review_session_status() IS '更新审核会话状态（根据审核意见）';

-- 4.5 创建人做最终决策的函数
CREATE OR REPLACE FUNCTION finalize_review_decision(
    p_session_id INTEGER,
    p_decision_by INTEGER,
    p_decision VARCHAR(50),
    p_comment TEXT DEFAULT NULL
)
RETURNS BOOLEAN AS $
DECLARE
    session_record RECORD;
    task_record RECORD;
    new_status VARCHAR(50);
BEGIN
    -- 获取审核会话信息
    SELECT * INTO session_record
    FROM review_sessions
    WHERE id = p_session_id;
    
    -- 获取任务信息
    SELECT * INTO task_record
    FROM tasks
    WHERE id = session_record.task_id;
    
    -- 验证决策人是否为任务创建人
    IF task_record.creator_id != p_decision_by THEN
        RAISE EXCEPTION '只有任务创建人可以做最终决策';
    END IF;
    
    -- 验证决策值
    IF p_decision NOT IN ('approved', 'rejected') THEN
        RAISE EXCEPTION '决策必须是 approved 或 rejected';
    END IF;
    
    -- 更新审核会话
    UPDATE review_sessions
    SET final_decision = p_decision,
        final_decision_by = p_decision_by,
        final_decision_at = CURRENT_TIMESTAMP,
        final_decision_comment = p_comment,
        status = p_decision,
        completed_at = CURRENT_TIMESTAMP,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_session_id;
    
    -- 根据审核类型和决策结果，更新任务状态
    IF session_record.review_type = 'goal_review' THEN
        IF p_decision = 'approved' THEN
            new_status := 'req_pending_plan';
        ELSE
            new_status := 'req_goal_rejected';
        END IF;
    ELSIF session_record.review_type = 'solution_review' THEN
        IF p_decision = 'approved' THEN
            new_status := 'req_pending_plan';
        ELSE
            new_status := 'req_goal_rejected';
        END IF;
    ELSIF session_record.review_type = 'plan_review' THEN
        IF p_decision = 'approved' THEN
            new_status := 'req_in_progress';
        ELSE
            new_status := 'req_plan_rejected';
        END IF;
    END IF;
    
    -- 更新任务状态
    UPDATE tasks
    SET status_code = new_status,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = session_record.task_id;
    
    -- 更新被审核对象的状态
    IF session_record.target_type = 'requirement_goals' THEN
        UPDATE requirement_goals
        SET status = p_decision,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = session_record.target_id;
    ELSIF session_record.target_type = 'requirement_solutions' THEN
        UPDATE requirement_solutions
        SET status = p_decision,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = session_record.target_id;
    ELSIF session_record.target_type = 'execution_plans' THEN
        UPDATE execution_plans
        SET status = p_decision,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = session_record.target_id;
    END IF;
    
    -- 发送通知给执行人
    INSERT INTO notifications (user_id, task_id, type, title, content)
    VALUES (
        task_record.executor_id,
        session_record.task_id,
        'review_result',
        '审核结果通知',
        '您的' || session_record.review_type || '审核' || 
        CASE WHEN p_decision = 'approved' THEN '已通过' ELSE '已驳回' END
    );
    
    RETURN TRUE;
END;
$ LANGUAGE plpgsql;

COMMENT ON FUNCTION finalize_review_decision(INTEGER, INTEGER, VARCHAR, TEXT) IS '创建人对审核做最终决策（陪审团模式）';

-- 4.6 触发器：更新被审核对象的状态
CREATE OR REPLACE FUNCTION update_review_target_status()
RETURNS TRIGGER AS $
BEGIN
    -- 根据审核会话的最终决策，更新被审核对象的状态
    IF NEW.final_decision IS NOT NULL AND OLD.final_decision IS NULL THEN
        CASE NEW.target_type
            WHEN 'requirement_goals' THEN
                UPDATE requirement_goals
                SET status = NEW.final_decision,
                    updated_at = CURRENT_TIMESTAMP
                WHERE id = NEW.target_id;
            
            WHEN 'requirement_solutions' THEN
                UPDATE requirement_solutions
                SET status = NEW.final_decision,
                    updated_at = CURRENT_TIMESTAMP
                WHERE id = NEW.target_id;
            
            WHEN 'execution_plans' THEN
                UPDATE execution_plans
                SET status = NEW.final_decision,
                    updated_at = CURRENT_TIMESTAMP
                WHERE id = NEW.target_id;
        END CASE;
    END IF;
    
    RETURN NEW;
END;
$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_review_target_status
    AFTER UPDATE OF final_decision ON review_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_review_target_status();

COMMENT ON FUNCTION update_review_target_status() IS '审核决策后更新被审核对象的状态';

-- ============================================
-- 更新原有表的外键关联
-- ============================================

-- 更新 requirement_goals 表，关联到审核会话
ALTER TABLE requirement_goals DROP COLUMN IF EXISTS status;
ALTER TABLE requirement_goals ADD COLUMN status VARCHAR(50) DEFAULT 'draft';

COMMENT ON COLUMN requirement_goals.status IS '目标状态：draft-草稿, pending-待审核, approved-已通过, rejected-已驳回';

-- 更新 requirement_solutions 表
ALTER TABLE requirement_solutions DROP COLUMN IF EXISTS status;
ALTER TABLE requirement_solutions ADD COLUMN status VARCHAR(50) DEFAULT 'draft';

COMMENT ON COLUMN requirement_solutions.status IS '方案状态：draft-草稿, pending-待审核, approved-已通过, rejected-已驳回';

-- 更新 execution_plans 表
ALTER TABLE execution_plans DROP COLUMN IF EXISTS status;
ALTER TABLE execution_plans ADD COLUMN status VARCHAR(50) DEFAULT 'draft';

COMMENT ON COLUMN execution_plans.status IS '计划状态：draft-草稿, pending-待审核, approved-已通过, rejected-已驳回';

-- ============================================
-- 综合查询视图
-- ============================================

-- 视图：任务完整信息（包含审核进度）
CREATE OR REPLACE VIEW v_task_full_details AS
SELECT 
    t.*,
    tt.name as task_type_name,
    ts.name as status_name,
    u1.username as creator_name,
    u2.username as executor_name,
    d.name as department_name,
    
    -- 当前审核会话信息
    rs.id as current_review_session_id,
    rs.review_type as current_review_type,
    rs.status as review_status,
    rs.review_mode,
    
    -- 审核统计
    rss.total_reviewers,
    rss.approve_count,
    rss.reject_count,
    rss.average_score,
    
    -- 受阻信息
    bt.id as blocked_id,
    bt.blocked_reason,
    bt.blocker_type,
    bt.status as blocked_status
    
FROM tasks t
LEFT JOIN task_types tt ON t.task_type_code = tt.code
LEFT JOIN task_statuses ts ON t.status_code = ts.code
LEFT JOIN users u1 ON t.creator_id = u1.id
LEFT JOIN users u2 ON t.executor_id = u2.id
LEFT JOIN departments d ON t.department_id = d.id
LEFT JOIN review_sessions rs ON t.id = rs.task_id AND rs.status IN ('pending', 'in_review')
LEFT JOIN v_review_session_stats rss ON rs.id = rss.session_id
LEFT JOIN blocked_tasks bt ON t.id = bt.task_id AND bt.status != 'resolved'
WHERE t.deleted_at IS NULL;

COMMENT ON VIEW v_task_full_details IS '任务完整信息视图（包含审核和受阻信息）';

-- ============================================
-- 实用查询示例
-- ============================================

-- 示例1：检查用户是否为某部门负责人
-- SELECT EXISTS (
--     SELECT 1 FROM department_leaders 
--     WHERE department_id = 1 AND user_id = 5
-- ) as is_leader;

-- 示例2：查询用户负责的所有部门
-- SELECT * FROM v_user_departments WHERE user_id = 5;

-- 示例3：创建受阻任务并生成解决任务
-- INSERT INTO blocked_tasks (task_id, blocked_reason, blocker_type, reported_by)
-- VALUES (10, '缺少API文档', 'resource', 3);
-- 
-- -- 创建解决任务
-- INSERT INTO tasks (task_no, title, task_type_code, status_code, creator_id, parent_task_id)
-- VALUES ('UNIT-2024-999', '编写API文档', 'unit_task', 'unit_draft', 3, 10);
-- 
-- -- 关联解决任务
-- UPDATE blocked_tasks SET resolution_task_id = LAST_INSERT_ID WHERE id = ...;

-- 示例4：查询允许的状态转换
-- SELECT * FROM get_allowed_transitions(10, 5);

-- 示例5：发起审核会话（目标审核）
-- INSERT INTO review_sessions (
--     task_id, review_type, target_type, target_id, 
--     initiated_by, review_mode, required_approvals
-- )
-- VALUES (1, 'goal_review', 'requirement_goals', 1, 3, 'jury', 2);
-- 
-- -- 邀请陪审团成员
-- INSERT INTO task_participants (task_id, user_id, role, invited_by, status)
-- VALUES 
--     (1, 5, 'jury', 1, 'pending'),
--     (1, 6, 'jury', 1, 'pending'),
--     (1, 7, 'jury', 1, 'pending');

-- 示例6：陪审团成员提交审核意见
-- INSERT INTO review_records (review_session_id, reviewer_id, reviewer_role, opinion, comment, score)
-- VALUES (1, 5, 'jury', 'approve', '目标清晰，方案可行', 4);

-- 示例7：查看审核会话统计
-- SELECT * FROM v_review_session_stats WHERE session_id = 1;

-- 示例8：创建人做最终决策
-- SELECT finalize_review_decision(1, 1, 'approved', '综合陪审团意见，目标通过');

-- 示例9：查询所有受阻任务
-- SELECT 
--     t.task_no,
--     t.title,
--     bt.blocked_reason,
--     bt.blocker_type,
--     u.username as reported_by_name,
--     bt.status
-- FROM blocked_tasks bt
-- JOIN tasks t ON bt.task_id = t.id
-- JOIN users u ON bt.reported_by = u.id
-- WHERE bt.status != 'resolved'
-- ORDER BY bt.blocked_at DESC;

-- 示例10：查询任务的完整审核历史
-- SELECT 
--     rs.review_type,
--     rs.status,
--     rs.final_decision,
--     u1.username as decision_by_name,
--     rs.final_decision_at,
--     rr.reviewer_id,
--     u2.username as reviewer_name,
--     rr.opinion,
--     rr.comment,
--     rr.score
-- FROM review_sessions rs
-- LEFT JOIN review_records rr ON rs.id = rr.review_session_id
-- LEFT JOIN users u1 ON rs.final_decision_by = u1.id
-- LEFT JOIN users u2 ON rr.reviewer_id = u2.id
-- WHERE rs.task_id = 1
-- ORDER BY rs.created_at DESC, rr.reviewed_at;