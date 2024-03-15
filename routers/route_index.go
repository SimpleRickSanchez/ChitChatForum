package router

import (
	"controller"
	"middleware"
	"model"

	"github.com/gin-gonic/gin"
)

func IdexRouterInit(router *gin.Engine) {
	sessionNames := []string{model.UserSession}
	indexRouter := router.Group("/",
		middleware.Timer,
		model.GetSessionFunc(sessionNames),
		middleware.CheckSessionAuth)
	{
		indexController := controller.IndexControllers{}
		indexRouter.GET("", indexController.Index)
		indexRouter.GET("thread", indexController.Thread)
		indexRouter.POST("create", indexController.Create)
	}

}
