package play_list

import (
	"math"
	"share_song/music"
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
const (
	StatusPlaying PlayStatus = "playing"
	StatusPaused  PlayStatus = "paused"

	ModeLoop     PlayMode = "loop"     // 循环播放
	ModeSequence PlayMode = "sequence" // 顺序播放
)

type PlayList struct {
	list        []*music.Song
	curPlaySong *currentPlaying // start from 0

	status int32
	mode   PlayMode
	mtx    sync.Mutex
	playChan chan struct{}
	pauseChan chan struct{}
}

func NewPlayList() *PlayList {
	pList := &PlayList{
		list:        make([]*music.Song, 0),
		status:      pause,
		mtx:         sync.Mutex{},
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
	if currentPos >= len(p.list) || currentPos < 0 {
		p.list = append(p.list, song)
	} else {
		newSlice := append([]*music.Song{song}, p.list[currentPos+1:]...)
		p.list = append(p.list[0:currentPos+1], newSlice...)
	}
}

func (p *PlayList) PlayNext() {
	p.playChan <- struct{}{}
}

func (p *PlayList) StartPlay() {
	if atomic.LoadInt32(&p.status) != play {
		p.mtx.Lock()
		if atomic.LoadInt32(&p.status) != play {
			{
				p.PlayNext()
				go p.playLoop()
			}
		}
		p.mtx.Unlock()
	}
}

func (p *PlayList) SetPause() {
	p.pauseChan <- struct{}{}
}

func (p *PlayList) playLoop() {
	// playChan 接收信号则播放
	// pauseChan 接收信号则暂停

	for {
		select {
		case <- p.pauseChan:
			{
				if atomic.CompareAndSwapInt32(&p.status, play, pause) {
					p.curPlaySong.pause()
				}
			}
		case <- p.playChan:
			{
				if atomic.CompareAndSwapInt32(&p.status, pause, play) {
				}

				p.mtx.Lock()
				if p.curPlaySong.Pos < 0 {
					p.curPlaySong.Pos++
					p.curPlaySong.PlayOffset = 0
				}
				var current int
				if p.mode == ModeLoop {
					current = (p.curPlaySong.Pos) % len(p.list)
				} else if p.mode == ModeSequence {
					current = int(math.Min(float64(p.curPlaySong.Pos), float64(len(p.list) - 1)))
				}
				currentSong := p.list[current]
				p.mtx.Unlock()

				p.curPlaySong.length = currentSong.Length
				p.curPlaySong.start()
			}
		default:
			{}
		}
	}
}

func (p *PlayList) GetCurrentPlay() (*music.Song, int) {
	curPlaySong := p.curPlaySong
	song := p.list[curPlaySong.Pos]
	return song, curPlaySong.Pos
}