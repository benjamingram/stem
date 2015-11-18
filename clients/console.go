package clients

import (
	"fmt"
	"github.com/benjamingram/stem"
	"log"
)

type Console struct {
	Hub     *stem.ChannelHub
	channel chan string
}

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

func (cc *Console) Stop() {
	// If the channel is already closed, nothing more to do
	if cc.channel == nil {
		return
	}

	cc.Hub.DeregisterChannel(&cc.channel)
	close(cc.channel)

	log.Println("Console Client Stopped")
}
