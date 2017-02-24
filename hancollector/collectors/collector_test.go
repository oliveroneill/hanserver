package collectors

import (
    "testing"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

type MockCollector struct {
    ImageCollector
    queryLimit int
}

func NewMockCollector(queryLimit int) ImageCollector {
    c := &MockCollector{
        ImageCollector: NewAPIRestrictedCollector(),
    }
    c.queryLimit = queryLimit
    return c
}

type MockConfig struct {
    config.CollectorConfig
}
func (c *MockConfig) IsEnabled() bool {
    return true
}
func (c *MockConfig) GetCollectorName() string {
    return ""
}

func (c *MockCollector) GetConfig() config.CollectorConfiguration {
    return &MockConfig{
        CollectorConfig: config.CollectorConfig{
            QueryLimit: c.queryLimit,
            QueryWindow: 60 * 60,
        },
    }
}

func (c *MockCollector) GetImages(lat float64, lng float64) ([]imagedata.ImageData, error) {
    return []imagedata.ImageData{}, nil
}

func TestAbleToQuery(t *testing.T) {
    queryLimit := 5
    collector := NewMockCollector(queryLimit)
    for i := 0; i < queryLimit; i++ {
        // make queries until reaching limit
        if !collector.ableToQuery(collector.GetConfig()) {
            t.Error("Expected to be able to query on", i)
        }
    }
    // check that we can't query anymore
    if collector.ableToQuery(collector.GetConfig()) {
        t.Error("Expected not to be able to query after reaching limit")
    }
}