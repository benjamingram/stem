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
	pageTemplate = template.Must(template.New("webSocketHost").Parse(websocketHostTemplate))
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

const websocketHostTemplate = `
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
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
  <script type="text/javascript">
      $(function() {
          var conn;
          var log = $("#log");

          function appendLog(msg) {
              var d = log[0]
              var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
              msg.appendTo(log)
              if (doScroll) {
                  d.scrollTop = d.scrollHeight - d.clientHeight;
              }
          }

          if (window["WebSocket"]) {
              conn = new WebSocket("ws://{{$}}/ws");
              conn.onclose = function(evt) {
                  appendLog($("<div><b>Connection closed.</b></div>"))
              }
              conn.onmessage = function(evt) {
                  appendLog($("<div/>").text(evt.data))
              }
          } else {
              appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
          }
      });
  </script>
  <style type="text/css">
    html, body {
      height: 100%;
      background-color: #333;
      overflow: hidden;
    }

    header {
      background-color: #5cb85c;
      -webkit-box-shadow: inset 0 -2px 5px rgba(0,0,0,.1);
      box-shadow: inset 0 -2px 5px rgba(0,0,0,.1);
      text-shadow: 0 1px 3px rgba(0,0,0,.5);
      color: white;
      font-size: 1.3em;
      padding: 10px;
      margin-bottom: 20px;
    }
    .sub-title { color: #d9d9d9; font-size: .8em; }

    #log {
        color: white;
        margin: 0;
        padding: 0.5em 0.5em 0.5em 0.5em;

        overflow: auto;
    }
  </style>
</head>
<body>
    <header>
        <span>Stem</span>
        -
        <span class="sub-title">WebSocket Viewer</span>
    </header>

    <div class="fluid">
      <div id="log"></div>
    </div>
</body>
</html>
`
