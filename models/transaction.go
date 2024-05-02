package models

import "time"

type Transcation struct {
	SenderId   string
	ReceiverId string
	Amount     float32
	Timestamp  time.Time
}
