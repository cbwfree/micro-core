package debug

import (
	"bytes"
	"fmt"
	"github.com/micro/go-micro/v2/util/log"
	"time"
)

// 性能分析
type Prof struct {
	label string
	t     time.Time
	start *CallerInfo
	end   *CallerInfo
}

func (p *Prof) Result() {
	p.end = GetCaller(3)
	var b = new(bytes.Buffer)
	b.WriteString(fmt.Sprintf("[Prof][%s] Runtime Prof Info, time: %s\n", p.label, time.Since(p.t)))
	b.WriteString(fmt.Sprintf(" -> Start [%s]: %s:%d\n", p.start.Func, p.start.File, p.start.Line))
	b.WriteString(fmt.Sprintf(" -> End   [%s]: %s:%d", p.end.Func, p.end.File, p.end.Line))
	log.Trace(b.String())
}

func NewProf(label ...string) *Prof {
	var key string
	if len(label) > 0 {
		key = label[0]
	} else {
		key = "Label"
	}
	return &Prof{
		label: key,
		t:     time.Now(),
		start: GetCaller(3),
	}
}
