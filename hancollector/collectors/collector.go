package collectors

import (
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

// QueryRange is the maximum radius of each query in metres
const QueryRange = 5000

// ImageCollector is an interface used for collecting images
// This should be implemented for each media source
type ImageCollector interface {
    // a configuration must be implemented for each collector
    GetConfig() config.CollectorConfiguration
    GetImages(Lat float64, Lng float64) ([]imagedata.ImageData, error)
}