package user_v2

import "github.com/gin-gonic/gin"

type Controller struct {
	s Service
}

func NewController(s Service) *Controller {
	return &Controller{
		s: s,
	}
}

func (c *Controller) Connect(ctx *gin.Context) {
	var connKey string = ctx.Query("online-token")
	c.s.Connect(connKey, ctx.Writer, ctx.Request)
}
