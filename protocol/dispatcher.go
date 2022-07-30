package protocol

import (
	"fmt"
	"share_song/global"
)

type RouteMapFunc = func(pro *Protocol) (map[string]interface{}, error)

type Dispatcher struct {
	routeMap map[string]RouteMapFunc
}

func (d *Dispatcher) Usable() bool {
	return d.routeMap != nil
}

func (d *Dispatcher) Key() string {
	return global.KeyDispatcher
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		routeMap: make(map[string]RouteMapFunc),
	}
}

func (d *Dispatcher) Register(path string, handle RouteMapFunc) {
	if _, ok := d.routeMap[path]; !ok {
		d.routeMap[path] = handle
	}
}

func (d *Dispatcher) Dispatch(pro *Protocol) (map[string]interface{}, error) {
	if routeFunc, ok := d.routeMap[pro.Path]; ok {
		response, err := routeFunc(pro)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, fmt.Errorf("path not found")
}
