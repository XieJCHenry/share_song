package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger  *zap.SugaredLogger
	service Service
}

func NewController(logger *zap.SugaredLogger, service Service) *Controller {
	return &Controller{
		logger:  logger,
		service: service,
	}
}

func (c *Controller) Login(ctx *gin.Context) {

	req := &LoginRequest{}
	err := ctx.BindJSON(req)
	if err != nil {
		c.logger.Errorf("login bind json failed, err=%s", err)
		return
	}

	user, err := c.service.Login(ctx.Request.Context(), req.UserName, req.Phone)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(200, &LoginResponse{
			InstanceId: user.InstanceId,
			UserName:   user.Name,
			SongList:   user.OperatedSongs,
		})
	}
}

func (c *Controller) Logout(ctx *gin.Context) {

	req := &LogoutRequest{}
	err := ctx.BindJSON(req)
	if err != nil {
		c.logger.Errorf("login bind json failed, err=%s", err)
		return
	}

	err = c.service.Logout(ctx.Request.Context(), req.InstanceId, req.UserName)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(200, &LogoutResponse{
			InstanceId: req.InstanceId,
		})
	}

}

func (c *Controller) RegisterAccount(ctx *gin.Context) {
	req := &RegisterAccountRequest{}
	err := ctx.BindJSON(req)
	if err != nil {
		c.logger.Errorf("login bind json failed, err=%s", err)
		return
	}

	user, err := c.service.RegisterAccount(ctx.Request.Context(), req.UserName, req.Phone)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, &RegisterAccountResponse{
			InstanceId: user.InstanceId,
			UserName:   user.Name,
		})
		return
	}
}

func (c *Controller) CancelAccount(ctx *gin.Context) {
	req := &CancelAccountRequest{}
	err := ctx.BindJSON(req)
	if err != nil {
		c.logger.Errorf("login bind json failed, err=%s", err)
		return
	}

	err = c.service.CancelAccount(ctx.Request.Context(), req.InstanceId, req.UserName, req.Phone)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(200, &CancelAccountResponse{
			InstanceId: req.InstanceId,
		})
	}
}

func (c *Controller) SearchOnlineUsers(ctx *gin.Context) {
	req := &SearchOnlineUsersRequest{}
	err := ctx.BindJSON(req)
	if err != nil {
		c.logger.Errorf("login bind json failed, err=%s", err)
		return
	}

	users, err := c.service.SearchOnlineUsers(ctx.Request.Context(), req.Query, req.WithTimeOut)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(200, &SearchOnlineUsersResponse{
			WithTimeOut: req.WithTimeOut,
			Users:       users,
		})
	}
}

func (c *Controller) GetOnlineToken(ctx *gin.Context) {
	loginKey := ctx.Query("login-key")
	onlineToken, err := c.service.GetOnlineToken(ctx.Request.Context(), loginKey)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, gin.H{
			"loginKey":    loginKey,
			"onlineToken": onlineToken,
		})
	}
}
