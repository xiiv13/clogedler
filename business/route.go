package business

import (
	"clogedler/framework"
	"clogedler/framework/middleware"
)

// RegisterRouter 注册路由规则
func RegisterRouter(core *framework.Core) {
	// 静态路由+HTTP方法匹配
	core.Get("/user/login", middleware.Test1(), UserLoginController)

	// 批量通用前缀
	subjectApi := core.Group("/subject")
	{
		subjectApi.Use(middleware.Test1())
		// 动态路由
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectUpdateController)
		subjectApi.Get("/:id", middleware.Test1(), SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", SubjectNameController)
		}
	}
}
