package collectors

import (
	"github.com/oliveroneill/hanserver/hanapi"
	"github.com/oliveroneill/hanserver/hancollector/collectors/config"
	"sync"
	"time"
)

// QueryRange is the maximum radius of each query in metres
const QueryRange = 5000

// ImageCollector is an interface used for collecting images
// This should be implemented for each media source
type ImageCollector interface {
	// a configuration must be implemented for each collector
	GetConfig() config.CollectorConfiguration
	GetImages(Lat float64, Lng float64) ([]hanapi.ImageData, error)
	ableToQuery(config config.CollectorConfiguration) bool
}

// APIRestrictedCollector is an implementation of ImageCollector that monitors
// query calls. This should be extended since it does not implement GetImages
// or GetConfig. See `instagramcollector.go` for example
type APIRestrictedCollector struct {
	receivedError bool
	queryCount    int
	lastQueryTime int64
	mutex         sync.Mutex
}

// NewAPIRestrictedCollector creates a simple implementation of ImageCollector
// that monitors API calls
func NewAPIRestrictedCollector() *APIRestrictedCollector {
	return &APIRestrictedCollector{
		receivedError: false,
		queryCount:    0,
		lastQueryTime: 0,
		mutex:         sync.Mutex{},
	}
}

// ableToQuery will return true if the collector hasn't reached its API limit
// This will assume that a query will be made if the call is true, therefore
// increasing the query count
func (c *APIRestrictedCollector) ableToQuery(config config.CollectorConfiguration) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.queryCount < config.GetQueryLimit() && !c.receivedError {
		c.queryCount++
		c.lastQueryTime = time.Now().Unix()
		return true
	}
	timeSinceLastQuery := time.Now().Unix() - c.lastQueryTime
	if timeSinceLastQuery > config.GetQueryWindow() {
		c.queryCount = 0
		c.receivedError = false
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
func (c *APIRestrictedCollector) GetImages(Lat float64, Lng float64) ([]hanapi.ImageData, error) {
	return []hanapi.ImageData{}, nil
}
