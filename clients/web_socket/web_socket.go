package web_socket

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	pingPeriod   = (pongWait * 9) / 10
	pongWait     = 60 * time.Second
	readDeadline = 60 * time.Second
	readLimit    = 512
	writeWait    = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocket struct {
	ws *websocket.Conn
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(readLimit)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func write(ws *websocket.Conn, mt int, payload []byte) error {
	if ws == nil {
		return errors.New("no web socket connection: ws")
	}
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteMessage(mt, payload)
}

func (s *WebSocket) HandleSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	defer ws.Close()

	s.ws = ws

	reader(ws)
}

func (s WebSocket) Write(p []byte) (n int, err error) {
	err = write(s.ws, websocket.TextMessage, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}
