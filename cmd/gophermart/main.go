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
	//ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here

	userSvc, err := cfg.BuildUserService(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println()

	if _, err := userSvc.Register(ctx, "vasya4", "God"); err != nil {
		log.Println(err)
	}

	fmt.Println()

	if _, err := userSvc.Register(ctx, "vasya2", "God"); err != nil {
		log.Println(err)
	}

	fmt.Println()

	user, err := userSvc.Login(ctx, "vasya2", "God")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(user)

	user1, err1 := userSvc.Login(ctx, "vasya3", "God!")
	if err1 != nil {
		log.Println(err1)
	}
	fmt.Println(user1)

	//user2, err2 := userSvc.Login(ctx, "vasya2", "God!")
	//if err2 != nil {
	//	fmt.Println(err2)
	//}
	//fmt.Println(user2)

}
