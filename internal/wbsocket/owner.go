package wbsocket

import (
	"share_song/global"
	logger2 "share_song/logger"
	"share_song/protocol"
)

type Owner struct {
	key  string
	conn *Connection
}

func NewOwner(key string, conn *Connection) *Owner {
	return &Owner{
		key:  key,
		conn: conn,
	}
}

func (o *Owner) Run() error {
	dispatcher := global.GetGlobalObject(global.KeyDispatcher).(*protocol.Dispatcher)
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger).Sugared()

	defer func() {
		logger.Infof("websocket close %s", o.key)
		o.stop()
	}()

	for {
		pro, err := o.conn.ReadMessage()
		if err != nil {
			logger.Errorf("owner read message %s failed, err=%s", pro.Path, err)
			return err
		}
		response, err := dispatcher.Dispatch(&pro)
		if err != nil {
			logger.Errorf("owner dispatch %s failed, err=%s", pro.Path, err)
			continue
		}
		err = o.conn.WriteMessage(protocol.Protocol{Body: response})
		if err != nil {
			logger.Errorf("owner write message %s failed, err=%s", pro.Path, err)
			return err
		}
	}
}

func (o *Owner) Key() string {
	return o.key
}

func (o *Owner) Conn() *Connection {
	return o.conn
}

func (o *Owner) OnLogin() {
	//playListService.GetCurrentSongs()
}

func (o *Owner) stop() {
	o.conn.Close()
	connPool := global.GetGlobalObject(global.KeyConnPool).(*Pool)
	connPool.Remove(o.Key())
}
