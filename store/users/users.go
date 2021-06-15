package users

import (
	"context"
	"github.com/azubkokshe/tg-group-promouter/models"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	Store(ctx context.Context, tx *sqlx.Tx, user *models.User) error
	Delete(ctx context.Context, tx *sqlx.Tx, id int64) error
}
