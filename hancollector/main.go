package main

import (
	"github.com/oliveroneill/hanserver/hancollector/imagepopulation"
	"github.com/oliveroneill/hanserver/hanapi/db"
)

func main() {
	// connect to mongo
	db := db.NewMongoInterface()

	populator := imagepopulation.NewImagePopulator()
	// call it once before starting the timer
	populator.PopulateImageDB(db)
}