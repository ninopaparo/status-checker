package main

import (
	"flag"
	"log"

	"github.com/ninopaparo/status-checker/cmd/slack"
)

var slackStatusCurrent = flag.Bool("current", false, "get Slack's current health status")
var slackStatusHistory = flag.Bool("history", false, "get Slack's history health status")
var slackDebugMode = flag.Bool("debug-mode", false, "enable debug mode")

func main() {
	flag.Parse()

	// Create Slack HTTP Client
	sc, err := slack.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	if *slackStatusCurrent {
		res, err := sc.GetCurrentStatus()
		if err != nil {
			log.Fatal(err)
		}

		if *slackDebugMode {
			err = res.DebugResponse()
			if err != nil {
				log.Fatal(err)
			}
		}
		res.CurrentStatus()
	}

	if *slackStatusHistory {
		res, err := sc.GetStatusHistory()
		if err != nil {
			log.Fatal(err)
		}

		if *slackDebugMode {
			err = slack.DebugResponse(res)
			if err != nil {
				log.Fatal(err)
			}
		}
		err = slack.DisplayIncidentHistory(res)
		if err != nil {
			log.Fatal(err)
		}
	}

}
