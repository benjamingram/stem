package hosts

import (
	"log"
	"net"
	"net/http"
	"sync"
	"text/template"

	"github.com/benjamingram/stem"
	"github.com/benjamingram/stem/clients"
	"github.com/gorilla/mux"
)

// HostStatus is used to track hosts that should be running
type HostStatus struct {
	API       bool
	Console   bool
	WebSocket bool
}

// Host provides configuration for Host and ClientHosts
type Host struct {
	console   clients.Console
	webSocket clients.WebSocketHost
	api       API

	hub        *stem.ChannelHub
	hostStatus HostStatus
	listener   net.Listener
	waitGroup  sync.WaitGroup

	Addr          string
	APIAddr       string
	WebSocketAddr string
}

var homepageTemplate = template.Must(template.ParseFiles("hosts/launcher.html"))

// Start initiates listening for new requests
func (h *Host) Start(initialStatus HostStatus) {
	var ch stem.ChannelHub

	// Initialize hosts
	h.console = clients.Console{Hub: &ch}
	h.webSocket = clients.WebSocketHost{Addr: h.WebSocketAddr, Hub: &ch}
	h.api = API{Addr: h.APIAddr, Hub: &ch}

	// Initialize hosts' states
	h.syncHostStatuses(initialStatus)

	handler := h.mapRoutes()

	log.Println("Web Host Started -", h.Addr)

	http.ListenAndServe(h.Addr, handler)
}

func (h *Host) syncHostStatuses(hs HostStatus) {
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

func (h *Host) mapRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", h.homepageHandler)

	r.HandleFunc("/api/start", h.startAPIHandler)
	r.HandleFunc("/api/stop", h.stopAPIHandler)

	r.HandleFunc("/websocket/start", h.startWebSocketHandler)
	r.HandleFunc("/websocket/stop", h.stopWebSocketHandler)

	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/login-oauth", loginOauthHandler)
	r.HandleFunc("/login-oauth-response", loginOauthResponseHandler)

	r.HandleFunc("/console/start", h.startConsoleHandler)
	r.HandleFunc("/console/stop", h.stopConsoleHandler)

	return r
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

	data := struct {
		API           bool
		APIAddr       string
		Console       bool
		WebSocket     bool
		WebSocketAddr string
	}{
		API:           h.hostStatus.API,
		APIAddr:       h.APIAddr,
		Console:       h.hostStatus.Console,
		WebSocket:     h.hostStatus.WebSocket,
		WebSocketAddr: h.WebSocketAddr,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homepageTemplate.Execute(w, data)
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
