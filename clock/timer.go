package clock

import (
	"github.com/micro/go-micro/v2/util/log"
	"time"
)

type TimerHandler func(t *Timer) error

// 定时器
type Timer struct {
	timer   *time.Timer   // 定时器
	name    string        // 名称
	delay   time.Duration // 延时时间
	handler TimerHandler  // 事件函数
	exit    chan bool     // 退出触发
}

// 获取定时器名称
func (t *Timer) Name() string {
	return t.name
}

// 获取定时器名称
func (t *Timer) String() string {
	return "Timer"
}

// 获取定时器延时时间
func (t *Timer) Time() time.Duration {
	return t.delay
}

// 启动定时器
func (t *Timer) Run() {
	t.timer = time.AfterFunc(t.delay, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("[Timer] the [%s] timer run panic: %s", t.name, err)
			}
			t.Stop()
		}()
		if err := t.handler(t); err != nil {
			log.Warn("[Timer] the [%s] timer run err: %s", t.name, err)
		}
	})
}

// 停止定时器
func (t *Timer) Stop() {
	if t.timer == nil {
		return
	}
	t.timer.Stop()
	t.timer = nil
	t.exit <- true
}

// 重载定时器
func (t *Timer) Reload(d ...time.Duration) {
	if t.timer == nil {
		return
	}
	if len(d) > 0 && d[0] > time.Duration(0) {
		t.delay = d[0]
	}
	t.timer.Reset(t.delay)
}

// 监听退出
func (t *Timer) WaitExit() {
	<-t.exit
}

// 新建定时器
func NewTimer(name string, delay time.Duration, handler TimerHandler) *Timer {
	t := &Timer{
		name:    name,
		delay:   delay,
		handler: handler,
		exit:    make(chan bool),
	}
	return t
}
