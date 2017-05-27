package main

import (
	"flag"
	"github.com/oliveroneill/hanserver/hanapi/dao"
	"github.com/oliveroneill/hanserver/hanapi/reporting"
	"github.com/oliveroneill/hanserver/hancollector/imagepopulation"
)

func main() {
	// connect to mongo
	db := dao.NewMongoInterface()

	slackAPITokenPtr := flag.String("slacktoken", "", "Specify the API token for logging through Slack")
	flag.Parse()
	logger := reporting.NewSlackLogger(*slackAPITokenPtr)
	populator := imagepopulation.NewImagePopulator(logger)
	// call it once before starting the timer
	populator.PopulateImageDB(db)
}