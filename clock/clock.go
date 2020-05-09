package clock

import (
	"sync"
	"time"
)

type Clock interface {
	Name() string
	Time() time.Duration
	Run()
	Stop()
	Reload(d ...time.Duration)
	WaitExit()
}

type ExitHandler func()

var (
	admin     *Admin
	adminOnce sync.Once
)

// 全局定时器管理
func Instance() *Admin {
	adminOnce.Do(func() {
		admin = NewAdmin("global")
	})
	return admin
}
