package pkg

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidLogin    = errors.New("invalid login")
	ErrInvalidPassword = errors.New("wrong password")
	ErrAlreadyExists   = errors.New("object exists in the DB")
	ErrNotExists       = errors.New("object not exists in the DB")
)

/*
200 — номер заказа уже был загружен этим пользователем;
202 — новый номер заказа принят в обработку;
400 — неверный формат запроса;
401 — пользователь не аутентифицирован;
409 — номер заказа уже был загружен другим пользователем;
422 — неверный формат номера заказа;
500 — внутренняя ошибка сервера.

*/
