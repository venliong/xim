package xim

import (
	"container/list"
	"sync"
)

type User struct {
	ID       string
	MsgChan  chan string
	messages *list.List
	lock     sync.RWMutex
}

func NewUser(userid string) *User {
	if userid == "" {
		return nil
	}

	return &User{
		ID:       userid,
		MsgChan:  make(chan string, 0),
		messages: list.New()}
}

func (p *User) Destroy() {
	p.messages.Init()
	close(p.MsgChan)
}

func (p *User) PushMessage(message string) {
	select {
	case p.MsgChan <- message:
	default:
		p.lock.Lock()
		p.messages.PushBack(message)
		p.lock.Unlock()
	}
}

func (p *User) HistoryMessage() string {
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
