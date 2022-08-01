package wbsocket

import (
	"fmt"
	"share_song/protocol"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type Connection struct {
	logger    *zap.SugaredLogger
	wsConn    *websocket.Conn
	readChan  chan protocol.Protocol
	writeChan chan protocol.Protocol
	closeChan chan struct{}
	closeMtx  sync.Mutex
	closed    *atomic.Bool
}

func NewJsonConnection(logger *zap.SugaredLogger, wcConn *websocket.Conn) *Connection {
	conn := &Connection{
		logger:    logger,
		wsConn:    wcConn,
		readChan:  make(chan protocol.Protocol),
		writeChan: make(chan protocol.Protocol),
		closeChan: make(chan struct{}),
		closeMtx:  sync.Mutex{},
		closed:    atomic.NewBool(false),
	}

	go conn.readLoop()
	go conn.writeLoop()

	return conn
}

func (c *Connection) Close() {
	c.wsConn.Close()

	c.closeMtx.Lock()
	defer c.closeMtx.Unlock()

	if !c.closed.Load() {
		c.closed.Store(true)
		close(c.closeChan)
	}
}

func (c *Connection) IsClosed() bool {
	return c.closed.Load()
}

func (c *Connection) ReadMessage() (protocol.Protocol, error) {
	select {
	case readData := <-c.readChan:
		return readData, nil
	case <-c.closeChan:
		return protocol.Protocol{}, fmt.Errorf("websocket closed")
	}
}

func (c *Connection) WriteMessage(pro protocol.Protocol) error {
	select {
	case c.writeChan <- pro:
	case <-c.closeChan:
		return fmt.Errorf("websocket closed")
	}
	return nil
}

func (c *Connection) readLoop() {
	defer func() {
		if err := recover(); err != nil {

		}
		c.Close()
	}()

	for {
		var pro protocol.Protocol
		err := c.wsConn.ReadJSON(&pro)
		if err != nil && c.closed.Load() {
			c.logger.Errorf("websocket read json failed, err=%s", err)
			return
		}
		c.logger.Debugf("websocket read json %s", pro)
		select {
		case c.readChan <- pro:
		case <-c.closeChan:
			c.logger.Info("websocket read chan close")
			return
		}
	}
}

func (c *Connection) writeLoop() {
	defer func() {
		if err := recover(); err != nil {

		}
		c.Close()
	}()

	for {
		select {
		case writeData := <-c.writeChan:
			c.logger.Debugf("websocket write json %s", writeData)
			err := c.wsConn.WriteJSON(writeData)
			if err != nil && c.closed.Load() {
				c.logger.Errorf("websocket write json failed, err=%s", err)
				return
			}

		case <-c.closeChan:
			c.logger.Info("websocket write chan close")
			return
		}
	}
}
