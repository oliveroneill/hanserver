package config

// TwitterConfiguration is a Configuration type specifying information about
// Twitter collection
type TwitterConfiguration struct {
    CollectorConfig
    APIKey          string
    APISecret       string
    AccessToken     string
    AccessSecret    string
}

// TwitterConfig is the current configuration
var TwitterConfig = &TwitterConfiguration{
    CollectorConfig: CollectorConfig{},
}

/*
 * Specify all configurable details needed to run this collector in here
 */
func init() {
    TwitterConfig.CollectorConfig.CollectorName = "twitter"
    // easily turn on or off each collector
    TwitterConfig.CollectorConfig.Enabled = false

    // update every minute
    TwitterConfig.CollectorConfig.UpdateFrequency = 1 * 60 * 60
    TwitterConfig.CollectorConfig.QueryWindow = 15 * 60
    TwitterConfig.CollectorConfig.QueryLimit = 150

    TwitterConfig.APIKey = ""
    TwitterConfig.APISecret = ""
    TwitterConfig.AccessToken = ""
    TwitterConfig.AccessSecret = ""
}
