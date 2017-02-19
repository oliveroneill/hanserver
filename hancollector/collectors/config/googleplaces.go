package config

// GooglePlacesConfiguration is a Configuration type specifying information about
// Google Places collection
type GooglePlacesConfiguration struct {
    Enabled        bool
    APIKey         string
    CollectorName  string
    // Use this since Google Places requires a separate request for each image
    // This offers the ability to proxy or store images and avoid
    // revealing your Google Maps API key
    // See the bottom of `hanhttpserver/hanserver.go` for an example
    // of a proxy server used to keep the API key safe.
    // Use `{ipaddress}:80/maps/api/place/photo` as the value of PhotoUrl for
    // this work with `hanserver.go`
    // Note: don't include http:// on the start of this url as it's added in
    // within the collector
    PhotoURL      string
}

// GooglePlacesConfig is the current configuration
var GooglePlacesConfig = new(GooglePlacesConfiguration)

func init() {
    GooglePlacesConfig.CollectorName = "google-places"
    GooglePlacesConfig.Enabled = false
    GooglePlacesConfig.APIKey = ""
    GooglePlacesConfig.PhotoURL = ""
}

// IsEnabled determines whether the collector is used or not
func (c GooglePlacesConfiguration) IsEnabled() bool {
    return c.Enabled
}