package main

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/cmd/gophermart/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	userSvc, err := cfg.BuildUserService()
	if err != nil {
		panic(err)
	}

	user, err := userSvc.CreateUser(context.Background(), "vasya", "God")
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
}
