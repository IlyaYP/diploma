package main

import (
	"context"
	"github.com/IlyaYP/diploma/cmd/gophermart/config"
	"os/signal"
	"syscall"
)

func main() {
	run(false)
}
func run(test bool) error {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	srv, err := cfg.BuildServer(ctx)
	if err != nil {
		panic(err)
	}

	if test {
		return nil
	}

	go srv.Serve(ctx)

	<-ctx.Done()
	srv.Close(ctx)

	// closing everything, do not care about errors
	for _, c := range cfg.Closer {
		c.Close()
	}

	return nil

}
