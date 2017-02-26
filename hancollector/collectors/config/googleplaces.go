package config

// GooglePlacesConfiguration is a Configuration type specifying information about
// Google Places collection
type GooglePlacesConfiguration struct {
    CollectorConfig
    APIKey         string
    // Use this since Google Places requires a separate request for each image
    // This offers the ability to proxy or store images and avoid
    // revealing your Google Maps API key
    // See the bottom of `hanhttpserver/hanserver.go` for an example
    // of a proxy server used to keep the API key safe.
    // Use `{protocol}://{ipaddress}:80/maps/api/place/photo` as the value of PhotoUrl for
    // this work with `hanserver.go`
    PhotoURL      string
}

// GooglePlacesConfig is the current configuration
var GooglePlacesConfig = &GooglePlacesConfiguration{
    CollectorConfig: CollectorConfig{},
}

func init() {
    GooglePlacesConfig.CollectorConfig.CollectorName = "google-places"
    GooglePlacesConfig.CollectorConfig.Enabled = false

    // update every 12 hours
    GooglePlacesConfig.CollectorConfig.UpdateFrequency = 12 * 60 * 60
    GooglePlacesConfig.CollectorConfig.QueryWindow = 24 * 60 * 60
    GooglePlacesConfig.CollectorConfig.QueryLimit = 900

    GooglePlacesConfig.APIKey = ""
    GooglePlacesConfig.PhotoURL = ""
}
