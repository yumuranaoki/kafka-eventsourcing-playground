package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	brokers             = []string{"localhost:9092"}
	group   goka.Group  = "collector"
	table   goka.Table  = goka.GroupTable(group)
	stream  goka.Stream = "message"
)

func main() {
	cb := func(ctx goka.Context, msg interface{}) {
		var state string

		if val := ctx.Value(); val != nil {
			fmt.Printf("val: %s", val)
			state = val.(string)
		}

		state = msg.(string)

		ctx.SetValue(state)
		log.Printf("key = %s, state = %v, msg = %v", ctx.Key(), state, msg)
	}

	g := goka.DefineGroup(group,
		goka.Input(stream, new(codec.String), cb),
		goka.Persist(new(codec.String)),
	)

	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)
	go func() {
		defer close(done)
		if err = p.Run(ctx); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait
	cancel()
	<-done
}
