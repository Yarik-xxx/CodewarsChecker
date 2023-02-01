package models

import "time"

type Challenge struct {
	Id         string
	Rank       string
	LastUpdate time.Time
}

type Message struct {
	UserName   string
	Text       string
	LastUpdate time.Time
}
