package main

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/cmd/gophermart/config"
	"github.com/IlyaYP/diploma/model"
	//"github.com/IlyaYP/diploma/provider/accrual/http"
	http "github.com/IlyaYP/diploma/provider/accrual/mock"
	"github.com/IlyaYP/diploma/service/order"
	"github.com/IlyaYP/diploma/service/user"
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
	st.Migrate(ctx)
	defer st.Close()

	//userSvc, err := cfg.BuildUserService(ctx)
	//if err != nil {
	//	return err
	//}

	// Build User Service
	userSvc, err := user.New(
		user.WithConfig(cfg.UserService),
		user.WithUserStorage(st),
	)
	if err != nil {
		return fmt.Errorf("building user service: %w", err)
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
		log.Println(err)
		return err
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

	//orderSvc, err := cfg.BuildOrderService(ctx)
	//if err != nil {
	//	return err
	//}

	// BuildAccrualProvider builds Accrual Provider dependency
	//accPr, err := cfg.BuildAccrualProvider(ctx)
	//if err != nil {
	//	return fmt.Errorf("accrual provider: %w", err)
	//}
	c := http.NewDefaultConfig()
	c.AccrualAddress = cfg.AccrualHTTPProvider.AccrualAddress

	accPr, err := http.New(http.WithConfig(&c))
	if err != nil {
		return fmt.Errorf("building provider: %w", err)
	}

	// Build Order Service
	orderSvc, err := order.New(
		order.WithOrderStorage(st),
		order.WithAccrualProvider(accPr),
	)
	if err != nil {
		return fmt.Errorf("building order service: %w", err)
	}

	if _, err := orderSvc.CreateOrder(ctx, model.Order{Number: "12345678903",
		Status: model.OrderStatusNew,
		User:   "vasya"}); err != nil {
		return err
	}

	if _, err := orderSvc.CreateOrder(ctx, model.Order{Number: "12345678911",
		Status: model.OrderStatusNew,
		User:   "vasya"}); err != nil {
		return err
	}

	if _, err := orderSvc.CreateOrder(ctx, model.Order{Number: "12345678929",
		Status: model.OrderStatusNew,
		User:   "kolya"}); err != nil {
		return err
	}

	if _, err := orderSvc.CreateOrder(ctx, model.Order{Number: "12345678937",
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

	if err := orderSvc.NewWithdrawal(ctx, model.Withdrawal{Order: "22345678902",
		Sum:  100,
		User: "vasya"}); err != nil {
		log.Println(err)
	}

	if err := orderSvc.NewWithdrawal(ctx, model.Withdrawal{Order: "22345678910",
		Sum:  100,
		User: "vasya"}); err != nil {
		log.Println(err)
	}

	if err := orderSvc.NewWithdrawal(ctx, model.Withdrawal{Order: "22345678928",
		Sum:  2000,
		User: "vasya"}); err != nil {
		log.Println(err)
	}

	if err := orderSvc.NewWithdrawal(ctx, model.Withdrawal{Order: "22345678936",
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
