package config

import (
	"encoding/json"
	"fmt"
	"time"
)

// CollectionConfig holds configuration information for each collector
type CollectionConfig struct {
	InstagramConfig *InstagramConfiguration `json:"instagram"`
	FlickrConfig    *FlickrConfiguration    `json:"flickr"`
	TwitterConfig   *TwitterConfiguration   `json:"twitter"`
}

// UnmarshalConfig will convert a json string into the CollectionConfig struct
func UnmarshalConfig(jsonString string) CollectionConfig {
	c := CollectionConfig{}
	c.InstagramConfig = InstagramConfig
	c.FlickrConfig = FlickrConfig
	c.TwitterConfig = TwitterConfig
	err := json.Unmarshal([]byte(jsonString), &c)
	if err != nil {
		// TODO: send error back
		fmt.Println(err)
	}
	return c
}

// CollectorConfiguration is the base configuration
type CollectorConfiguration interface {
	IsEnabled() bool
	GetCollectorName() string
	// The frequency in seconds at which the collector should be updated for
	// all regions
	GetUpdateFrequency() time.Duration
	// Limit on query per GetQueryWindow seconds
	GetQueryLimit() int
	// in seconds
	GetQueryWindow() int64
}

// CollectorConfig is a type used for CollectorConfiguration interface
type CollectorConfig struct {
	CollectorConfiguration
	CollectorName   string
	Enabled         bool          `json:"enabled"`
	UpdateFrequency time.Duration `json:"update_frequency"`
	QueryLimit      int           `json:"query_limit"`
	QueryWindow     int64         `json:"query_window"`
}

// IsEnabled if this collector should be used
func (c CollectorConfig) IsEnabled() bool {
	return c.Enabled
}

// GetCollectorName returns the name of the collector for logging purposes
func (c CollectorConfig) GetCollectorName() string {
	return c.CollectorName
}

// GetUpdateFrequency returns the frequency which the collector should be
// updated
func (c CollectorConfig) GetUpdateFrequency() time.Duration {
	return c.UpdateFrequency
}

// GetQueryLimit returns the limit of queries to be made per query window
func (c CollectorConfig) GetQueryLimit() int {
	return c.QueryLimit
}

// GetQueryWindow returns the window to be used with GetQueryLimit
func (c CollectorConfig) GetQueryWindow() int64 {
	return c.QueryWindow
}
