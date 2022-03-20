package handler

import (
	"github.com/go-chi/render"
	"net/http"
)

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error  `json:"-"`               // low-level runtime error
	HTTPStatusCode int    `json:"-"`               // http response status code
	StatusText     string `json:"status"`          // user-level status message
	AppCode        int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var (
	NewOrderAccepted       = &ErrResponse{HTTPStatusCode: 202, StatusText: "New order accepted"}
	ErrNoData              = &ErrResponse{HTTPStatusCode: 204, StatusText: "No data"}
	ErrNotFound            = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
	ErrMethodNotAllowed    = &ErrResponse{HTTPStatusCode: 405, StatusText: "Method not allowed"}
	ErrBadRequest          = &ErrResponse{HTTPStatusCode: 400, StatusText: "Bad request"}
	ErrAlreadyExists       = &ErrResponse{HTTPStatusCode: 409, StatusText: "Already exists"}
	ErrInvalidLogin        = &ErrResponse{HTTPStatusCode: 401, StatusText: "Invalid login"}
	ErrInvalidOrderNum     = &ErrResponse{HTTPStatusCode: 422, StatusText: "Invalid order number"} //422 — неверный формат номера заказа;
	ErrInsufficientBalance = &ErrResponse{HTTPStatusCode: 402, StatusText: "Insufficient balance"}
)

/*
402 — на счету недостаточно средств;
200 — номер заказа уже был загружен этим пользователем;
202 — новый номер заказа принят в обработку;
400 — неверный формат запроса;
401 — пользователь не аутентифицирован;
409 — номер заказа уже был загружен другим пользователем;
422 — неверный формат номера заказа;
500 — внутренняя ошибка сервера.

*/
