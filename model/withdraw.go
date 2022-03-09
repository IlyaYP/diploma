package model

import (
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

// Withdraw keeps withdraw data.
type (
	Withdraw struct {
		Order       uint64    `json:"order"`
		Sum         int       `json:"sum"`
		ProcessedAt time.Time `json:"processed_at"`
		User        string    `json:"-"`
	}
	Withdrawals []Withdraw

	Balance struct {
		Current   int `json:"current"`
		Withdrawn int `json:"withdrawn"`
	}
)

func (*Balance) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
