package model

import (
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

/*
	{
        "number": "9278923470",
        "status": "PROCESSED",
        "accrual": 500,
        "uploaded_at": "2020-12-10T15:15:45+03:00"
    }
*/

// Order keeps order data.
type Order struct {
	Number     int         `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    int         `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
	user       string
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"        // NEW — заказ загружен в систему, но не попал в обработку;
	OrderStatusProcessing OrderStatus = "PROCESSING" // PROCESSING — вознаграждение за заказ рассчитывается;
	OrderStatusInvalid    OrderStatus = "INVALID"    // INVALID — система расчёта вознаграждений отказала в расчёте;
	OrderStatusProcessed  OrderStatus = "PROCESSED"  // PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
)

var (
	// orderStatusMap maps OrderStatus value to its int representation.
	orderStatusToIntMap = map[OrderStatus]int{
		OrderStatusNew:        1,
		OrderStatusProcessing: 2,
		OrderStatusInvalid:    3,
		OrderStatusProcessed:  4,
	}

	// orderStatusToStrMap maps OrderStatus value to its string representation.
	orderStatusToStrMap = map[int]OrderStatus{
		1: OrderStatusNew,
		2: OrderStatusProcessing,
		3: OrderStatusInvalid,
		4: OrderStatusProcessed,
	}
)

// NewOrderStatusFromInt returns OrderStatus by its int representation (might be invalid).
func NewOrderStatusFromInt(v int) OrderStatus {
	return orderStatusToStrMap[v]
}

// String implements the fmt.Stringer interface.
func (s OrderStatus) String() string {
	return string(s)
}

// Int returns enum value int representation.
func (s OrderStatus) Int() int {
	return orderStatusToIntMap[s]
}

// Validate validates enum value.
func (s OrderStatus) Validate() error {
	_, found := orderStatusToIntMap[s]
	if !found {
		return fmt.Errorf("unknown value: %v", s)
	}

	return nil
}

// GetLoggerContext enriches logger context with essential Order fields.
func (o Order) GetLoggerContext(logCtx zerolog.Context) zerolog.Context {
	logCtx = logCtx.Str("login", o.user)
	logCtx = logCtx.Int("ordernum", o.Number)
	return logCtx
}
