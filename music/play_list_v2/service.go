package play_list_v2

import (
	"share_song/global"
	logger2 "share_song/logger"
	"share_song/music/play_list"
	"share_song/protocol"
)

func Register(dispatcher *protocol.Dispatcher) {
	dispatcher.Register("/play_list_v2/api/v1/current_list", GetCurrentList)
	dispatcher.Register("/play_list_v2/api/v1/set_next", SetNextPlay)
}

func SetNextPlay(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(play_list.PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	userId := pro.Body["userId"].(string)
	songId := pro.Body["songId"].(string)
	songs, err := playListService.SetNextPlay(nil, userId, songId)
	if err != nil {
		logger.Sugared().Errorf("set next play failed, err=%s", err)
		return map[string]interface{}{
			"error": err.Error(),
		}, nil
	}
	return map[string]interface{}{
		"songs": songs,
	}, nil
}

func GetCurrentList(pro *protocol.Protocol) (map[string]interface{}, error) {
	playListService := global.GetGlobalObject(global.KeyPlayListService).(play_list.PlayerListService)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	songs, err := playListService.GetCurrentSongs(nil)
	if err != nil {
		logger.Sugared().Errorf("get current songs failed, err=%s", err)
		return map[string]interface{}{
			"error": err.Error(),
		}, nil
	}

	return map[string]interface{}{
		"songs": songs,
	}, nil

}
