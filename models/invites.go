package models

type Invites struct {
	ChannelID int64  `db:"channel_id"`
	FromID    string `db:"from_id"`
	MemberID  string `db:"member_id"`
}
