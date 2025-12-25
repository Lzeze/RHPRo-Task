SET session_replication_role = 'replica';
-- public DDL
-- CREATE SCHEMA "public";
COMMENT ON SCHEMA "public" IS 'standard public schema';
ALTER SCHEMA "public" OWNER TO "pg_database_owner";

-- public.blocked_tasks DDL
DROP TABLE IF EXISTS "public"."blocked_tasks";
CREATE SEQUENCE IF NOT EXISTS "public"."blocked_tasks_id_seq";
CREATE TABLE "public"."blocked_tasks" (
"id" int4 NOT NULL DEFAULT nextval('blocked_tasks_id_seq'::regclass),
"task_id" int4 NOT NULL,
"blocked_reason" text NOT NULL,
"blocker_type" varchar(50) NOT NULL,
"blocking_task_id" int4,
"solution_description" text,
"resolution_task_id" int4,
"status" varchar(50) DEFAULT 'open'::character varying,
"blocked_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"resolved_at" timestamptz(6),
"reported_by" int4 NOT NULL,
"assigned_to" int4,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.department_leaders DDL
DROP TABLE IF EXISTS "public"."department_leaders";
CREATE SEQUENCE IF NOT EXISTS "public"."department_leaders_id_seq";
CREATE TABLE "public"."department_leaders" (
"id" int4 NOT NULL DEFAULT nextval('department_leaders_id_seq'::regclass),
"department_id" int4 NOT NULL,
"user_id" int4 NOT NULL,
"appointed_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"appointed_by" int4,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"deleted_at" timestamptz(6),
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.departments DDL
DROP TABLE IF EXISTS "public"."departments";
CREATE SEQUENCE IF NOT EXISTS "public"."departments_id_seq";
CREATE TABLE "public"."departments" (
"id" int4 NOT NULL DEFAULT nextval('departments_id_seq'::regclass),
"name" varchar(100) NOT NULL,
"description" text,
"parent_id" int4,
"status" int4 DEFAULT 1,
"created_at" timestamptz(6),
"updated_at" timestamptz(6),
"deleted_at" timestamptz(6),
PRIMARY KEY ("id"));

-- public.execution_plans DDL
DROP TABLE IF EXISTS "public"."execution_plans";
CREATE SEQUENCE IF NOT EXISTS "public"."execution_plans_id_seq";
CREATE TABLE "public"."execution_plans" (
"id" int4 NOT NULL DEFAULT nextval('execution_plans_id_seq'::regclass),
"task_id" int4 NOT NULL,
"version" int4 DEFAULT 1,
"tech_stack" text NOT NULL,
"implementation_steps" jsonb NOT NULL,
"resource_requirements" text,
"risk_assessment" text,
"status" varchar(50) DEFAULT 'pending'::character varying,
"submitted_by" int4,
"submitted_at" timestamptz(6),
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"title" varchar(500) NOT NULL DEFAULT ''::character varying,
PRIMARY KEY ("id"));

-- public.notifications DDL
DROP TABLE IF EXISTS "public"."notifications";
CREATE SEQUENCE IF NOT EXISTS "public"."notifications_id_seq";
CREATE TABLE "public"."notifications" (
"id" int4 NOT NULL DEFAULT nextval('notifications_id_seq'::regclass),
"user_id" int4 NOT NULL,
"task_id" int4,
"type" varchar(50) NOT NULL,
"title" varchar(255) NOT NULL,
"content" text,
"is_read" bool DEFAULT false,
"read_at" timestamptz(6),
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.permissions DDL
DROP TABLE IF EXISTS "public"."permissions";
CREATE SEQUENCE IF NOT EXISTS "public"."permissions_id_seq";
CREATE TABLE "public"."permissions" (
"id" int4 NOT NULL DEFAULT nextval('permissions_id_seq'::regclass),
"name" varchar(100) NOT NULL,
"description" varchar(255),
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"deleted_at" timestamptz(6),
PRIMARY KEY ("id"));

-- public.requirement_goals DDL
DROP TABLE IF EXISTS "public"."requirement_goals";
CREATE SEQUENCE IF NOT EXISTS "public"."requirement_goals_id_seq";
CREATE TABLE "public"."requirement_goals" (
"id" int4 NOT NULL DEFAULT nextval('requirement_goals_id_seq'::regclass),
"execution_plan_id" int4 NOT NULL,
"goal_no" int4 NOT NULL,
"title" varchar(255) NOT NULL,
"description" text NOT NULL,
"success_criteria" text,
"priority" int4 DEFAULT 2,
"status" varchar(50) DEFAULT 'pending'::character varying,
"sort_order" int4 DEFAULT 0,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"start_date" timestamp(6),
"end_date" timestamp(6),
PRIMARY KEY ("id"));

-- public.requirement_solutions DDL
DROP TABLE IF EXISTS "public"."requirement_solutions";
CREATE SEQUENCE IF NOT EXISTS "public"."requirement_solutions_id_seq";
CREATE TABLE "public"."requirement_solutions" (
"id" int4 NOT NULL DEFAULT nextval('requirement_solutions_id_seq'::regclass),
"task_id" int4 NOT NULL,
"version" int4 DEFAULT 1,
"content" text,
"mindmap_url" varchar(500),
"file_name" varchar(255),
"file_size" int8,
"status" varchar(50) DEFAULT 'pending'::character varying,
"submitted_by" int4,
"submitted_at" timestamptz(6),
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"title" varchar(500) NOT NULL DEFAULT ''::character varying,
"mindmap_markdown" text,
PRIMARY KEY ("id"));

-- public.review_records DDL
DROP TABLE IF EXISTS "public"."review_records";
CREATE SEQUENCE IF NOT EXISTS "public"."review_records_id_seq";
CREATE TABLE "public"."review_records" (
"id" int4 NOT NULL DEFAULT nextval('review_records_id_seq'::regclass),
"review_session_id" int4 NOT NULL,
"reviewer_id" int4 NOT NULL,
"reviewer_role" varchar(50),
"opinion" varchar(50) NOT NULL,
"comment" text,
"score" int4,
"attachments" jsonb,
"vote_weight" numeric(3,2) DEFAULT 1.0,
"reviewed_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"deleted_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.review_sessions DDL
DROP TABLE IF EXISTS "public"."review_sessions";
CREATE SEQUENCE IF NOT EXISTS "public"."review_sessions_id_seq";
CREATE TABLE "public"."review_sessions" (
"id" int4 NOT NULL DEFAULT nextval('review_sessions_id_seq'::regclass),
"task_id" int4 NOT NULL,
"review_type" varchar(50) NOT NULL,
"target_type" varchar(50) NOT NULL,
"target_id" int4 NOT NULL,
"initiated_by" int4 NOT NULL,
"initiated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"status" varchar(50) DEFAULT 'pending'::character varying,
"review_mode" varchar(50) NOT NULL,
"required_approvals" int4 DEFAULT 1,
"final_decision" varchar(50),
"final_decision_by" int4,
"final_decision_at" timestamptz(6),
"final_decision_comment" text,
"completed_at" timestamptz(6),
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"deleted_at" timestamptz(6),
PRIMARY KEY ("id"));

-- public.role_permissions DDL
DROP TABLE IF EXISTS "public"."role_permissions";
CREATE TABLE "public"."role_permissions" (
"role_id" int4 NOT NULL,
"permission_id" int4 NOT NULL,
PRIMARY KEY ("role_id",
"permission_id"));

-- public.roles DDL
DROP TABLE IF EXISTS "public"."roles";
CREATE SEQUENCE IF NOT EXISTS "public"."roles_id_seq";
CREATE TABLE "public"."roles" (
"id" int4 NOT NULL DEFAULT nextval('roles_id_seq'::regclass),
"name" varchar(50) NOT NULL,
"description" varchar(255),
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"deleted_at" timestamptz(6),
PRIMARY KEY ("id"));

-- public.task_attachments DDL
DROP TABLE IF EXISTS "public"."task_attachments";
CREATE SEQUENCE IF NOT EXISTS "public"."task_attachments_id_seq";
CREATE TABLE "public"."task_attachments" (
"id" int4 NOT NULL DEFAULT nextval('task_attachments_id_seq'::regclass),
"task_id" int4 NOT NULL,
"file_name" varchar(255) NOT NULL,
"file_url" varchar(500) NOT NULL,
"file_type" varchar(100),
"file_size" int8,
"uploaded_by" int4 NOT NULL,
"attachment_type" varchar(50),
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.task_change_logs DDL
DROP TABLE IF EXISTS "public"."task_change_logs";
CREATE SEQUENCE IF NOT EXISTS "public"."task_change_logs_id_seq";
CREATE TABLE "public"."task_change_logs" (
"id" int4 NOT NULL DEFAULT nextval('task_change_logs_id_seq'::regclass),
"task_id" int4 NOT NULL,
"user_id" int4 NOT NULL,
"change_type" varchar(50) NOT NULL,
"field_name" varchar(100),
"old_value" text,
"new_value" text,
"comment" text,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.task_comments DDL
DROP TABLE IF EXISTS "public"."task_comments";
CREATE SEQUENCE IF NOT EXISTS "public"."task_comments_id_seq";
CREATE TABLE "public"."task_comments" (
"id" int4 NOT NULL DEFAULT nextval('task_comments_id_seq'::regclass),
"task_id" int4 NOT NULL,
"user_id" int4 NOT NULL,
"content" text NOT NULL,
"parent_comment_id" int4,
"attachments" jsonb,
"is_private" bool DEFAULT false,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"deleted_at" timestamptz(6),
PRIMARY KEY ("id"));

-- public.task_milestones DDL
DROP TABLE IF EXISTS "public"."task_milestones";
CREATE SEQUENCE IF NOT EXISTS "public"."task_milestones_id_seq";
CREATE TABLE "public"."task_milestones" (
"id" int4 NOT NULL DEFAULT nextval('task_milestones_id_seq'::regclass),
"task_id" int4 NOT NULL,
"name" varchar(255) NOT NULL,
"description" text,
"target_date" date NOT NULL,
"actual_date" date,
"status" varchar(50) DEFAULT 'pending'::character varying,
"sort_order" int4 DEFAULT 0,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.task_participants DDL
DROP TABLE IF EXISTS "public"."task_participants";
CREATE SEQUENCE IF NOT EXISTS "public"."task_participants_id_seq";
CREATE TABLE "public"."task_participants" (
"id" int4 NOT NULL DEFAULT nextval('task_participants_id_seq'::regclass),
"task_id" int4 NOT NULL,
"user_id" int4 NOT NULL,
"role" varchar(50) NOT NULL,
"status" varchar(50) DEFAULT 'pending'::character varying,
"invited_by" int4,
"invited_at" timestamptz(6),
"response_at" timestamptz(6),
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.task_status_transitions DDL
DROP TABLE IF EXISTS "public"."task_status_transitions";
CREATE SEQUENCE IF NOT EXISTS "public"."task_status_transitions_id_seq";
CREATE TABLE "public"."task_status_transitions" (
"id" int4 NOT NULL DEFAULT nextval('task_status_transitions_id_seq'::regclass),
"task_type_code" varchar(50) NOT NULL,
"from_status_code" varchar(50) NOT NULL,
"to_status_code" varchar(50) NOT NULL,
"required_role" varchar(50),
"requires_approval" bool DEFAULT false,
"is_allowed" bool DEFAULT true,
"description" text,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.task_statuses DDL
DROP TABLE IF EXISTS "public"."task_statuses";
CREATE SEQUENCE IF NOT EXISTS "public"."task_statuses_id_seq";
CREATE TABLE "public"."task_statuses" (
"id" int4 NOT NULL DEFAULT nextval('task_statuses_id_seq'::regclass),
"code" varchar(50) NOT NULL,
"name" varchar(100) NOT NULL,
"task_type_code" varchar(50),
"sort_order" int4 DEFAULT 0,
"description" text,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.task_tag_rel DDL
DROP TABLE IF EXISTS "public"."task_tag_rel";
CREATE TABLE "public"."task_tag_rel" (
"task_id" int4 NOT NULL,
"tag_id" int4 NOT NULL,
PRIMARY KEY ("task_id",
"tag_id"));

-- public.task_tags DDL
DROP TABLE IF EXISTS "public"."task_tags";
CREATE SEQUENCE IF NOT EXISTS "public"."task_tags_id_seq";
CREATE TABLE "public"."task_tags" (
"id" int4 NOT NULL DEFAULT nextval('task_tags_id_seq'::regclass),
"name" varchar(50) NOT NULL,
"color" varchar(20),
"description" text,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.task_types DDL
DROP TABLE IF EXISTS "public"."task_types";
CREATE SEQUENCE IF NOT EXISTS "public"."task_types_id_seq";
CREATE TABLE "public"."task_types" (
"id" int4 NOT NULL DEFAULT nextval('task_types_id_seq'::regclass),
"code" varchar(50) NOT NULL,
"name" varchar(100) NOT NULL,
"description" text,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY ("id"));

-- public.tasks DDL
DROP TABLE IF EXISTS "public"."tasks";
CREATE SEQUENCE IF NOT EXISTS "public"."tasks_id_seq";
CREATE TABLE "public"."tasks" (
"id" int8 NOT NULL DEFAULT nextval('tasks_id_seq'::regclass),
"created_at" timestamptz(6),
"updated_at" timestamptz(6),
"deleted_at" timestamptz(6),
"task_no" varchar(50) NOT NULL,
"title" varchar(255) NOT NULL,
"description" text,
"task_type_code" varchar(50) NOT NULL,
"status_code" varchar(50) NOT NULL,
"creator_id" int8 NOT NULL,
"executor_id" int8,
"department_id" int8,
"parent_task_id" int8,
"root_task_id" int8,
"task_level" int8 DEFAULT 0,
"task_path" varchar(500),
"child_sequence" int8 DEFAULT 0,
"total_subtasks" int8 DEFAULT 0,
"completed_subtasks" int8 DEFAULT 0,
"expected_start_date" timestamptz(6),
"expected_end_date" timestamptz(6),
"actual_start_date" timestamptz(6),
"actual_end_date" timestamptz(6),
"priority" int8 DEFAULT 2,
"progress" int8 DEFAULT 0,
"is_cross_department" bool DEFAULT false,
"is_in_pool" bool DEFAULT false,
"is_template" bool DEFAULT false,
"split_from_plan_id" int8,
"split_at" timestamptz(6),
"solution_deadline" int4 DEFAULT 0,
PRIMARY KEY ("id"));

-- public.user_roles DDL
DROP TABLE IF EXISTS "public"."user_roles";
CREATE TABLE "public"."user_roles" (
"user_id" int4 NOT NULL,
"role_id" int4 NOT NULL,
PRIMARY KEY ("user_id",
"role_id"));

-- public.users DDL
DROP TABLE IF EXISTS "public"."users";
CREATE SEQUENCE IF NOT EXISTS "public"."users_id_seq";
CREATE TABLE "public"."users" (
"id" int4 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
"username" varchar(50) NOT NULL,
"email" varchar(100),
"password" varchar(255) NOT NULL,
"mobile" varchar(20) NOT NULL,
"status" int4 DEFAULT 1,
"created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
"deleted_at" timestamptz(6),
"is_department_leader" bool DEFAULT false,
"job_title" varchar(100),
"department_id" int4,
"nickname" varchar(50),
"wechat_unionid" varchar(64),
"wechat_openid" varchar(64),
"avatar" varchar(500),
PRIMARY KEY ("id"));

-- public.v_user_departments DDL
CREATE VIEW "public"."v_user_departments" AS  SELECT dl.user_id,
    u.username,
    d.id AS department_id,
    d.name AS department_name,
    dl.appointed_at
   FROM ((department_leaders dl
     JOIN users u ON ((dl.user_id = u.id)))
     JOIN departments d ON ((dl.department_id = d.id)))
  WHERE (d.deleted_at IS NULL);
COMMENT ON VIEW "public"."v_user_departments" IS '用户负责的部门视图';
-- public.v_department_leaders_view DDL
CREATE VIEW "public"."v_department_leaders_view" AS  SELECT dl.department_id,
    d.name AS department_name,
    dl.user_id,
    u.username,
    u.email,
    dl.appointed_at
   FROM ((department_leaders dl
     JOIN users u ON ((dl.user_id = u.id)))
     JOIN departments d ON ((dl.department_id = d.id)))
  WHERE (d.deleted_at IS NULL)
  ORDER BY dl.appointed_at;
COMMENT ON VIEW "public"."v_department_leaders_view" IS '部门负责人列表视图';
-- public.get_all_subtasks DDL
DROP FUNCTION IF EXISTS "public"."get_all_subtasks";
CREATE OR REPLACE FUNCTION public.get_all_subtasks(task_id_param integer)
 RETURNS TABLE(task_id integer, task_no character varying, title character varying, task_level integer, status_code character varying)
 LANGUAGE plpgsql
AS $function$
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
$function$
;
-- public.get_task_ancestors DDL
DROP FUNCTION IF EXISTS "public"."get_task_ancestors";
CREATE OR REPLACE FUNCTION public.get_task_ancestors(task_id_param integer)
 RETURNS TABLE(task_id integer, task_no character varying, title character varying, task_level integer)
 LANGUAGE plpgsql
AS $function$
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
$function$
;
-- public.update_updated_at_column DDL
DROP FUNCTION IF EXISTS "public"."update_updated_at_column";
CREATE OR REPLACE FUNCTION public.update_updated_at_column()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$function$
;
-- public.blocked_tasks Indexes
COMMENT ON TABLE "public"."blocked_tasks" IS '受阻任务表';
CREATE INDEX "idx_blocked_tasks_blocking_task_id" ON "public"."blocked_tasks" USING btree ("blocking_task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_blocked_tasks_resolution_task_id" ON "public"."blocked_tasks" USING btree ("resolution_task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_blocked_tasks_status" ON "public"."blocked_tasks" USING btree ("status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_blocked_tasks_task_id" ON "public"."blocked_tasks" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."blocked_tasks" ADD CONSTRAINT "blocked_tasks_reported_by_fkey" FOREIGN KEY ("reported_by") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."blocked_tasks" ADD CONSTRAINT "blocked_tasks_assigned_to_fkey" FOREIGN KEY ("assigned_to") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."blocked_tasks"."id" IS '主键ID';
COMMENT ON COLUMN "public"."blocked_tasks"."task_id" IS '被阻塞的任务ID';
COMMENT ON COLUMN "public"."blocked_tasks"."blocked_reason" IS '受阻原因详细描述';
COMMENT ON COLUMN "public"."blocked_tasks"."blocker_type" IS '受阻类型：dependency-依赖阻塞, resource-资源不足, technical-技术难题, external-外部因素';
COMMENT ON COLUMN "public"."blocked_tasks"."blocking_task_id" IS '阻塞任务ID（如果是依赖其他任务导致的阻塞）';
COMMENT ON COLUMN "public"."blocked_tasks"."solution_description" IS '解决方案描述';
COMMENT ON COLUMN "public"."blocked_tasks"."resolution_task_id" IS '为解决阻塞而创建的任务ID';
COMMENT ON COLUMN "public"."blocked_tasks"."status" IS '受阻状态：open-未解决, in_progress-解决中, resolved-已解决';
COMMENT ON COLUMN "public"."blocked_tasks"."blocked_at" IS '受阻时间';
COMMENT ON COLUMN "public"."blocked_tasks"."resolved_at" IS '解决时间';
COMMENT ON COLUMN "public"."blocked_tasks"."reported_by" IS '报告人用户ID';
COMMENT ON COLUMN "public"."blocked_tasks"."assigned_to" IS '指派解决人用户ID';
COMMENT ON COLUMN "public"."blocked_tasks"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."blocked_tasks"."updated_at" IS '更新时间';

-- public.department_leaders Indexes
COMMENT ON TABLE "public"."department_leaders" IS '部门负责人关联表（支持一人多部门、一部门多负责人）';
CREATE UNIQUE INDEX "department_leaders_department_id_user_id_key" ON "public"."department_leaders" USING btree ("department_id"  "pg_catalog"."int4_ops" ASC NULLS LAST,"user_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_department_leaders_deleted_at" ON "public"."department_leaders" USING btree ("deleted_at"  "pg_catalog"."timestamptz_ops" ASC NULLS LAST);
CREATE INDEX "idx_department_leaders_department_id" ON "public"."department_leaders" USING btree ("department_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_department_leaders_user_id" ON "public"."department_leaders" USING btree ("user_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."department_leaders" ADD CONSTRAINT "department_leaders_department_id_fkey" FOREIGN KEY ("department_id") REFERENCES "public"."departments" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."department_leaders" ADD CONSTRAINT "department_leaders_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."department_leaders" ADD CONSTRAINT "department_leaders_appointed_by_fkey" FOREIGN KEY ("appointed_by") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."department_leaders"."id" IS '主键ID';
COMMENT ON COLUMN "public"."department_leaders"."department_id" IS '部门ID';
COMMENT ON COLUMN "public"."department_leaders"."user_id" IS '负责人用户ID';
COMMENT ON COLUMN "public"."department_leaders"."appointed_at" IS '任命时间';
COMMENT ON COLUMN "public"."department_leaders"."appointed_by" IS '任命人用户ID';
COMMENT ON COLUMN "public"."department_leaders"."created_at" IS '创建时间';

-- public.departments Indexes
COMMENT ON TABLE "public"."departments" IS '部门表';
CREATE INDEX "idx_departments_parent_id" ON "public"."departments" USING btree ("parent_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."departments" ADD CONSTRAINT "departments_parent_id_fkey" FOREIGN KEY ("parent_id") REFERENCES "public"."departments" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."departments"."id" IS '主键ID';
COMMENT ON COLUMN "public"."departments"."name" IS '部门名称';
COMMENT ON COLUMN "public"."departments"."description" IS '部门描述';
COMMENT ON COLUMN "public"."departments"."parent_id" IS '父部门ID（支持多级部门）';
COMMENT ON COLUMN "public"."departments"."status" IS '状态：1-正常，0-禁用';
COMMENT ON COLUMN "public"."departments"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."departments"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."departments"."deleted_at" IS '软删除时间';
CREATE TRIGGER "update_departments_updated_at"
    BEFORE UPDATE
    ON "public"."departments"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.execution_plans Indexes
COMMENT ON TABLE "public"."execution_plans" IS '执行计划表';
CREATE INDEX "idx_execution_plans_task_id" ON "public"."execution_plans" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."execution_plans" ADD CONSTRAINT "execution_plans_submitted_by_fkey" FOREIGN KEY ("submitted_by") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."execution_plans"."id" IS '主键ID';
COMMENT ON COLUMN "public"."execution_plans"."task_id" IS '关联的任务ID';
COMMENT ON COLUMN "public"."execution_plans"."version" IS '计划版本号（支持多次修改）';
COMMENT ON COLUMN "public"."execution_plans"."tech_stack" IS '技术栈选型说明';
COMMENT ON COLUMN "public"."execution_plans"."implementation_steps" IS '实施步骤JSON：[{step:1, name:"步骤名", description:"描述", duration:3}]';
COMMENT ON COLUMN "public"."execution_plans"."resource_requirements" IS '资源需求说明';
COMMENT ON COLUMN "public"."execution_plans"."risk_assessment" IS '风险评估说明';
COMMENT ON COLUMN "public"."execution_plans"."status" IS '计划状态：pending-待审核，approved-已通过，rejected-已驳回';
COMMENT ON COLUMN "public"."execution_plans"."submitted_by" IS '提交人用户ID';
COMMENT ON COLUMN "public"."execution_plans"."submitted_at" IS '提交时间';
COMMENT ON COLUMN "public"."execution_plans"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."execution_plans"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."execution_plans"."title" IS '执行计划标题（用于在列表中快速识别）';
CREATE TRIGGER "update_execution_plans_updated_at"
    BEFORE UPDATE
    ON "public"."execution_plans"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.notifications Indexes
COMMENT ON TABLE "public"."notifications" IS '通知消息表';
CREATE INDEX "idx_notifications_created_at" ON "public"."notifications" USING btree ("created_at"  "pg_catalog"."timestamptz_ops" ASC NULLS LAST);
CREATE INDEX "idx_notifications_is_read" ON "public"."notifications" USING btree ("is_read"  "pg_catalog"."bool_ops" ASC NULLS LAST);
CREATE INDEX "idx_notifications_user_id" ON "public"."notifications" USING btree ("user_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."notifications" ADD CONSTRAINT "notifications_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."notifications"."id" IS '主键ID';
COMMENT ON COLUMN "public"."notifications"."user_id" IS '接收通知的用户ID';
COMMENT ON COLUMN "public"."notifications"."task_id" IS '关联的任务ID';
COMMENT ON COLUMN "public"."notifications"."type" IS '通知类型：task_assigned-任务指派，review_request-审核请求，status_change-状态变更，comment-评论通知';
COMMENT ON COLUMN "public"."notifications"."title" IS '通知标题';
COMMENT ON COLUMN "public"."notifications"."content" IS '通知内容';
COMMENT ON COLUMN "public"."notifications"."is_read" IS '是否已读';
COMMENT ON COLUMN "public"."notifications"."read_at" IS '阅读时间';
COMMENT ON COLUMN "public"."notifications"."created_at" IS '创建时间';

-- public.permissions Indexes
COMMENT ON TABLE "public"."permissions" IS '权限表';
CREATE INDEX "idx_permissions_deleted_at" ON "public"."permissions" USING btree ("deleted_at"  "pg_catalog"."timestamptz_ops" ASC NULLS LAST);
CREATE INDEX "idx_permissions_name" ON "public"."permissions" USING btree ("name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "permissions_name_key" ON "public"."permissions" USING btree ("name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
COMMENT ON COLUMN "public"."permissions"."id" IS '权限ID';
COMMENT ON COLUMN "public"."permissions"."name" IS '权限名称（如：user:read）';
COMMENT ON COLUMN "public"."permissions"."description" IS '权限描述';
COMMENT ON COLUMN "public"."permissions"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."permissions"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."permissions"."deleted_at" IS '软删除时间';
CREATE TRIGGER "update_permissions_updated_at"
    BEFORE UPDATE
    ON "public"."permissions"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.requirement_goals Indexes
COMMENT ON TABLE "public"."requirement_goals" IS '需求目标表（支持多目标）';
CREATE INDEX "idx_requirement_goals_execution_plan_id" ON "public"."requirement_goals" USING btree ("execution_plan_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "requirement_goals_execution_plan_id_goal_no_key" ON "public"."requirement_goals" USING btree ("execution_plan_id"  "pg_catalog"."int4_ops" ASC NULLS LAST,"goal_no"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."requirement_goals" ADD CONSTRAINT "requirement_goals_execution_plan_id_fkey" FOREIGN KEY ("execution_plan_id") REFERENCES "public"."execution_plans" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."requirement_goals"."id" IS '主键ID';
COMMENT ON COLUMN "public"."requirement_goals"."execution_plan_id" IS '关联的执行计划ID';
COMMENT ON COLUMN "public"."requirement_goals"."goal_no" IS '目标编号（同一任务内的序号）';
COMMENT ON COLUMN "public"."requirement_goals"."title" IS '目标标题';
COMMENT ON COLUMN "public"."requirement_goals"."description" IS '目标描述';
COMMENT ON COLUMN "public"."requirement_goals"."success_criteria" IS '成功标准/验收标准';
COMMENT ON COLUMN "public"."requirement_goals"."priority" IS '目标优先级：1-低，2-中，3-高，4-紧急';
COMMENT ON COLUMN "public"."requirement_goals"."status" IS '目标状态：pending-待审核，approved-已通过，rejected-已驳回';
COMMENT ON COLUMN "public"."requirement_goals"."sort_order" IS '排序顺序';
COMMENT ON COLUMN "public"."requirement_goals"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."requirement_goals"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."requirement_goals"."start_date" IS '目标开始时间';
COMMENT ON COLUMN "public"."requirement_goals"."end_date" IS '目标结束时间';
CREATE TRIGGER "update_requirement_goals_updated_at"
    BEFORE UPDATE
    ON "public"."requirement_goals"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.requirement_solutions Indexes
COMMENT ON TABLE "public"."requirement_solutions" IS '需求思路方案表';
CREATE INDEX "idx_requirement_solutions_task_id" ON "public"."requirement_solutions" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."requirement_solutions" ADD CONSTRAINT "requirement_solutions_submitted_by_fkey" FOREIGN KEY ("submitted_by") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."requirement_solutions"."id" IS '主键ID';
COMMENT ON COLUMN "public"."requirement_solutions"."task_id" IS '关联的任务ID';
COMMENT ON COLUMN "public"."requirement_solutions"."version" IS '方案版本号（支持多次修改）';
COMMENT ON COLUMN "public"."requirement_solutions"."content" IS '方案文字说明';
COMMENT ON COLUMN "public"."requirement_solutions"."mindmap_url" IS '脑图文件存储URL';
COMMENT ON COLUMN "public"."requirement_solutions"."file_name" IS '脑图文件名';
COMMENT ON COLUMN "public"."requirement_solutions"."file_size" IS '文件大小（字节）';
COMMENT ON COLUMN "public"."requirement_solutions"."status" IS '方案状态：pending-待审核，approved-已通过，rejected-已驳回';
COMMENT ON COLUMN "public"."requirement_solutions"."submitted_by" IS '提交人用户ID';
COMMENT ON COLUMN "public"."requirement_solutions"."submitted_at" IS '提交时间';
COMMENT ON COLUMN "public"."requirement_solutions"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."requirement_solutions"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."requirement_solutions"."title" IS '方案标题（用于在列表中快速识别）';
COMMENT ON COLUMN "public"."requirement_solutions"."mindmap_markdown" IS '脑图 Markdown 文本';
CREATE TRIGGER "update_requirement_solutions_updated_at"
    BEFORE UPDATE
    ON "public"."requirement_solutions"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.review_records Indexes
COMMENT ON TABLE "public"."review_records" IS '审核记录表（存储每个审核人的意见）';
CREATE INDEX "idx_review_records_reviewer_id" ON "public"."review_records" USING btree ("reviewer_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_review_records_session_id" ON "public"."review_records" USING btree ("review_session_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "review_records_review_session_id_reviewer_id_key" ON "public"."review_records" USING btree ("review_session_id"  "pg_catalog"."int4_ops" ASC NULLS LAST,"reviewer_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."review_records" ADD CONSTRAINT "review_records_review_session_id_fkey" FOREIGN KEY ("review_session_id") REFERENCES "public"."review_sessions" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."review_records" ADD CONSTRAINT "review_records_reviewer_id_fkey" FOREIGN KEY ("reviewer_id") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."review_records"."id" IS '主键ID';
COMMENT ON COLUMN "public"."review_records"."review_session_id" IS '关联的审核会话ID';
COMMENT ON COLUMN "public"."review_records"."reviewer_id" IS '审核人用户ID';
COMMENT ON COLUMN "public"."review_records"."reviewer_role" IS '审核人角色：creator-创建人, jury-陪审团成员, expert-专家';
COMMENT ON COLUMN "public"."review_records"."opinion" IS '审核意见：approve-同意, reject-拒绝, abstain-弃权';
COMMENT ON COLUMN "public"."review_records"."comment" IS '审核意见详细说明';
COMMENT ON COLUMN "public"."review_records"."score" IS '评分（1-5分，可选）';
COMMENT ON COLUMN "public"."review_records"."attachments" IS '附件外链信息JSON';
COMMENT ON COLUMN "public"."review_records"."vote_weight" IS '投票权重（陪审团模式，默认1.0）';
COMMENT ON COLUMN "public"."review_records"."reviewed_at" IS '审核时间';
COMMENT ON COLUMN "public"."review_records"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."review_records"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."review_records"."deleted_at" IS '删除时间';
CREATE TRIGGER "update_review_records_updated_at"
    BEFORE UPDATE
    ON "public"."review_records"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.review_sessions Indexes
COMMENT ON TABLE "public"."review_sessions" IS '审核会话表（管理整个审核流程）';
CREATE INDEX "idx_review_sessions_deleted_at" ON "public"."review_sessions" USING btree ("deleted_at"  "pg_catalog"."timestamptz_ops" ASC NULLS LAST);
CREATE INDEX "idx_review_sessions_status" ON "public"."review_sessions" USING btree ("status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_review_sessions_target" ON "public"."review_sessions" USING btree ("target_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,"target_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_review_sessions_task_id" ON "public"."review_sessions" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."review_sessions" ADD CONSTRAINT "review_sessions_initiated_by_fkey" FOREIGN KEY ("initiated_by") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."review_sessions" ADD CONSTRAINT "review_sessions_final_decision_by_fkey" FOREIGN KEY ("final_decision_by") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."review_sessions"."id" IS '主键ID';
COMMENT ON COLUMN "public"."review_sessions"."task_id" IS '关联的任务ID';
COMMENT ON COLUMN "public"."review_sessions"."review_type" IS '审核类型：solution_review-目标方案审核,  plan_review-计划审核';
COMMENT ON COLUMN "public"."review_sessions"."target_type" IS '被审核对象的表名：requirement_goals, requirement_solutions, execution_plans';
COMMENT ON COLUMN "public"."review_sessions"."target_id" IS '被审核对象的ID';
COMMENT ON COLUMN "public"."review_sessions"."initiated_by" IS '发起审核的用户ID（通常是执行人）';
COMMENT ON COLUMN "public"."review_sessions"."initiated_at" IS '发起审核时间';
COMMENT ON COLUMN "public"."review_sessions"."status" IS '审核状态：pending-待审核, in_review-审核中, approved-已通过, rejected-已驳回, cancelled-已取消';
COMMENT ON COLUMN "public"."review_sessions"."review_mode" IS '审核模式：single-单人审核（创建人）, jury-陪审团审核（多人投票）';
COMMENT ON COLUMN "public"."review_sessions"."required_approvals" IS '需要的通过票数（陪审团模式）';
COMMENT ON COLUMN "public"."review_sessions"."final_decision" IS '最终决策：approved-通过, rejected-驳回';
COMMENT ON COLUMN "public"."review_sessions"."final_decision_by" IS '最终决策人用户ID（通常是任务创建人）';
COMMENT ON COLUMN "public"."review_sessions"."final_decision_at" IS '最终决策时间';
COMMENT ON COLUMN "public"."review_sessions"."final_decision_comment" IS '最终决策说明';
COMMENT ON COLUMN "public"."review_sessions"."completed_at" IS '审核完成时间';
COMMENT ON COLUMN "public"."review_sessions"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."review_sessions"."updated_at" IS '更新时间';
CREATE TRIGGER "update_review_sessions_updated_at"
    BEFORE UPDATE
    ON "public"."review_sessions"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.role_permissions Indexes
COMMENT ON TABLE "public"."role_permissions" IS '角色权限关联表（多对多）';
CREATE INDEX "idx_role_permissions_permission_id" ON "public"."role_permissions" USING btree ("permission_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_role_permissions_role_id" ON "public"."role_permissions" USING btree ("role_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."role_permissions" ADD CONSTRAINT "role_permissions_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."role_permissions" ADD CONSTRAINT "role_permissions_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."role_permissions"."role_id" IS '角色ID';
COMMENT ON COLUMN "public"."role_permissions"."permission_id" IS '权限ID';

-- public.roles Indexes
COMMENT ON TABLE "public"."roles" IS '角色表';
CREATE INDEX "idx_roles_deleted_at" ON "public"."roles" USING btree ("deleted_at"  "pg_catalog"."timestamptz_ops" ASC NULLS LAST);
CREATE INDEX "idx_roles_name" ON "public"."roles" USING btree ("name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "roles_name_key" ON "public"."roles" USING btree ("name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
COMMENT ON COLUMN "public"."roles"."id" IS '角色ID';
COMMENT ON COLUMN "public"."roles"."name" IS '角色名称';
COMMENT ON COLUMN "public"."roles"."description" IS '角色描述';
COMMENT ON COLUMN "public"."roles"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."roles"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."roles"."deleted_at" IS '软删除时间';
CREATE TRIGGER "update_roles_updated_at"
    BEFORE UPDATE
    ON "public"."roles"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.task_attachments Indexes
COMMENT ON TABLE "public"."task_attachments" IS '任务附件表';
CREATE INDEX "idx_task_attachments_task_id" ON "public"."task_attachments" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."task_attachments" ADD CONSTRAINT "task_attachments_uploaded_by_fkey" FOREIGN KEY ("uploaded_by") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."task_attachments"."id" IS '主键ID';
COMMENT ON COLUMN "public"."task_attachments"."task_id" IS '关联的任务ID';
COMMENT ON COLUMN "public"."task_attachments"."file_name" IS '文件名';
COMMENT ON COLUMN "public"."task_attachments"."file_url" IS '文件存储URL';
COMMENT ON COLUMN "public"."task_attachments"."file_type" IS '文件类型（MIME类型）';
COMMENT ON COLUMN "public"."task_attachments"."file_size" IS '文件大小（字节）';
COMMENT ON COLUMN "public"."task_attachments"."uploaded_by" IS '上传人用户ID';
COMMENT ON COLUMN "public"."task_attachments"."attachment_type" IS '附件类型：requirement-需求相关，solution-方案相关，plan-计划相关，general-通用附件';
COMMENT ON COLUMN "public"."task_attachments"."created_at" IS '上传时间';

-- public.task_change_logs Indexes
COMMENT ON TABLE "public"."task_change_logs" IS '任务变更历史表';
CREATE INDEX "idx_task_change_logs_created_at" ON "public"."task_change_logs" USING btree ("created_at"  "pg_catalog"."timestamptz_ops" ASC NULLS LAST);
CREATE INDEX "idx_task_change_logs_task_id" ON "public"."task_change_logs" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."task_change_logs" ADD CONSTRAINT "task_change_logs_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."task_change_logs"."id" IS '主键ID';
COMMENT ON COLUMN "public"."task_change_logs"."task_id" IS '关联的任务ID';
COMMENT ON COLUMN "public"."task_change_logs"."user_id" IS '操作人用户ID';
COMMENT ON COLUMN "public"."task_change_logs"."change_type" IS '变更类型：status_change-状态变更，assign-指派变更，update-信息更新，comment-评论';
COMMENT ON COLUMN "public"."task_change_logs"."field_name" IS '变更字段名称';
COMMENT ON COLUMN "public"."task_change_logs"."old_value" IS '变更前的值';
COMMENT ON COLUMN "public"."task_change_logs"."new_value" IS '变更后的值';
COMMENT ON COLUMN "public"."task_change_logs"."comment" IS '变更说明';
COMMENT ON COLUMN "public"."task_change_logs"."created_at" IS '变更时间';

-- public.task_comments Indexes
COMMENT ON TABLE "public"."task_comments" IS '任务评论表';
CREATE INDEX "idx_task_comments_parent_comment_id" ON "public"."task_comments" USING btree ("parent_comment_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_task_comments_task_id" ON "public"."task_comments" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_task_comments_user_id" ON "public"."task_comments" USING btree ("user_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."task_comments" ADD CONSTRAINT "task_comments_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."task_comments" ADD CONSTRAINT "task_comments_parent_comment_id_fkey" FOREIGN KEY ("parent_comment_id") REFERENCES "public"."task_comments" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."task_comments"."id" IS '主键ID';
COMMENT ON COLUMN "public"."task_comments"."task_id" IS '关联的任务ID';
COMMENT ON COLUMN "public"."task_comments"."user_id" IS '评论人用户ID';
COMMENT ON COLUMN "public"."task_comments"."content" IS '评论内容';
COMMENT ON COLUMN "public"."task_comments"."parent_comment_id" IS '父评论ID（用于回复功能）';
COMMENT ON COLUMN "public"."task_comments"."attachments" IS '附件信息JSON：[{name:"文件名", url:"地址", size:123}]';
COMMENT ON COLUMN "public"."task_comments"."is_private" IS '是否为私密评论（仅部分人可见）';
COMMENT ON COLUMN "public"."task_comments"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."task_comments"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."task_comments"."deleted_at" IS '软删除时间';
CREATE TRIGGER "update_task_comments_updated_at"
    BEFORE UPDATE
    ON "public"."task_comments"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.task_milestones Indexes
COMMENT ON TABLE "public"."task_milestones" IS '任务时间节点/里程碑表';
CREATE INDEX "idx_task_milestones_task_id" ON "public"."task_milestones" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
COMMENT ON COLUMN "public"."task_milestones"."id" IS '主键ID';
COMMENT ON COLUMN "public"."task_milestones"."task_id" IS '关联的任务ID';
COMMENT ON COLUMN "public"."task_milestones"."name" IS '节点名称';
COMMENT ON COLUMN "public"."task_milestones"."description" IS '节点描述';
COMMENT ON COLUMN "public"."task_milestones"."target_date" IS '目标完成日期';
COMMENT ON COLUMN "public"."task_milestones"."actual_date" IS '实际完成日期';
COMMENT ON COLUMN "public"."task_milestones"."status" IS '节点状态：pending-待完成，completed-已完成，delayed-延期';
COMMENT ON COLUMN "public"."task_milestones"."sort_order" IS '排序顺序';
COMMENT ON COLUMN "public"."task_milestones"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."task_milestones"."updated_at" IS '更新时间';
CREATE TRIGGER "update_task_milestones_updated_at"
    BEFORE UPDATE
    ON "public"."task_milestones"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.task_participants Indexes
COMMENT ON TABLE "public"."task_participants" IS '任务参与人表';
CREATE INDEX "idx_task_participants_task_id" ON "public"."task_participants" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_task_participants_user_id" ON "public"."task_participants" USING btree ("user_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "task_participants_task_id_user_id_role_key" ON "public"."task_participants" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST,"user_id"  "pg_catalog"."int4_ops" ASC NULLS LAST,"role" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
ALTER TABLE "public"."task_participants" ADD CONSTRAINT "task_participants_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."task_participants" ADD CONSTRAINT "task_participants_invited_by_fkey" FOREIGN KEY ("invited_by") REFERENCES "public"."users" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."task_participants"."id" IS '主键ID';
COMMENT ON COLUMN "public"."task_participants"."task_id" IS '关联的任务ID';
COMMENT ON COLUMN "public"."task_participants"."user_id" IS '参与人用户ID';
COMMENT ON COLUMN "public"."task_participants"."role" IS '参与角色：creator-创建人，executor-执行人，reviewer-审核人，jury-陪审团，observer-观察者';
COMMENT ON COLUMN "public"."task_participants"."status" IS '参与状态：pending-待确认，accepted-已接受，rejected-已拒绝';
COMMENT ON COLUMN "public"."task_participants"."invited_by" IS '邀请人用户ID';
COMMENT ON COLUMN "public"."task_participants"."invited_at" IS '邀请时间';
COMMENT ON COLUMN "public"."task_participants"."response_at" IS '响应时间';
COMMENT ON COLUMN "public"."task_participants"."created_at" IS '创建时间';

-- public.task_types Indexes
COMMENT ON TABLE "public"."task_types" IS '任务类型表';
CREATE UNIQUE INDEX "task_types_code_key" ON "public"."task_types" USING btree ("code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
COMMENT ON COLUMN "public"."task_types"."id" IS '主键ID';
COMMENT ON COLUMN "public"."task_types"."code" IS '任务类型编码（requirement-需求任务, unit_task-最小单元任务）';
COMMENT ON COLUMN "public"."task_types"."name" IS '任务类型名称';
COMMENT ON COLUMN "public"."task_types"."description" IS '任务类型描述';
COMMENT ON COLUMN "public"."task_types"."created_at" IS '创建时间';

-- public.task_statuses Indexes
COMMENT ON TABLE "public"."task_statuses" IS '任务状态表';
CREATE UNIQUE INDEX "task_statuses_code_key" ON "public"."task_statuses" USING btree ("code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
ALTER TABLE "public"."task_statuses" ADD CONSTRAINT "task_statuses_task_type_code_fkey" FOREIGN KEY ("task_type_code") REFERENCES "public"."task_types" ("code")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."task_statuses"."id" IS '主键ID';
COMMENT ON COLUMN "public"."task_statuses"."code" IS '状态编码（唯一标识）';
COMMENT ON COLUMN "public"."task_statuses"."name" IS '状态名称';
COMMENT ON COLUMN "public"."task_statuses"."task_type_code" IS '所属任务类型编码';
COMMENT ON COLUMN "public"."task_statuses"."sort_order" IS '排序顺序';
COMMENT ON COLUMN "public"."task_statuses"."description" IS '状态描述';
COMMENT ON COLUMN "public"."task_statuses"."created_at" IS '创建时间';



-- public.task_status_transitions Indexes
COMMENT ON TABLE "public"."task_status_transitions" IS '任务状态转换规则表（状态机配置）';
CREATE INDEX "idx_status_transitions_from_status" ON "public"."task_status_transitions" USING btree ("from_status_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_status_transitions_task_type" ON "public"."task_status_transitions" USING btree ("task_type_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "task_status_transitions_task_type_code_from_status_code_to__key" ON "public"."task_status_transitions" USING btree ("task_type_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,"from_status_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,"to_status_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
ALTER TABLE "public"."task_status_transitions" ADD CONSTRAINT "task_status_transitions_task_type_code_fkey" FOREIGN KEY ("task_type_code") REFERENCES "public"."task_types" ("code")ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."task_status_transitions" ADD CONSTRAINT "task_status_transitions_from_status_code_fkey" FOREIGN KEY ("from_status_code") REFERENCES "public"."task_statuses" ("code")ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."task_status_transitions" ADD CONSTRAINT "task_status_transitions_to_status_code_fkey" FOREIGN KEY ("to_status_code") REFERENCES "public"."task_statuses" ("code")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."task_status_transitions"."id" IS '主键ID';
COMMENT ON COLUMN "public"."task_status_transitions"."task_type_code" IS '任务类型编码';
COMMENT ON COLUMN "public"."task_status_transitions"."from_status_code" IS '源状态编码';
COMMENT ON COLUMN "public"."task_status_transitions"."to_status_code" IS '目标状态编码';
COMMENT ON COLUMN "public"."task_status_transitions"."required_role" IS '需要的角色：creator-创建人, executor-执行人, reviewer-审核人';
COMMENT ON COLUMN "public"."task_status_transitions"."requires_approval" IS '是否需要审核批准';
COMMENT ON COLUMN "public"."task_status_transitions"."is_allowed" IS '是否允许此转换';
COMMENT ON COLUMN "public"."task_status_transitions"."description" IS '转换说明';
COMMENT ON COLUMN "public"."task_status_transitions"."created_at" IS '创建时间';


-- public.task_tag_rel Indexes
COMMENT ON TABLE "public"."task_tag_rel" IS '任务与标签关系表';
CREATE INDEX "idx_task_tag_rel_tag_id" ON "public"."task_tag_rel" USING btree ("tag_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_task_tag_rel_task_id" ON "public"."task_tag_rel" USING btree ("task_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."task_tag_rel" ADD CONSTRAINT "task_tag_rel_task_id_fkey" FOREIGN KEY ("task_id") REFERENCES "public"."tasks" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."task_tag_rel" ADD CONSTRAINT "task_tag_rel_tag_id_fkey" FOREIGN KEY ("tag_id") REFERENCES "public"."task_tags" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."task_tag_rel"."task_id" IS '关联 tasks.id';
COMMENT ON COLUMN "public"."task_tag_rel"."tag_id" IS '关联 task_tags.id';

-- public.task_tags Indexes
COMMENT ON TABLE "public"."task_tags" IS '任务标签表';
CREATE UNIQUE INDEX "task_tags_name_key" ON "public"."task_tags" USING btree ("name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
COMMENT ON COLUMN "public"."task_tags"."id" IS '主键ID';
COMMENT ON COLUMN "public"."task_tags"."name" IS '标签名称（唯一）';
COMMENT ON COLUMN "public"."task_tags"."color" IS '标签颜色（用于前端显示）';
COMMENT ON COLUMN "public"."task_tags"."description" IS '标签描述';
COMMENT ON COLUMN "public"."task_tags"."created_at" IS '创建时间';


-- public.tasks Indexes
COMMENT ON TABLE "public"."tasks" IS '任务主表（统一管理所有类型任务及其层级关系）';
CREATE INDEX "idx_tasks_creator_id" ON "public"."tasks" USING btree ("creator_id"  "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE INDEX "idx_tasks_deleted_at" ON "public"."tasks" USING btree ("deleted_at"  "pg_catalog"."timestamptz_ops" ASC NULLS LAST);
CREATE INDEX "idx_tasks_department_id" ON "public"."tasks" USING btree ("department_id"  "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE INDEX "idx_tasks_executor_id" ON "public"."tasks" USING btree ("executor_id"  "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE INDEX "idx_tasks_parent_task_id" ON "public"."tasks" USING btree ("parent_task_id"  "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE INDEX "idx_tasks_root_task_id" ON "public"."tasks" USING btree ("root_task_id"  "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE INDEX "idx_tasks_split_from_plan_id" ON "public"."tasks" USING btree ("split_from_plan_id"  "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "idx_tasks_task_no" ON "public"."tasks" USING btree ("task_no" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
ALTER TABLE "public"."tasks" ADD CONSTRAINT "fk_tasks_split_from_plan" FOREIGN KEY ("split_from_plan_id") REFERENCES "public"."execution_plans" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."tasks"."id" IS '主键ID';
COMMENT ON COLUMN "public"."tasks"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."tasks"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."tasks"."deleted_at" IS '软删除时间';
COMMENT ON COLUMN "public"."tasks"."task_no" IS '任务编号（唯一，如：REQ-2024-001, UNIT-2024-001）';
COMMENT ON COLUMN "public"."tasks"."title" IS '任务标题';
COMMENT ON COLUMN "public"."tasks"."description" IS '任务描述';
COMMENT ON COLUMN "public"."tasks"."task_type_code" IS '任务类型编码（requirement-需求任务, unit_task-单元任务）';
COMMENT ON COLUMN "public"."tasks"."status_code" IS '任务状态编码';
COMMENT ON COLUMN "public"."tasks"."creator_id" IS '创建人用户ID';
COMMENT ON COLUMN "public"."tasks"."executor_id" IS '执行人/负责人用户ID';
COMMENT ON COLUMN "public"."tasks"."department_id" IS '所属部门ID';
COMMENT ON COLUMN "public"."tasks"."parent_task_id" IS '父任务ID（用于建立父子关系）';
COMMENT ON COLUMN "public"."tasks"."root_task_id" IS '根任务ID（顶层任务的ID，方便追溯到最初的需求）';
COMMENT ON COLUMN "public"."tasks"."task_level" IS '任务层级：0-顶层任务，1-一级子任务，2-二级子任务...';
COMMENT ON COLUMN "public"."tasks"."task_path" IS '任务路径（如：1/5/12，表示任务1的子任务5的子任务12）';
COMMENT ON COLUMN "public"."tasks"."child_sequence" IS '在父任务中的序号（用于子任务排序，从1开始）';
COMMENT ON COLUMN "public"."tasks"."total_subtasks" IS '直接子任务总数（冗余字段，提升查询性能）';
COMMENT ON COLUMN "public"."tasks"."completed_subtasks" IS '已完成的直接子任务数（冗余字段）';
COMMENT ON COLUMN "public"."tasks"."expected_start_date" IS '期望开始日期';
COMMENT ON COLUMN "public"."tasks"."expected_end_date" IS '期望完成日期';
COMMENT ON COLUMN "public"."tasks"."actual_start_date" IS '实际开始日期';
COMMENT ON COLUMN "public"."tasks"."actual_end_date" IS '实际完成日期';
COMMENT ON COLUMN "public"."tasks"."priority" IS '优先级：1-低，2-中，3-高，4-紧急';
COMMENT ON COLUMN "public"."tasks"."progress" IS '任务进度百分比（0-100）';
COMMENT ON COLUMN "public"."tasks"."is_cross_department" IS '是否跨部门任务';
COMMENT ON COLUMN "public"."tasks"."is_in_pool" IS '是否在待领池中（未指派执行人）';
COMMENT ON COLUMN "public"."tasks"."is_template" IS '是否为模板任务（用于快速创建相似任务）';
COMMENT ON COLUMN "public"."tasks"."split_from_plan_id" IS '从哪个执行计划拆分出来的（关联execution_plans表）';
COMMENT ON COLUMN "public"."tasks"."split_at" IS '任务拆分时间';
COMMENT ON COLUMN "public"."tasks"."solution_deadline" IS '思路方案截止天数（需求类任务创建时可设定，表示执行人接受任务后需在N天内提交方案，0表示不限制）';
CREATE TRIGGER "update_tasks_updated_at"
    BEFORE UPDATE
    ON "public"."tasks"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- public.user_roles Indexes
COMMENT ON TABLE "public"."user_roles" IS '用户角色关联表（多对多）';
CREATE INDEX "idx_user_roles_role_id" ON "public"."user_roles" USING btree ("role_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_user_roles_user_id" ON "public"."user_roles" USING btree ("user_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
ALTER TABLE "public"."user_roles" ADD CONSTRAINT "user_roles_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."user_roles" ADD CONSTRAINT "user_roles_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id")ON DELETE CASCADE ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."user_roles"."user_id" IS '用户ID';
COMMENT ON COLUMN "public"."user_roles"."role_id" IS '角色ID';

-- public.users Indexes
COMMENT ON TABLE "public"."users" IS '用户表';
CREATE INDEX "idx_users_deleted_at" ON "public"."users" USING btree ("deleted_at"  "pg_catalog"."timestamptz_ops" ASC NULLS LAST);
CREATE INDEX "idx_users_department_id" ON "public"."users" USING btree ("department_id"  "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_users_email" ON "public"."users" USING btree ("email" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_users_username" ON "public"."users" USING btree ("username" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "users_email_key" ON "public"."users" USING btree ("email" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "users_mobile_key" ON "public"."users" USING btree ("mobile" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "users_wechat_openid_key" ON "public"."users" USING btree ("wechat_openid" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "users_wechat_unionid_key" ON "public"."users" USING btree ("wechat_unionid" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
ALTER TABLE "public"."users" ADD CONSTRAINT "users_department_id_fkey" FOREIGN KEY ("department_id") REFERENCES "public"."departments" ("id")ON DELETE NO ACTION ON UPDATE NO ACTION;
COMMENT ON COLUMN "public"."users"."id" IS '用户ID';
COMMENT ON COLUMN "public"."users"."username" IS '用户名';
COMMENT ON COLUMN "public"."users"."email" IS '邮箱';
COMMENT ON COLUMN "public"."users"."password" IS '密码（加密）';
COMMENT ON COLUMN "public"."users"."mobile" IS '手机号';
COMMENT ON COLUMN "public"."users"."status" IS '状态：1-正常，3-禁用，2-待审核';
COMMENT ON COLUMN "public"."users"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."users"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."users"."deleted_at" IS '软删除时间';
COMMENT ON COLUMN "public"."users"."is_department_leader" IS '是否为部门负责人';
COMMENT ON COLUMN "public"."users"."job_title" IS '职位名称';
COMMENT ON COLUMN "public"."users"."department_id" IS '所属部门ID';
COMMENT ON COLUMN "public"."users"."nickname" IS '昵称';
COMMENT ON COLUMN "public"."users"."wechat_unionid" IS '微信全局唯一ID（跨应用）';
COMMENT ON COLUMN "public"."users"."wechat_openid" IS '微信OpenID（单应用内唯一）';
COMMENT ON COLUMN "public"."users"."avatar" IS '用户头像URL（可存微信头像）';
CREATE TRIGGER "update_users_updated_at"
    BEFORE UPDATE
    ON "public"."users"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

SET session_replication_role = 'origin';
