package user_v2

import (
	"net/http"
	"share_song/global"
	"share_song/internal/wbsocket"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Service interface {
	Connect(connKey string, w http.ResponseWriter, r *http.Request)
}

type service struct {
	logger *zap.SugaredLogger
}

func NewUserServiceV2(logger *zap.SugaredLogger) Service {

	return &service{
		logger: logger,
	}
}

func (s *service) Connect(connKey string, w http.ResponseWriter, r *http.Request) {
	connPool := global.GetGlobalObject(global.KeyConnPool).(*wbsocket.Pool)

	if !connPool.Contains(connKey) {
		var upGrader = websocket.Upgrader{}
		c, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			s.logger.Errorf("upgrade connection failed, err=%s", err)
			return
		}
		con := wbsocket.NewJsonConnection(s.logger, c)
		owner := wbsocket.NewOwner(connKey, con)
		owner.OnLogin() // todo 在登录时推送一次当前的播放状态
		connPool.Add(connKey, owner)
		go owner.Run()
	}
}
