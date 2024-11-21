package sse

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"gopkg.in/antage/eventsource.v1"
)

type SSESender struct {
	ctx    context.Context
	cancel context.CancelFunc
	es     eventsource.EventSource
}

func NewEventSource(
	writeTimeout time.Duration,
	idleConnTimeout time.Duration,
	closeOnTimeout bool,
) *SSESender {
	ctx, cancel := context.WithCancel(context.Background())
	return &SSESender{
		ctx:    ctx,
		cancel: cancel,
		es: eventsource.New(
			&eventsource.Settings{
				Timeout:        writeTimeout,
				IdleTimeout:    idleConnTimeout,
				CloseOnTimeout: closeOnTimeout,
			}, nil),
	}
}

func (sse SSESender) GetES() eventsource.EventSource {
	return sse.es
}

func (sse SSESender) Close() {
	sse.es.Close() // Close consumers
	sse.cancel()   // Close producer
}

func (sse SSESender) Run() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-sse.ctx.Done():
				fmt.Println("shutdown es daemon")
				break
			case <-ticker.C:
				sse.es.SendEventMessage("tick", "tick-event", strconv.Itoa(rand.Intn(10)))
			}
		}
	}()
}
