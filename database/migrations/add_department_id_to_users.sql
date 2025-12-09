-- 添加 department_id 到 users 表
ALTER TABLE "public"."users" ADD COLUMN "department_id" int4;
CREATE INDEX "idx_users_department_id" ON "public"."users" ("department_id");
ALTER TABLE "public"."users" ADD CONSTRAINT "users_department_id_fkey" FOREIGN KEY ("department_id") REFERENCES "public"."departments" ("id");
COMMENT ON COLUMN "public"."users"."department_id" IS '所属部门ID';
