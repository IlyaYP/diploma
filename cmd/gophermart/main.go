package main

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/cmd/gophermart/config"
	"github.com/IlyaYP/diploma/pkg/logging"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here

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

	user1, err1 := userSvc.Login(ctx, "vasya3", "God!")
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(user1)

	user2, err2 := userSvc.Login(ctx, "vasya2", "God!")
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(user2)

}
