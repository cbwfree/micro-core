package event

import (
	"sync"
)

// 事件管理器
type Events struct {
	sync.RWMutex
	name   string
	events map[string]*Event // 事件列表
}

func (e *Events) String() string {
	return e.name
}

func (e *Events) Exist(name string) bool {
	e.RLock()
	defer e.RUnlock()
	_, ok := e.events[name]
	return ok
}

// Bind 绑定事件
func (e *Events) Bind(name string, handler Handler, maxCount ...int64) error {
	if e.Exist(name) {
		return ErrExists
	}

	e.Lock()
	e.events[name] = NewEvent(name, handler, maxCount...)
	e.Unlock()

	return nil
}

// Bind 绑定事件 (仅执行一次)
func (e *Events) Once(name string, handler Handler) error {
	if e.Exist(name) {
		return ErrExists
	}

	e.Lock()
	e.events[name] = NewEvent(name, handler, 1)
	e.Unlock()

	return nil
}

// Delete 删除事件
func (e *Events) Delete(name string) {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.events[name]; ok {
		delete(e.events, name)
	}
}

// Emit 执行事件
func (e *Events) Call(name string, data interface{}) error {
	e.RLock()
	h, ok := e.events[name]
	e.RUnlock()

	if !ok {
		return ErrNotFound
	}

	if !h.IsAllowExec() {
		e.Delete(name)
		return ErrMaxCount
	}

	return h.Call(data)
}

func NewEvents(name string) *Events {
	return &Events{
		name:   name,
		events: make(map[string]*Event),
	}
}
