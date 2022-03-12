package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/provider/accrual"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
)

var _ accrual.Provider = (*service)(nil)

const (
	serviceName = "http-provider"
)

/*
GET /api/orders/{number} HTTP/1.1
Content-Length: 0

Возможные коды ответа:

200 — успешная обработка запроса.

Формат ответа:
200 OK HTTP/1.1
Content-Type: application/json
...

{
    "order": "<number>",
    "status": "PROCESSED",
    "accrual": 500
}

Поля объекта ответа:

order — номер заказа;

status — статус расчёта начисления:

REGISTERED — заказ зарегистрирован, но не начисление не рассчитано;
INVALID — заказ не принят к расчёту, и вознаграждение не будет начислено;
PROCESSING — расчёт начисления в процессе;
PROCESSED — расчёт начисления окончен;
accrual — рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.

429 — превышено количество запросов к сервису.

Формат ответа:
429 Too Many Requests HTTP/1.1
Content-Type: text/plain
Retry-After: 60

No more than N requests per minute allowed

500 — внутренняя ошибка сервера.


*/

type (
	service struct {
		config *Config
		client *http.Client
	}

	option func(svc *service) error
)

// WithConfig sets Config.
func WithConfig(cfg *Config) option {
	return func(svc *service) error {
		svc.config = cfg
		return nil
	}
}

// New creates a new service.
func New(opts ...option) (*service, error) {
	svc := &service{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	if svc.config == nil {
		return nil, fmt.Errorf("config: nil")
	}
	if err := svc.config.validate(); err != nil {
		return nil, fmt.Errorf("Config validation: %w", err)
	}

	svc.client = &http.Client{Timeout: svc.config.timeout}

	return svc, nil
}

// ProcessOrder request accrual for order
func (svc *service) ProcessOrder(ctx context.Context, order model.Order) (model.Order, error) {

	return model.Order{}, nil
}

// sendRequest sends request
// GET /api/orders/{number} HTTP/1.1
func (svc *service) sendRequest(ctx context.Context, order model.Order) (model.Order, error) {
	//ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here
	logger := svc.Logger(ctx)
	logger.Info().Msg("sendRequest")

	endpoint := svc.config.AccrualAddress + "/api/orders/" + order.Number
	jsonData, err := json.Marshal(order)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Err(err).Msg("sendRequest:http.NewRequest")
		return model.Order{}, err
	}

	response, err := svc.client.Do(req)
	if err != nil {
		logger.Err(err).Msg("sendRequest:Error sending request to API endpoint.")
		return model.Order{}, err
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Err(err).Msg("sendRequest:ReadAll:Couldn't read response body.")
		return model.Order{}, err
	}

	if response.StatusCode == 500 {
		logger.Err(pkg.ErrServerError).Msg("sendRequest:error on server side")
		return model.Order{}, err
	}

	if response.StatusCode == 429 { // TODO: Get Retry-After: 60
		logger.Err(pkg.ErrTooManyRequests).Msgf("sendRequest:Too Many Requests: %s", string(body))
		return model.Order{}, err
	}

	if response.StatusCode == 200 {
		if err := json.Unmarshal(body, &order); err != nil {
			logger.Err(err).Msg("sendRequest:Unmarshal:Couldn't parse response body.")
			return model.Order{}, err
		}
	}

	return order, nil
}

// Logger returns logger with ServiceKey field set.
func (svc *service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}
