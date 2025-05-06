package event

import (
	. "biathlon/internal/models"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseEvent(line string) (Event, error) {
	var event Event

	re := regexp.MustCompile(`\[(.*?)] (\d+) (\d+)(?: (.*))?`)

	if strings.TrimSpace(line) == "" {
		panic("invalid event")
	}

	m := re.FindStringSubmatch(line)
	if len(m) < 4 {
		panic("invalid event")
	}
	t, _ := time.Parse("15:04:05.000", m[1])
	eid, _ := strconv.Atoi(m[2])
	var params []string
	if len(m) > 4 && m[4] != "" {
		params = strings.Fields(m[4])
	}

	event = Event{
		Time:        t,
		EventID:     eid,
		Competitor:  m[3],
		ExtraParams: params,
	}

	return event, nil
}
