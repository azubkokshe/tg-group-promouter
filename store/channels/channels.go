package channels

import (
	"context"
	"github.com/azubkokshe/tg-group-promouter/models"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	Store(ctx context.Context, tx *sqlx.Tx, invite *models.Channel) error
	Delete(ctx context.Context, tx *sqlx.Tx, id int64) error
	GetByID(ctx context.Context, id int64) (*models.Channel, error)
}
