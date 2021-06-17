package channels_pg

import (
	"context"
	"github.com/azubkokshe/tg-group-promouter/models"
	"github.com/azubkokshe/tg-group-promouter/store/channels"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB          *sqlx.DB
	insertQuery string
	getQuery    string
}

func NewRepository(db *sqlx.DB) channels.Store {
	return &Repository{
		DB: db,
		insertQuery: `INSERT INTO tbl_channel (id, title)
                      VALUES(:id, :title)
                      ON CONFLICT (id)
                      DO NOTHING`,
		getQuery: `SELECT id, title FROM tbl_channel WHERE id=$1`,
	}
}

func (r *Repository) Store(ctx context.Context, tx *sqlx.Tx, channel *models.Channel) error {
	query, args, err := sqlx.Named(r.insertQuery, channel)
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

func (r *Repository) Delete(ctx context.Context, tx *sqlx.Tx, id int64) error {
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*models.Channel, error) {
	c := models.Channel{}

	err := r.DB.Get(&c, r.getQuery, id)
	if err != nil {
		return &models.Channel{}, err
	}
	return &c, nil
}
