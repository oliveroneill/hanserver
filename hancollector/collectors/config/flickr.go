package config

// FlickrConfiguration is a Configuration type specifying information about
// Flickr collection
type FlickrConfiguration struct {
	CollectorConfig
	APIKey string `json:"api_key"`
	Secret string `json:"secret"`
}

// FlickrConfig is the current configuration
var FlickrConfig = &FlickrConfiguration{
	CollectorConfig: CollectorConfig{},
}

func init() {
	FlickrConfig.CollectorConfig.CollectorName = "flickr"
	// easily turn on or off each collector
	FlickrConfig.CollectorConfig.Enabled = false

	// update every hour
	FlickrConfig.CollectorConfig.UpdateFrequency = 1 * 60 * 60
	FlickrConfig.CollectorConfig.QueryWindow = 1 * 60 * 60
	FlickrConfig.CollectorConfig.QueryLimit = 3000
}
