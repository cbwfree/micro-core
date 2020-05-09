package web

import (
	"github.com/cbwfree/micro-core/conv"
	"sync"
)

type SocketConns struct {
	sync.RWMutex
	conns map[string]*SocketConn
}

// 获取客户端列表
func (s *SocketConns) All() map[string]*SocketConn {
	s.RLock()
	defer s.RUnlock()

	conns := make(map[string]*SocketConn, len(s.conns))
	for k, v := range s.conns {
		conns[k] = v
	}

	return conns
}

// 连接总数
func (s *SocketConns) Count() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.conns)
}

// 添加新的客户端
func (s *SocketConns) Put(sc *SocketConn) {
	s.Lock()
	defer s.Unlock()

	s.conns[sc.Id()] = sc
}

// 获取客户端信息
func (s *SocketConns) Get(id string) *SocketConn {
	s.RLock()
	defer s.RUnlock()

	if sc, ok := s.conns[id]; ok {
		return sc
	}

	return nil
}

// 移除客户端连接
func (s *SocketConns) Del(id string) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.conns[id]; ok {
		delete(s.conns, id)
	}
}

// 通过Meta Key获取
func (s *SocketConns) GetByMeta(key string, value interface{}) []*SocketConn {
	s.RLock()
	defer s.RUnlock()

	var conns []*SocketConn

	valStr := conv.String(value)

	for _, sess := range s.conns {
		if sess.GetMeta(key) == valStr {
			conns = append(conns, sess)
		}
	}

	return conns
}

// 通过Meta Key删除
func (s *SocketConns) etByMeta(key string, value interface{}) {
	s.RLock()
	defer s.Unlock()

	valStr := conv.String(value)

	for id, sess := range s.conns {
		if sess.GetMeta(key) == valStr {
			delete(s.conns, id)
		}
	}
}

// 清空所有连接
func (s *SocketConns) Clean() {
	s.Lock()
	defer s.Unlock()

	for id, sc := range s.conns {
		sc.Destroy()
		delete(s.conns, id)
	}
}

// 实例化客户端连接管理器
func newSocketConns() *SocketConns {
	return &SocketConns{
		conns: make(map[string]*SocketConn),
	}
}
