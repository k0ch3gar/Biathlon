package main

import (
	. "biathlon/internal/api"
	processor "biathlon/internal/domain"
	"flag"
)

func init() {
	flag.StringVar(&PathToLogs, "logs", PathToLogs, "path to log file")
	flag.StringVar(&PathToInEvents, "in", PathToInEvents, "path to input events file")
	flag.StringVar(&PathToOutEvents, "out", PathToOutEvents, "path to output events file")
	flag.StringVar(&PathToFinalReport, "report", PathToFinalReport, "path to final report file")
	flag.StringVar(&PathToConfig, "config", PathToConfig, "path to config file")
}

func main() {
	flag.Parse()
	system := processor.NewCompetitionSystem()
	system.Start()
}
