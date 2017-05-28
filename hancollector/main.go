package main

import (
	"io"
	"os"
	"fmt"
	"flag"
	"bytes"
	"github.com/oliveroneill/hanserver/hanapi/dao"
	"github.com/oliveroneill/hanserver/hanapi/reporting"
	"github.com/oliveroneill/hanserver/hancollector/imagepopulation"
)

func main() {
	slackAPITokenPtr := flag.String("slacktoken", "", "Specify the API token for logging through Slack")
	flag.Parse()

	flag.Usage = printUsage
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// connect to mongo
	db := dao.NewMongoInterface()

	// parse config
	config := configToString(flag.Arg(0))

	logger := reporting.NewSlackLogger(*slackAPITokenPtr)
	populator := imagepopulation.NewImagePopulator(config, logger)
	// call it once before starting the timer
	populator.PopulateImageDB(db)
}

func printUsage() {
	fmt.Printf("Usage: %s config_file ...\n", os.Args[0])
	flag.PrintDefaults()
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