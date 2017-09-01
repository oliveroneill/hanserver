package collectors

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/oliveroneill/hanserver/hanapi/imagedata"
	"github.com/oliveroneill/hanserver/hancollector/collectors/config"
	"github.com/oliveroneill/hanserver/hancollector/util"
	"time"
)

// TwitterCollector implements the collector interface for Twitter
type TwitterCollector struct {
	ImageCollector
	config *config.TwitterConfiguration
}

// NewTwitterCollector creates a new `TwitterCollector`
func NewTwitterCollector(config *config.TwitterConfiguration) *TwitterCollector {
	c := &TwitterCollector{
		ImageCollector: NewAPIRestrictedCollector(),
		config:         config,
	}
	return c
}

// GetConfig returns the configuration for the Twitter source
// Use this to store api keys and enable/disable collectors
func (c *TwitterCollector) GetConfig() config.CollectorConfiguration {
	return config.TwitterConfig
}

// GetImages returns new images queried by location on Twitter
func (c *TwitterCollector) GetImages(lat float64, lng float64) ([]imagedata.ImageData, error) {
	if !c.GetConfig().IsEnabled() {
		return []imagedata.ImageData{}, nil
	}
	// Twitter client setup
	// TODO: couldn't get app auth using oauth2 working
	conf := oauth1.NewConfig(c.config.APIKey, c.config.APISecret)
	token := oauth1.NewToken(c.config.AccessToken, c.config.AccessSecret)
	// http.Client will automatically authorize Requests
	httpClient := conf.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return c.getImagesWithClient(client, lat, lng)
}

func (c *TwitterCollector) getImagesWithClient(client *twitter.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
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

func (c *TwitterCollector) queryImages(client *twitter.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
	// check that we haven't reached query limits
	if !c.ableToQuery(c.GetConfig()) {
		return []imagedata.ImageData{}, nil
	}
	includeEntities := true
	params := &twitter.SearchTweetParams{
		Query:           "filter:images",
		Geocode:         fmt.Sprintf("%f,%f,%dkm", lat, lng, QueryRange/1000),
		IncludeEntities: &includeEntities,
	}
	media, _, err := client.Search.Tweets(params)
	if err != nil {
		// we failed so just return the error
		return []imagedata.ImageData{}, err
	}

	images := []imagedata.ImageData{}
	for _, m := range media.Statuses {
		// if it doesn't have an image then ignore
		if len(m.Entities.Media) == 0 {
			continue
		}
		// parse date to timestamp
		const longForm = "Mon Jan 2 15:04:05 -0700 2006"
		t, err := time.Parse(longForm, m.CreatedAt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// if it's not geotagged then we ignore it
		// TODO: these images are still within the region but due
		// to missing information this data is potentially unhelpful
		if m.Coordinates == nil || len(m.Coordinates.Coordinates) < 2 {
			continue
		}
		newImage := imagedata.NewImage(m.Text, t.Unix(),
			m.Entities.Media[0].MediaURL,
			m.Entities.Media[0].MediaURL, m.IDStr,
			m.Coordinates.Coordinates[1],
			m.Coordinates.Coordinates[0],
			m.Entities.Media[0].DisplayURL,
			m.User.Name, m.User.ProfileImageURL,
			c.config.CollectorName)
		images = append(images, *newImage)
	}
	return images, nil
}
