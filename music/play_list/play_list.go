package play_list

import (
	"math"
	"share_song/global"
	"share_song/internal/wbsocket"
	logger2 "share_song/logger"
	"share_song/music"
	"share_song/protocol"
	"sync"
	"sync/atomic"
)

/*
播放服务：
1、状态有几种？
播放中、暂停、切换中

服务端的基本播放逻辑
每一首播放时，都开启一个定时器，定时为歌曲时长，结束后切换到下一首。
*/

const (
	play  int32 = 0
	pause int32 = 1
)

type PlayStatus string
type PlayMode string
type Action string

const (
	StatusPlaying PlayStatus = "playing"
	StatusPaused  PlayStatus = "paused"

	ModeLoop     PlayMode = "loop"     // 循环播放
	ModeSequence PlayMode = "sequence" // 顺序播放
)

type PlayList struct {
	list        []*music.Song
	curPlaySong *currentPlaying // start from 0

	status    int32
	mode      PlayMode
	mtx       sync.Mutex
	playChan  chan struct{}
	pauseChan chan struct{}
}

func NewPlayList() *PlayList {
	pList := &PlayList{
		list:      make([]*music.Song, 0),
		status:    pause,
		mode:      ModeLoop,
		mtx:       sync.Mutex{},
		playChan:  make(chan struct{}, 1),
		pauseChan: make(chan struct{}, 1),
	}
	curPlay := NewCurrentPlaying(0, pList)
	pList.curPlaySong = curPlay
	return pList
}

func (p *PlayList) Append(song *music.Song) bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.list = append(p.list, song)
	return true
}

func (p *PlayList) Remove(song *music.Song) {

}

func (p *PlayList) BatchRemove(songs []*music.Song) []string {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	deletedMap := make(map[string]struct{})
	for _, song := range songs {
		deletedMap[song.InstanceId] = struct{}{}
	}

	var newList []*music.Song
	var deletedIds []string
	for _, song := range p.list {
		if _, ok := deletedMap[song.InstanceId]; !ok {
			newList = append(newList, song)
		} else {
			deletedIds = append(deletedIds, song.InstanceId)
		}
	}
	p.list = newList
	return deletedIds
}

func (p *PlayList) GetAll() []*music.Song {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.list
}

func (p *PlayList) SetNext(song *music.Song) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	currentPos := p.curPlaySong.Pos
	if currentPos >= len(p.list) {
		p.list = append(p.list, song)
	} else {
		newSlice := append([]*music.Song{song}, p.list[currentPos+1:]...)
		p.list = append(p.list[0:currentPos+1], newSlice...)
	}
}

func (p *PlayList) PlayNext() {
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	if atomic.LoadInt32(&p.status) == play {
		p.mtx.Lock()
		if atomic.LoadInt32(&p.status) == play {
			p.curPlaySong.reset()
			logger.Sugared().Infof("song play next")
			p.playChan <- struct{}{}
		}
		p.mtx.Unlock()
	}
}

func (p *PlayList) StartPlay() (*music.Song, int) {

	wg := sync.WaitGroup{}

	if len(p.list) > 0 {
		if atomic.LoadInt32(&p.status) != play {
			p.mtx.Lock()
			if atomic.LoadInt32(&p.status) != play {
				{
					wg.Add(1)
					go p.playLoop(&wg)
				}
			}
			p.mtx.Unlock()
		}
	}

	wg.Wait()
	return p.GetCurrentPlay()
}

func (p *PlayList) SetPause() {
	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)
	if atomic.CompareAndSwapInt32(&p.status, play, pause) {
		logger.Sugared().Infof("song pause")
		p.pauseChan <- struct{}{}
	}
}

func (p *PlayList) playLoop(wg *sync.WaitGroup) {
	// playChan 接收信号则播放
	// pauseChan 接收信号则暂停

	logger := global.GetGlobalObject(global.KeyLogger).(*logger2.Logger)

	for {
		select {
		case <-p.pauseChan:
			{
				logger.Sugared().Debugf("pause chan signal")
				if atomic.CompareAndSwapInt32(&p.status, play, pause) {
					logger.Sugared().Infof("status swap from play to pause")
					p.curPlaySong.pause()
					p.onPause()
				}
			}
		case <-p.playChan:
			{
				if atomic.CompareAndSwapInt32(&p.status, pause, play) {
					logger.Sugared().Infof("status swap from pause to play")
				}

				p.mtx.Lock()
				p.curPlaySong.Pos++
				p.curPlaySong.PlayOffset = 0

				if p.curPlaySong.Pos >= len(p.list) && p.mode == ModeSequence {
					p.SetPause()
				} else {
					var current int
					if p.mode == ModeLoop {
						current = (p.curPlaySong.Pos) % len(p.list)
					} else if p.mode == ModeSequence {
						current = int(math.Min(float64(p.curPlaySong.Pos), float64(len(p.list)-1)))
					}
					currentSong := p.list[current]
					logger.Sugared().Infof("currentSong is %v", currentSong)
					p.mtx.Unlock()

					p.curPlaySong.length = currentSong.Length
					p.curPlaySong.Pos = current
					p.curPlaySong.start()
					p.onPlay()
				}
			}
		default:
			{
				if atomic.CompareAndSwapInt32(&p.status, pause, play) {
					p.mtx.Lock()

					logger.Sugared().Debugf("send start signal at default")
					p.playChan <- struct{}{}
					wg.Done()

					p.mtx.Unlock()
				}
			}
		}
	}
}

func (p *PlayList) GetCurrentPlay() (*music.Song, int) {
	curPlaySong := p.curPlaySong
	if curPlaySong.Pos < 0 || len(p.list) <= 0 {
		return nil, -1
	}
	song := p.list[curPlaySong.Pos]
	return song, curPlaySong.Pos
}

func (p *PlayList) GetCurrentStatus() (*music.Song, int, PlayStatus) {
	curSong, pos := p.GetCurrentPlay()
	if p.status == play {
		return curSong, pos, StatusPlaying
	} else {
		return nil, -1, StatusPaused
	}
}

func (p *PlayList) onPlay() {
	connPool := global.GetGlobalObject(global.KeyConnPool).(*wbsocket.Pool)
	currentSong, pos := p.GetCurrentPlay()
	if currentSong != nil && pos >= 0 {
		connPool.ForEach(func(o *wbsocket.Owner) error {
			err := o.Conn().WriteMessage(protocol.Protocol{
				Body: map[string]interface{}{
					"key":        o.Key(),
					"path":       PathStartPlay,
					"playStatus": StatusPlaying,
					"current": map[string]interface{}{
						"instance_id": currentSong.InstanceId,
						"pos":         pos,
					},
				},
			})
			return err
		})
	}
}

func (p *PlayList) onPause() {
	connPool := global.GetGlobalObject(global.KeyConnPool).(*wbsocket.Pool)
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
}
