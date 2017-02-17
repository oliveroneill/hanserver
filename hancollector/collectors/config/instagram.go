package config

// InstagramConfiguration is a Configuration type specifying information about
// Instagram collection
type InstagramConfiguration struct {
    Enabled        bool
    AccessToken    string
    CollectorName  string
}

// InstagramConfig is the current configuration
var InstagramConfig = new(InstagramConfiguration)

/*
 * Specify all configurable details needed to run this collector in here
 */
func init() {
    InstagramConfig.CollectorName = "instagram"
    // easily turn on or off each collector
    InstagramConfig.Enabled = false
    // could be retrieved via json etc.
    InstagramConfig.AccessToken = ""
}

// IsEnabled determines whether the collector is used or not
// Unfortunately you have to implement this method every time. Go does not
// allow you to return an inherited struct as the same type as the original.
// So I was forced to use an interface instead.
func (c InstagramConfiguration) IsEnabled() bool {
    return c.Enabled
}