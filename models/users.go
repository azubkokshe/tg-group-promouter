package models

type User struct {
	ID        int64  `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	UserName  string `db:"username"`
	IsBot     bool   `db:"is_bot"`
}
