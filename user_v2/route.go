package user_v2

import "github.com/gin-gonic/gin"

func (c *Controller) Register(e *gin.Engine) {
	e.GET("user/api/v1/connect", c.Connect)
}
