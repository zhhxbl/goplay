package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/zhhxbl/goplay"
)

type TcpInstance struct {
	info      goplay.InstanceInfo
	hook      goplay.IServerHook
	ctrl      *goplay.InstanceCtrl
	transport goplay.ITransport
}

func NewTcpInstance(name string, addr string, transport goplay.ITransport, hook goplay.IServerHook) (*TcpInstance, error) {
	if transport == nil {
		return nil, errors.New("tcp instance transport must not be nil")
	}
	if hook == nil {
		return nil, errors.New("tcp instance server hook must not be nil")
	}
	return &TcpInstance{info: goplay.InstanceInfo{Name: name, Address: addr, Type: TypeTcp}, transport: transport, hook: hook, ctrl: new(goplay.InstanceCtrl)}, nil
}

func (i *TcpInstance) accept(s *goplay.Session) {
	var err error
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = fmt.Errorf("panic: %v\n%v", panicInfo, string(debug.Stack()))
		}
		_ = s.Conn.Tcp.Conn.Close()
		i.hook.OnClose(s, err)
	}()

	i.hook.OnConnect(s, nil)
	err = i.onReady(s)
}

func (i *TcpInstance) onReady(s *goplay.Session) (err error) {
	var n int
	var buffer = make([]byte, 4096)
	var request *goplay.Request
	var conn = s.Conn.Tcp.Conn

	for {
		sessContext := s.Context()
		select {
		case <-sessContext.Done():
			return sessContext.Err()
		default:
			if n, err = conn.Read(buffer); err != nil {
				return
			}
			s.Conn.Tcp.Surplus = append(s.Conn.Tcp.Surplus, buffer[:n]...)
			if true {
				if request, err = i.transport.Receive(s.Conn); err != nil {
					return
				}
				if request == nil {
					continue
				} else {
					s.Conn.Tcp.Version = request.Version
					if err = doRequest(s, request); err != nil {
						return err
					}
				}
			}
		}
	}
}

func (i *TcpInstance) Info() goplay.InstanceInfo {
	return i.info
}

func (i *TcpInstance) Hook() goplay.IServerHook {
	return i.hook
}

func (i *TcpInstance) Transport() goplay.ITransport {
	return i.transport
}

func (i *TcpInstance) Ctrl() *goplay.InstanceCtrl {
	return i.ctrl
}

func (i *TcpInstance) Run(listener net.Listener) error {
	for {
		if conn, err := listener.Accept(); err != nil {
			i.hook.OnConnect(goplay.NewSession(context.Background(), nil, i), err)
			continue
		} else {
			go func() {
				s := goplay.NewSession(context.Background(), new(goplay.Conn), i)
				s.Conn.Tcp.Conn = conn
				i.accept(s)
			}()
		}
	}
}

func (i *TcpInstance) Close() {
	i.ctrl.WaitTask()
}
