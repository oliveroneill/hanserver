package main

import (
	"flag"
	"fmt"
	"github.com/oliveroneill/hanserver/hanapi/dao"
	"time"
)

// DefaultImageCountLimit is the default maximum amount of images allowed in
// the database before cleaning up
const DefaultImageCountLimit = 500000

// DefaultClearanceCount is the default amount of images to be cleared if the
// maximum is reached
const DefaultClearanceCount = 100000

// Watch the database and clear old images when it starts reaching a max size
func main() {
	// connect to mongo
	db := dao.NewMongoInterface()

	// parse arguments
	limitUsageString := "Specify the maximum amount of images allowed in the database"
	imageCountLimitPtr := flag.Int("imagelimit", DefaultImageCountLimit, limitUsageString)
	clearanceUsageString := "Specify how many images should be cleared when reaching the maximum"
	clearanceCountPtr := flag.Int("clear", DefaultClearanceCount, clearanceUsageString)
	flag.Parse()
	imageLimit := *imageCountLimitPtr
	clearanceCount := *clearanceCountPtr

	checkAndClean(db, imageLimit, clearanceCount)
	// every hour the database is checked and old images are cleared out
	freq := 60 * 60 * time.Second
	for _ = range time.NewTicker(freq).C {
		checkAndClean(db, imageLimit, clearanceCount)
	}
}

func checkAndClean(db dao.DatabaseInterface, limit int, clear int) {
	if db.Size() >= limit {
		fmt.Println("Cleaning up", clear, "images")
		db.DeleteOldImages(clear)
	}
}
