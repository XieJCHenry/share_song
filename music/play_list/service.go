package play_list

import (
	"context"
	"share_song/global"
	"share_song/music"

	"github.com/gin-gonic/gin"
)

type PlayerListService interface {
	global.UsableObject
	AddSongs(ctx *gin.Context, userId string, songIds []string) ([]music.Song, error)
	RemoveSongs(ctx *gin.Context, userId string, songIds []string) ([]music.Song, error)
	GetCurrentSongs(ctx context.Context) ([]music.Song, error)

	SetNextPlay(ctx *gin.Context, userId string, songId string) ([]music.Song, error)
	SetPause(ctx *gin.Context, userId string) error
	StartPlay(ctx *gin.Context, userId string) (*music.Song, int, error)
}
