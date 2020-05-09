package multi

import (
	"bytes"
	"fmt"
	"github.com/cbwfree/micro-core/debug"
	"github.com/micro/go-micro/v2/util/log"
	"sync"
	"time"
)

type TaskHandler func() error

type Task struct {
	sync.Mutex
	wg     *sync.WaitGroup
	tasks  []TaskHandler
	errors *sync.Map
	level  int
}

func (mt *Task) Do(work ...TaskHandler) {
	mt.tasks = append(mt.tasks, work...)
}

// 执行并发任务
func (mt *Task) Run(timeout ...time.Duration) error {
	var trace = new(bytes.Buffer)
	var caller = debug.GetCaller(mt.level)

	mt.Lock()
	defer mt.Unlock()

	now := time.Now()
	wLen := len(mt.tasks)
	done := make(chan struct{})

	mt.errors = new(sync.Map)
	mt.wg = new(sync.WaitGroup)

	mt.wg.Add(wLen)

	trace.WriteString("[Task] ")
	if caller != nil {
		trace.WriteString(fmt.Sprintf("[%s] %s:%d, ", caller.Func, caller.File, caller.Line))
	}
	trace.WriteString(fmt.Sprintf("Start %d Task ...\n", wLen))

	for i, fn := range mt.tasks {
		go func(i int, fn TaskHandler) {
			defer mt.wg.Done()
			st := time.Now()
			err := fn()
			mt.errors.Store(i, err)
			if trace != nil {
				diffTime := time.Since(st)
				if err != nil {
					trace.WriteString(fmt.Sprintf(" -> Run %d task time: %s, Error: %v\n", i, diffTime, err))
				} else if diffTime > DefaultLongTime {
					trace.WriteString(fmt.Sprintf(" -> Run %d task time: %s\n", i, diffTime))
				}
			}
		}(i, fn)
	}

	go func() {
		mt.wg.Wait()
		done <- struct{}{}
	}()

	var afterTime time.Duration
	if len(timeout) > 0 {
		afterTime = timeout[0]
	} else {
		afterTime = DefaultTimeout
	}

	var err error
	select {
	case <-done:
		err = nil
	case <-time.After(afterTime):
		for i := 0; i < len(mt.tasks); i++ {
			_, ok := mt.errors.Load(i)
			if !ok && trace != nil {
				trace.WriteString(fmt.Sprintf(" -> Warning: Run %d task timeout\n", i))
			}
		}
		err = ErrTimeout
	}

	if trace != nil {
		trace.WriteString(fmt.Sprintf(" -> Total time spent: %s", time.Since(now)))
		log.Trace(trace)
	}

	return err
}

// 获取任务执行结果
func (mt *Task) Errors() []error {
	var errs []error
	for i := 0; i < len(mt.tasks); i++ {
		if err, ok := mt.errors.Load(i); ok {
			if err != nil {
				errs = append(errs, err.(error))
			} else {
				errs = append(errs, nil)
			}
		} else {
			errs = append(errs, ErrNotFound)
		}
	}
	return errs
}

// 实例化并发任务
func NewTasks(task ...TaskHandler) *Task {
	w := &Task{
		level: 3,
	}
	w.Do(task...)
	return w
}

// 执行并发任务
func RunTasks(tasks []TaskHandler, timeout ...time.Duration) error {
	w := &Task{
		level: 4,
	}
	w.Do(tasks...)
	return w.Run(timeout...)
}
