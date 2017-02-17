package collectors

import (
    "fmt"
    "time"
<<<<<<< HEAD
    "regexp"
=======
>>>>>>> added: google places collector
    "googlemaps.github.io/maps"
    "golang.org/x/net/context"
    "github.com/kellydunn/golang-geo"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

// GooglePlacesCollector implements the collector interface for Google Places
type GooglePlacesCollector struct {
    timeSinceLastQuery int64
}

// NewGooglePlacesCollector creates a new `GooglePlacesCollector`
func NewGooglePlacesCollector() *GooglePlacesCollector {
    c := new(GooglePlacesCollector)
    c.timeSinceLastQuery = 0
    return c
}

// GetConfig returns the configuration for the Google Places source
// Use this to store api keys and enable/disable collectors
<<<<<<< HEAD
func (c *GooglePlacesCollector) GetConfig() config.CollectorConfiguration {
=======
func (c GooglePlacesCollector) GetConfig() config.CollectorConfiguration {
>>>>>>> added: google places collector
    return config.GooglePlacesConfig
}

// GetImages returns new images queried by location on Google Places
<<<<<<< HEAD
func (c *GooglePlacesCollector) GetImages(lat float64, lng float64) ([]imagedata.ImageData, error) {
    if !c.GetConfig().IsEnabled() {
        return []imagedata.ImageData{}, nil
    }
    timeSinceLastUpdate := time.Now().Unix() - c.timeSinceLastQuery
    // due to Google Maps strict query limits, we'll only query every 12 hours
    // Google Places updates very slowly anyway, so this should be fine
    if timeSinceLastUpdate < 12 * 60 * 60 && timeSinceLastUpdate > 1 {
=======
func (c GooglePlacesCollector) GetImages(lat float64, lng float64) ([]imagedata.ImageData, error) {
    if !c.GetConfig().IsEnabled() {
        return []imagedata.ImageData{}, nil
    }
    // due to Google Maps strict query limits, we'll only query every 12 hours
    // Google Places updates very slowly anyway, so this should be fine
    if time.Now().Unix() - c.timeSinceLastQuery < 12 * 60 * 60 {
>>>>>>> added: google places collector
        return []imagedata.ImageData{}, nil
    }
    c.timeSinceLastQuery = time.Now().Unix()
    client, err := maps.NewClient(maps.WithAPIKey(config.GooglePlacesConfig.APIKey))
    if err != nil {
        return []imagedata.ImageData{}, nil
    }
    return c.getImagesWithClient(client, lat, lng)
}

<<<<<<< HEAD
func (c *GooglePlacesCollector) getImagesWithClient(client *maps.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
=======
func (c GooglePlacesCollector) getImagesWithClient(client *maps.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
>>>>>>> added: google places collector
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

<<<<<<< HEAD
func (c *GooglePlacesCollector) queryImages(client *maps.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
=======
func (c GooglePlacesCollector) queryImages(client *maps.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
>>>>>>> added: google places collector
    r := &maps.NearbySearchRequest{
        Location: &maps.LatLng{lat, lng},
        Radius:   QueryRange,
    }

    media, err := client.NearbySearch(context.Background(), r)
    if err != nil {
        // we failed so just return the error
        return []imagedata.ImageData {}, err
    }

    images := []imagedata.ImageData {}
    for _, m := range media.Results {
<<<<<<< HEAD
        if len(m.Photos) == 0 {
            continue
        }
        // Nearby Search only returns one image
        // To retrieve more you can use a Place Details request, but it appears
        // that there are no unique identifiers for images and different
        // photo_reference values can point to the same image
        image := m.Photos[0]
        userName := ""
        link := ""
        // retrieve user name and link for attributions within the app
        if len(image.HTMLAttributions) > 0 {
            contrib := image.HTMLAttributions[0]
            r, err := regexp.Compile(">(.*)<")
            if err == nil {
                results := r.FindStringSubmatch(contrib)
                if len(results) > 1 {
                    userName = results[1]
                }
            }
            r, err = regexp.Compile("href=\"(.*)\"")
            if err == nil {
                results := r.FindStringSubmatch(contrib)
                if len(results) > 1 {
                    link = results[1]
                }
            }
        }
        // format the url to include the photo reference
        // use `config/googleplaces.go` to set up where the URL should point to
        url := fmt.Sprintf("%s?photoreference=%s", config.GooglePlacesConfig.PhotoURL, image.PhotoReference)
        thumbnailURL := fmt.Sprintf("%s&maxwidth=%d", url, 200)
        // choose maximum val for standard res
        standardURL := fmt.Sprintf("%s&maxwidth=%d", url, 1599)
        // For now the place name is used for the identifier to avoid duplicates
        id := m.Name
        // TODO: using the current timestamp may give a bias to google places
        newImage := imagedata.NewImage(m.Name, time.Now().Unix(),
            standardURL, thumbnailURL, id,
            m.Geometry.Location.Lat, m.Geometry.Location.Lng,
            link, userName, "", config.GooglePlacesConfig.CollectorName)
        images = append(images, *newImage)
=======
        for _, image := range m.Photos {
            // format the url to include the photo reference
            // use config/googleplaces.go to set up where the URL should point to
            url := fmt.Sprintf("%s?photoreference=%s", config.GooglePlacesConfig.PhotoURL, image.PhotoReference)
            // TODO: using the current timestamp may give a bias to google places
            newImage := imagedata.NewImage(m.Name, time.Now().Unix(),
                url, url, image.PhotoReference, m.Geometry.Location.Lat,
                m.Geometry.Location.Lng, "", "", "",
                config.GooglePlacesConfig.CollectorName)
            images = append(images, *newImage)
        }
>>>>>>> added: google places collector
    }
    return images, nil
}
