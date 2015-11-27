package main

import (
	"flag"

	"github.com/benjamingram/stem/hosts"
)

// Command Line Parameters
var webAddr = flag.String("web-addr", ":8877", "http web service address")

var initConsole = flag.Bool("console", false, "start console client")

var initAPI = flag.Bool("api", false, "start http API service")
var apiAddr = flag.String("api-addr", ":9988", "http api service address")

var initWebSocket = flag.Bool("websocket", false, "start http web socket service")
var webSocketAddr = flag.String("websocket-addr", ":7766", "web socket service address")

var secretFile = flag.String("secret-file", "control-panel-secret.json", "secret file for control panel")

func main() {
	flag.Parse()

	hostStatus := hosts.HostStatus{API: *initAPI,
		Console:   *initConsole,
		WebSocket: *initWebSocket}

	web := hosts.Host{Addr: *webAddr,
		APIAddr:       *apiAddr,
		WebSocketAddr: *webSocketAddr}

	hosts.LoadSecretsFromFile(*secretFile)
	web.Start(hostStatus)
}
