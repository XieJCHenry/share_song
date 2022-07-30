package play_list

import "github.com/gin-gonic/gin"

func (c *PlayListController) Register(e *gin.Engine) {
	e.POST("/play_list/api/v1/add", c.AddSongs)
	e.DELETE("/play_list/api/v1/delete", c.RemoveSongs)
	e.GET("/play_list/api/v1/current_list", c.GetCurrentSongs)
}
