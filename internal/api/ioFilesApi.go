package api

import (
	"biathlon/internal/config"
	. "biathlon/internal/event"
	. "biathlon/internal/models"
	"bufio"
	"encoding/json"
	"os"
)

var PathToLogs string = "output/latest.log"
var PathToInEvents string = "input/events_in.txt"
var PathToOutEvents string = "output/events_out.txt"
var PathToFinalReport string = "output/final_report.txt"
var PathToConfig string = "config/config.json"

func ReadEvents() []Event {
	file, err := os.Open(PathToInEvents)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var events []Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		event, err := ParseEvent(scanner.Text())
		if err == nil {
			events = append(events, event)
		}
	}
	return events
}

func WriteLog(log string) error {
	file, err := os.OpenFile(PathToLogs, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(log + "\n")
	if err != nil {
		return err
	}

	return nil
}

func ReadConfig() (config.Config, error) {
	file, err := os.ReadFile(PathToConfig)
	if err != nil {
		return config.Config{}, err
	}

	var result config.Config
	err = json.Unmarshal(file, &result)

	return result, nil
}

func WriteEvent(event string) error {
	file, err := os.OpenFile(PathToOutEvents, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	_, err = file.WriteString(event)
	if err != nil {
		return err
	}

	return nil
}

func WriteReport(report []string) error {
	file, err := os.OpenFile(PathToFinalReport, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, line := range report {
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}

		_, err = file.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}

	return nil
}
