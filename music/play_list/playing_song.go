package play_list

import (
	"share_song/global"
	logger2 "share_song/logger"
	"time"
)

type currentPlaying struct {
	Pos        int    // 在播放列表中的下标
	PlayOffset uint32 // 当前播放进度
	length     uint32
	pauseChan  chan struct{}
	pList      *PlayList
}

func NewCurrentPlaying(length uint32, pList *PlayList) *currentPlaying {
	return &currentPlaying{
		Pos:        -1,
		PlayOffset: 0,
		length:     length,
		pauseChan:  make(chan struct{}, 1),
		pList:      pList,
	}
}

func (c *currentPlaying) pause() {
	c.pauseChan <- struct{}{}
}

func (c *currentPlaying) start() {
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	dur := c.length - c.PlayOffset
	if dur <= 0 {
		c.pList.PlayNext()
		return
	}
	tm := time.NewTimer(time.Second * time.Duration(dur))

	go func() {
		defer func() {
			if c.PlayOffset >= c.length {
				c.pList.PlayNext()
			}
		}()

		for {
			select {
			case <-c.pauseChan:
				{
					logger.Sugared().Debugf("song play pause")
					return
				}
			case <-tm.C:
				{
					c.PlayOffset = c.length
					logger.Sugared().Debugf("song play shutdown by timer")
					return
				}
			default:
				{
					time.Sleep(time.Second * 1)
					c.PlayOffset++
					if c.PlayOffset%10 == 0 {
						logger.Sugared().Debugf("song offset = %d", c.PlayOffset)
					}
					if c.PlayOffset >= c.length {
						logger.Sugared().Infof("song %d play to the end", c.PlayOffset)
						return
					}
				}
			}
		}
	}()
}

func (c *currentPlaying) reset() {
	c.length = 0
	c.PlayOffset = 0
}
