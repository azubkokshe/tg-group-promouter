package invites_pg

import (
	"context"
	"github.com/azubkokshe/tg-group-promouter/models"
	"github.com/azubkokshe/tg-group-promouter/store/invites"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB          *sqlx.DB
	insertQuery string
}

func NewRepository(db *sqlx.DB) invites.Store {
	return &Repository{
		DB: db,
		insertQuery: `INSERT INTO tbl_invite (channel_id, from_id, member_id)
                      VALUES(:channel_id, :from_id, :member_id)
                      ON CONFLICT ON CONSTRAINT tbl_invite_pk DO NOTHING;`,
	}
}

func (r *Repository) Store(ctx context.Context, tx *sqlx.Tx, invite *models.Invites) error {
	query, args, err := sqlx.Named(r.insertQuery, invite)
	if err != nil {
		return err
	}
	query = tx.Rebind(query)
	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, tx *sqlx.Tx, channelID int64, fromID int64, memberID int64) error {
	return nil
}
