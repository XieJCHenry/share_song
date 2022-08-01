package play_list_v2

import (
	"share_song/global"
	logger2 "share_song/logger"
	"share_song/music/play_list"
	"share_song/protocol"
)

const (
	PathCurrentList string = "/play_list_v2/api/v1/current_list"
	PathSetNext string = "/play_list_v2/api/v1/set_next"
	PathStartPlay string = "/play_list_v2/api/v1/start_play"
	PathSetPause string = "/play_list_v2/api/v1/set_pause"
)

func Register(dispatcher *protocol.Dispatcher) {
	dispatcher.Register(PathCurrentList, GetCurrentList)
	dispatcher.Register(PathSetNext, SetNextPlay)
	dispatcher.Register(PathStartPlay, StartPlay)
	dispatcher.Register(PathSetPause, SetPause)
}

func StartPlay(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(play_list.PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	curPlaySong, pos, err := playListService.StartPlay(nil, userId)
	if err != nil {
		logger.Sugared().Errorf("start play failed, err=%s", err)
		return wrapperError(err)
	} else {
		return map[string]interface{}{
			"key":   userId,
			"playStatus": play_list.StatusPlaying,
			"current": map[string]interface{}{
				"instance_id": curPlaySong.InstanceId,
				"pos": pos,
			},
		}, nil
	}
}

func SetNextPlay(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(play_list.PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	songId := pro.Body["songId"].(string)
	songs, err := playListService.SetNextPlay(nil, userId, songId)
	if err != nil {
		logger.Sugared().Errorf("set next play failed, err=%s", err)
		return wrapperError(err)
	}
	return map[string]interface{}{
		"songs": songs,
	}, nil
}

func SetPause(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(play_list.PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	err := playListService.SetPause(nil, userId)
	if err != nil {
		logger.Sugared().Errorf("start play failed, err=%s", err)
		return wrapperError(err)
	} else {
		return map[string]interface{}{
			"key":   userId,
			"playStatus": play_list.StatusPaused,
		}, nil
	}
}

func GetCurrentList(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(play_list.PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	songs, err := playListService.GetCurrentSongs(nil)
	if err != nil {
		logger.Sugared().Errorf("get current songs failed, err=%s", err)
		return wrapperError(err)
	}

	return map[string]interface{}{
		"key":   userId,
		"songs": songs,
	}, nil

}

func wrapperError(err error) (map[string]interface{}, error) {
	return map[string]interface{}{
		"error": err.Error(),
	}, nil
}
