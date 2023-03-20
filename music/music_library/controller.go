package music_library

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MusicLibraryController struct {
	logger  *zap.SugaredLogger
	service MusicLibraryService
}

func NewController(logger *zap.SugaredLogger, service MusicLibraryService) *MusicLibraryController {
	return &MusicLibraryController{
		logger:  logger,
		service: service,
	}
}

func (c *MusicLibraryController) BatchAddSong(ctx *gin.Context) {
	request := &BatchAddSongRequest{}
	err := ctx.BindJSON(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("参数错误：%s", err.Error())})
		c.logger.Errorf("search songs bind json failed, err=%s", err)
		return
	}

	result, err := c.service.BatchAddSong(request.Songs)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, &BatchAddSongResponse{
			Result: result,
		})
	}
}

func (c *MusicLibraryController) BatchRemoveSong(ctx *gin.Context) {
	request := &BatchRemoveSongRequest{}
	err := ctx.BindJSON(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("参数错误：%s", err.Error())})
		c.logger.Errorf("search songs bind json failed, err=%s", err)
		return
	}

	result, err := c.service.BatchRemoveSong(request.SongIds)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, &BatchRemoveSongResponse{
			Result: result,
		})
	}
}

func (c *MusicLibraryController) SearchSongs(ctx *gin.Context) {

	request := &SearchSongsRequest{}
	err := ctx.BindJSON(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("参数错误：%s", err.Error())})
		c.logger.Errorf("search songs bind json failed, err=%s", err)
		return
	}

	searchSongs, err := c.service.SearchSongs(request.Fields, request.Page, request.PageSize)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(200, &SearchSongsResponse{
			Songs:    searchSongs,
			Page:     request.Page,
			PageSize: request.PageSize,
		})
	}
}

func (c *MusicLibraryController) SearchSongByCondition(ctx *gin.Context) {

}
