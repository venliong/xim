package common

import (
	"container/list"
	"sync"
)

type UserMessage struct {
	ID       string
	MsgChan  chan string
	messages *list.List
	lock     sync.RWMutex
}

func NewUserMessage(userid string) *UserMessage {
	if userid == "" {
		return nil
	}

	return &UserMessage{
		ID:       userid,
		MsgChan:  make(chan string, 0),
		messages: list.New()}
}

func (p *UserMessage) Destroy() {
	p.messages.Init()
	close(p.MsgChan)
}

func (p *UserMessage) PushMessage(message string) {
	select {
	case p.MsgChan <- message:
	default:
		p.lock.Lock()
		p.messages.PushBack(message)
		p.lock.Unlock()
	}
}

func (p *UserMessage) HistoryMessage() string {
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
