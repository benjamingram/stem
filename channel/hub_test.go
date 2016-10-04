package channel

import "testing"

func TestRegisterChannelRequiresChannel(t *testing.T) {
	var ch Hub

	err := ch.RegisterChannel(nil, []string{"*"})

	if err == nil || err.Error() != "no channel specified" {
		t.Errorf("RegisterChannel(nil, []string{\"*\"}) = %v, want %v", err, "no channel specified")
	}
}

func TestRegisterChannelRequiresTopic(t *testing.T) {
	var ch Hub

	c := make(chan string)

	err := ch.RegisterChannel(&c, nil)

	if err == nil || err.Error() != "no topics specified" {
		t.Errorf("RegisterChannel(c, nil) = %v, want %v", err, "no topics specified")
	}
}

func TestSendMessageSendsToRegisteredChannel(t *testing.T) {
	var ch Hub

	c := make(chan string)
	ch.RegisterChannel(&c, []string{"*"})

	expected := "howdy doody"

	go func() {
		ch.SendMessage(expected, "c1")
	}()

	received := <-c

	if received != expected {
		t.Errorf("SendMessage(\"%v\", \"*\"); Received = %v, want = %v", expected, received, expected)
	}
}

func TestSendMessageSendsToMatchedTopic(t *testing.T) {
	var ch Hub

	c := make(chan string)
	ch.RegisterChannel(&c, []string{"c1"})

	expected := "howdy doody"

	go func() {
		ch.SendMessage(expected, "c1")
	}()

	received := <-c

	if received != expected {
		t.Errorf("SendMessage(\"%v\", \"c1\"); Received = %v, want = %v", expected, received, expected)
	}
}

func TestSendMessageDoesNotSendToUnmatchedTopic(t *testing.T) {
	var ch Hub

	c := make(chan string)
	ch.RegisterChannel(&c, []string{"c3"})

	go func() {
		ch.SendMessage("howdy doody", "c1")
		close(c)
	}()

	received := <-c

	if received != "" {
		t.Errorf("SendMessage(\"howdy doody\", \"c1\"); Received = %v, want = \"\"", received)
	}
}

func TestDeregisterChannelRemovesChannelFromReceivingMessages(t *testing.T) {
	var ch Hub

	c := make(chan string)

	ch.RegisterChannel(&c, []string{"c1"})
	ch.DeregisterChannel(&c)

	go func() {
		ch.SendMessage("howdy doody", "c1")
		close(c)
	}()

	received := <-c

	if received != "" {
		t.Errorf("SendMessage(\"howdy doody\", \"c1\"); Received = %v, want = \"\"", received)
	}
}
