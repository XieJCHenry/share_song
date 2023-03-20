package hello

import "share_song/protocol"

func Register(dispatcher *protocol.Dispatcher) {
	dispatcher.Register("/api/v1/hello", HelloWorld)
}

func HelloWorld(pro *protocol.Protocol) (map[string]interface{}, error) {

	return map[string]interface{}{
		"msg": "hello world!",
	}, nil
}
