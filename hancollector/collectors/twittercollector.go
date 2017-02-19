package collectors

import (
    "fmt"
    "time"
    "github.com/dghubble/oauth1"
    "github.com/dghubble/go-twitter/twitter"
    "github.com/kellydunn/golang-geo"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

// TwitterCollector implements the collector interface for Twitter
type TwitterCollector struct {
    timeSinceLastQuery int64
}

// NewTwitterCollector creates a new `TwitterCollector`
func NewTwitterCollector() *TwitterCollector {
    c := new(TwitterCollector)
    c.timeSinceLastQuery = 0
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
    timeSinceLastUpdate := time.Now().Unix() - c.timeSinceLastQuery
    // Twitter is rate limited to 15 minute windows. We'll limit to an hour to
    // be safe
    if timeSinceLastUpdate < 1 * 60 * 60 && timeSinceLastUpdate > 1 {
        return []imagedata.ImageData{}, nil
    }
    c.timeSinceLastQuery = time.Now().Unix()
    // Twitter client setup
    // TODO: couldn't get app auth using oauth2 working
    conf := oauth1.NewConfig(config.TwitterConfig.APIKey, config.TwitterConfig.APISecret)
    token := oauth1.NewToken(config.TwitterConfig.AccessToken, config.TwitterConfig.AccessSecret)
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

func (c *TwitterCollector) queryImages(client *twitter.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
    includeEntities := true
    params := &twitter.SearchTweetParams{
        Query: "filter:images",
        Geocode: fmt.Sprintf("%f,%f,%dkm", lat, lng, QueryRange / 1000),
        IncludeEntities: &includeEntities,
    }
    media, _, err := client.Search.Tweets(params)
    if err != nil {
        // we failed so just return the error
        return []imagedata.ImageData {}, err
    }

    images := []imagedata.ImageData {}
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
            config.TwitterConfig.CollectorName)
        images = append(images, *newImage)
    }
    return images, nil
}
