package collectors

import (
    "fmt"
    "time"
    "strconv"
    "net/http"
    "github.com/manki/flickgo"
    "github.com/kellydunn/golang-geo"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

// FlickrCollector implements the collector interface for Flickr
type FlickrCollector struct {
    lastUpdateTime int64
}

// NewFlickrCollector creates a new `FlickrCollector`
func NewFlickrCollector() *FlickrCollector {
    c := new(FlickrCollector)
    c.lastUpdateTime = 0
    return c
}

// GetConfig returns the configuration for the Flickr source
// Use this to store api keys and enable/disable collectors
func (c FlickrCollector) GetConfig() config.CollectorConfiguration {
    return config.FlickrConfig
}

// GetImages returns new images queried by location on Flickr
func (c *FlickrCollector) GetImages(lat float64, lng float64) ([]imagedata.ImageData, error) {
    if !c.GetConfig().IsEnabled() {
        return []imagedata.ImageData{}, nil
    }
    // Only update every hour, due to having to request each image location separately
    timeSinceLastUpdate := time.Now().Unix() - c.lastUpdateTime
    // here we allow 1 second overlap, in case one region has just started updating
    if timeSinceLastUpdate < 1 * 60 * 60 && timeSinceLastUpdate > 1 {
        return []imagedata.ImageData{}, nil
    }
    c.lastUpdateTime = time.Now().Unix()
    client := flickgo.New(config.FlickrConfig.APIKey, config.FlickrConfig.Secret, http.DefaultClient)
    return c.getImagesWithClient(client, lat, lng)
}

func (c *FlickrCollector) getImagesWithClient(client *flickgo.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
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

func (c *FlickrCollector) queryImages(client *flickgo.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
    request := map[string]string {
        "api_key": config.FlickrConfig.APIKey,
        "lat": fmt.Sprintf("%f", lat),
        "lon": fmt.Sprintf("%f", lng),
        "per_page": fmt.Sprintf("%d", 500),
    }
    response, err := client.Search(request)
    if err != nil {
        // we failed so just return the error
        return []imagedata.ImageData {}, err
    }

    images := []imagedata.ImageData {}
    for _, m := range response.Photos {

        // we then have to request the exact location
        // TODO: this extra request slows everything down
        location, err := client.GetLocation(
            map[string]string {
                "api_key": config.FlickrConfig.APIKey,
                "photo_id": m.ID,
            },
        )
        if err != nil {
            fmt.Println(err)
            continue
        }

        url := fmt.Sprintf("https://farm%s.staticflickr.com/%s/%s_%s",
            m.Farm, m.Server, m.ID, m.Secret)
        // add the extension, this will be formatted using Sprintf
        url += "_%s.jpg"
        userLink := fmt.Sprintf("https://www.flickr.com/photos/%s/%s", m.Owner, m.ID)
        // convert location to floats
        lat, err := strconv.ParseFloat(location.Location.Latitude, 64)
        if err != nil {
            continue
        }
        lng, err := strconv.ParseFloat(location.Location.Longitude, 64)
        if err != nil {
            continue
        }
        newImage := imagedata.NewImage(m.Title, time.Now().Unix(),
            fmt.Sprintf(url, "b"), fmt.Sprintf(url, "t"), m.ID,
            lat, lng, userLink, "Flickr", "",
            config.FlickrConfig.CollectorName)
        images = append(images, *newImage)
    }
    return images, nil
}
