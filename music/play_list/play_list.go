package play_list

import (
	"share_song/music"
	"share_song/utils/collection/linked_list"
	"sync"
)

type PlayList struct {
	list     []*music.Song
	linklist *linked_list.List[*music.Song]
	current  int
	mtx      sync.Mutex
}

func NewPlayList() *PlayList {
	return &PlayList{
		list:     make([]*music.Song, 0),
		linklist: linked_list.NewList[*music.Song](),
		current:  -1,
		mtx:      sync.Mutex{},
	}
}

func (p *PlayList) Append(song *music.Song) bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.list = append(p.list, song)
	p.linklist.Append(song)
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

	if p.current >= len(p.list) || p.current < 0 {
		p.list = append(p.list, song)
	} else {
		newSlice := append([]*music.Song{song}, p.list[p.current+1:]...)
		p.list = append(p.list[0:p.current+1], newSlice...)
	}
}
