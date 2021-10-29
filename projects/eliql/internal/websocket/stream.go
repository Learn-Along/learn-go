package websocket

import (
	"github.com/Learn-Along/learn-go/projects/eliql/internal"
	"github.com/gorilla/websocket"
)

/*
The stream of map[string]interface{} that comes from the websocket URL
*/
type WebsocketStream struct {
	conn *websocket.Conn
	Url  string
}

func NewWebsocketStream(url string) *WebsocketStream {
	return &WebsocketStream{Url: url}
}

func (w *WebsocketStream) Start(recv chan internal.Record) error {
	return nil
}

func (w *WebsocketStream) Close() error {
	return nil
}

func (w *WebsocketStream) Equals(other interface{}) bool {
	otherStream, ok := other.(WebsocketStream)
	if !ok {
		return false
	}

	return otherStream.Url == w.Url
}
