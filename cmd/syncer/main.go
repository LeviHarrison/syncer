package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/leviharrison/syncer/terminal"
	"github.com/pkg/errors"
)

func main() {
	t, err := terminal.Init()
	if err != nil {
		fmt.Printf("%v\n", errors.Wrap(err, "Initializing terminal."))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go waitCancel(cancel)

	var wg sync.WaitGroup
	print := make(chan string)

	wg.Add(1)
	go t.Run(ctx, wg, print)

	wg.Wait()
}

func waitCancel(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	cancel()
}
