package web

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/micro/go-micro/v2/util/log"
	"net/http"
	"sync"
	"time"
)

// Receive WebSocket Message Event
type OnReceiveHandler func(s *Socket, sess *SocketConn, data []byte) error

// WebSocket Conn Close Event
type OnConnCloseHandler func(s *Socket, sess *SocketConn) error

type Socket struct {
	wg          sync.WaitGroup
	web         *Server             // Web Server
	path        string              // WebSocket Path
	upgrader    *websocket.Upgrader //
	conns       *SocketConns
	onReceive   OnReceiveHandler
	onConnClose OnConnCloseHandler
}

func (s *Socket) Web() *Server {
	return s.web
}

func (s *Socket) Path() string {
	return s.path
}

func (s *Socket) Conns() *SocketConns {
	return s.conns
}

// WebSocket Handler
func (s *Socket) Handler(c echo.Context) error {
	// 将HTTP请求升级为WebSocket
	conn, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	s.wg.Add(1)
	defer s.wg.Done()

	// 创建客户端连接对象
	sc := NewSocketConn(conn)
	s.conns.Put(sc)

	log.Debugf("[%s][%s] successfully connected...", sc.Id(), sc.RemoteAddr().String())

	// 接收消息处理
	for {
		// 接收消息
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("[%s] Read Message Failure: %s", sc.Id(), err.Error())
			break
		}

		if err := s.onReceive(s, sc, data); err != nil {
			log.Errorf("%s", err)
			break
		}
	}

	// 读消息失败后清理客户端
	sc.Destroy()
	s.conns.Del(sc.Id())     // 清理session
	_ = s.onConnClose(s, sc) // 连接断开处理

	log.Debugf("[%s][%s] disconnected ...", sc.Id(), sc.RemoteAddr().String())

	return nil
}

func (s *Socket) Close() {
	s.conns.Clean()
}

func NewSocket(web *Server, path string, timeout time.Duration) *Socket {
	ws := &Socket{
		web:   web,
		path:  path,
		conns: newSocketConns(),
		upgrader: &websocket.Upgrader{
			HandshakeTimeout: timeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}
	return ws
}
