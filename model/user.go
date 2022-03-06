package model

import (
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
)

// User keeps user data.
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// GetLoggerContext enriches logger context with essential User fields.
func (u User) GetLoggerContext(logCtx zerolog.Context) zerolog.Context {
	return logCtx.Str("login", u.Login)
}

func (u *User) Bind(r *http.Request) error {
	if u.Login == "" {
		return fmt.Errorf("Login is a required field")
	}
	if u.Password == "" {
		return fmt.Errorf("Password is a required field")
	}
	return nil
}

func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
