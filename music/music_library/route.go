package music_library

import "github.com/gin-gonic/gin"

func (c *MusicLibraryController) Register(e *gin.Engine) {
	e.POST("/music_library/api/v1/add", c.BatchAddSong)
	//e.DELETE("/music_library/api/v1/remove", c.BatchRemoveSong)
	e.POST("/music_library/api/v1/search", c.SearchSongs)
}
