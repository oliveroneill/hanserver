package collectors

import (
    "github.com/gedex/go-instagram/instagram"
    "github.com/kellydunn/golang-geo"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

// QueryRange is the maximum radius of Instagram queries
const QueryRange = 5000

// InstagramCollector implements the collector interface for Instagram
type InstagramCollector struct {
}

// NewInstagramCollector creates a new `InstagramCollector`
func NewInstagramCollector() *InstagramCollector {
    c := new(InstagramCollector)
    return c
}

// GetConfig returns the configuration for the Instagram source
// Use this to store api keys and enable/disable collectors
func (c InstagramCollector) GetConfig() config.CollectorConfiguration {
    return config.InstagramConfig
}

// GetImages returns new images queried by location on Instagram
func (c InstagramCollector) GetImages(lat float64, lng float64) ([]imagedata.ImageData, error) {
    if !c.GetConfig().IsEnabled() {
        return []imagedata.ImageData{}, nil
    }
    client := instagram.NewClient(nil)
    client.AccessToken = config.InstagramConfig.AccessToken
    return c.getImagesWithClient(client, lat, lng)
}

func (c InstagramCollector) getImagesWithClient(client *instagram.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
    images, err := c.queryImages(client, lat, lng)
    if err != nil {
        return images, err
    }
    // continue search until we have at least 100 images
    for degrees := float64(0); degrees < 360 && len(images) < 100; degrees += 90 {
        // search 5 kilometers in each direction
        p := geo.NewPoint(lat, lng)
        // find another point that's at the edge of the previous query
        newPoint := p.PointAtDistanceAndBearing(QueryRange / 1000, degrees)
        queryResponse, err := c.queryImages(client, newPoint.Lat(), newPoint.Lng())
        if err != nil {
            continue
        }
        images = append(images, queryResponse...)
    }
    return images, nil
}

func (c InstagramCollector) queryImages(client *instagram.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
    opt := &instagram.Parameters{
        Lat: lat,
        Lng: lng,
        Distance: QueryRange,
    }
    media, _, err := client.Media.Search(opt)
    if err != nil {
        // we failed so just return the error
        return []imagedata.ImageData {}, err
    }

    images := []imagedata.ImageData {}
    for _, m := range media {
        // make sure the caption is not nil
        text := ""
        if m.Caption != nil {
            text = m.Caption.Text
        }
        newImage := imagedata.NewImage(text, m.CreatedTime,
            m.Images.StandardResolution.URL, m.Images.Thumbnail.URL, m.ID,
            m.Location.Latitude, m.Location.Longitude, m.Link,
            m.User.Username, m.User.ProfilePicture)
        images = append(images, *newImage)
    }
    return images, nil
}

