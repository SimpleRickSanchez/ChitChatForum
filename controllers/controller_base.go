package controller

import (
	"model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// func (con BaseController) success(c *gin.Context) {
// 	c.String(http.StatusOK, "sucess")
// }

// func (con BaseController) error(c *gin.Context) {
// 	c.String(http.StatusOK, "error")
// }

func (con BaseController) getSession(c *gin.Context) sessions.Session {
	return model.GetSession(c)
}
