package processor

import (
	"biathlon/internal/api"
	. "biathlon/internal/config"
	. "biathlon/internal/models"
	"fmt"
	"strings"
	"time"
)

type CompetitionSystem struct {
	Config      Config
	Competitors map[string]*CompetitorState
	Logs        []string
	BaseStart   time.Time
	StartDelta  time.Duration
}

func NewCompetitionSystem() *CompetitionSystem {
	cfg, err := api.ReadConfig()
	if err != nil {
		panic(err)
	}

	baseStart, _ := time.Parse("15:04:05", cfg.Start)
	delta, _ := time.ParseDuration(cfg.StartDelta)
	return &CompetitionSystem{
		Config:      cfg,
		Competitors: make(map[string]*CompetitorState),
		BaseStart:   baseStart,
		StartDelta:  delta,
	}
}

func (c *CompetitionSystem) Start() {
	c.ProcessEvents(api.ReadEvents())
	err := api.WriteReport(c.GenerateReport())
	if err != nil {
		panic(err)
	}
}

func (c *CompetitionSystem) ProcessEvents(events []Event) {
	for _, e := range events {
		competitor := e.Competitor
		if _, exists := c.Competitors[competitor]; !exists {
			c.Competitors[competitor] = &CompetitorState{
				ID: competitor,
			}
		}
		comp := c.Competitors[competitor]

		switch e.EventID {
		case 1:
			c.log(e.Time, fmt.Sprintf("The competitor(%s) registered", competitor))
			comp.NotStarted = true
		case 2:
			t, _ := time.Parse("15:04:05.000", e.ExtraParams[0])
			comp.StartPlanned = t
			c.log(e.Time, fmt.Sprintf("The start time for the competitor(%s) was set by a draw to %s", competitor, t.Format("15:04:05.000")))
		case 3:
			c.log(e.Time, fmt.Sprintf("The competitor(%s) is on the start line", competitor))
		case 4:
			comp.StartActual = &e.Time
			comp.NotStarted = false
			c.log(e.Time, fmt.Sprintf("The competitor(%s) has started", competitor))
		case 5:
			c.log(e.Time, fmt.Sprintf("The competitor(%s) is on the firing range(%s)", competitor, e.ExtraParams[0]))
		case 6:
			comp.Hits++
			comp.Shots++
			c.log(e.Time, fmt.Sprintf("The target(%s) has been hit by competitor(%s)", e.ExtraParams[0], competitor))
		case 7:
			c.log(e.Time, fmt.Sprintf("The competitor(%s) left the firing range", competitor))
		case 8:
			c.log(e.Time, fmt.Sprintf("The competitor(%s) entered the penalty laps", competitor))
			comp.LastEventTime = e.Time
		case 9:
			delta := e.Time.Sub(comp.LastEventTime)
			comp.PenaltyTime += delta
			c.log(e.Time, fmt.Sprintf("The competitor(%s) left the penalty laps", competitor))
		case 10:
			lapTime := e.Time.Sub(comp.LastEventTime)
			comp.LapTimes = append(comp.LapTimes, lapTime)
			comp.LapsCompleted++
			if comp.LapsCompleted == c.Config.Laps {
				err := api.WriteEvent(fmt.Sprintf("[%s] %s", e.Time.Format("15:04:05.000"), fmt.Sprintf("The competitor(%s) has finished", competitor)))
				if err != nil {
					panic(err)
				}

			}

			c.log(e.Time, fmt.Sprintf("The competitor(%s) ended the main lap", competitor))
		case 11:
			comp.NotFinished = true
			c.log(e.Time, fmt.Sprintf("The competitor(%s) can`t continue: %s", competitor, strings.Join(e.ExtraParams, " ")))
		}
		comp.LastEventTime = e.Time
	}

}

func (c *CompetitionSystem) GenerateReport() []string {
	var report []string
	for _, comp := range c.Competitors {
		var total string
		if comp.NotStarted {
			total = "[NotStarted]"
		} else if comp.NotFinished {
			total = "[NotFinished]"
		} else if comp.StartActual != nil {
			diff := comp.LastEventTime.Sub(*comp.StartActual)
			total = fmt.Sprintf("{%s}", formatDuration(diff))
		}
		var laps []string
		for _, t := range comp.LapTimes {
			avg := float64(c.Config.LapLen) / t.Seconds()
			laps = append(laps, fmt.Sprintf("{%s, %.3f}", formatDuration(t), avg))
		}
		penaltyAvg := 0.0
		if comp.PenaltyTime > 0 {
			penaltyAvg = float64(c.Config.PenaltyLen) / comp.PenaltyTime.Seconds()
		}
		report = append(report,
			fmt.Sprintf("%s %s %v {%s, %.3f} %d/%d",
				total, comp.ID, laps, formatDuration(comp.PenaltyTime), penaltyAvg, comp.Hits, comp.Shots))
	}
	return report
}

func (c *CompetitionSystem) log(t time.Time, msg string) {
	err := api.WriteLog(fmt.Sprintf("[%s] %s", t.Format("15:04:05.000"), msg))
	if err != nil {
		panic(err)
	}
}

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%06.3f",
		int(d.Hours()), int(d.Minutes())%60, d.Seconds()-float64(d.Minutes()*60))
}
