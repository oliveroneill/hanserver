package hanapi

import (
    "fmt"
    "os"
    "sort"
    "math"
    "github.com/nlopes/slack"
    "github.com/kellydunn/golang-geo"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hanapi/db"
    "github.com/oliveroneill/hanserver/hanapi/feedsort"
)

// RegionSize is radius of a region in meters
const RegionSize = 5000

// ContainsRegion - determines whether a point is within a specific region
func ContainsRegion(db db.DatabaseInterface, lat float64, lng float64) bool {
    return GetRegion(db, lat, lng) != nil
}

// GetRegion - returns the region which the specified lat, lng lies in
func GetRegion(db db.DatabaseInterface, lat float64, lng float64) *imagedata.Location {
    regions := db.GetRegions()
    currentPoint := geo.NewPoint(lat, lng)
    // loop through each region and return the first one that the point is
    // enclosed in
    for _, r := range regions {
        p := geo.NewPoint(r.Lat, r.Lng)
        if p.GreatCircleDistance(currentPoint) <= RegionSize / 1000 {
            return &r
        }
    }
    return nil
}

// GetRegions - returns the currently used regions
func GetRegions(db db.DatabaseInterface) []imagedata.Location {
    return db.GetRegions()
}

// AddRegion - adds a new region for image population
func AddRegion(db db.DatabaseInterface, lat float64, lng float64) {
    if !ContainsRegion(db, lat, lng) {
        db.AddRegion(lat, lng)
    }
}

// GetImages - get images near the location sorted by distance and recency
func GetImages(db db.DatabaseInterface, lat float64, lng float64) []imagedata.ImageData {
    return GetImagesWithRange(db, lat, lng, -1, -1)
}

// GetImagesWithStart - get images starting at a certain point
func GetImagesWithStart(db db.DatabaseInterface, lat float64, lng float64,
    start int) []imagedata.ImageData {
    return GetImagesWithRange(db, lat, lng, start, -1)
}

// GetImagesWithEnd - get images from the beginning to the specified end
func GetImagesWithEnd(db db.DatabaseInterface, lat float64, lng float64,
    end int) []imagedata.ImageData {
    return GetImagesWithRange(db, lat, lng, -1, end)
}

// GetImagesWithRange - Specify a range, so that you can query a portion of the image list
// @param start - start is optional, use -1 to signify no value, indexing starts at zero
// @param end - end is optional, use -1 to signify no value
func GetImagesWithRange(db db.DatabaseInterface, lat float64, lng float64,
    start int, end int) []imagedata.ImageData {
    // 100 images will be sorted at a time
    return getImagesWithRangeAndSampleSize(db, lat, lng, start, end, 100)
}

// SampleSize will determine the bias on distance, as we aren't sorting the
// whole database on recency on distance. Only the closest images are then
// sorted
func getImagesWithRangeAndSampleSize(db db.DatabaseInterface, lat float64,
    lng float64, start int, end int, sampleSize int) []imagedata.ImageData {
    // fix input values
    if start < 0 {
        start = 0
    }
    if end < 0 {
        end = start + sampleSize
    }
    // remove incorrect requests
    if end < start {
        return []imagedata.ImageData{}
    }
    startSort, endSort := getRange(sampleSize, start, end)

    images := []imagedata.ImageData{}
    // if the request is larger than our sample size then we need to sort
    // multiple sets of our sample size individually to avoid duplicate
    // images
    if endSort - startSort > sampleSize {
        // calculate the range queries that need to be made to only ever sort
        // arrays of sample size
        rangeEnd := end
        rangeStart := start
        // edge case will not return the previous sampleSize boundary, so we
        // subtract and add 1 to the range to get new start and end values
        // which specify where to query and sort
        if end % sampleSize == 0 {
            rangeEnd--
        }
        if start % sampleSize == 0 {
            rangeStart++
        }
        // work out the nearest sample size from the start position and the
        // next smallest sample size from the end position
        closestEnd, closestStart := getRange(sampleSize, rangeEnd, rangeStart)
        if start == 0 {
            closestEnd = sampleSize
        }
        // call image range recursively
        // get the first portion
        images = append(images, getImagesWithRangeAndSampleSize(db, lat, lng, start, closestEnd, sampleSize)...)
        // go through all sample size chunks in between the start and end portion
        for i := closestEnd; i < closestStart; i += sampleSize {
            images = append(images, getImagesWithRangeAndSampleSize(db, lat, lng, i, i+sampleSize, sampleSize)...)
        }
        // get the last portion
        images = append(images, getImagesWithRangeAndSampleSize(db, lat, lng, closestStart, end, sampleSize)...)
        return images
    }
    // this is the base case where we get the images and sort
    images = db.GetImages(lat, lng, startSort, endSort)
    // sort
    sort.Sort(feedsort.BySum(images))
    // figure out where to slice the array
    sliceStart := start - startSort
    sliceEnd := sliceStart + (end - start)
    // return empty if we're out of range
    if sliceStart > len(images) {
        return []imagedata.ImageData{}
    }
    // set relevant values since start and end can be optional
    if sliceStart < 0 {
        sliceStart = 0
    }
    if sliceEnd < 0 || sliceEnd > len(images) {
        sliceEnd = len(images)
    }
    return images[sliceStart:sliceEnd]
}

func getRange(sampleSize int, start int, end int) (int, int) {
    startSort := int(math.Floor(float64(start)/float64(sampleSize)) * float64(sampleSize))
    endSort := int(math.Ceil(float64(end)/float64(sampleSize)) * float64(sampleSize))
    return startSort, endSort
}

// ReportImage - report an image to be removed
// @param id - the image ID which should match one in imagedata.ImageData
// @param reason - reason for reporting
func ReportImage(db db.DatabaseInterface, id string, reason string) {
    db.SoftDelete(id, reason)

    // notify through Slack bot
    apiToken := os.Getenv("SLACK_API_TOKEN")
    if len(apiToken) == 0 {
        return
    }
    channelName := "hanserver"
    api := slack.New(apiToken)
    params := slack.PostMessageParameters{}
    _, _, err := api.PostMessage(channelName, fmt.Sprintf("Image %s reported because: %s", id, reason), params)
    if err != nil {
        fmt.Println(err)
        return
    }
}
