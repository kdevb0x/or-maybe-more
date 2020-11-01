package signal

import (
	"context"
	"fmt"
	"net/mail"
	"time"
)

type MessageType int
const (
	_ MessageType = iota
	Email
)

func sendNotification(target string, msgtype MessageType) error {
	switch msgtype {
	case Email:
		ctx := context.Background()
		timeOp(byEmail(ctx, target)
	}
}

func timeOp(op func()) time.Duration {
	var done = make(chan struct{})
	start := time.Now()
	go func(d chan<- struct{}) {
		op()
		d <- struct{}{}
	}
	_ := <-done
	return start.Sub(time.Now())

}

func byEmail(ctx context.Context, emailaddr string) {

}
