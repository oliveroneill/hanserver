package config

// FlickrConfiguration is a Configuration type specifying information about
// Flickr collection
type FlickrConfiguration struct {
    Enabled        bool
    APIKey         string
    Secret         string
    CollectorName  string
}

// FlickrConfig is the current configuration
var FlickrConfig = new(FlickrConfiguration)

/*
 * Specify all configurable details needed to run this collector in here
 */
func init() {
    FlickrConfig.CollectorName = "flickr"
    // easily turn on or off each collector
    FlickrConfig.Enabled = false
    // could be retrieved via json etc.
    FlickrConfig.APIKey = ""
    FlickrConfig.Secret = ""
}

// IsEnabled if flickr collection should be used
func (c FlickrConfiguration) IsEnabled() bool {
    return c.Enabled
}

// GetCollectorName returns the name of the collector for logging purposes
func (c FlickrConfiguration) GetCollectorName() string {
    return c.CollectorName
}