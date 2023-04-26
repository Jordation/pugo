package main

import (
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
)

type QC_SM struct {
	sync.RWMutex
	m map[string]*QueueChannel
}
type L_SM struct {
	sync.RWMutex
	m map[string]*Lobby
}

func (sm *QC_SM) Make() {
	sm.m = make(map[string]*QueueChannel)
}
func (sm *QC_SM) Get(id string) (*QueueChannel, bool) {
	sm.Lock()
	defer sm.Unlock()
	v, ok := sm.m[id]
	if ok {
		return v, ok
	}
	_, f, l, _ := runtime.Caller(1)
	log.Error("[QUEUE MAP GET]: Key not found ", f, l)
	return nil, ok
}

func (sm *L_SM) Make() {
	sm.m = make(map[string]*Lobby)
}
func (sm *QC_SM) Set(id string, val *QueueChannel) {
	sm.Lock()
	defer sm.Unlock()
	sm.m[id] = val
}

func (sm *L_SM) Get(id string) (*Lobby, bool) {
	sm.Lock()
	defer sm.Unlock()
	v, ok := sm.m[id]
	if ok {
		return v, ok
	}
	_, f, l, _ := runtime.Caller(1)
	log.Error("[LOBBY MAP GET]: Key not found ", f, l)
	return nil, ok
}

func (sm *L_SM) Set(id string, val *Lobby) {
	sm.Lock()
	defer sm.Unlock()
	sm.m[id] = val
}
