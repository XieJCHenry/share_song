package play_list

import (
	"share_song/global"
	logger2 "share_song/logger"
	"share_song/protocol"
)

const (
	PathCurPlayList   string = "/play_list_v2/api/v1/current_list"
	PathSetNext       string = "/play_list_v2/api/v1/set_next"
	PathStartPlay     string = "/play_list_v2/api/v1/start_play"
	PathSetPause      string = "/play_list_v2/api/v1/set_pause"
	PathCurrentStatus string = "/play_list_v2/api/v1/current_status"
)

func Register(dispatcher *protocol.Dispatcher) {
	dispatcher.Register(PathCurPlayList, GetCurrentList)
	dispatcher.Register(PathCurrentStatus, GetCurrentStatus)
	dispatcher.Register(PathSetNext, SetNextPlay)
	dispatcher.Register(PathStartPlay, StartPlay)
	dispatcher.Register(PathSetPause, SetPause)
}

func StartPlay(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	curPlaySong, pos, err := playListService.StartPlay(nil, userId)
	if err != nil {
		logger.Sugared().Errorf("start play failed, err=%s", err)
		return wrapperError(err)
	} else {
		var instanceId = ""
		var actualPos = -1
		if curPlaySong != nil {
			instanceId = curPlaySong.InstanceId
			actualPos = pos
		}
		return map[string]interface{}{
			"path":       PathStartPlay,
			"key":        userId,
			"playStatus": StatusPlaying,
			"current": map[string]interface{}{
				"instance_id": instanceId,
				"pos":         actualPos,
			},
		}, nil
	}
}

func SetNextPlay(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	songId := pro.Body["songId"].(string)
	songs, err := playListService.SetNextPlay(nil, userId, songId)
	if err != nil {
		logger.Sugared().Errorf("set next play failed, err=%s", err)
		return wrapperError(err)
	}
	return map[string]interface{}{
		"path":  PathSetNext,
		"key":   userId,
		"songs": songs,
	}, nil
}

func SetPause(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	err := playListService.SetPause(nil, userId)
	if err != nil {
		logger.Sugared().Errorf("start play failed, err=%s", err)
		return wrapperError(err)
	} else {
		return map[string]interface{}{
			"path":       PathSetPause,
			"key":        userId,
			"playStatus": StatusPaused,
		}, nil
	}
}

func GetCurrentList(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	songs, err := playListService.GetCurrentSongs(nil)
	if err != nil {
		logger.Sugared().Errorf("get current songs failed, err=%s", err)
		return wrapperError(err)
	}

	return map[string]interface{}{
		"path":  PathCurPlayList,
		"key":   userId,
		"songs": songs,
	}, nil
}

func GetCurrentStatus(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	song, pos, playStatus, err := playListService.GetCurrentStatus(nil)
	if err != nil {
		logger.Sugared().Errorf("get current status failed, err=%s", err)
		return wrapperError(err)
	}

	return map[string]interface{}{
		"path": PathCurrentStatus,
		"key":  userId,
		"status": map[string]interface{}{
			"status":      playStatus,
			"currentSong": song,
			"pos":         pos,
		},
	}, nil
}

func wrapperError(err error) (map[string]interface{}, error) {
	return map[string]interface{}{
		"error": err.Error(),
	}, nil
}
