package models

type CompetitorResult struct {
	ID           string     `json:"id"`
	TotalTime    string     `json:"total_time"`
	Laps         []LapStats `json:"laps"`
	PenaltyTime  string     `json:"penalty_time"`
	PenaltySpeed float64    `json:"penalty_speed"`
	HitSummary   string     `json:"hit_summary"`
	Status       string     `json:"status"` // "Finished", "NotStarted", "NotFinished"
}
