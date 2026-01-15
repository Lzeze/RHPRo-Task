-- ============================================
-- 年度规划系统数据库迁移脚本
-- Annual Planning System Database Migration
-- ============================================

-- ============================================
-- 1. 部门准则表 (department_guidelines)
-- ============================================
DROP TABLE IF EXISTS "public"."department_guidelines";
CREATE SEQUENCE IF NOT EXISTS "public"."department_guidelines_id_seq";
CREATE TABLE "public"."department_guidelines" (
    "id" int4 NOT NULL DEFAULT nextval('department_guidelines_id_seq'::regclass),
    "department_id" int4 NOT NULL,
    "file_name" varchar(255) NOT NULL,
    "file_path" varchar(500) NOT NULL,
    "file_type" varchar(50) NOT NULL,
    "file_size" int8 NOT NULL,
    "version" int4 DEFAULT 1,
    "is_current" bool DEFAULT true,
    "uploaded_by" int4 NOT NULL,
    "remark" text,
    "created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz(6),
    PRIMARY KEY ("id")
);

COMMENT ON TABLE "public"."department_guidelines" IS '部门行为准则表';
COMMENT ON COLUMN "public"."department_guidelines"."id" IS '主键ID';
COMMENT ON COLUMN "public"."department_guidelines"."department_id" IS '部门ID';
COMMENT ON COLUMN "public"."department_guidelines"."file_name" IS '文件名';
COMMENT ON COLUMN "public"."department_guidelines"."file_path" IS '文件存储路径';
COMMENT ON COLUMN "public"."department_guidelines"."file_type" IS '文件类型（pdf/doc/docx/png/jpg等）';
COMMENT ON COLUMN "public"."department_guidelines"."file_size" IS '文件大小（字节）';
COMMENT ON COLUMN "public"."department_guidelines"."version" IS '版本号';
COMMENT ON COLUMN "public"."department_guidelines"."is_current" IS '是否当前生效版本';
COMMENT ON COLUMN "public"."department_guidelines"."uploaded_by" IS '上传人ID';
COMMENT ON COLUMN "public"."department_guidelines"."remark" IS '备注';
COMMENT ON COLUMN "public"."department_guidelines"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."department_guidelines"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."department_guidelines"."deleted_at" IS '软删除时间';

