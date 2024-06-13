package transport

import (
	"github.com/gorilla/websocket"
	"github.com/zhhOceanfly/goplay"
	"github.com/zhhOceanfly/goplay/binder"
	"github.com/zhhOceanfly/goplay/library/golang/json"
)

type WsJsonTransport struct {
}

func NewWsJsonTransport() *WsJsonTransport {
	return new(WsJsonTransport)
}

func (m *WsJsonTransport) Receive(c *goplay.Conn) (*goplay.Request, error) {
	var request goplay.Request
	request.Respond = true
	request.ActionName, request.Render = ParseHttpPath(c.Http.Request.URL.Path)

	if len(c.Websocket.Message) > 0 {
		request.InputBinder = binder.NewJsonBinder(c.Websocket.Message)
	} else {
		request.InputBinder = ParseHttpInput(c.Http.Request, 4096)
	}

	return &request, nil
}

func (m *WsJsonTransport) Response(c *goplay.Conn, res *goplay.Response) error {
	var err error
	var data []byte
	var messageType = c.Websocket.MessageType

	if messageType == 0 {
		messageType = websocket.TextMessage
	}

	if data, err = json.MarshalEscape(res.Output.All(), false, false); err != nil {
		return err
	}

	if err := c.Websocket.WebsocketConn.WriteMessage(messageType, data); err != nil {
		return err
	}

	return nil
}
