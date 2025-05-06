package main

import (
	processor "biathlon/internal/domain"
)

func main() {
	system := processor.NewCompetitionSystem()
	system.Start()
}
