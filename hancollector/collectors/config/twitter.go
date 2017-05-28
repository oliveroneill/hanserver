package config

// TwitterConfiguration is a Configuration type specifying information about
// Twitter collection
type TwitterConfiguration struct {
	CollectorConfig
	APIKey		 string `json:"api_key"`
	APISecret	 string `json:"api_secret"`
	AccessToken	 string `json:"access_token"`
	AccessSecret string `json:"access_secret"`
}

// TwitterConfig is the current configuration
var TwitterConfig = &TwitterConfiguration{
	CollectorConfig: CollectorConfig{},
}

func init() {
	TwitterConfig.CollectorConfig.CollectorName = "twitter"
	// easily turn on or off each collector
	TwitterConfig.CollectorConfig.Enabled = false

	// update every minute
	TwitterConfig.CollectorConfig.UpdateFrequency = 1 * 60 * 60
	TwitterConfig.CollectorConfig.QueryWindow = 15 * 60
	TwitterConfig.CollectorConfig.QueryLimit = 150
}
