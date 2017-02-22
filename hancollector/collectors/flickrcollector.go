package collectors

import (
    "fmt"
    "time"
    "strconv"
    "net/http"
    "github.com/oliveroneill/flickgo"
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
    return c.queryImages(client, lat, lng)
}

func (c *FlickrCollector) queryImages(client *flickgo.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
    request := flickgo.PhotosSearchParams{
        Lat: fmt.Sprintf("%f", lat),
        Lon: fmt.Sprintf("%f", lng),
        PerPage: 500,
    }
    response, err := client.PhotosSearch(request)
    if err != nil {
        // we failed so just return the error
        return []imagedata.ImageData {}, err
    }

    images := []imagedata.ImageData {}
    for _, m := range response.Photos {
        secret, err := strconv.Atoi(m.Secret)
        // we then have to request for user info and licensing info
        // TODO: this extra request slows everything down
        res, err := client.PhotosGetInfo(flickgo.PhotosGetInfoParams {
            PhotoID: m.ID,
            Secret: secret,
        })
        photoInfo := res.PhotoInfo
        license, err := strconv.ParseFloat(photoInfo.License, 64)
        if err != nil {
            fmt.Println(err)
            continue
        }
        // ensure the license allows us to show it
        if license == 0 {
            continue
        }
        // we then have to request the exact location
        // TODO: so many requests...
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
        createdAt, err := strconv.ParseInt(photoInfo.DateUploaded, 0, 64)
        if err != nil {
            continue
        }
        newImage := imagedata.NewImage(m.Title, createdAt,
            fmt.Sprintf(url, "b"), fmt.Sprintf(url, "t"), m.ID,
            lat, lng, userLink, photoInfo.Owner.UserName, "",
            config.FlickrConfig.CollectorName)
        images = append(images, *newImage)
    }
    return images, nil
}

