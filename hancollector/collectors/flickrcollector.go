package collectors

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/oliveroneill/flickgo"
	"github.com/oliveroneill/hanserver/hanapi/imagedata"
	"github.com/oliveroneill/hanserver/hancollector/util"
	"github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

// FlickrCollector implements the collector interface for Flickr
type FlickrCollector struct {
	ImageCollector
	config *config.FlickrConfiguration
}

// NewFlickrCollector creates a new `FlickrCollector`
func NewFlickrCollector(config *config.FlickrConfiguration) *FlickrCollector {
	c := &FlickrCollector{
		ImageCollector: NewAPIRestrictedCollector(),
		config: config,
	}
	return c
}

// GetConfig returns the configuration for the Flickr source
// Use this to store api keys and enable/disable collectors
func (c *FlickrCollector) GetConfig() config.CollectorConfiguration {
	return c.config
}

// GetImages returns new images queried by location on Flickr
func (c *FlickrCollector) GetImages(lat float64, lng float64) ([]imagedata.ImageData, error) {
	if !c.GetConfig().IsEnabled() {
		return []imagedata.ImageData{}, nil
	}
	client := flickgo.New(c.config.APIKey, c.config.Secret, http.DefaultClient)
	return c.getImagesWithClient(client, lat, lng)
}

func (c *FlickrCollector) getImagesWithClient(client *flickgo.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
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

func (c *FlickrCollector) queryImages(client *flickgo.Client, lat float64, lng float64) ([]imagedata.ImageData, error) {
	// check that we haven't reached query limits
	if !c.ableToQuery(c.GetConfig()) {
		return []imagedata.ImageData {}, nil
	}
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
		// check that we haven't reached query limits
		if !c.ableToQuery(c.GetConfig()) {
			break
		}
		secret, err := strconv.Atoi(m.Secret)
		// we then have to request for user info and licensing info
		// TODO: this extra request slows everything down
		res, err := client.PhotosGetInfo(flickgo.PhotosGetInfoParams {
			PhotoID: m.ID,
			Secret: secret,
		})
		if err != nil {
			continue
		}
		photoInfo := res.PhotoInfo
		license, err := strconv.ParseFloat(photoInfo.License, 64)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// ensure the license allows us to show it
		// 0 = All Rights Reserved
		if license == 0 {
			continue
		}
		// check that we haven't reached query limits
		if !c.ableToQuery(c.GetConfig()) {
			break
		}
		// we then have to request the exact location
		// TODO: so many requests...
		location, err := client.GetLocation(
			map[string]string {
				"api_key": c.config.APIKey,
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
			c.config.CollectorName)
		images = append(images, *newImage)
	}
	return images, nil
}

