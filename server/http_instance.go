package server

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"

	"github.com/zhhOceanfly/goplay"
)

type httpInstance struct {
	info      goplay.InstanceInfo
	hook      goplay.IServerHook
	ctrl      *goplay.InstanceCtrl
	transport goplay.ITransport

	tlsConfig  *tls.Config
	httpServer http.Server
	ws         *wsInstance
	sse        *sseInstance
}

func NewHttpInstance(name string, addr string, transport goplay.ITransport, hook goplay.IServerHook) (*httpInstance, error) {
	if transport == nil {
		return nil, errors.New("tcp instance transport must not be nil")
	}
	if hook == nil {
		return nil, errors.New("tcp instance server hook must not be nil")
	}
	return &httpInstance{info: goplay.InstanceInfo{Name: name, Address: addr, Type: TypeHttp}, transport: transport, hook: hook, ctrl: new(goplay.InstanceCtrl)}, nil
}

func (i *httpInstance) Run(listener net.Listener) error {
	i.httpServer.Handler = i
	if i.tlsConfig != nil {
		listener = tls.NewListener(listener, i.tlsConfig)
	}
	var err = i.httpServer.Serve(listener)
	return err
}

func (i *httpInstance) Close() {
	i.ctrl.WaitTask()
}

func (i *httpInstance) Info() goplay.InstanceInfo {
	return i.info
}

func (i *httpInstance) Hook() goplay.IServerHook {
	return i.hook
}

func (i *httpInstance) Transport() goplay.ITransport {
	return i.transport
}

func (i *httpInstance) Ctrl() *goplay.InstanceCtrl {
	return i.ctrl
}

func (i *httpInstance) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var request *goplay.Request
	var sess = goplay.NewSession(r.Context(), new(goplay.Conn), i)
	sess.Conn.Http.Request, sess.Conn.Http.ResponseWriter = r, w

	if i.ws != nil {
		if conn, _ := i.ws.update(w, r); conn != nil {
			sess.Server = i.ws
			sess.Conn.Websocket.WebsocketConn = conn
			i.ws.accept(sess)
			return
		}
	}

	if i.sse != nil {
		if err = i.sse.update(r); err == nil {
			sess.Server = i.sse
			i.sse.accept(sess)
			return
		}
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = fmt.Errorf("panic: %v\n%v", panicInfo, string(debug.Stack()))
		}
		i.hook.OnClose(sess, err)
	}()

	i.hook.OnConnect(sess, nil)
	request, err = i.transport.Receive(sess.Conn)
	err = doRequest(sess, request)
}

func (i *httpInstance) SetWSInstance(ws *wsInstance) {
	i.ws = ws
}

func (i *httpInstance) SetSSEInstance(sse *sseInstance) {
	i.sse = sse
}

func (i *httpInstance) WithCertificate(cert tls.Certificate) *httpInstance {
	if i.tlsConfig == nil {
		i.tlsConfig = &tls.Config{}
	}
	i.tlsConfig.Certificates = []tls.Certificate{cert}
	i.tlsConfig.Rand = rand.Reader
	return i
}
