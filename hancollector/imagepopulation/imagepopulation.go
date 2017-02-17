package imagepopulation

import (
    "os"
    "fmt"
    "sync"
    "github.com/nlopes/slack"
    "github.com/oliveroneill/hanserver/hanapi"
    "github.com/oliveroneill/hanserver/hanapi/db"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors"
)

// ImagePopulator is a type that will populate images from its set of
// collectors
type ImagePopulator struct {
    collectorsList []collectors.ImageCollector
}

// NewImagePopulator creates a new `ImagePopulator`
func NewImagePopulator() *ImagePopulator {
    p := new(ImagePopulator)
    p.collectorsList = []collectors.ImageCollector {
        collectors.NewTwitterCollector(),
        collectors.NewInstagramCollector(),
        collectors.NewFlickrCollector(),
        collectors.NewGooglePlacesCollector(),
    }
    return p
}

func (p ImagePopulator) getCollectors() []collectors.ImageCollector {
    return p.collectorsList
}

// PopulateImageDBWithLoc will populate the database with images at this
// specific location
func (p *ImagePopulator) PopulateImageDBWithLoc(db db.DatabaseInterface, lat float64, lng float64) {
    populateImageDBWithCollectors(db, p.getCollectors(), lat, lng)
}

// PopulateImageDB will populate the database with images using the regions
// set in the database. This will return once each region has new images from
// at least one collector
func (p *ImagePopulator) PopulateImageDB(db db.DatabaseInterface) {
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
        go func(region imagedata.ImageLocation) {
            defer wg.Done()
            p.PopulateImageDBWithLoc(db, region.Lat, region.Lng)
        }(region)
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
    successChannel := make(chan int)
    failureChannel := make(chan int)
    atLeastOneEnabled := false
    region := imagedata.NewImageLocation(lat, lng)
    for _, collector := range collectorArr {
        if !collector.GetConfig().IsEnabled() {
            continue
        }
        atLeastOneEnabled = true
        go func(c collectors.ImageCollector) {
            images, err := c.GetImages(lat, lng)
            if err != nil {
                reportError(err, c.GetConfig().GetCollectorName())
                failureChannel <- 1
                return
            }
            for _, img := range images {
                img.Region = region
                db.AddImage(img)
            }
            successChannel <- 1
        }(collector)
    }

    if !atLeastOneEnabled {
        panic(`No collectors enabled. Please go to hancollector/collectors/config and set
            Enabled to true on at least one`)
    }
    failures := 0
    for {
        select {
            case <-successChannel:
                return
            case <-failureChannel:
                failures++
                // wait for all failures until we give up
                if failures >= len(collectorArr) {
                    return
                }
        }
    }
}

// reports errors through Slack
func reportError(err error, collectorName string) {
    fmt.Fprintln(os.Stderr, collectorName, "Error:", err)
    apiToken := os.Getenv("SLACK_API_TOKEN")
    if len(apiToken) == 0 {
        fmt.Println("Slack support not set up. Please set SLACK_API_TOKEN environment variable")
        return
    }
    channelName := "hanserver"
    api := slack.New(apiToken)
    params := slack.PostMessageParameters{}
    _, _, err = api.PostMessage(channelName, fmt.Sprintf("%s Error: %s", collectorName, err), params)
    if err != nil {
        fmt.Println("%s", err)
        return
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
func (p ImagePopulator) CleanImages(db db.DatabaseInterface) {
    // go func() {
    //     for _, img := range db.GetAllImages() {
    //         // test image
    //     }
    // }()
}