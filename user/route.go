package user

import "github.com/gin-gonic/gin"

func (c *Controller) Register(e *gin.Engine) {
	e.POST("/user/api/v1/login", c.Login)
	e.POST("/user/api/v1/logout", c.Logout)
	e.POST("/user/api/v1/register", c.RegisterAccount)
	e.POST("/user/api/v1/cancel", c.CancelAccount)
	e.POST("/user/api/v1/search", c.SearchOnlineUsers)
	e.POST("/user/api/v1/get-online-token", c.GetOnlineToken)
}
