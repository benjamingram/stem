package hosts

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/benjamingram/stem"
)

// API is used to specify configuration for the API Host
type API struct {
	listener  net.Listener
	waitGroup sync.WaitGroup

	Addr string
	Hub  *stem.ChannelHub
}

// Start begins listening for new requests
func (api *API) Start() {
	api.Stop()

	// Setup listener
	l, err := net.Listen("tcp", api.Addr)
	if err != nil {
		log.Fatal(err)
	}

	api.listener = l

	mux := http.NewServeMux()
	mux.HandleFunc("/", api.rootHandler)

	go func() {
		api.waitGroup.Add(1)
		defer api.waitGroup.Done()

		http.Serve(l, mux)
	}()

	log.Println("API Host Started -", api.Addr)
}

// Stop ends listening for new requests
func (api *API) Stop() {
	if api.listener == nil {
		return
	}

	log.Println("Stopping API Host...")

	api.listener.Close()

	api.waitGroup.Wait()

	api.listener = nil

	log.Println("API Host Stopped")
}

func (api *API) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	val, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.Hub.SendMessage(string(val), "*")

	w.WriteHeader(http.StatusOK)
}
