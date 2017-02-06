package hanapi

import (
    "sort"
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
func GetRegion(db db.DatabaseInterface, lat float64, lng float64) *imagedata.ImageLocation {
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
func GetRegions(db db.DatabaseInterface) []imagedata.ImageLocation {
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
    images := db.GetImages(lat, lng)
    // sort
    sort.Sort(feedsort.BySum(images))
    return images
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
    // we have to get all images and shrink, because the sort happens after
    // the query unfortunately
    images := GetImages(db, lat, lng)
    // return empty if we're out of range
    if start > len(images) {
        return []imagedata.ImageData{}
    }
    // set relevant values since start and end can be optional
    if start < 0 {
        start = 0
    }
    if end < 0 || end > len(images) {
        end = len(images)
    }
    return images[start:end]
}