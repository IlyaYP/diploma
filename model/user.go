package model

import "github.com/rs/zerolog"

// User keeps user data.
type User struct {
	Login    string `json:"login" yaml:"login"`
	Password string `json:"password" yaml:"password"`
}

// GetLoggerContext enriches logger context with essential User fields.
func (u User) GetLoggerContext(logCtx zerolog.Context) zerolog.Context {
	return logCtx.Str("login", u.Login)
}
