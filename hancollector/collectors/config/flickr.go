package config

// FlickrConfiguration is a Configuration type specifying information about
// Flickr collection
type FlickrConfiguration struct {
    CollectorConfig
    APIKey          string
    Secret          string
}

// FlickrConfig is the current configuration
var FlickrConfig = &FlickrConfiguration{
    CollectorConfig: CollectorConfig{},
}

/*
 * Specify all configurable details needed to run this collector in here
 */
func init() {
    FlickrConfig.CollectorConfig.CollectorName = "flickr"
    // easily turn on or off each collector
    FlickrConfig.CollectorConfig.Enabled = true

    // update every hour
    FlickrConfig.CollectorConfig.UpdateFrequency = 1 * 60 * 60
    FlickrConfig.CollectorConfig.QueryWindow = 1 * 60 * 60
    FlickrConfig.CollectorConfig.QueryLimit = 3000

    FlickrConfig.APIKey = ""
    FlickrConfig.Secret = ""
}