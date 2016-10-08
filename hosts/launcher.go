package hosts

import (
  "log"
  "net"
	"net/http"
	"sync"
	"text/template"

	"github.com/benjamingram/stem/channel"
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
	initialized bool
	console     clients.Console
	webSocket   clients.WebSocketHost
	api         API

	hub        *channel.Hub
	hostStatus HostStatus
	listener   net.Listener
	waitGroup  sync.WaitGroup

	Addr          string
	APIAddr       string
	WebSocketAddr string
}

var homepageTemplate = template.Must(template.New("launcherTemplate").Parse(launcherTemplate))

// Initialize initializes the host with initialStatus
func (h *Host) Initialize(initialStatus HostStatus) {
	var ch channel.Hub

	// Initialize hosts
	h.console = clients.Console{Hub: &ch}
	h.webSocket = clients.WebSocketHost{Addr: h.WebSocketAddr, Hub: &ch}
	h.api = API{Addr: h.APIAddr, Hub: &ch}

	// Initialize hosts' states
	h.syncHostStatuses(initialStatus)

	handler := h.mapRoutes()
	http.Handle("/", handler)

	h.initialized = true
	log.Println("Web Host Initialized")
}

// Start initiates listening for new requests
func (h *Host) Start() {
	log.Println("Web Host Started -", h.Addr)
	http.ListenAndServe(h.Addr, nil)
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

const launcherTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>stem</title>
  <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css" rel="stylesheet" integrity="sha256-MfvZlkHCEqatNoGiOXveE8FIwMzZg4W85qfrfIFBfYc= sha512-dTfge/zgoMYpP7QbHy4gWMEGsbsdZeCXz7irItjcC3sPUFtf0kuFbDz/ixG7ArTxmDjLXDmezHubeNikyKGVyQ==" crossorigin="anonymous">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap-theme.min.css" integrity="sha384-aUGj/X2zp5rLCbBxumKTCw2Z50WgIr1vs/PFN4praOTvYXWlVyh2UtNUU0KAUhAX" crossorigin="anonymous">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.4.0/css/font-awesome.min.css">
  <style type="text/css">
  html, body {
    height: 100%;
    background-color: #333;
  }

  header {
    background-color: #428bca;
    -webkit-box-shadow: inset 0 -2px 5px rgba(0,0,0,.1);
    box-shadow: inset 0 -2px 5px rgba(0,0,0,.1);
    text-shadow: 0 1px 3px rgba(0,0,0,.5);
    color: white;
    font-size: 1.3em;
    padding: 10px;
    margin-bottom: 20px;
  }
  .sub-title { color: #d9d9d9; font-size: .8em; }
  .panel { width: 200px; margin-left: 20px; }
  .panel-body { text-align: center; min-height: 100px; }
  .panel-body .host-location { display: block; margin-bottom: 10px; }
  </style>
</head>
<body>
    <header>
        <span>Stem</span>
        -
        <span class="sub-title">Server Host Control Panel</span>
    </header>

    <div class="fluid">
      <div class="panel panel-default pull-left">
        {{ if .API }}
        <div class="panel-heading">
          <h3 class="panel-title">
            API
            <span class="label label-success pull-right">Running</span>
          </h3>
        </div>
        <div class="panel-body">
          <span class="host-location">http://localhost{{ .APIAddr }}</span>
          <a class="btn btn-default" href="api/stop">
            Stop API
            &nbsp;
            <i class="fa fa-power-off"></i>
          </a>
        </div>
        {{ else }}
        <div class="panel-heading">
          <h3 class="panel-title">
            API
            <span class="label label-danger pull-right">Stopped</span>
          </h3>
        </div>
        <div class="panel-body">
          <span class="host-location">&nbsp;</span>
          <a class="btn btn-default" href="api/start">
            Start API
            &nbsp;
            <i class="fa fa-power-off"></i>
          </a>
        </div>
        {{ end }}
      </div>

      <div class="panel panel-default pull-left">
        {{ if .WebSocket }}
        <div class="panel-heading">
          <h3 class="panel-title">
            WebSocket
            <span class="label label-success pull-right">Running</span>
          </h3>
        </div>
        <div class="panel-body">
          <a class="host-location" href="http://localhost{{ .WebSocketAddr }}" target="_blank">http://localhost{{ .WebSocketAddr }}</a>
          <a class="btn btn-default" href="websocket/stop">
            Stop WebSocket
            &nbsp;
            <i class="fa fa-power-off"></i>
          </a>
        </div>
        {{ else }}
        <div class="panel-heading">
          <h3 class="panel-title">
            WebSocket
            <span class="label label-danger pull-right">Stopped</span>
          </h3>
        </div>
        <div class="panel-body">
          <span class="host-location">&nbsp;</span>
          <a class="btn btn-default" href="websocket/start">
            Start WebSocket
            &nbsp;
            <i class="fa fa-power-off"></i>
          </a>
        </div>
        {{ end }}
      </div>

      <div class="panel panel-default pull-left">
        {{ if .Console }}
        <div class="panel-heading">
          <h3 class="panel-title">
            Console
            <span class="label label-success pull-right">Running</span>
          </h3>
        </div>
        <div class="panel-body">
          <span class="host-location">&nbsp;</span>
          <a class="btn btn-default" href="console/stop">
            Stop Console
            &nbsp;
            <i class="fa fa-power-off"></i>
          </a>
        </div>
        {{ else }}
        <div class="panel-heading">
          <h3 class="panel-title">
            Console
            <span class="label label-danger pull-right">Stopped</span>
          </h3>
        </div>
        <div class="panel-body">
          <span class="host-location">&nbsp;</span>
          <a class="btn btn-default" href="console/start">
            Start Console
            &nbsp;
            <i class="fa fa-power-off"></i>
          </a>
        </div>
        {{ end }}
      </div>

    </div>
    <!-- <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js" integrity="sha256-Sk3nkD6mLTMOF0EOpNtsIry+s1CsaqQC1rVLTAy+0yc= sha512-K1qjQ+NcF2TYO/eI3M6v8EiNYZfA95pQumfvcVrTHtwQVDG+aHRqLi/ETn2uB+1JqwYqVG3LIvdm9lj6imS/pQ==" crossorigin="anonymous"></script> -->
</body>
</html>
`