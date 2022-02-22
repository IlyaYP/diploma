package main

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/cmd/gophermart/config"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	userSvc, err := cfg.BuildUserService(ctx)
	if err != nil {
		panic(err)
	}

	if _, err := userSvc.CreateUser(ctx, "vasya2", "God"); err != nil {
		log.Println(err)
	}

	user, err := userSvc.Login(ctx, "vasya2", "God")
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
}
