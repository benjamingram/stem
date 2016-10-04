package clients

import (
	"log"
	"net"
	"net/http"
	"sync"
	"text/template"

	"github.com/benjamingram/stem/channel"
	"github.com/benjamingram/stem/clients/websocket"
)

var (
	pageTemplate = template.Must(template.ParseFiles("clients/websocket/main.html"))
)

// WebSocketHost is a wrapper http server to host the websocket client UI
type WebSocketHost struct {
	listener  net.Listener
	logger    *log.Logger
	waitGroup sync.WaitGroup

	Addr string
	Hub  *channel.Hub
}

// Start the WebSocketHost listening for incoming requests
func (wh *WebSocketHost) Start() {
	c := make(chan string)

	var ws websocket.WebSocket

	wh.Hub.RegisterChannel(&c, []string{"*"})

	wh.logger = log.New(&ws, "", log.LstdFlags)

	l, err := net.Listen("tcp", wh.Addr)
	if err != nil {
		log.Fatal(err)
	}

	wh.listener = l

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHomepage)
	mux.HandleFunc("/ws", ws.HandleSocket)

	go func() {
		wh.waitGroup.Add(1)
		defer wh.waitGroup.Done()

		http.Serve(l, mux)
	}()

	go func() {
		for {
			message, more := <-c

			if more {
				wh.logger.Println(message)
			} else {
				wh.Hub.DeregisterChannel(&c)
			}
		}
	}()

	log.Println("Web Socket Host Started -", wh.Addr)
}

// Stop the WebSocketHost from taking new requests
func (wh *WebSocketHost) Stop() {
	if wh.listener == nil {
		return
	}

	log.Println("Stopping Web Socket Host...")
	wh.listener.Close()
	wh.waitGroup.Wait()
	wh.listener = nil
	log.Println("Web Socket Host Stopped")
}

func handleHomepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", 404)
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pageTemplate.Execute(w, r.Host)
}
