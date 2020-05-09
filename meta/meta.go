// meta 支持的数据类型: string bool int int32 int64 []int []string []int32 []int64
// 不支持任何浮点数
package meta

import (
	"context"
	"errors"
	"github.com/cbwfree/micro-core/srv"
	"github.com/micro/go-micro/v2/metadata"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrNoContext   = errors.New("no meta context")
	ErrInvalidMeta = errors.New("invalid service meta")
)

type NodeMeta interface {
	NodeId() string
	Context() context.Context
}

// Meta 基础类
type Meta struct {
	sync.RWMutex
	meta metadata.Metadata
}

func (m *Meta) Get(key string) string {
	m.RLock()
	defer m.RUnlock()

	return m.meta[key]
}

func (m *Meta) Int(key string) int {
	v, _ := strconv.Atoi(m.Get(key))
	return v
}

func (m *Meta) Int32(key string) int32 {
	v, _ := strconv.Atoi(m.Get(key))
	return int32(v)
}

func (m *Meta) Int64(key string) int64 {
	v, _ := strconv.ParseInt(m.Get(key), 10, 64)
	return v
}

func (m *Meta) Uint32(key string) uint32 {
	v, _ := strconv.ParseUint(m.Get(key), 10, 32)
	return uint32(v)
}

func (m *Meta) Uint64(key string) uint64 {
	v, _ := strconv.ParseUint(m.Get(key), 10, 64)
	return v
}

func (m *Meta) Bool(key string) bool {
	v, _ := strconv.ParseBool(m.Get(key))
	return v
}

func (m *Meta) Slice(key string) []string {
	m.RLock()
	defer m.RUnlock()

	val := m.meta[key]
	if val == "" {
		return nil
	}

	return strings.Split(val, ",")
}

// 创建新的 context
func (m *Meta) Context() context.Context {
	m.RLock()
	defer m.RUnlock()
	return metadata.NewContext(context.TODO(), m.meta)
}

// 获取meta数据
func (m *Meta) Data() metadata.Metadata {
	m.RLock()
	defer m.RUnlock()
	return m.meta
}

// 获取meta数据 (map[string]interface{})
func (m *Meta) Interface() map[string]interface{} {
	m.RLock()
	defer m.RUnlock()

	result := make(map[string]interface{}, len(m.meta))
	for k, v := range m.meta {
		result[k] = v
	}

	return result
}

// 更新单个属性值
func (m *Meta) SetValue(name string, value interface{}) {
	m.Lock()
	m.meta[name] = toMetaValue(value)
	m.Unlock()
}

// 设置属性值
func (m *Meta) SetValues(meta map[string]interface{}) {
	m.Lock()
	defer m.Unlock()

	for k, v := range meta {
		m.meta[k] = toMetaValue(v)
	}
}

func NewMeta(meta metadata.Metadata) *Meta {
	if meta == nil {
		meta = make(metadata.Metadata)
	}
	return &Meta{
		meta: meta,
	}
}

func FromMeta(ctx context.Context) (*Meta, error) {
	meta, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, ErrNoContext
	}
	return NewMeta(meta), nil
}

func toMetaValue(value interface{}) string {
	switch val := value.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case int:
		return strconv.Itoa(val)
	case int32:
		return strconv.Itoa(int(val))
	case int64:
		return strconv.FormatInt(val, 10)
	case bool:
		return strconv.FormatBool(val)
	case []string:
		return strings.Join(val, ",")
	case []int:
		var ss []string
		for _, v := range val {
			ss = append(ss, strconv.Itoa(v))
		}
		return strings.Join(ss, ",")
	case []int32:
		var ss []string
		for _, v := range val {
			ss = append(ss, strconv.Itoa(int(v)))
		}
		return strings.Join(ss, ",")
	case []int64:
		var ss []string
		for _, v := range val {
			ss = append(ss, strconv.FormatInt(v, 10))
		}
		return strings.Join(ss, ",")
	default:
		return ""
	}
}

// 检查节点是否有效
func CheckMetaNodeIsValid(name string, mt NodeMeta) bool {
	if mt.NodeId() == "" {
		return false
	}
	return srv.CheckServiceNode(name, mt.NodeId())
}
