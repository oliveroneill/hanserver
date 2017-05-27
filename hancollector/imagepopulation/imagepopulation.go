package imagepopulation

import (
	"os"
	"fmt"
	"time"
	"sync"
	"github.com/oliveroneill/hanserver/hanapi"
	"github.com/oliveroneill/hanserver/hanapi/db"
	"github.com/oliveroneill/hanserver/hanapi/reporting"
	"github.com/oliveroneill/hanserver/hanapi/imagedata"
	"github.com/oliveroneill/hanserver/hancollector/collectors"
)

// ImagePopulator is a type that will populate images from its set of
// collectors
type ImagePopulator struct {
	collectorsList []collectors.ImageCollector
	logger         reporting.Logger
}

// NewImagePopulator creates a new `ImagePopulator`
// @param logger - optional logging support
func NewImagePopulator(logger reporting.Logger) *ImagePopulator {
	p := new(ImagePopulator)
	p.collectorsList = []collectors.ImageCollector {
		collectors.NewTwitterCollector(),
		collectors.NewInstagramCollector(),
		collectors.NewFlickrCollector(),
	}
	p.logger = logger
	return p
}

func (p *ImagePopulator) getCollectors() []collectors.ImageCollector {
	return p.collectorsList
}

// PopulateImageDBWithLoc will populate the database with images at this
// specific location
func (p *ImagePopulator) PopulateImageDBWithLoc(db db.DatabaseInterface, lat float64, lng float64) {
	populateImageDBWithCollectors(db, p.getCollectors(), lat, lng, p.logger)
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
		go p.startPopulating(db, collector, regions)
	}
	if !atLeastOneEnabled {
		panic(`No collectors enabled. Please go to hancollector/collectors/config and set
			Enabled to true on at least one`)
	}
	// wait forever
	wg.Wait()
}

func (p *ImagePopulator) startPopulating(db db.DatabaseInterface,
					 					 c collectors.ImageCollector,
					 					 regions []imagedata.Location) {
	p.populate(db, c, regions)
	// update the collector at its configured frequency
	freq := c.GetConfig().GetUpdateFrequency() * time.Second
	for _ = range time.NewTicker(freq).C {
		p.populate(db, c, regions)
	}
}

func (p *ImagePopulator) populate(db db.DatabaseInterface,
					 			  c collectors.ImageCollector,
					 			  regions []imagedata.Location) {
	fmt.Println("Populating", c.GetConfig().GetCollectorName())
	// update once at the start
	for _, region := range regions {
		// populate the image db for this collector
		populateImageDBWithCollectors(db,
							   		  []collectors.ImageCollector{c},
							          region.Lat, region.Lng, p.logger)
	}
}

/*
	Search through each collector at a specific location and
	add them to the database
	This will return when at least one image in this region is found
	OR if all collectors fail
*/
func populateImageDBWithCollectors(db db.DatabaseInterface,
	collectorArr []collectors.ImageCollector, lat float64, lng float64,
	logger reporting.Logger) {
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
				reportError(err, c.GetConfig().GetCollectorName(), logger)
				failureChannel <- 1
				return
			}
			db.AddBulkImagesToRegion(images, region)
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

func reportError(err error, collectorName string, logger reporting.Logger) {
	fmt.Fprintln(os.Stderr, collectorName, "Error:", err)
	if logger != nil {
		logger.Log(fmt.Sprintf("%s Error: %s", collectorName, err))
	}
}
