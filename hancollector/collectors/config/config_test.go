package config

import (
	"fmt"
	"testing"
)

func TestUnmarshalConfig(t *testing.T) {
	queryLimit := 140
	accessToken := "testApiToken"
	json := `{"instagram": {"enabled": true, "query_limit": %d, "access_token": "%s"}}`
	testStr := fmt.Sprintf(json, queryLimit, accessToken)
	result := UnmarshalConfig(testStr)
	instagram := result.InstagramConfig
	if !instagram.IsEnabled() {
		t.Error("Expected instagram config to be enabled")
	}
	q := instagram.GetQueryLimit()
	if q != queryLimit {
		t.Error("Expected query limit to be", queryLimit, "but was", q)
	}
	if instagram.AccessToken != accessToken {
		t.Error("Expected query limit to be", accessToken, "but was", instagram.AccessToken)
	}
	// ensure that default values are used for missing fields
	if InstagramConfig.UpdateFrequency != instagram.GetUpdateFrequency() {
		t.Error("Expected update frequency to be",
			InstagramConfig.UpdateFrequency, "but was",
			instagram.GetUpdateFrequency())
	}
	// test that no other configs are enabled since they're missing from
	// the json
	if result.FlickrConfig.IsEnabled() {
		t.Error("Expected flickr config to be disabled")
	}
}
