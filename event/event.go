package event

import (
	"errors"
	"sync"
)

var (
	ErrMaxCount = errors.New("run max count")
	ErrNotFound = errors.New("not found event")
	ErrExists   = errors.New("event is exists")
)

// Handler 事件函数
type Handler func(data interface{}) error

type Event struct {
	sync.Mutex
	name     string
	handler  Handler
	nowCount int64
	maxCount int64
}

func (e *Event) Name() string {
	return e.name
}

func (e *Event) NowCount() int64 {
	return e.nowCount
}

func (e *Event) MaxCount() int64 {
	return e.maxCount
}

// IsAllowExec 是否可以执行
func (e *Event) IsAllowExec() bool {
	if e.maxCount == 0 {
		return true
	}
	return e.nowCount < e.maxCount
}

func (e *Event) Call(obj interface{}) error {
	if e.maxCount > 0 && e.nowCount >= e.maxCount {
		return ErrMaxCount
	}

	e.Lock()
	defer e.Unlock()

	e.nowCount++

	return e.handler(obj)
}

// NewEvent 实例化事件对象
func NewEvent(name string, handler Handler, maxCount ...int64) *Event {
	e := &Event{
		name:    name,
		handler: handler,
	}
	if len(maxCount) > 0 {
		e.maxCount = maxCount[0]
	}
	return e
}
