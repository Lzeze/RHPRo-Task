package routes

import (
	"RHPRo-Task/controllers"
	"RHPRo-Task/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "RHPRo-Task/docs" // Swagger 文档
)

func SetupRoutes() *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(middlewares.CORSMiddleware()) // CORS 中间件 - 必须在最前面
	router.Use(middlewares.RecoveryMiddleware())
	router.Use(middlewares.LoggerMiddleware())
	router.Use(middlewares.ErrorHandlerMiddleware())
	router.Use(middlewares.ValidatorMiddleware())

	// 初始化控制器
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	adminController := controllers.NewAdminController()
	taskController := controllers.NewTaskController()
	detailController := controllers.NewTaskDetailController()
	deptController := controllers.NewDepartmentController()

	// 公开路由
	public := router.Group("/api/v1")
	{
		// 认证相关
		public.POST("/auth/register", authController.Register)
		public.POST("/auth/login", authController.Login)

		// 健康检查
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Swagger API 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 需要认证的路由
	auth := router.Group("/api/v1")
	auth.Use(middlewares.AuthMiddleware())
	{
		// 用户信息
		auth.GET("/profile", userController.GetProfile)
	}
	// 用户管理路由（需要user:read权限）
	userRoutes := router.Group("/api/v1/users")
	userRoutes.Use(middlewares.AuthMiddleware())
	{
		// 查看用户列表和详情（需要user:read权限）
		userRoutes.GET("",
			// middlewares.PermissionMiddleware("user:read"),
			userController.GetUserList)
		userRoutes.GET("/:id",
			// middlewares.PermissionMiddleware("user:read"),
			userController.GetUserByID)

		// 更新用户（需要user:update权限）
		userRoutes.PUT("/:id",
			// middlewares.PermissionMiddleware("user:update"),
			userController.UpdateUser)

		// 分配角色（需要role:manage权限）
		userRoutes.POST("/:id/roles",
			// middlewares.PermissionMiddleware("role:manage"),
			userController.AssignRoles)

		// 管理员创建用户（需要user:create权限）
		userRoutes.POST("",
			// middlewares.PermissionMiddleware("user:create"),
			userController.CreateUser)

		// 审核用户（需要user:approve权限）
		userRoutes.POST("/:id/approve",
			// middlewares.PermissionMiddleware("user:approve"),
			userController.ApproveUser)

		// 获取可指派的执行人列表（用于任务分配，只需要当前用户有效即可）
		userRoutes.GET("/assignable",
			userController.GetAssignableUsers)
	}

	// 部门管理路由
	deptRoutes := router.Group("/api/v1/departments")
	deptRoutes.Use(middlewares.AuthMiddleware())
	{
		// deptRoutes.POST("", middlewares.PermissionMiddleware("dept:create"), deptController.CreateDepartment)
		// deptRoutes.PUT("/:id", middlewares.PermissionMiddleware("dept:update"), deptController.UpdateDepartment)
		// deptRoutes.DELETE("/:id", middlewares.PermissionMiddleware("dept:delete"), deptController.DeleteDepartment)

		// 获取当前用户负责的部门列表（必须放在 /:id 之前）
		deptRoutes.GET("/my-departments", deptController.GetUserDepartments)

		// 设置默认部门（必须放在 /:id 之前）
		deptRoutes.PUT("/default", deptController.SetDefaultDepartment)

		//创建部门
		deptRoutes.POST("", deptController.CreateDepartment)
		//更新部门
		deptRoutes.PUT("/:id", deptController.UpdateDepartment)
		//删除部门
		deptRoutes.DELETE("/:id", deptController.DeleteDepartment)
		// 获取部门列表和详情
		deptRoutes.GET("", deptController.GetDepartmentList)
		deptRoutes.GET("/:id", deptController.GetDepartmentDetail)

		// 负责人管理
		// deptRoutes.POST("/:id/leaders", middlewares.PermissionMiddleware("dept:manage"), deptController.AddLeader)
		// deptRoutes.DELETE("/:id/leaders/:userId", middlewares.PermissionMiddleware("dept:manage"), deptController.RemoveLeader)
		//添加负责人
		deptRoutes.POST("/:id/leaders", deptController.AddLeader)
		//移除负责人
		deptRoutes.DELETE("/:id/leaders/:userId", deptController.RemoveLeader)

		// 人员分配
		// deptRoutes.POST("/:id/users", middlewares.PermissionMiddleware("dept:manage"), deptController.AssignUsers)
		deptRoutes.POST("/:id/users", deptController.AssignUsers)
	}

	// 任务管理路由
	taskRoutes := router.Group("/api/v1/tasks")
	taskRoutes.Use(middlewares.AuthMiddleware())
	{
		// 任务 CRUD
		taskRoutes.POST("", taskController.CreateTask)
		//所有任务列表
		taskRoutes.GET("", taskController.GetTaskList)
		// 我的任务列表（必须放在 /:id 之前，避免路径匹配冲突）
		taskRoutes.GET("/my", taskController.GetMyTasks)
		// 任务详情（包含最新版本的方案和计划）
		taskRoutes.GET("/:id", detailController.GetTaskDetail)
		// 更新任务
		taskRoutes.PUT("/:id", taskController.UpdateTask)
		//删除任务
		taskRoutes.DELETE("/:id", taskController.DeleteTask)

		// 任务状态转换
		taskRoutes.POST("/:id/transit", taskController.TransitStatus)

		// 任务分配
		taskRoutes.POST("/:id/assign", taskController.AssignExecutor)

		// 任务专用访问 (Task Token)
		taskRoutes.GET("/current", middlewares.TaskTokenMiddleware(), taskController.GetTaskInfo)

		// 任务详情相关接口
		// 获取任务的所有方案版本
		taskRoutes.GET("/:id/solutions", detailController.GetTaskSolutions)
		// 获取任务的所有执行计划版本
		taskRoutes.GET("/:id/execution-plans", detailController.GetTaskExecutionPlans)
		// 获取任务的审核历史
		taskRoutes.GET("/:id/reviews", detailController.GetTaskReviewHistory)
		// 获取任务的变更日志
		taskRoutes.GET("/:id/change-logs", detailController.GetTaskChangeLogs)
		// 获取任务的时间轴
		taskRoutes.GET("/:id/timeline", detailController.GetTaskTimeline)
	}

	// 任务流程路由
	flowController := controllers.NewTaskFlowController()
	flowRoutes := router.Group("/api/v1/tasks")
	// flowRoutes.Use(middlewares.AuthMiddleware())
	{
		//认领任务、接受任务、
		flowRoutes.POST("/:id/accept", flowController.AcceptTask)
		//拒绝任务
		flowRoutes.POST("/:id/reject", flowController.RejectTask)
		//提交目标 (待用,目前目标和执行计划合并提交)
		flowRoutes.POST("/:id/goals", flowController.SubmitGoals)
		//发起审核(待用)
		flowRoutes.POST("/:id/review", flowController.InitiateReview)
		// 提交解决方案（第一步：方案审核）
		flowRoutes.POST("/:id/solution", flowController.SubmitSolution)
		// 提交执行计划+目标（第二步：计划审核）
		flowRoutes.POST("/:id/execution-plan", flowController.SubmitExecutionPlanWithGoals)
	}

	// 审核会话路由
	reviewRoutes := router.Group("/api/v1/review-sessions")
	// reviewRoutes.Use(middlewares.AuthMiddleware())
	{
		// 获取审核会话详情
		reviewRoutes.GET("/:sessionId", flowController.GetReviewSession)
		// 提交审核意见
		reviewRoutes.POST("/:sessionId/opinion", flowController.SubmitReviewOpinion)
		// 最终决策
		reviewRoutes.POST("/:sessionId/finalize", flowController.FinalizeReview)
		// 邀请陪审团成员
		reviewRoutes.POST("/:sessionId/invite-jury", flowController.InviteJury)
		// 移除陪审团成员
		reviewRoutes.DELETE("/:sessionId/jury/:juryMemberId", flowController.RemoveJuryMember)
	}

	// 管理员路由（需要permission:manage权限）
	adminRoutes := router.Group("/api/v1/admin")
	adminRoutes.Use(middlewares.AuthMiddleware())
	adminRoutes.Use(middlewares.PermissionMiddleware("permission:manage"))
	{
		adminRoutes.GET("/roles", adminController.GetRoleList)
		adminRoutes.GET("/permissions", adminController.GetPermissionList)
	}

	return router
}
