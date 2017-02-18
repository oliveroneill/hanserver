package main

import (
    "log"
    "time"
    "github.com/oliveroneill/hanserver/hancollector/imagepopulation"
    "github.com/oliveroneill/hanserver/hanapi/db"
)

func main() {
    // connect to mongo
    db := db.NewMongoInterface()

    populator := imagepopulation.NewImagePopulator()
    // call it once before starting the timer
    populator.PopulateImageDB(db)
    // start getting images every 30 seconds
    for _ = range time.NewTicker(30 * time.Second).C {
        log.Println("Populating...")
        populator.CleanImages(db)
        populator.PopulateImageDB(db)
    }
}