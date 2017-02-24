package imagepopulation

import (
    "os"
    "fmt"
    "time"
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
    regions := hanapi.GetRegions(db)
    if len(regions) == 0 {
        fmt.Println(`Warning: There are no specified regions. Either query
            hanhttpserver or set a region in the database`)
        return
    }

    collectors := p.getCollectors()
    var wg sync.WaitGroup
    wg.Add(len(collectors))

    atLeastOneEnabled := false
    for _, collector := range collectors {
        if !collector.GetConfig().IsEnabled() {
            continue
        }
        atLeastOneEnabled = true
        go startPopulating(db, collector, regions)
    }
    if !atLeastOneEnabled {
        panic(`No collectors enabled. Please go to hancollector/collectors/config and set
            Enabled to true on at least one`)
    }
    // wait forever
    wg.Wait()
}

func startPopulating(db db.DatabaseInterface,
                     c collectors.ImageCollector,
                     regions []imagedata.Location) {
    populate(db, c, regions)
    // update the collector at its configured frequency
    for _ = range time.NewTicker(c.GetConfig().GetUpdateFrequency() * time.Second).C {
        populate(db, c, regions)
    }
}

func populate(db db.DatabaseInterface,
                     c collectors.ImageCollector,
                     regions []imagedata.Location) {
    fmt.Println("Populating", c.GetConfig().GetCollectorName())
    // update once at the start
    for _, region := range regions {
        // populate the image db for this collector
        populateImageDBWithCollectors(db,
                               []collectors.ImageCollector{c},
                               region.Lat, region.Lng)
    }
}

/*
    Search through each collector at a specific location and
    add them to the database
    This will return when at least one image in this region is found
    OR if all collectors fail
*/
func populateImageDBWithCollectors(db db.DatabaseInterface,
                                   collectorArr []collectors.ImageCollector,
                                   lat float64, lng float64) {
    // use a channel to wait for first response, so that we can return without
    // unnecessarily waiting for all collector
    successChannel := make(chan int)
    failureChannel := make(chan int)
    atLeastOneEnabled := false
    region := imagedata.NewLocation(lat, lng)
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
            // only succeed if at least one image was found
            if len(images) > 0 {
                successChannel <- 1
            } else {
                // consider retrieving no images a failure
                failureChannel <- 1
            }
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
