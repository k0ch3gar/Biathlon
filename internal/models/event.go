package models

import "time"

type Event struct {
	Time        time.Time
	EventID     int
	Competitor  string
	ExtraParams []string
}
