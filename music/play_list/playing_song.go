package play_list

import "time"

type currentPlaying struct {
	Pos int	// 在播放列表中的下标
	PlayOffset uint32 // 当前播放进度
	length uint32
	pauseChan chan struct{}
	pList *PlayList
}

func NewCurrentPlaying(length uint32, pList *PlayList) *currentPlaying {
	return &currentPlaying{
		Pos:        -1,
		PlayOffset: 0,
		length:     length,
		pauseChan:  make(chan struct{}),
		pList:      pList,
	}
}

func (c *currentPlaying) pause() {
	c.pauseChan <- struct{}{}
}

func (c *currentPlaying) start() {
	dur := c.length - c.PlayOffset
	if dur <= 0 {
		c.pList.PlayNext()
		return
	}
	tm := time.NewTimer(time.Second * time.Duration(dur))

	go func() {
		for {
			select {
			case <- c.pauseChan:
				return
			case <- tm.C:
				{
					c.pList.PlayNext()
					return
				}
			default:
				{
					time.Sleep(time.Second * 1)
					c.PlayOffset++
				}
			}
		}
	}()
}

