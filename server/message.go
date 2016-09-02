package main

import (
	"github.com/twinj/uuid"
	"time"
)

type Message struct {
	Id        uuid.Uuid
	Type      string
	Name      string
	Timestamp time.Time
	Text      string
}
