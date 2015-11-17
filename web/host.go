package web

import (
	"log"
	"net"
	"net/http"
	"stem"
	"stem/clients"
	"sync"
	"text/template"
)

type HostStatus struct {
    API bool
    Console bool
    WebSocket bool
}

type Host struct {
    console clients.Console
    webSocket clients.WebSocketHost
    api stem.API

    hub *stem.ChannelHub
    hostStatus HostStatus
    listener net.Listener
    waitGroup sync.WaitGroup

	Addr string
    APIAddr string
    WebSocketAddr string
}

var homepageTemplate = template.Must(template.ParseFiles("web/home.html"))

func (h *Host) Start(initialStatus HostStatus) {
    var ch stem.ChannelHub

    // Initialize hosts
    h.console = clients.Console { Hub: &ch }
    h.webSocket = clients.WebSocketHost { Addr: h.WebSocketAddr, Hub: &ch}
    h.api = stem.API { Addr: h.APIAddr, Hub: &ch }

    // Initialize hosts' states
    h.syncHostStatuses(initialStatus)

    mux := http.NewServeMux()

		h.mapRoutes(mux)

    log.Println("Web Host Started -", h.Addr)

    http.ListenAndServe(h.Addr, mux)
}

func (h *Host) syncHostStatuses(hs HostStatus){
    if h.hostStatus.API != hs.API {
        if hs.API {
            h.api.Start()
        } else {
            h.api.Stop()
        }

        h.hostStatus.API = hs.API
    }

    if h.hostStatus.Console != hs.Console {
        if hs.Console {
            h.console.Start()
        } else {
            h.console.Stop()
        }

        h.hostStatus.Console = hs.Console
    }

    if h.hostStatus.WebSocket != hs.WebSocket {
        if hs.WebSocket {
            h.webSocket.Start()
        } else {
            h.webSocket.Stop()
        }

        h.hostStatus.WebSocket = hs.WebSocket
    }
}

func (h *Host) mapRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.homepageHandler)

	mux.HandleFunc("/api/start", h.startAPIHandler)
	mux.HandleFunc("/api/stop", h.stopAPIHandler)

	mux.HandleFunc("/websocket/start", h.startWebSocketHandler)
	mux.HandleFunc("/websocket/stop", h.stopWebSocketHandler)

	mux.HandleFunc("/console/start", h.startConsoleHandler)
	mux.HandleFunc("/console/stop", h.stopConsoleHandler)
}

func (h *Host) homepageHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.Error(w, "Page not found", 404)
		return
  }

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homepageTemplate.Execute(w, h.hostStatus)
}

func (h *Host) startAPIHandler(w http.ResponseWriter, r *http.Request) {
	s := h.hostStatus
	s.API = true
	h.syncHostStatuses(s)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Host) stopAPIHandler(w http.ResponseWriter, r *http.Request) {
	s := h.hostStatus
	s.API = false
	h.syncHostStatuses(s)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Host) startWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	s := h.hostStatus
	s.WebSocket = true
	h.syncHostStatuses(s)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Host) stopWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	s := h.hostStatus
	s.WebSocket = false
	h.syncHostStatuses(s)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Host) startConsoleHandler(w http.ResponseWriter, r *http.Request) {
	s := h.hostStatus
	s.Console = true
	h.syncHostStatuses(s)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Host) stopConsoleHandler(w http.ResponseWriter, r *http.Request) {
	s := h.hostStatus
	s.Console = false
	h.syncHostStatuses(s)
	http.Redirect(w, r, "/", http.StatusFound)
}
