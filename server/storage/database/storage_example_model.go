package database

import "time"

type ExampleModel struct {
	ID          int       `db:"id"`
	ChannelID   string    `db:"channel_id"`
	MainTs      string    `db:"main_thread_ts"`
	IsEmbedded  bool      `db:"is_embedded"`
	ChatSummary string    `db:"chat_summary"`
	ChatHistory string    `db:"chat_history"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
