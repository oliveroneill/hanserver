package imagepopulation

import (
    "fmt"
    "sync"
    "github.com/oliveroneill/hanserver/hanapi"
    "github.com/oliveroneill/hanserver/hanapi/db"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors"
)

func getCollectors() []collectors.ImageCollector {
    return []collectors.ImageCollector{collectors.NewInstagramCollector()}
}

// PopulateImageDBWithLoc will populate the database with images at this
// specific location
func PopulateImageDBWithLoc(db db.DatabaseInterface, lat float64, lng float64) {
    populateImageDBWithCollectors(db, getCollectors(), lat, lng)
}

// PopulateImageDB will populate the database with images using the regions
// set in the database. This will return once each region has new images from
// at least one collector
func PopulateImageDB(db db.DatabaseInterface) {
    var wg sync.WaitGroup
    regions := hanapi.GetRegions(db)
    if len(regions) == 0 {
        fmt.Println(`Warning: There are no specified regions. Either query
            hanhttpserver or set a region in the database`)
        return
    }
    wg.Add(len(regions))
    for _, region := range regions {
        // we'll wait for all collectors to complete, so that everything
        // completes
        go func() {
            defer wg.Done()
            PopulateImageDBWithLoc(db, region.Lat, region.Lng)
        }()
    }
    wg.Wait()
}

/*
    Search through each collector at a specific location and
    add them to the database
    This will return when one collector has found images
*/
func populateImageDBWithCollectors(db db.DatabaseInterface,
                                   collectorArr []collectors.ImageCollector,
                                   lat float64, lng float64) {
    // use a channel to wait for first response, so that we can return without
    // unnecessarily waiting for all collector
    channel := make(chan int)
    atLeastOneEnabled := false
    region := imagedata.NewImageLocation(lat, lng)
    for _, collector := range collectorArr {
        go func(c collectors.ImageCollector) {
            if !c.GetConfig().IsEnabled() {
                // set value so that we don't wait forever
                channel <- 1
                return
            }
            atLeastOneEnabled = true
            for _, img := range c.GetImages(lat, lng) {
                img.Region = region
                db.AddImage(img)
            }
            channel <- 1
        }(collector)
    }
    // wait for first to return
    <-channel
    if !atLeastOneEnabled {
        panic(`No collectors enabled. Please go to hancollector/collectors/config and set
            Enabled to true on at least one`)
    }
}

// CleanImages is an asynchronous method that makes sure every image points to
// an actual image, this will clean up images deleted from their original
// source
//
// TODO: not sure whether to go through every image and monitor response codes
// (costly), or just delete old images. Or the client could send a message when images
// are dead, this seems exploitable but we could check upon receiving the message to
// confirm
func CleanImages(db db.DatabaseInterface) {
    // go func() {
    //     for _, img := range db.GetAllImages() {
    //         // test image
    //     }
    // }()
}