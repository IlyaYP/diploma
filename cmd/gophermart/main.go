package main

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/service/user"
)

func main() {
	userSvc, err := user.New()
	if err != nil {
		panic(err)
	}

	user, err := userSvc.CreateUser(context.Background(), "vasya", "God")
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
}
