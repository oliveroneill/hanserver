package collectors

import (
	"sync"
	"time"
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
	ableToQuery(config config.CollectorConfiguration) bool
}

// APIRestrictedCollector is an implementation of ImageCollector that monitors
// query calls. This should be extended since it does not implement GetImages
// or GetConfig. See `instagramcollector.go` for example
type APIRestrictedCollector struct {
	queryCount	  int
	lastQueryTime int64
	mutex		  sync.Mutex
}

// NewAPIRestrictedCollector creates a simple implementation of ImageCollector
// that monitors API calls
func NewAPIRestrictedCollector() ImageCollector {
	return &APIRestrictedCollector{
		queryCount: 0,
		lastQueryTime: 0,
		mutex: sync.Mutex{},
	}
}

// ableToQuery will return true if the collector hasn't reached its API limit
// This will assume that a query will be made if the call is true, therefore
// increasing the query count
func (c *APIRestrictedCollector) ableToQuery(config config.CollectorConfiguration) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.queryCount < config.GetQueryLimit() {
		c.queryCount++
		c.lastQueryTime = time.Now().Unix()
		return true
	}
	timeSinceLastQuery := time.Now().Unix() - c.lastQueryTime
	if timeSinceLastQuery > config.GetQueryWindow() {
		c.queryCount = 0
		c.lastQueryTime = time.Now().Unix()
		return true
	}
	return false
}

// GetConfig placeholder method to be overriden
func (c *APIRestrictedCollector) GetConfig() config.CollectorConfiguration {
	return config.CollectorConfig{}
}

// GetImages placeholder method to be overriden
func (c *APIRestrictedCollector) GetImages(Lat float64, Lng float64) ([]imagedata.ImageData, error) {
	return []imagedata.ImageData{}, nil
}