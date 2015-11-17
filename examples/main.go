package main

import (
	"flag"
	"stem/web"
)

// Command Line Parameters
var webAddr = flag.String("web-addr", ":8877", "http web service address")

var initConsole = flag.Bool("console", false, "start console client")

var initAPI = flag.Bool("api", false, "start http API service")
var apiAddr = flag.String("api-addr", ":9988", "http api service address")

var initWebSocket = flag.Bool("websocket", false, "start http web socket service")
var webSocketAddr = flag.String("websocket-addr", ":7766", "web socket service address")

func main() {
	flag.Parse()

	hostStatus := web.HostStatus{API: *initAPI,
		Console:   *initConsole,
		WebSocket: *initWebSocket}

	web := web.Host{Addr: *webAddr,
		APIAddr:       *apiAddr,
		WebSocketAddr: *webSocketAddr}

	web.Start(hostStatus)
}
