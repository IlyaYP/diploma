package psql

import (
	"context"
	"errors"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// CreateUser creates a new model.User.
// Returns ErrAlreadyExists if user exists.
func (svc *service) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	logger := svc.Logger(ctx)
	logger.UpdateContext(user.GetLoggerContext)

	_, err := svc.pool.Exec(ctx, `insert into users(login, password) values ($1, $2)`,
		user.Login, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				logger.Err(err).Msg("Error creating user")
				return model.User{}, pkg.ErrAlreadyExists
			}
		}

		logger.Err(err).Msg("Error creating user")
		return model.User{}, err
	}

	logger.Info().Msg("Successfully created user")

	return user, nil
}

// GetUserByLogin returns model.User by its login if exists.
func (svc *service) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	logger := svc.Logger(ctx)
	// Build input
	user := model.User{
		Login: login,
	}

	logger.UpdateContext(user.GetLoggerContext)
	logger.Info().Msg("GetUserByLogin")

	err := svc.pool.QueryRow(ctx, "select password from users where login=$1", login).Scan(&user.Password)
	switch err {
	case nil:
		return &user, nil
	case pgx.ErrNoRows:
		logger.Err(pkg.ErrNotExists).Msg("Error GetUserByLogin")
		return nil, pkg.ErrInvalidLogin
	default:
		logger.Err(err).Msg("Error GetUserByLogin")
		return nil, err
	}
}
