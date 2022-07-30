package wbsocket

import (
	"share_song/global"
	"sync"

	"go.uber.org/zap"
)

type Pool struct {
	logger *zap.SugaredLogger
	data   map[string]*Owner
	mtx    sync.Mutex
}

func (p *Pool) Usable() bool {
	if p.data == nil || p.logger == nil {
		return false
	}
	return true
}

func (p *Pool) Key() string {
	return global.KeyConnPool
}

func NewConnectionPool(logger *zap.SugaredLogger) *Pool {

	return &Pool{
		logger: logger,
		data:   make(map[string]*Owner),
		mtx:    sync.Mutex{},
	}
}

func (p *Pool) Add(key string, con *Owner) {
	if _, exists := p.data[key]; !exists {
		p.mtx.Lock()
		if _, exists2 := p.data[key]; !exists2 {
			p.data[key] = con
		}
		p.mtx.Unlock()
	}
}

func (p *Pool) Remove(key string) {
	if _, exists := p.data[key]; exists {
		p.mtx.Lock()
		if _, exists2 := p.data[key]; exists2 {
			delete(p.data, key)
		}
		p.mtx.Unlock()
	}
}

func (p *Pool) Get(key string) *Owner {
	var result *Owner
	if _, exists := p.data[key]; exists {
		p.mtx.Lock()
		if con, exists2 := p.data[key]; exists2 {
			result = con
		}
		p.mtx.Unlock()
	}
	return result
}

func (p *Pool) Contains(key string) bool {
	return p.Get(key) != nil
}

func (p *Pool) ForEach(f func(o *Owner) error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	for _, ow := range p.data {
		if err := f(ow); err != nil {
			p.logger.Errorf("for each func failed, err=%s", err)
		}
	}
}
