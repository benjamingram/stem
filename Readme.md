# Stem
A small foray into using the Go language

## Goal
Stem is a simplistic Web Host that takes input from an API and distributes the data across a modular system of clients.

## Supported Modules

### API
The API receives simple values as input through HTTP Post requests.

### WebSockets
The WebSockets Host is a Web UI for streaming the incoming results from the API.

### Console
The Console streams the input from the API data to os.Stderr
