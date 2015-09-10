package common

import (
	"sync"
)

type Session struct {
	data map[string]interface{}

	lock sync.Mutex
}

func NewSession() *Session {
	return &Session{data: make(map[string]interface{})}
}

func (p *Session) Get(key string) interface{} {
	return p.data[key]
}

func (p *Session) Set(key string, val interface{}) {
	p.lock.Lock()
	p.data[key] = val
	p.lock.Unlock()

	return
}
