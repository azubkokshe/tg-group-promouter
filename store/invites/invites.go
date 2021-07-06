package invites

import (
	"context"
	"github.com/azubkokshe/tg-group-promouter/models"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	Store(ctx context.Context, tx *sqlx.Tx, invite *models.Invites) error
	Journal(ctx context.Context, tx *sqlx.Tx, invite *models.Journal) error
	Delete(ctx context.Context, tx *sqlx.Tx, channelID int64, fromID int64, memberID int64) error
	SelectRating(ctx context.Context, channelID int64) *[]models.AllRating
}
