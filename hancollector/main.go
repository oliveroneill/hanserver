package main

import (
	"flag"
	"github.com/oliveroneill/hanserver/hanapi/db"
	"github.com/oliveroneill/hanserver/hanapi/reporting"
	"github.com/oliveroneill/hanserver/hancollector/imagepopulation"
)

func main() {
	// connect to mongo
	db := db.NewMongoInterface()

	slackAPITokenPtr := flag.String("slacktoken", "", "Specify the API token for logging through Slack")
	flag.Parse()
	logger := reporting.NewSlackLogger(*slackAPITokenPtr)
	populator := imagepopulation.NewImagePopulator(logger)
	// call it once before starting the timer
	populator.PopulateImageDB(db)
}