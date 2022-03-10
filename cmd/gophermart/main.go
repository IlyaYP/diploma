package main

import (
	"context"
	"github.com/IlyaYP/diploma/cmd/gophermart/config"
)

func main() {
	run(false)
}
func run(test bool) error {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	srv, err := cfg.BuildServer(ctx)
	if err != nil {
		panic(err)
	}

	if !test {
		return srv.ListenAndServe()
	}
	return nil

}