CREATE INDEX "idx_department_guidelines_department_id" ON "public"."department_guidelines" USING btree ("department_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_department_guidelines_is_current" ON "public"."department_guidelines" USING btree ("is_current" "pg_catalog"."bool_ops" ASC NULLS LAST);
CREATE INDEX "idx_department_guidelines_deleted_at" ON "public"."department_guidelines" USING btree ("deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST);

ALTER TABLE "public"."department_guidelines" ADD CONSTRAINT "department_guidelines_department_id_fkey" 
    FOREIGN KEY ("department_id") REFERENCES "public"."departments" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."department_guidelines" ADD CONSTRAINT "department_guidelines_uploaded_by_fkey" 
    FOREIGN KEY ("uploaded_by") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

CREATE TRIGGER "update_department_guidelines_updated_at"
    BEFORE UPDATE ON "public"."department_guidelines"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 2. 产品主线表 (product_lines)
-- ============================================
DROP TABLE IF EXISTS "public"."product_lines";
CREATE SEQUENCE IF NOT EXISTS "public"."product_lines_id_seq";
CREATE TABLE "public"."product_lines" (
    "id" int4 NOT NULL DEFAULT nextval('product_lines_id_seq'::regclass),
    "product_no" varchar(50) NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" text,
    "creator_department_id" int4 NOT NULL,
    "creator_id" int4 NOT NULL,
    "status" varchar(50) DEFAULT 'active',
    "created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz(6),
    PRIMARY KEY ("id")
);

COMMENT ON TABLE "public"."product_lines" IS '产品主线表（全局可见，跨部门）';
COMMENT ON COLUMN "public"."product_lines"."id" IS '主键ID';
COMMENT ON COLUMN "public"."product_lines"."product_no" IS '产品编号（系统自动生成，格式：PRD-2026-001）';
COMMENT ON COLUMN "public"."product_lines"."name" IS '产品名称';
COMMENT ON COLUMN "public"."product_lines"."description" IS '产品描述';
COMMENT ON COLUMN "public"."product_lines"."creator_department_id" IS '创建部门ID（首次创建的部门）';
COMMENT ON COLUMN "public"."product_lines"."creator_id" IS '创建人ID';
COMMENT ON COLUMN "public"."product_lines"."status" IS '状态：active-活跃，archived-已归档';
COMMENT ON COLUMN "public"."product_lines"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."product_lines"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."product_lines"."deleted_at" IS '软删除时间';

CREATE UNIQUE INDEX "product_lines_product_no_key" ON "public"."product_lines" USING btree ("product_no" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_product_lines_status" ON "public"."product_lines" USING btree ("status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_product_lines_creator_department_id" ON "public"."product_lines" USING btree ("creator_department_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_product_lines_deleted_at" ON "public"."product_lines" USING btree ("deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST);

ALTER TABLE "public"."product_lines" ADD CONSTRAINT "product_lines_creator_department_id_fkey" 
    FOREIGN KEY ("creator_department_id") REFERENCES "public"."departments" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."product_lines" ADD CONSTRAINT "product_lines_creator_id_fkey" 
    FOREIGN KEY ("creator_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

CREATE TRIGGER "update_product_lines_updated_at"
    BEFORE UPDATE ON "public"."product_lines"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 3. 年度计划表 (annual_plans)
-- ============================================
DROP TABLE IF EXISTS "public"."annual_plans";
CREATE SEQUENCE IF NOT EXISTS "public"."annual_plans_id_seq";
CREATE TABLE "public"."annual_plans" (
    "id" int4 NOT NULL DEFAULT nextval('annual_plans_id_seq'::regclass),
    "plan_no" varchar(50) NOT NULL,
    "name" varchar(255) NOT NULL,
    "year" int4 NOT NULL,
    "department_id" int4 NOT NULL,
    "description" text,
    "status" varchar(50) DEFAULT 'draft',
    "creator_id" int4 NOT NULL,
    "published_at" timestamptz(6),
    "archived_at" timestamptz(6),
    "created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz(6),
    PRIMARY KEY ("id")
);

COMMENT ON TABLE "public"."annual_plans" IS '年度计划表';
COMMENT ON COLUMN "public"."annual_plans"."id" IS '主键ID';
COMMENT ON COLUMN "public"."annual_plans"."plan_no" IS '计划编号（系统自动生成，格式：AP-2026-001）';
COMMENT ON COLUMN "public"."annual_plans"."name" IS '计划名称';
COMMENT ON COLUMN "public"."annual_plans"."year" IS '年份';
COMMENT ON COLUMN "public"."annual_plans"."department_id" IS '部门ID';
COMMENT ON COLUMN "public"."annual_plans"."description" IS '计划描述';
COMMENT ON COLUMN "public"."annual_plans"."status" IS '状态：draft-草稿，active-进行中，archived-已归档';
COMMENT ON COLUMN "public"."annual_plans"."creator_id" IS '创建人ID';
COMMENT ON COLUMN "public"."annual_plans"."published_at" IS '发布时间';
COMMENT ON COLUMN "public"."annual_plans"."archived_at" IS '归档时间';
COMMENT ON COLUMN "public"."annual_plans"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."annual_plans"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."annual_plans"."deleted_at" IS '软删除时间';

CREATE UNIQUE INDEX "annual_plans_plan_no_key" ON "public"."annual_plans" USING btree ("plan_no" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE UNIQUE INDEX "annual_plans_department_year_key" ON "public"."annual_plans" USING btree ("department_id" "pg_catalog"."int4_ops" ASC NULLS LAST, "year" "pg_catalog"."int4_ops" ASC NULLS LAST) WHERE deleted_at IS NULL;
CREATE INDEX "idx_annual_plans_year" ON "public"."annual_plans" USING btree ("year" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_annual_plans_department_id" ON "public"."annual_plans" USING btree ("department_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_annual_plans_status" ON "public"."annual_plans" USING btree ("status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_annual_plans_deleted_at" ON "public"."annual_plans" USING btree ("deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST);

ALTER TABLE "public"."annual_plans" ADD CONSTRAINT "annual_plans_department_id_fkey" 
    FOREIGN KEY ("department_id") REFERENCES "public"."departments" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."annual_plans" ADD CONSTRAINT "annual_plans_creator_id_fkey" 
    FOREIGN KEY ("creator_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

CREATE TRIGGER "update_annual_plans_updated_at"
    BEFORE UPDATE ON "public"."annual_plans"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- ============================================
-- 4. 计划节点表 (plan_nodes)
-- ============================================
DROP TABLE IF EXISTS "public"."plan_nodes";
CREATE SEQUENCE IF NOT EXISTS "public"."plan_nodes_id_seq";
CREATE TABLE "public"."plan_nodes" (
    "id" int4 NOT NULL DEFAULT nextval('plan_nodes_id_seq'::regclass),
    "node_no" varchar(50) NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" text,
    "annual_plan_id" int4 NOT NULL,
    "product_line_id" int4 NOT NULL,
    "stage" varchar(50) NOT NULL,
    "parent_node_id" int4,
    "root_node_id" int4,
    "node_level" int4 DEFAULT 0,
    "node_path" varchar(500),
    "sort_order" int4 DEFAULT 0,
    "owner_id" int4,
    "expected_start_date" timestamptz(6),
    "expected_end_date" timestamptz(6),
    "actual_start_date" timestamptz(6),
    "actual_end_date" timestamptz(6),
    "status" varchar(50) DEFAULT 'pending',
    "total_tasks" int4 DEFAULT 0,
    "completed_tasks" int4 DEFAULT 0,
    "creator_id" int4 NOT NULL,
    "created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz(6),
    PRIMARY KEY ("id")
);

COMMENT ON TABLE "public"."plan_nodes" IS '计划节点表（年度计划拆分的具体工作项）';
COMMENT ON COLUMN "public"."plan_nodes"."id" IS '主键ID';
COMMENT ON COLUMN "public"."plan_nodes"."node_no" IS '节点编号（系统自动生成，格式：PN-2026-001）';
COMMENT ON COLUMN "public"."plan_nodes"."name" IS '节点名称';
COMMENT ON COLUMN "public"."plan_nodes"."description" IS '节点描述';
COMMENT ON COLUMN "public"."plan_nodes"."annual_plan_id" IS '所属年度计划ID';
COMMENT ON COLUMN "public"."plan_nodes"."product_line_id" IS '所属产品主线ID';
COMMENT ON COLUMN "public"."plan_nodes"."stage" IS '计划阶段：germination-萌芽期，experiment-试验期，maturity-成熟期，promotion-推广期';
COMMENT ON COLUMN "public"."plan_nodes"."parent_node_id" IS '父节点ID（支持多级嵌套）';
COMMENT ON COLUMN "public"."plan_nodes"."root_node_id" IS '根节点ID';
COMMENT ON COLUMN "public"."plan_nodes"."node_level" IS '节点层级（根节点为0）';
COMMENT ON COLUMN "public"."plan_nodes"."node_path" IS '节点路径（如：1/2/3）';
COMMENT ON COLUMN "public"."plan_nodes"."sort_order" IS '同级排序序号';
COMMENT ON COLUMN "public"."plan_nodes"."owner_id" IS '负责人ID';
COMMENT ON COLUMN "public"."plan_nodes"."expected_start_date" IS '期望开始日期';
COMMENT ON COLUMN "public"."plan_nodes"."expected_end_date" IS '期望结束日期';
COMMENT ON COLUMN "public"."plan_nodes"."actual_start_date" IS '实际开始日期';
COMMENT ON COLUMN "public"."plan_nodes"."actual_end_date" IS '实际结束日期';
COMMENT ON COLUMN "public"."plan_nodes"."status" IS '状态：pending-待开始，in_progress-进行中，completed-已完成，cancelled-已取消';
COMMENT ON COLUMN "public"."plan_nodes"."total_tasks" IS '任务总数（冗余字段，优化查询）';
COMMENT ON COLUMN "public"."plan_nodes"."completed_tasks" IS '已完成任务数';
COMMENT ON COLUMN "public"."plan_nodes"."creator_id" IS '创建人ID';
COMMENT ON COLUMN "public"."plan_nodes"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."plan_nodes"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."plan_nodes"."deleted_at" IS '软删除时间';

CREATE UNIQUE INDEX "plan_nodes_node_no_key" ON "public"."plan_nodes" USING btree ("node_no" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_plan_nodes_annual_plan_id" ON "public"."plan_nodes" USING btree ("annual_plan_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_plan_nodes_product_line_id" ON "public"."plan_nodes" USING btree ("product_line_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_plan_nodes_stage" ON "public"."plan_nodes" USING btree ("stage" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_plan_nodes_parent_node_id" ON "public"."plan_nodes" USING btree ("parent_node_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_plan_nodes_root_node_id" ON "public"."plan_nodes" USING btree ("root_node_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_plan_nodes_status" ON "public"."plan_nodes" USING btree ("status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_plan_nodes_deleted_at" ON "public"."plan_nodes" USING btree ("deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST);
-- 跨部门协作聚合索引：通过 product_line_id + stage 查询同产品同阶段的所有部门节点
CREATE INDEX "idx_plan_nodes_product_stage" ON "public"."plan_nodes" USING btree ("product_line_id" "pg_catalog"."int4_ops" ASC NULLS LAST, "stage" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);

ALTER TABLE "public"."plan_nodes" ADD CONSTRAINT "plan_nodes_annual_plan_id_fkey" 
    FOREIGN KEY ("annual_plan_id") REFERENCES "public"."annual_plans" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."plan_nodes" ADD CONSTRAINT "plan_nodes_product_line_id_fkey" 
    FOREIGN KEY ("product_line_id") REFERENCES "public"."product_lines" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."plan_nodes" ADD CONSTRAINT "plan_nodes_parent_node_id_fkey" 
    FOREIGN KEY ("parent_node_id") REFERENCES "public"."plan_nodes" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."plan_nodes" ADD CONSTRAINT "plan_nodes_owner_id_fkey" 
    FOREIGN KEY ("owner_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."plan_nodes" ADD CONSTRAINT "plan_nodes_creator_id_fkey" 
    FOREIGN KEY ("creator_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
-- 阶段值约束
ALTER TABLE "public"."plan_nodes" ADD CONSTRAINT "plan_nodes_stage_check" 
    CHECK (stage IN ('germination', 'experiment', 'maturity', 'promotion'));

CREATE TRIGGER "update_plan_nodes_updated_at"
    BEFORE UPDATE ON "public"."plan_nodes"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 5. 计划目标表 (plan_goals)
-- ============================================
DROP TABLE IF EXISTS "public"."plan_goals";
CREATE SEQUENCE IF NOT EXISTS "public"."plan_goals_id_seq";
CREATE TABLE "public"."plan_goals" (
    "id" int4 NOT NULL DEFAULT nextval('plan_goals_id_seq'::regclass),
    "plan_node_id" int4 NOT NULL,
    "goal_no" int4 NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" text,
    "acceptance_criteria" text,
    "status" varchar(50) DEFAULT 'pending',
    "completed_at" timestamptz(6),
    "completed_by" int4,
    "sort_order" int4 DEFAULT 0,
    "created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz(6),
    PRIMARY KEY ("id")
);

COMMENT ON TABLE "public"."plan_goals" IS '计划目标表（计划节点的具体目标）';
COMMENT ON COLUMN "public"."plan_goals"."id" IS '主键ID';
COMMENT ON COLUMN "public"."plan_goals"."plan_node_id" IS '所属计划节点ID';
COMMENT ON COLUMN "public"."plan_goals"."goal_no" IS '目标编号（节点内序号）';
COMMENT ON COLUMN "public"."plan_goals"."name" IS '目标名称';
COMMENT ON COLUMN "public"."plan_goals"."description" IS '目标描述';
COMMENT ON COLUMN "public"."plan_goals"."acceptance_criteria" IS '验收标准';
COMMENT ON COLUMN "public"."plan_goals"."status" IS '完成状态：pending-待完成，completed-已完成';
COMMENT ON COLUMN "public"."plan_goals"."completed_at" IS '完成时间';
COMMENT ON COLUMN "public"."plan_goals"."completed_by" IS '完成人ID';
COMMENT ON COLUMN "public"."plan_goals"."sort_order" IS '排序序号';
COMMENT ON COLUMN "public"."plan_goals"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."plan_goals"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."plan_goals"."deleted_at" IS '软删除时间';

CREATE UNIQUE INDEX "plan_goals_node_goal_no_key" ON "public"."plan_goals" USING btree ("plan_node_id" "pg_catalog"."int4_ops" ASC NULLS LAST, "goal_no" "pg_catalog"."int4_ops" ASC NULLS LAST) WHERE deleted_at IS NULL;
CREATE INDEX "idx_plan_goals_plan_node_id" ON "public"."plan_goals" USING btree ("plan_node_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_plan_goals_status" ON "public"."plan_goals" USING btree ("status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "idx_plan_goals_deleted_at" ON "public"."plan_goals" USING btree ("deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST);

ALTER TABLE "public"."plan_goals" ADD CONSTRAINT "plan_goals_plan_node_id_fkey" 
    FOREIGN KEY ("plan_node_id") REFERENCES "public"."plan_nodes" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."plan_goals" ADD CONSTRAINT "plan_goals_completed_by_fkey" 
    FOREIGN KEY ("completed_by") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

CREATE TRIGGER "update_plan_goals_updated_at"
    BEFORE UPDATE ON "public"."plan_goals"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 6. 节点关联表 (node_links) - 仅用于阶段递进
-- ============================================
DROP TABLE IF EXISTS "public"."node_links";
CREATE SEQUENCE IF NOT EXISTS "public"."node_links_id_seq";
CREATE TABLE "public"."node_links" (
    "id" int4 NOT NULL DEFAULT nextval('node_links_id_seq'::regclass),
    "source_node_id" int4 NOT NULL,
    "target_node_id" int4 NOT NULL,
    "link_type" varchar(50) NOT NULL DEFAULT 'stage_progression',
    "creator_id" int4 NOT NULL,
    "created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz(6),
    PRIMARY KEY ("id")
);

COMMENT ON TABLE "public"."node_links" IS '节点关联表（仅用于阶段递进关联）';
COMMENT ON COLUMN "public"."node_links"."id" IS '主键ID';
COMMENT ON COLUMN "public"."node_links"."source_node_id" IS '源节点ID（前一阶段）';
COMMENT ON COLUMN "public"."node_links"."target_node_id" IS '目标节点ID（后一阶段）';
COMMENT ON COLUMN "public"."node_links"."link_type" IS '关联类型：stage_progression-阶段递进';
COMMENT ON COLUMN "public"."node_links"."creator_id" IS '创建人ID';
COMMENT ON COLUMN "public"."node_links"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."node_links"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."node_links"."deleted_at" IS '软删除时间';

CREATE UNIQUE INDEX "node_links_source_target_key" ON "public"."node_links" USING btree ("source_node_id" "pg_catalog"."int4_ops" ASC NULLS LAST, "target_node_id" "pg_catalog"."int4_ops" ASC NULLS LAST) WHERE deleted_at IS NULL;
CREATE INDEX "idx_node_links_source_node_id" ON "public"."node_links" USING btree ("source_node_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_node_links_target_node_id" ON "public"."node_links" USING btree ("target_node_id" "pg_catalog"."int4_ops" ASC NULLS LAST);
CREATE INDEX "idx_node_links_deleted_at" ON "public"."node_links" USING btree ("deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST);

ALTER TABLE "public"."node_links" ADD CONSTRAINT "node_links_source_node_id_fkey" 
    FOREIGN KEY ("source_node_id") REFERENCES "public"."plan_nodes" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."node_links" ADD CONSTRAINT "node_links_target_node_id_fkey" 
    FOREIGN KEY ("target_node_id") REFERENCES "public"."plan_nodes" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."node_links" ADD CONSTRAINT "node_links_creator_id_fkey" 
    FOREIGN KEY ("creator_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
-- 关联类型约束
ALTER TABLE "public"."node_links" ADD CONSTRAINT "node_links_link_type_check" 
    CHECK (link_type = 'stage_progression');

CREATE TRIGGER "update_node_links_updated_at"
    BEFORE UPDATE ON "public"."node_links"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 7. 扩展现有 tasks 表 - 添加计划节点绑定字段
-- ============================================
ALTER TABLE "public"."tasks" ADD COLUMN IF NOT EXISTS "plan_node_id" int4;
ALTER TABLE "public"."tasks" ADD COLUMN IF NOT EXISTS "bound_at" timestamptz(6);
ALTER TABLE "public"."tasks" ADD COLUMN IF NOT EXISTS "bound_by" int4;

COMMENT ON COLUMN "public"."tasks"."plan_node_id" IS '绑定的计划节点ID';
COMMENT ON COLUMN "public"."tasks"."bound_at" IS '绑定时间';
COMMENT ON COLUMN "public"."tasks"."bound_by" IS '绑定人ID';

CREATE INDEX IF NOT EXISTS "idx_tasks_plan_node_id" ON "public"."tasks" USING btree ("plan_node_id" "pg_catalog"."int4_ops" ASC NULLS LAST);

ALTER TABLE "public"."tasks" ADD CONSTRAINT "tasks_plan_node_id_fkey" 
    FOREIGN KEY ("plan_node_id") REFERENCES "public"."plan_nodes" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;

-- ============================================
-- 8. 视图：产品主线各阶段进度汇总
-- ============================================
DROP VIEW IF EXISTS "public"."v_product_line_stage_progress";
CREATE VIEW "public"."v_product_line_stage_progress" AS
SELECT 
    pl.id AS product_line_id,
    pl.product_no,
    pl.name AS product_name,
    pn.stage,
    d.id AS department_id,
    d.name AS department_name,
    COUNT(pn.id) AS node_count,
    SUM(pn.total_tasks) AS total_tasks,
    SUM(pn.completed_tasks) AS completed_tasks,
    CASE 
        WHEN SUM(pn.total_tasks) > 0 
        THEN ROUND((SUM(pn.completed_tasks)::numeric / SUM(pn.total_tasks)::numeric) * 100, 2)
        ELSE 0 
    END AS completion_rate
FROM product_lines pl
LEFT JOIN plan_nodes pn ON pl.id = pn.product_line_id AND pn.deleted_at IS NULL
LEFT JOIN annual_plans ap ON pn.annual_plan_id = ap.id AND ap.deleted_at IS NULL
LEFT JOIN departments d ON ap.department_id = d.id AND d.deleted_at IS NULL
WHERE pl.deleted_at IS NULL
GROUP BY pl.id, pl.product_no, pl.name, pn.stage, d.id, d.name
ORDER BY pl.id, 
    CASE pn.stage 
        WHEN 'germination' THEN 1 
        WHEN 'experiment' THEN 2 
        WHEN 'maturity' THEN 3 
        WHEN 'promotion' THEN 4 
    END,
    d.id;

COMMENT ON VIEW "public"."v_product_line_stage_progress" IS '产品主线各阶段各部门进度汇总视图';

-- ============================================
-- 9. 视图：年度计划统计汇总
-- ============================================
DROP VIEW IF EXISTS "public"."v_annual_plan_statistics";
CREATE VIEW "public"."v_annual_plan_statistics" AS
SELECT 
    ap.id AS annual_plan_id,
    ap.plan_no,
    ap.name AS plan_name,
    ap.year,
    ap.department_id,
    d.name AS department_name,
    ap.status,
    COUNT(DISTINCT pn.id) AS node_count,
    COALESCE(SUM(pn.total_tasks), 0) AS total_tasks,
    COALESCE(SUM(pn.completed_tasks), 0) AS completed_tasks,
    CASE 
        WHEN COALESCE(SUM(pn.total_tasks), 0) > 0 
        THEN ROUND((COALESCE(SUM(pn.completed_tasks), 0)::numeric / SUM(pn.total_tasks)::numeric) * 100, 2)
        ELSE 0 
    END AS completion_rate,
    COUNT(DISTINCT pn.id) FILTER (WHERE pn.status = 'in_progress') AS in_progress_nodes,
    COUNT(DISTINCT pn.id) FILTER (WHERE pn.status = 'completed') AS completed_nodes
FROM annual_plans ap
LEFT JOIN departments d ON ap.department_id = d.id AND d.deleted_at IS NULL
LEFT JOIN plan_nodes pn ON ap.id = pn.annual_plan_id AND pn.deleted_at IS NULL
WHERE ap.deleted_at IS NULL
GROUP BY ap.id, ap.plan_no, ap.name, ap.year, ap.department_id, d.name, ap.status
ORDER BY ap.year DESC, ap.department_id;

COMMENT ON VIEW "public"."v_annual_plan_statistics" IS '年度计划统计汇总视图';

-- ============================================
-- 10. 视图：跨部门协作节点（同产品同阶段）
-- ============================================
DROP VIEW IF EXISTS "public"."v_cross_department_nodes";
CREATE VIEW "public"."v_cross_department_nodes" AS
SELECT 
    pn.product_line_id,
    pl.product_no,
    pl.name AS product_name,
    pn.stage,
    pn.id AS node_id,
    pn.node_no,
    pn.name AS node_name,
    ap.department_id,
    d.name AS department_name,
    pn.status,
    pn.total_tasks,
    pn.completed_tasks,
    pn.expected_start_date,
    pn.expected_end_date
FROM plan_nodes pn
JOIN product_lines pl ON pn.product_line_id = pl.id AND pl.deleted_at IS NULL
JOIN annual_plans ap ON pn.annual_plan_id = ap.id AND ap.deleted_at IS NULL
JOIN departments d ON ap.department_id = d.id AND d.deleted_at IS NULL
WHERE pn.deleted_at IS NULL
ORDER BY pn.product_line_id, 
    CASE pn.stage 
        WHEN 'germination' THEN 1 
        WHEN 'experiment' THEN 2 
        WHEN 'maturity' THEN 3 
        WHEN 'promotion' THEN 4 
    END,
    ap.department_id;

COMMENT ON VIEW "public"."v_cross_department_nodes" IS '跨部门协作节点视图（按产品主线+阶段聚合）';

-- ============================================
-- 11. 初始化权限数据
-- ============================================
INSERT INTO permissions (id, name, description, created_at, updated_at) VALUES
-- 年度计划权限
(21, 'annual_plan:read', '查看年度计划', NOW(), NOW()),
(22, 'annual_plan:create', '创建年度计划', NOW(), NOW()),
(23, 'annual_plan:update', '更新年度计划', NOW(), NOW()),
(24, 'annual_plan:delete', '删除年度计划', NOW(), NOW()),
(25, 'annual_plan:publish', '发布年度计划', NOW(), NOW()),
(26, 'annual_plan:archive', '归档年度计划', NOW(), NOW()),
-- 产品主线权限
(27, 'product_line:read', '查看产品主线', NOW(), NOW()),
(28, 'product_line:create', '创建产品主线', NOW(), NOW()),
(29, 'product_line:update', '更新产品主线', NOW(), NOW()),
(30, 'product_line:delete', '删除产品主线', NOW(), NOW()),
-- 计划节点权限
(31, 'plan_node:read', '查看计划节点', NOW(), NOW()),
(32, 'plan_node:create', '创建计划节点', NOW(), NOW()),
(33, 'plan_node:update', '更新计划节点', NOW(), NOW()),
(34, 'plan_node:delete', '删除计划节点', NOW(), NOW()),
-- 部门准则权限
(35, 'guideline:read', '查看部门准则', NOW(), NOW()),
(36, 'guideline:upload', '上传部门准则', NOW(), NOW()),
(37, 'guideline:delete', '删除部门准则', NOW(), NOW()),
-- 统计权限
(38, 'statistics:read', '查看统计数据', NOW(), NOW()),
(39, 'statistics:export', '导出统计数据', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 重置序列
SELECT setval('permissions_id_seq', (SELECT MAX(id) FROM permissions) + 1, false);

-- 管理员拥有所有新权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions WHERE id >= 21 AND id <= 39
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 部门经理权限
INSERT INTO role_permissions (role_id, permission_id) VALUES
(2, 21), -- annual_plan:read
(2, 22), -- annual_plan:create
(2, 23), -- annual_plan:update
(2, 25), -- annual_plan:publish
(2, 26), -- annual_plan:archive
(2, 27), -- product_line:read
(2, 28), -- product_line:create
(2, 29), -- product_line:update
(2, 31), -- plan_node:read
(2, 32), -- plan_node:create
(2, 33), -- plan_node:update
(2, 34), -- plan_node:delete
(2, 35), -- guideline:read
(2, 36), -- guideline:upload
(2, 38)  -- statistics:read
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 普通用户权限
INSERT INTO role_permissions (role_id, permission_id) VALUES
(3, 21), -- annual_plan:read
(3, 27), -- product_line:read
(3, 31), -- plan_node:read
(3, 35), -- guideline:read
(3, 38)  -- statistics:read
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ============================================
-- 迁移完成
-- ============================================
-- 新增表：
-- 1. department_guidelines - 部门行为准则表
-- 2. product_lines - 产品主线表
-- 3. annual_plans - 年度计划表
-- 4. plan_nodes - 计划节点表
-- 5. plan_goals - 计划目标表
-- 6. node_links - 节点关联表（阶段递进）
-- 
-- 扩展表：
-- - tasks 表新增 plan_node_id, bound_at, bound_by 字段
--
-- 新增视图：
-- - v_product_line_stage_progress - 产品主线各阶段进度汇总
-- - v_annual_plan_statistics - 年度计划统计汇总
-- - v_cross_department_nodes - 跨部门协作节点
--
-- 编号生成逻辑（在 Service 层实现）：
-- - 产品编号：PRD-{年份}-{序号}，如 PRD-2026-001
-- - 年度计划编号：AP-{年份}-{序号}，如 AP-2026-001
-- - 计划节点编号：PN-{年份}-{序号}，如 PN-2026-001
-- ============================================
