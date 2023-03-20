package play_list

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PlayListController struct {
	logger  *zap.SugaredLogger
	service PlayerListService
}

func NewController(logger *zap.SugaredLogger, service PlayerListService) *PlayListController {
	return &PlayListController{
		logger:  logger,
		service: service,
	}
}

func (c *PlayListController) AddSongs(ctx *gin.Context) {

	request := &AddSongsRequest{}
	err := ctx.BindJSON(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("参数错误：%s", err.Error())})
		c.logger.Errorf("search songs bind json failed, err=%s", err)
		return
	}

	songs, err := c.service.AddSongs(ctx, request.UserId, request.SongIds)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		c.logger.Errorf("usr %s add songs failed, err=%s", request.UserId, err)
		return
	} else {
		ctx.JSON(200, &AddSongsResponse{
			Songs: songs,
		})
	}
}

func (c *PlayListController) RemoveSongs(ctx *gin.Context) {

	request := &RemoveSongsRequest{}
	err := ctx.BindJSON(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("参数错误：%s", err.Error())})
		c.logger.Errorf("search songs bind json failed, err=%s", err)
		return
	}

	songs, err := c.service.RemoveSongs(ctx, request.UserId, request.SongIds)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		c.logger.Errorf("usr %s add songs failed, err=%s", request.UserId, err)
		return
	} else {
		ctx.JSON(200, &RemoveSongsResponse{
			Songs: songs,
		})
	}
}

func (c *PlayListController) GetCurrentSongs(ctx *gin.Context) {
	currentSongs, _ := c.service.GetCurrentSongs(ctx.Request.Context())
	ctx.JSON(200, &GetCurrentSongsResponse{
		Songs: currentSongs,
	})
}
