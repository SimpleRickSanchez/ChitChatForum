package router

import (
	"controller"
	"middleware"
	"model"

	"github.com/gin-gonic/gin"
)

func UserRouterInit(router *gin.Engine) {
	sessionNames := []string{model.UserSession}
	userRouter := router.Group("/user",
		middleware.Timer,
		model.GetSessionFunc(sessionNames),
		middleware.CheckSessionAuth)
	{
		userController := controller.UserController{}
		userRouter.GET("/login", userController.Login)
		userRouter.POST("/check", userController.Check)
		userRouter.POST("/checkexists", userController.CheckExists)
		userRouter.POST("/authenticate", userController.Auth)
		userRouter.GET("/logout", userController.Logout)
		userRouter.GET("/signup", userController.SignUp)
		userRouter.POST("/signup_account", userController.DoSignUp)
		userRouter.POST("/salt", userController.Salt)
		// userRouter.GET("/test", userController.Test)
		// userRouter.POST("/test", userController.PTest)
		userRouter.GET("/forge", userController.ForgeUser)
		// userRouter.GET("/killsleeper", userController.Killsleeper)
	}

}
