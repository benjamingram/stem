package channel

import (
	"errors"
	"sync"
)

// Hub is responsible for piping messages to all registered channels
type Hub struct {
	sync.RWMutex
	channels map[*chan string]map[string]struct{}
}

// SendMessage publishes a message to all matching channels registered to the topic
func (ch *Hub) SendMessage(message string, topic string) {
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

// RegisterChannel registers the specified channel with the specified topics
func (ch *Hub) RegisterChannel(channel *chan string, topics []string) (err error) {
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

// DeregisterChannel removes the channel from the list of channels to write to
func (ch *Hub) DeregisterChannel(channel *chan string) {
	ch.Lock()
	defer ch.Unlock()

	if _, ok := ch.channels[channel]; ok {
		delete(ch.channels, channel)
	}
}
