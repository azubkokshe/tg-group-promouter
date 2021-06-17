package models

type Channel struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
}
