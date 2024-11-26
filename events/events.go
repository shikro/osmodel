package events

import (
	"context"
	"fmt"
	"time"
)

type PingEvent struct {
	ping chan struct{}
	freq time.Duration
}

func NewPingEvent(freq time.Duration) *PingEvent {
	return &PingEvent{
		ping: make(chan struct{}),
		freq: freq,
	}
}

func (e *PingEvent) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(e.freq)
				fmt.Println("event: ping")
				e.ping <- struct{}{}
			}
		}
	}()
}

func (e *PingEvent) Event() chan struct{} {
	return e.ping
}
