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

    // call it once before starting the timer
    imagepopulation.PopulateImageDB(db)
    // start getting images every 30 seconds
    for _ = range time.NewTicker(30 * time.Second).C {
        log.Println("Populating...")
        imagepopulation.CleanImages(db)
        imagepopulation.PopulateImageDB(db)
    }
}