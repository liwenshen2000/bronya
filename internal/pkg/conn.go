package pkg

import (
	"fmt"
	"io"

	"github.com/gorilla/websocket"
)

type TCPWebsocketConn struct {
	*websocket.Conn
	reader io.Reader
}

func (w *TCPWebsocketConn) Read(p []byte) (n int, err error) {
	for {
		for w.reader == nil {
			t, r, e := w.NextReader()
			if e != nil {
				err = fmt.Errorf("[Websocket] NextReader Error: %v", e)
				return
			}
			if t != websocket.BinaryMessage {
				continue
			}
			w.reader = r
		}
		m, e := w.reader.Read(p)
		if m != 0 {
			return m, nil
		}
		if e == io.EOF {
			w.reader = nil
		}
	}
}

func (w *TCPWebsocketConn) Write(p []byte) (n int, err error) {
	if err = w.WriteMessage(websocket.BinaryMessage, p); err != nil {
		return
	}
	n += len(p)
	return
}

func NewTCPWebsocketConn(conn *websocket.Conn) *TCPWebsocketConn {
	return &TCPWebsocketConn{Conn: conn, reader: nil}
}
