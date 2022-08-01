package play_list

import (
	"context"
	"fmt"
	"share_song/global"
	"share_song/internal/wbsocket"
	"share_song/music"
	"share_song/music/music_library"
	"share_song/protocol"
	"share_song/user"
	_set "share_song/utils/collection/set"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// todo 需要实现“播放”逻辑，每一首曲子播放完后要更新playList的current，并且通知每一个client

var (
	songSearchFields = []string{"instance_id", "name", "artists", "length"}
)

type service struct {
	logger         *zap.SugaredLogger
	userService    user.Service
	libraryService music_library.MusicLibraryService
	playList       *PlayList
	connPool       *wbsocket.Pool
}

func (s *service) Usable() bool {
	return true
}

func (s *service) Key() string {
	return global.KeyPlayListService
}

func NewService(logger *zap.SugaredLogger, userService user.Service, libraryService music_library.MusicLibraryService) PlayerListService {

	return &service{
		logger:         logger,
		userService:    userService,
		libraryService: libraryService,
		playList:       NewPlayList(),
		connPool:       wbsocket.NewConnectionPool(logger),
	}
}

func (s *service) AddSongs(ctx *gin.Context, userId string, songIds []string) ([]music.Song, error) {
	return s.addSongsV1(ctx, userId, songIds)
}

func (s *service) addSongsV1(ctx *gin.Context, userId string, songIds []string) ([]music.Song, error) {

	var connKey string = getConnectionKey(ctx)
	if !s.connPool.Contains(connKey) {
		return nil, fmt.Errorf("未知连接，禁止操作")
	}

	query := map[string]interface{}{
		"instance_id": songIds,
	}
	songs, err := s.libraryService.SearchSongByCondition(songSearchFields, query, 1, 3000)
	if err != nil {
		s.logger.Errorf("search songs failed, err=%s", err)
		return nil, fmt.Errorf("查找歌曲失败，错误：%s", err)
	}

	songSet := _set.NewSetWith[string](songIds...)

	var remainSongs []music.Song
	for _, song := range songs {
		if songSet.Contains(song.InstanceId) {
			remainSongs = append(remainSongs, song)
		}
	}

	currentPlayList := s.playList.GetAll()
	go func() {
		s.connPool.ForEach(func(o *wbsocket.Owner) error {
			if o.Key() != userId {
				err1 := o.Conn().WriteMessage(protocol.Protocol{
					Body: map[string]interface{}{
						"songs": currentPlayList,
					},
				})
				if err1 != nil {
					s.logger.Errorf("websocket closed")
				}
			}
			return nil
		})
	}()

	return nil, nil
}

func getConnectionKey(ctx *gin.Context) string {
	return ctx.Query("online-token")
}

func (s *service) SetNextPlay(ctx *gin.Context, userId string, songId string) ([]music.Song, error) {

	query := map[string]interface{}{
		"instance_id": songId,
	}
	songs, err := s.libraryService.SearchSongByCondition(songSearchFields, query, 1, 1)
	if err != nil {
		s.logger.Errorf("search songs failed, err=%s", err)
		return nil, fmt.Errorf("查找歌曲失败，错误：%s", err)
	}

	if len(songs) <= 0 {
		return []music.Song{}, nil
	}

	song := &songs[0]

	s.playList.SetNext(song)
	currentPlayList := s.playList.GetAll()

	result := make([]music.Song, 0, len(currentPlayList))
	for _, song := range currentPlayList {
		result = append(result, *song)
	}

	connPool := global.GetGlobalObject(global.KeyConnPool).(*wbsocket.Pool)
	connPool.ForEach(func(o *wbsocket.Owner) error {
		if o.Key() != userId {
			err1 := o.Conn().WriteMessage(protocol.Protocol{
				Body: map[string]interface{}{
					"key":   o.Key(),
					"path":  PathSetNext,
					"songs": result,
				},
			})
			if err1 != nil {
				return err1
			}
		}
		return nil
	})

	return result, nil
}

func (s *service) RemoveSongs(ctx *gin.Context, userId string, songIds []string) ([]music.Song, error) {
	//user, err := s.searchOnlineUserByInstanceId(ctx, userId)
	//if err != nil {
	//	s.logger.Errorf("search online user %s failed, err=%s", userId, err)
	//	return nil, fmt.Errorf("查找用户失败：该用户不在线")
	//}

	fields := []string{"instanceId"}
	query := map[string]interface{}{
		"instanceId": songIds,
	}
	songs, err := s.libraryService.SearchSongByCondition(fields, query, 1, 3000)
	if err != nil {
		s.logger.Errorf("search songs failed, err=%s", err)
		return nil, fmt.Errorf("查找歌曲失败，错误：%s", err)
	}

	songsMap := make(map[string]struct{})
	for _, id := range songIds {
		songsMap[id] = struct{}{}
	}

	filterMap := make(map[string]int)
	var remainSongs []*music.Song
	for i, song := range songs {
		if _, ok := songsMap[song.InstanceId]; ok {
			filterMap[song.InstanceId] = i
			remainSongs = append(remainSongs, &song)
		}
	}

	var resultList []music.Song
	if len(remainSongs) > 0 {
		go func(resultSongList *[]music.Song, remainSongList []*music.Song) {
			if len(remainSongs) > 0 {
				s.playList.BatchRemove(remainSongs)
				current := s.playList.GetAll()
				for _, song := range current {
					*resultSongList = append(*resultSongList, *song)
				}
			}
			//user.UpdateOperatedSongs(deletedIds, "delete")
		}(&resultList, remainSongs)
	}

	return resultList, nil
}

func (s *service) GetCurrentSongs(ctx context.Context) ([]music.Song, error) {

	current := s.playList.GetAll()

	result := make([]music.Song, 0, len(current))
	for _, song := range current {
		result = append(result, *song)
	}

	return result, nil
}

func (s *service) SetPause(ctx *gin.Context, userId string) error {
	connPool := global.GetGlobalObject(global.KeyConnPool).(*wbsocket.Pool)

	s.playList.SetPause()

	connPool.ForEach(func(o *wbsocket.Owner) error {
		err := o.Conn().WriteMessage(protocol.Protocol{
			Body: map[string]interface{}{
				"key":        o.Key(),
				"path":       PathSetPause,
				"playStatus": StatusPaused,
			},
		})
		return err
	})
	return nil
}

func (s *service) StartPlay(ctx *gin.Context, userId string) (*music.Song, int, error) {
	connPool := global.GetGlobalObject(global.KeyConnPool).(*wbsocket.Pool)

	var currentPlay *music.Song
	var pos = -1

	// todo 需要测试正确性
	go func() {
		for {
			select {
			case curPlaySong := <-s.playList.CurrentChan:
				pos = s.playList.curPlaySong.Pos
				if curPlaySong != nil && pos >= 0 {
					connPool.ForEach(func(o *wbsocket.Owner) error {
						err := o.Conn().WriteMessage(protocol.Protocol{
							Body: map[string]interface{}{
								"key":        o.Key(),
								"path":       PathStartPlay,
								"playStatus": StatusPlaying,
								"current": map[string]interface{}{
									"instance_id": currentPlay.InstanceId,
									"pos":         pos,
								},
							},
						})
						return err
					})
				}
			}
		}
	}()

	s.playList.StartPlay()

	s.logger.Debugf("current is %v, pos = %d", currentPlay, pos)

	return currentPlay, pos, nil
}

func (s *service) searchOnlineUserByInstanceId(ctx context.Context, userId string) (*user.User, error) {
	onlineUsers, err := s.userService.SearchOnlineUsers(ctx, map[string]interface{}{
		"instanceId": userId,
	}, true)
	if err != nil {
		s.logger.Errorf("sarch user %s failed, err=%s", userId, err)
		return nil, err
	}

	if len(onlineUsers) == 1 {
		return onlineUsers[0], nil
	}

	s.logger.Errorf("user %s not found", userId)
	return nil, nil
}
