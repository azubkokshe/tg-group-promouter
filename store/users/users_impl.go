package users

import (
	"context"
	"github.com/azubkokshe/tg-group-promouter/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB          *sqlx.DB
	insertQuery string
}

func NewRepository(db *sqlx.DB) Store {
	return &Repository{
		DB: db,
		insertQuery: `INSERT INTO tbl_user (id, first_name, last_name, username, is_bot)
                      VALUES(:id, :first_name, :last_name, :username, :is_bot)
                      ON CONFLICT (id)
                      DO UPDATE SET first_name = excluded.first_name,
                                    last_name = excluded.last_name,
                                    username = excluded.username,
                                    is_bot = excluded.is_bot;`,
	}
}

func (r *Repository) Store(_ context.Context, tx *sqlx.Tx, user *models.User) error {
	query, args, err := sqlx.Named(r.insertQuery, user)
	if err != nil {
		return err
	}
	query = tx.Rebind(query)
	var id int64
	err = tx.Get(&id, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(_ context.Context, tx *sqlx.Tx, id int64) error {
	return nil
}
