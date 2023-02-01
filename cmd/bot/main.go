package main

import (
	"CodeWarsCheckerAnalysis/internal/app"
)

func main() {
	var bot app.Bot

	// bot fall protection
	defer func() {
		if err := recover(); err != nil {
			main()
		}
	}()

	bot.Start()
}
