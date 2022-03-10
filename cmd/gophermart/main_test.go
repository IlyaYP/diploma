package main

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/cmd/gophermart/config"
	"github.com/IlyaYP/diploma/model"
	"log"
	"testing"
)

//func TestGeneralTests(t *testing.T) {
//	tests := []struct {
//		name    string
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//		{"1", false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := GeneralTests(); (err != nil) != tt.wantErr {
//				t.Errorf("GeneralTests() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestGeneralTests(t *testing.T) {
	t.Run("functional test", func(t *testing.T) {
		wantErr := false
		if err := GeneralTests(); (err != nil) != wantErr {
			t.Errorf("GeneralTests() error = %v, wantErr %v", err, wantErr)
		}
	})
}

func GeneralTests() error {

	cfg, err := config.New()
	if err != nil {
		return err
	}
	ctx := context.Background()

	st, _ := cfg.BuildPsqlStorage(ctx)
	st.Destroy(ctx)
	st.Close()

	userSvc, err := cfg.BuildUserService(ctx)
	if err != nil {
		return err
	}

	fmt.Println()

	if _, err := userSvc.Register(ctx, "vasya", "God"); err != nil {
		return err
	}

	fmt.Println()

	if _, err := userSvc.Register(ctx, "kolya", "God"); err != nil {
		return err
	}

	fmt.Println()

	if user, err := userSvc.Login(ctx, "vasya", "God"); err != nil {
		return err
		log.Println(err)
	} else {
		fmt.Println(user)
	}

	fmt.Println()

	if user, err := userSvc.Login(ctx, "kolya", "God"); err != nil {
		return err
	} else {
		fmt.Println(user)
	}

	fmt.Println()

	orderSvc, err := cfg.BuildOrderService(ctx)
	if err != nil {
		return err
	}

	if _, err := orderSvc.CreateOrder(ctx, model.Order{Number: 12345678903,
		Status: model.OrderStatusNew,
		User:   "vasya"}); err != nil {
		return err
	}

	if _, err := orderSvc.CreateOrder(ctx, model.Order{Number: 12345678911,
		Status: model.OrderStatusNew,
		User:   "vasya"}); err != nil {
		return err
	}

	if _, err := orderSvc.CreateOrder(ctx, model.Order{Number: 12345678929,
		Status: model.OrderStatusNew,
		User:   "kolya"}); err != nil {
		return err
	}

	if _, err := orderSvc.CreateOrder(ctx, model.Order{Number: 12345678937,
		Status: model.OrderStatusNew,
		User:   "kolya"}); err != nil {
		return err
	}

	//if err := orderSvc.ProcessNewOrders(ctx); err != nil {
	//	log.Println(err)
	//}
	//
	//if err := orderSvc.ProcessOrders(ctx); err != nil {
	//	log.Println(err)
	//}

	if orders, err := orderSvc.GetOrdersByUser(ctx, "vasya"); err != nil {
		return err
	} else {
		fmt.Println(orders)
	}

	if err := orderSvc.NewWithdrawal(ctx, model.Withdrawal{Order: 22345678902,
		Sum:  100,
		User: "vasya"}); err != nil {
		log.Println(err)
	}

	if err := orderSvc.NewWithdrawal(ctx, model.Withdrawal{Order: 22345678910,
		Sum:  100,
		User: "vasya"}); err != nil {
		log.Println(err)
	}

	if err := orderSvc.NewWithdrawal(ctx, model.Withdrawal{Order: 22345678928,
		Sum:  2000,
		User: "vasya"}); err != nil {
		log.Println(err)
	}

	if err := orderSvc.NewWithdrawal(ctx, model.Withdrawal{Order: 22345678936,
		Sum:  100,
		User: "vasya"}); err != nil {
		log.Println(err)
	}

	if withdrawals, err := orderSvc.GetWithdrawalsByUser(ctx, "vasya"); err != nil {
		return err
	} else {
		fmt.Println(withdrawals)
	}

	if balance, err := orderSvc.GetBalanceByUser(ctx, "vasya"); err != nil {
		return err
	} else {
		log.Println(balance)
	}

	return nil
}
