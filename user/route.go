package user

import "github.com/gin-gonic/gin"

func (c *Controller) Register(e *gin.Engine) {
	e.POST("/api/v1/login", c.Login)
	e.POST("/api/v1/logout", c.Logout)
	e.POST("/api/v1/register", c.RegisterAccount)
	e.POST("/api/v1/cancel", c.CancelAccount)
	e.POST("/api/v1/search", c.SearchOnlineUsers)
}
