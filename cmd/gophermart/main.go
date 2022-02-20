package main

import (
	"context"
	"fmt"
	userService "github.com/IlyaYP/diploma/service/user"
	userStorage "github.com/IlyaYP/diploma/storage/psql"
)

func main() {
	st, err := userStorage.New()
	if err != nil {
		panic(err)
	}
	userSvc, err := userService.New(userService.WithUserStorage(st))
	if err != nil {
		panic(err)
	}

	user, err := userSvc.CreateUser(context.Background(), "vasya", "God")
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
}
