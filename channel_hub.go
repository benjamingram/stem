package stem

import (
	"errors"
	"sync"
)

type ChannelHub struct {
	sync.RWMutex
	channels map[*chan string]map[string]struct{}
}

func (ch *ChannelHub) SendMessage(message string, topic string) {
	ch.RLock()
	defer ch.RUnlock()

	for channel := range ch.channels {
		_, allTopics := ch.channels[channel]["*"]
		_, topicMatch := ch.channels[channel][topic]

		// If we do not have a topic match, skip the channel
		if !(allTopics || topicMatch) {
			continue
		}

		// Send the message
		*channel <- message
	}
}

func (ch *ChannelHub) RegisterChannel(channel *chan string, topics []string) (err error) {
	// Validate params
	if channel == nil {
		return errors.New("no channel specified")
	}

	if topics == nil || len(topics) == 0 {
		return errors.New("no topics specified")
	}

	// Maps are not thread-safe, let's make sure we are!
	ch.Lock()
	defer ch.Unlock()

	// Initialize channels map if it has not already been init
	if ch.channels == nil {
		ch.channels = make(map[*chan string]map[string]struct{})
	}

	// If channel is not already initialized, initialize topic list for channel
	if _, ok := ch.channels[channel]; !ok {
		ch.channels[channel] = make(map[string]struct{})
	}

	// Register each topic on the channel
	for _, topic := range topics {
		ch.channels[channel][topic] = struct{}{}
	}

	return nil
}

func (ch *ChannelHub) DeregisterChannel(channel *chan string) {
	ch.Lock()
	defer ch.Unlock()

	if _, ok := ch.channels[channel]; ok {
		delete(ch.channels, channel)
	}
}
