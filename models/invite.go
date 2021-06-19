package models

type Invites struct {
	ChannelID int64 `db:"channel_id"`
	FromID    int64 `db:"from_id"`
	MemberID  int64 `db:"member_id"`
}

type AllRating struct {
	ChannelID int64  `db:"channel_id"`
	UserID    int64  `db:"user_id"`
	UserName  string `db:"user_name"`
	Count     int64  `db:"count"`
}
