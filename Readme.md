# Stem
A small foray into using the Go language

## Goal
Stem is a simplistic Web Host that takes input from an API and distributes the data across a modular system of clients.

## Installation
```
$ go get github.com/benjamingram/stem
```

## Run Control panel
```
$ cd `go list -f '{{.Dir}}' github.com/benjamingram/stem`
$ go run examples/main.go
```

## Launch Servers
1. Open the Control Panel [http://localhost:8877](http://localhost:8877)
1. Click "Start API"
1. Click "Start Console"
1. In a new terminal
```
$ curl --data "Hello Stem" http://localhost:9988
```

## Supported Modules
### Web Host
The host provides a Web UI Control panel for turning on and off different supported modules.

### API
The API receives simple values as input through HTTP Post requests.

### WebSockets
The WebSockets Host is a Web UI for streaming the incoming results from the API.

### Console
The Console streams the input from the API data to os.Stderr
