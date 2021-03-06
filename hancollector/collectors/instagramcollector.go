package collectors

import (
	"github.com/gedex/go-instagram/instagram"
	"github.com/oliveroneill/hanserver/hanapi"
	"github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

// InstagramCollector implements the collector interface for Instagram
type InstagramCollector struct {
	*APIRestrictedCollector
	config *config.InstagramConfiguration
}

// NewInstagramCollector creates a new `InstagramCollector`
func NewInstagramCollector(config *config.InstagramConfiguration) *InstagramCollector {
	c := &InstagramCollector{
		APIRestrictedCollector: NewAPIRestrictedCollector(),
		config:                 config,
	}
	return c
}

// GetConfig returns the configuration for the Instagram source
// Use this to store api keys and enable/disable collectors
func (c *InstagramCollector) GetConfig() config.CollectorConfiguration {
	return c.config
}

// GetImages returns new images queried by location on Instagram
func (c *InstagramCollector) GetImages(lat float64, lng float64) ([]hanapi.ImageData, error) {
	if !c.GetConfig().IsEnabled() {
		return []hanapi.ImageData{}, nil
	}
	client := instagram.NewClient(nil)
	client.AccessToken = c.config.AccessToken
	return c.getImagesWithClient(client, lat, lng)
}

func (c *InstagramCollector) getImagesWithClient(client *instagram.Client, lat float64, lng float64) ([]hanapi.ImageData, error) {
	images, err := c.queryImages(client, lat, lng)
	if err != nil {
		return images, err
	}
	points := GetSurroundingPoints(lat, lng, QueryRange)
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

func (c *InstagramCollector) queryImages(client *instagram.Client, lat float64, lng float64) ([]hanapi.ImageData, error) {
	// check that we haven't reached query limits
	if !c.ableToQuery(c.GetConfig()) {
		return []hanapi.ImageData{}, nil
	}
	opt := &instagram.Parameters{
		Lat:      lat,
		Lng:      lng,
		Distance: QueryRange,
	}
	media, _, err := client.Media.Search(opt)
	if err != nil {
		c.APIRestrictedCollector.receivedError = true
		// we failed so just return the error
		return []hanapi.ImageData{}, err
	}

	images := []hanapi.ImageData{}
	for _, m := range media {
		// make sure the caption is not nil
		text := ""
		if m.Caption != nil {
			text = m.Caption.Text
		}
		newImage := hanapi.NewImage(text, m.CreatedTime,
			m.Images.StandardResolution.URL, m.Images.Thumbnail.URL, m.ID,
			m.Location.Latitude, m.Location.Longitude, m.Link,
			m.User.Username, m.User.ProfilePicture,
			c.config.CollectorName)
		images = append(images, *newImage)
	}
	return images, nil
}
