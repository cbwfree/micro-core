package clock

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

var (
	defaultAdmin = NewAdmin("default")
)

func A() *Admin {
	return defaultAdmin
}

// 定时器管理
type Admin struct {
	sync.RWMutex // 锁
	name         string
	cron         *cron.Cron // 计划任务
	task         map[string]cron.EntryID
	clock        map[string]Clock // 定时器
}

// 定时器名称
func (ca *Admin) Name() string {
	return ca.name
}

// 计划任务
func (ca *Admin) Cron() *cron.Cron {
	return ca.cron
}

// 添加计划任务
func (ca *Admin) AddTask(name, spec string, task func()) error {
	ca.Lock()
	defer ca.Unlock()

	id, err := ca.cron.AddFunc(spec, task)
	if err != nil {
		return err
	}

	ca.task[name] = id

	return nil
}

// 添加计划任务
func (ca *Admin) AddJob(name, spec string, job cron.Job) error {
	ca.Lock()
	defer ca.Unlock()

	id, err := ca.cron.AddJob(spec, job)
	if err != nil {
		return err
	}

	ca.task[name] = id

	return nil
}

// 检查任务是否有效
func (ca *Admin) CheckTaskValid(name string) bool {
	ca.RLock()
	defer ca.RUnlock()

	if tId, ok := ca.task[name]; ok {
		return ca.cron.Entry(tId).Valid()
	}

	return false
}

// 删除计划任务
func (ca *Admin) DelTask(name string) {
	ca.RLock()
	defer ca.RUnlock()

	if tId, ok := ca.task[name]; ok {
		ca.cron.Remove(tId)
		delete(ca.task, name)
	}
}

// 新建延时器
func (ca *Admin) NewTimer(name string, delay time.Duration, handler TimerHandler, exit ...ExitHandler) *Timer {
	ca.Lock()
	defer ca.Unlock()

	t := NewTimer(name, delay, handler)

	go func() {
		defer func() {
			ca.DelClock(t.Name())
		}()

		t.WaitExit()
		log.Debugf("[Timer][%s] is run finish ...", t.Name())

		if len(exit) > 0 {
			exit[0]()
		}
	}()

	// 注册
	ca.clock[t.Name()] = t

	t.Run()

	return t
}

// 新建断续器(不限执行次数)
func (ca *Admin) NewTicker(name string, interval time.Duration, handler TickerHandler, exit ...ExitHandler) *Ticker {
	ca.Lock()
	defer ca.Unlock()

	t := NewTicker(name, interval, 0, handler)

	// 监听退出
	go func() {
		defer func() {
			ca.DelClock(t.Name())
		}()

		t.WaitExit()
		log.Debugf("[Ticker][%s] is run finish, executed %d, max execute: %d ...", t.Name(), t.NowNum(), t.MaxNum())

		if len(exit) > 0 {
			exit[0]()
		}
	}()

	// 注册
	ca.clock[t.Name()] = t

	t.Run()

	return t
}

// 新建断续器 (限制最大可执行次数)
func (ca *Admin) NewMaxTicker(name string, interval time.Duration, max int64, handler TickerHandler, exit ...ExitHandler) *Ticker {
	ca.Lock()
	defer ca.Unlock()

	t := NewTicker(name, interval, max, handler)

	// 监听退出
	go func() {
		defer func() {
			ca.DelClock(t.Name())
		}()

		t.WaitExit()
		log.Debugf("[Ticker][%s] is run finish, executed %d, max execute: %d ...", t.Name(), t.NowNum(), t.MaxNum())

		if len(exit) > 0 {
			exit[0]()
		}
	}()

	// 注册
	ca.clock[t.Name()] = t

	t.Run()

	return t
}

// 获取定时器
func (ca *Admin) GetClock(name string) Clock {
	if c, ok := ca.clock[name]; ok {
		return c
	}
	return nil
}

// 获取延时器
func (ca *Admin) GetTimer(name string) *Timer {
	c := ca.GetClock(name)
	if c == nil {
		return nil
	}
	if t, ok := c.(*Timer); ok {
		return t
	}
	return nil
}

// 获取断续器
func (ca *Admin) GetTicker(name string) *Ticker {
	c := ca.GetClock(name)
	if c == nil {
		return nil
	}
	if t, ok := c.(*Ticker); ok {
		return t
	}
	return nil
}

// 删除定时器
func (ca *Admin) DelClock(id string) {
	ca.Lock()
	defer ca.Unlock()

	delete(ca.clock, id)
}

// 销毁所有定时器
func (ca *Admin) CleanAll() {
	ca.Lock()
	defer ca.Unlock()

	// 移除所有计划任务
	for _, v := range ca.task {
		ca.cron.Remove(v)
	}

	// 等待计划任务执行完毕
	<-ca.cron.Stop().Done()

	for _, c := range ca.clock {
		c.Stop()
	}

	ca.clock = make(map[string]Clock)
}

// 实例化信息的定时器管理
func NewAdmin(name string) *Admin {
	a := &Admin{
		name:  name,
		cron:  cron.New(),
		task:  make(map[string]cron.EntryID),
		clock: make(map[string]Clock),
	}
	a.cron.Start()

	return a
}
