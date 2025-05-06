package models

import "time"

type CompetitorState struct {
	ID            string
	StartPlanned  time.Time
	StartActual   *time.Time
	Disqualified  bool
	NotStarted    bool
	NotFinished   bool
	LapsCompleted int
	LapTimes      []time.Duration
	PenaltyTime   time.Duration
	Shots         int
	Hits          int
	LastEventTime time.Time
}
