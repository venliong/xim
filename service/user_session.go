package service

import (
	"container/list"
	"sync"
)

type UserSession struct {
	ID       string
	MsgChan  chan string
	messages *list.List
	lock     sync.RWMutex
}

func NewUserSession(id string) *UserSession {
	if id == "" {
		return nil
	}

	return &UserSession{
		ID:       id,
		MsgChan:  make(chan string, 0),
		messages: list.New()}
}

func (p *UserSession) Destroy() {
	p.messages.Init()
	close(p.MsgChan)
}

func (p *UserSession) PushMessage(message string) {
	select {
	case p.MsgChan <- message:
	default:
		p.lock.Lock()
		p.messages.PushBack(message)
		p.lock.Unlock()
	}
}

func (p *UserSession) HistoryMessage() string {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.messages.Len() == 0 {
		return ""
	}

	rst := "[" + p.messages.Front().Value.(string)
	for e := p.messages.Front().Next(); e != nil; e = e.Next() {
		rst += "," + e.Value.(string)

	}
	p.messages.Init()
	rst += "]"

	return rst
}
