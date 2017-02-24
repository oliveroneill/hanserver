package collectors

import (
    "github.com/gedex/go-instagram/instagram"
    "github.com/oliveroneill/hanserver/hancollector/util"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

// InstagramCollector implements the collector interface for Instagram
type InstagramCollector struct {
    ImageCollector
}

// NewInstagramCollector creates a new `InstagramCollector`
func NewInstagramCollector() *InstagramCollector {
    c := &InstagramCollector{
        ImageCollector: NewAPIRestrictedCollector(),
    }
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
    points := util.GetSurroundingPoints(lat, lng, QueryRange)
    // continue search until we have at least 100 images
    for i := 0; i < len(points) && len(images) < 100; i++ {
        queryResponse, err := c.queryImages(client, points[i].Lat, points[i].Lng)
        if err != nil {
            continue
        }
        images = append(images, queryResponse...)
    }
    return images, nil
}

func (c InstagramCollector) queryImages(client *instagram.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
    // check that we haven't reached query limits
    if !c.ableToQuery(c.GetConfig()) {
        return []imagedata.ImageData {}, nil
    }
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
            m.User.Username, m.User.ProfilePicture,
            config.InstagramConfig.CollectorName)
        images = append(images, *newImage)
    }
    return images, nil
}
