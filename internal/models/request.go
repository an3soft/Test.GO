package models

import "time"

type Request struct {
	UserId    int64
	ChatId    int64
	MessageID int
	UserName  string
	Text      string
	Received  time.Time
	Updated   time.Time
}
