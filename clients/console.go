package clients

import (
	"fmt"
	"log"

	"github.com/benjamingram/stem/channel"
)

// Console represents the client that sends output to the console
type Console struct {
	Hub     *channel.Hub
	channel chan string
}

// Start begins listening for new messages on the Hub
func (cc *Console) Start() {
	c := make(chan string)
	cc.channel = c

	cc.Hub.RegisterChannel(&cc.channel, []string{"*"})

	go func() {
		for {
			message, more := <-cc.channel

			if more {
				fmt.Println(message)
			} else {
				cc.channel = nil
				return
			}
		}
	}()

	log.Println("Console Client Started")
}

// Stop ends listening for new messages on the Hub
func (cc *Console) Stop() {
	// If the channel is already closed, nothing more to do
	if cc.channel == nil {
		return
	}

	cc.Hub.DeregisterChannel(&cc.channel)
	close(cc.channel)

	log.Println("Console Client Stopped")
}
