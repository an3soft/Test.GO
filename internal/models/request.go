package models

import "time"

type Request struct {
	Id       int64
	Received time.Time
}
