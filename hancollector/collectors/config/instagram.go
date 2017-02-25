package config

// InstagramConfiguration is a Configuration type specifying information about
// Instagram collection
type InstagramConfiguration struct {
    CollectorConfig
    AccessToken     string
}

// InstagramConfig is the current configuration
var InstagramConfig = &InstagramConfiguration{
    CollectorConfig: CollectorConfig{},
}

/*
 * Specify all configurable details needed to run this collector in here
 */
func init() {
    InstagramConfig.CollectorConfig.CollectorName = "instagram"
    // easily turn on or off each collector
    InstagramConfig.CollectorConfig.Enabled = false

    // update every minute
    InstagramConfig.CollectorConfig.UpdateFrequency = 1 * 60
    InstagramConfig.CollectorConfig.QueryWindow = 1 * 60 * 60
    InstagramConfig.CollectorConfig.QueryLimit = 4500

    InstagramConfig.AccessToken = ""
}
