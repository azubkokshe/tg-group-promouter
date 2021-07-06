package invites_pg

import (
	"context"
	"github.com/azubkokshe/tg-group-promouter/models"
	"github.com/azubkokshe/tg-group-promouter/store/invites"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB           *sqlx.DB
	insertQuery  string
	ratingQuery  string
	journalQueue string
}

func NewRepository(db *sqlx.DB) invites.Store {
	return &Repository{
		DB: db,
		insertQuery: `
						INSERT INTO tbl_invite
									(
												channel_id,
												from_id,
												member_id
									)
									VALUES
									(
												:channel_id,
												:from_id,
												:member_id
									)
						on conflict
						ON CONSTRAINT tbl_invite_pk do nothing;`,
		ratingQuery: `
						WITH tbl
							 AS (SELECT channel_id,
										title,
										from_id        AS from_user_id,
										tu.first_name  AS from_first_name,
										tu.last_name   AS from_last_name,
										member_id      AS member_user_id,
										tu2.first_name AS member_first_name,
										tu2.last_name  AS member_last_name
								 FROM   tbl_invite
										INNER JOIN tbl_channel tc
												ON tbl_invite.channel_id = tc.id
										INNER JOIN tbl_user tu
												ON tbl_invite.from_id = tu.id
										INNER JOIN tbl_user tu2
												ON tbl_invite.member_id = tu2.id
								WHERE channel_id=$1)
						SELECT channel_id,
							   from_user_id      AS user_id,
							   from_first_name
							   || ' '
							   || from_last_name AS user_name,
							   Count(*)          AS count
						FROM   tbl
						GROUP  BY channel_id,
								  user_id,
								  user_name
						ORDER BY count DESC;`,
		journalQueue: `
						INSERT INTO tbl_journal
									(
												record
									)
									VALUES
									(
												:record
									);
						`,
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

func (r *Repository) Journal(ctx context.Context, tx *sqlx.Tx, journal *models.Journal) error {
	query, args, err := sqlx.Named(r.journalQueue, journal)
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

func (r *Repository) SelectRating(ctx context.Context, channelID int64) *[]models.AllRating {
	var rating []models.AllRating
	_ = r.DB.Select(&rating, r.ratingQuery, channelID)
	return &rating
}

func (r *Repository) Delete(ctx context.Context, tx *sqlx.Tx, channelID int64, fromID int64, memberID int64) error {
	return nil
}
