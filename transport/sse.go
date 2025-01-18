package transport

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/zhhxbl/goplay"
	"github.com/zhhxbl/goplay/library/golang/json"
)

type SseTransport struct {
}

func NewSSETransport() *SseTransport {
	return new(SseTransport)
}

func (p *SseTransport) Receive(c *goplay.Conn) (*goplay.Request, error) {
	var request goplay.Request
	request.Respond = true
	request.ActionName, request.Render = ParseHttpPath(c.Http.Request.URL.Path)
	request.InputBinder = ParseHttpInput(c.Http.Request, 1024*4)
	request.Render = "json"
	return &request, nil
}

func (p *SseTransport) Response(c *goplay.Conn, res *goplay.Response) error {
	var err error
	var data []byte
	var w = c.Http.ResponseWriter
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = fmt.Errorf("panic: %v\n%v", panicInfo, string(debug.Stack()))
		}
	}()

	if res.Render != "json" {
		return errors.New("undefined " + res.Render + " sse response render")
	}
	if data, err = json.MarshalEscape(res.Output.All(), false, false); err != nil {
		return err
	}

	if _, err = fmt.Fprintf(w, "data: %s\n\n", string(data)); err != nil {
		return err
	}
	w.(http.Flusher).Flush()
	return err
}
