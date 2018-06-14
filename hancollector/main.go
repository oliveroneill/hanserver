package main

import (
	"bytes"
	"fmt"
	"github.com/oliveroneill/hanserver/hanapi"
	"github.com/oliveroneill/hanserver/hanapi/reporting"
	"github.com/oliveroneill/hanserver/hancollector/imagepopulation"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
)

func main() {
	configPath := kingpin.Arg("config", "Config file for data collection.").Required().String()
	slackAPIToken := kingpin.Flag("slacktoken", "Specify the API token for logging through Slack").String()
	kingpin.Parse()

	// connect to mongo
	db := hanapi.NewMongoInterface()

	// parse config
	config := configToString(*configPath)

	logger := reporting.NewSlackLogger(*slackAPIToken)
	populator := imagepopulation.NewImagePopulator(config, logger)
	// call it once before starting the timer
	populator.PopulateImageDB(db)
}

func configToString(path string) string {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	io.Copy(buf, f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(buf.Bytes())
}
