package model

import "context"

// Q: Какой из способов лучше?

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// userKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
//var userKey key

const (
	userKey = key(iota)
	orderKey
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "gmart context value " + k.name
}

var (
	UserCtxKey  = &contextKey{"user"}
	OrderCtxKey = &contextKey{"order"}
)

// UserNewContext returns a new Context that carries value u.
func UserNewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
	//return context.WithValue(ctx, UserCtxKey, u)
}

// UserFromContext returns the User value stored in ctx, if any.
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	//u, ok := ctx.Value(UserCtxKey).(*User)
	return u, ok
}

// OrderNewContext returns a new Context that carries value u.
func OrderNewContext(ctx context.Context, o *Order) context.Context {
	return context.WithValue(ctx, orderKey, o)
	//return context.WithValue(ctx, OrderCtxKey, o)
}

// OrderFromContext returns the User value stored in ctx, if any.
func OrderFromContext(ctx context.Context) (*Order, bool) {
	u, ok := ctx.Value(orderKey).(*Order)
	//u, ok := ctx.Value(OrderCtxKey).(*Order)
	return u, ok
}
