package clock

import (
	"github.com/micro/go-micro/v2/util/log"
	"time"
)

type TickerHandler func(t *Ticker) error

// 定时器
type Ticker struct {
	reload   bool          // 重载
	ticker   *time.Ticker  // 断续器
	name     string        // 名称
	interval time.Duration // 时间间隔
	handler  TickerHandler // 事件函数
	nowNum   int64         // 当前计数
	maxNum   int64         // 最大计数
	exit     chan bool     // 退出触发
}

// 获取定时器名称
func (t *Ticker) Name() string {
	return t.name
}

// 获取定时器名称
func (t *Ticker) String() string {
	return "Ticker"
}

// 获取定时器执行时间间隔
func (t *Ticker) Time() time.Duration {
	return t.interval
}

// 当前执行次数
func (t *Ticker) NowNum() int64 {
	return t.nowNum
}

// 最大执行次数
func (t *Ticker) MaxNum() int64 {
	return t.maxNum
}

// 启动定时器
func (t *Ticker) Run() {
	if t.ticker != nil {
		return
	}
	if t.interval <= time.Duration(0) {
		return
	}

	t.ticker = time.NewTicker(t.interval)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("[Ticker][%d] the [%s] ticker run panic: %s", t.nowNum, t.name, err)
			}
			t.exit <- true
		}()
		for range t.ticker.C {
			t.nowNum++
			if t.nowNum > t.maxNum && t.maxNum != 0 {
				t.Stop()
				return
			}
			if err := t.handler(t); err != nil {
				log.Warn("[Ticker][%d] the [%s] ticker run err: %s", t.nowNum, t.name, err)
			}
			if t.reload {
				t.Stop()
				t.reload = false
				t.Run()
			}
		}
	}()
}

// 停止定时器
func (t *Ticker) Stop() {
	if t.ticker == nil {
		return
	}
	t.ticker.Stop()
	t.ticker = nil
}

// 重载定时器
func (t *Ticker) Reload(d ...time.Duration) {
	if t.ticker == nil {
		return
	}
	if len(d) > 0 && d[0] > time.Duration(0) {
		t.interval = d[0]
	}
	t.nowNum = 0
	t.reload = true
}

// 监听退出
func (t *Ticker) WaitExit() {
	<-t.exit
	close(t.exit)
}

// 新建定时器
func NewTicker(name string, interval time.Duration, maxNum int64, handler TickerHandler) *Ticker {
	t := &Ticker{
		name:     name,
		interval: interval,
		handler:  handler,
		maxNum:   maxNum,
		exit:     make(chan bool),
	}
	return t
}
