package model

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
)

// User keeps user data.
type User struct {
	Login    string `json:"login" yaml:"login"`
	Password string `json:"password" yaml:"password"`
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// userKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
var userKey key

// NewContext returns a new Context that carries value u.
func UserNewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the User value stored in ctx, if any.
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
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
