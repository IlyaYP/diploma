package model

import (
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

/*
[
    {
        "order": "2377225624",
        "sum": 500,
        "processed_at": "2020-12-09T16:09:57+03:00"
    }
]

{
	"current": 500.5,
	"withdrawn": 42
}
*/

// Withdrawal keeps withdraw data.
type (
	Withdrawal struct {
		Order       uint64    `json:"order"`
		Sum         float64   `json:"sum"`
		ProcessedAt time.Time `json:"processed_at,omitempty"`
		User        string    `json:"-"`
	}
	Withdrawals []Withdrawal

	Balance struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}
)

func (w *Withdrawal) Bind(r *http.Request) error {
	if w.Order == 0 {
		return fmt.Errorf("Order is a required field")
	}
	if w.Sum == 0 {
		return fmt.Errorf("Sum is a required field")
	}
	return nil
}

func (*Withdrawal) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Withdrawals) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Balance) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GetLoggerContext enriches logger context with essential Order fields.
func (w *Withdrawal) GetLoggerContext(logCtx zerolog.Context) zerolog.Context {
	logCtx = logCtx.Str("login", w.User)
	logCtx = logCtx.Uint64("ordernum", w.Order)
	logCtx = logCtx.Float64("sum", w.Sum)
	return logCtx
}
