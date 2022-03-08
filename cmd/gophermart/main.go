package main

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/cmd/gophermart/config"
	"log"
)

func main() {
	run(false)
	//GeneralTests()

}
func run(test bool) error {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	//ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here

	srv, err := cfg.BuildServer(ctx)
	if err != nil {
		panic(err)
	}

	if !test {
		return srv.ListenAndServe()
	}
	return nil

}

func GeneralTests() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	ctx := context.Background()

	userSvc, err := cfg.BuildUserService(ctx)
	if err != nil {
		return err
	}

	fmt.Println()

	if _, err := userSvc.Register(ctx, "vasya", "God"); err != nil {
		log.Println(err)
	}

	fmt.Println()

	if _, err := userSvc.Register(ctx, "kolya", "God"); err != nil {
		log.Println(err)
	}

	fmt.Println()

	if user, err := userSvc.Login(ctx, "vasya", "God"); err != nil {
		log.Println(err)
	} else {
		fmt.Println(user)
	}

	fmt.Println()

	if user, err := userSvc.Login(ctx, "kolya", "God"); err != nil {
		log.Println(err)
	} else {
		fmt.Println(user)
	}

	fmt.Println()

	orderSvc, err := cfg.BuildOrderService(ctx)
	if err != nil {
		return err
	}

	fmt.Println()

	if orders, err := orderSvc.GetOrdersByUser(ctx, "vasya"); err != nil {
		log.Println(err)
	} else {
		fmt.Println(orders)
	}

	return nil
}
