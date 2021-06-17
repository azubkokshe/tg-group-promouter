package models

type Invites struct {
	ChannelID int64 `db:"channel_id"`
	FromID    int64 `db:"from_id"`
	MemberID  int64 `db:"member_id"`
}
