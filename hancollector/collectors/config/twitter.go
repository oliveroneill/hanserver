package config

// TwitterConfiguration is a Configuration type specifying information about
// Twitter collection
type TwitterConfiguration struct {
    Enabled       bool
    APIKey        string
    APISecret     string
    AccessToken   string
    AccessSecret  string
    CollectorName string
}

// TwitterConfig is the current configuration
var TwitterConfig = new(TwitterConfiguration)

/*
 * Specify all configurable details needed to run this collector in here
 */
func init() {
    TwitterConfig.CollectorName = "twitter"
    // easily turn on or off each collector
    TwitterConfig.Enabled = false
    // could be retrieved via json etc.
    TwitterConfig.APIKey = ""
    TwitterConfig.APISecret = ""
    TwitterConfig.AccessToken = ""
    TwitterConfig.AccessSecret = ""
}

// IsEnabled determines whether the collector is used or not
// Unfortunately you have to implement this method every time. Go does not
// allow you to return an inherited struct as the same type as the original.
// So I was forced to use an interface instead.
func (c TwitterConfiguration) IsEnabled() bool {
    return c.Enabled
}

// GetCollectorName returns the name of the collector for logging purposes
func (c TwitterConfiguration) GetCollectorName() string {
    return c.CollectorName
}