package web

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/util/log"
	"net"
	"sync"
)

const (
	MetaClientId = "Ws-Client-Id" // 客户端ID
)

// 客户端连接
type SocketConn struct {
	sync.RWMutex
	id        string            // 客户端ID
	meta      metadata.Metadata // metadata
	conn      *websocket.Conn   // Socket连接
	writeChan chan []byte       // 写入消息缓冲
	isClose   bool              // 是否已关闭
	isLinger  bool              // 是否丢弃未发送的数据
}

// 获取客户端ID
func (s *SocketConn) Id() string {
	return s.id
}

// 获取客户端meta信息
func (s *SocketConn) Meta() metadata.Metadata {
	s.RLock()
	defer s.RUnlock()

	return s.meta
}

// 获取客户端meta信息
func (s *SocketConn) GetMeta(key string) string {
	s.RLock()
	defer s.RUnlock()

	return s.meta[key]
}

// 设置meta信息
func (s *SocketConn) SetMeta(key string, value string) {
	s.Lock()
	defer s.Unlock()

	s.meta[key] = value
}

// 获取客户端metadata信息
func (s *SocketConn) MetaData() context.Context {
	s.RLock()
	defer s.RUnlock()

	return metadata.NewContext(context.TODO(), s.meta)
}

// 获取客户端IP地址
func (s *SocketConn) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

// 获取服务器IP地址
func (s *SocketConn) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

// 判断是否关闭
func (s *SocketConn) IsClose() bool {
	return s.isClose
}

// 写入消息
func (s *SocketConn) Write(payload []byte) error {
	s.Lock()
	defer s.Unlock()
	if s.isClose {
		return nil
	}

	s.doWrite(payload)
	return nil
}

// 关闭连接
func (s *SocketConn) Close() {
	s.Lock()
	defer s.Unlock()
	if s.isClose {
		return
	}

	s.doWrite(nil)
	s.isClose = true
}

// 销毁连接
func (s *SocketConn) Destroy() {
	s.Lock()
	defer s.Unlock()

	s.doDestroy()
}

// 执行写入消息
func (s *SocketConn) doWrite(buf []byte) {
	s.writeChan <- buf
}

// 关闭操作
func (s *SocketConn) doDestroy() {
	_ = s.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	_ = s.conn.Close()

	if !s.isClose {
		close(s.writeChan)
		s.isClose = true
	}
}

// 实例化客户端连接
func NewSocketConn(conn *websocket.Conn) *SocketConn {
	s := new(SocketConn)
	s.conn = conn
	s.id = uuid.New().String()
	s.meta = map[string]string{
		MetaClientId: s.id,
	}
	s.writeChan = make(chan []byte)

	// 异步处理推送消息
	go func() {
		for b := range s.writeChan {
			if b == nil {
				break
			}

			err := conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				break
			}
		}

		_ = s.conn.Close()

		s.Lock()
		s.isClose = true
		s.Unlock()

		log.Debugf("[%s] SocketConn Write Chan is Closed ...", s.id)
	}()

	return s
}
