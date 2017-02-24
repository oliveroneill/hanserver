package imagepopulation

import (
    "testing"
    "time"
    "reflect"
    "sync"
    "errors"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
    "github.com/oliveroneill/hanserver/hancollector/collectors"
    "github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

type MockDB struct {
    Images []imagedata.ImageData
    regions []imagedata.Location
    lock sync.Mutex
}

func NewMockDB(regions []imagedata.Location) *MockDB {
    c := new(MockDB)
    c.regions = regions
    c.Images = []imagedata.ImageData{}
    return c
}

func (c *MockDB) GetRegions() []imagedata.Location {
    return c.regions
}

func (c *MockDB) AddRegion(lat float64, lng float64) {
}

func (c *MockDB) AddImage(image imagedata.ImageData) {
    c.lock.Lock()
    c.Images = append(c.Images, image)
    c.lock.Unlock()
}

func (c *MockDB) GetImages(lat float64, lng float64) []imagedata.ImageData {
    return []imagedata.ImageData{}
}

func (c *MockDB) GetAllImages() []imagedata.ImageData {
    return []imagedata.ImageData{}
}

func (c *MockDB) SoftDelete(id string, reason string) {}

func (c *MockDB) Close() {}

type MockCollector struct {
    collectors.ImageCollector
    images []imagedata.ImageData
    sleepDelay time.Duration
    shouldError bool
}

/**
 * Sleep delay, the amount that GetImages should delay for
 */
func NewMockCollector(sleepDelay time.Duration, images []imagedata.ImageData, shouldError bool) *MockCollector {
    c := &MockCollector{
        ImageCollector: collectors.NewAPIRestrictedCollector(),
    }
    c.images = images
    c.sleepDelay = sleepDelay
    c.shouldError = shouldError
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
        CollectorConfig: config.CollectorConfig{},
    }
}

func (c *MockCollector) GetImages(lat float64, lng float64) ([]imagedata.ImageData, error) {
    if c.sleepDelay > 0 {
        time.Sleep(c.sleepDelay)
    }
    if c.shouldError {
        return nil, errors.New("Mock error")
    }
    return c.images, nil
}

func TestPopulateImageDB(t *testing.T) {
    firstImages := []imagedata.ImageData{
        *imagedata.NewImage("caption string", 10, "", "", "", 56, 33, "", "", "", ""),
        *imagedata.NewImage("dsgjsdk", 104, "", "", "", 5336, 3, "", "", "", ""),
    }
    secondImages := []imagedata.ImageData{
        *imagedata.NewImage("caption string2", 12, "", "", "", 532, 33, "", "", "", ""),
        *imagedata.NewImage("dsgjsdk2", 14, "", "", "", 5336, 3, "", "", "", ""),
    }
    thirdImages := []imagedata.ImageData{
        *imagedata.NewImage("caption string3", 10, "", "", "", 56, 233, "", "", "", ""),
        *imagedata.NewImage("dsgjsdk3", 104, "", "", "", 56, 32, "", "", "", ""),
    }
    collectorArray := []collectors.ImageCollector{
        NewMockCollector(1 * time.Millisecond, []imagedata.ImageData{}, true),
        NewMockCollector(4 * time.Millisecond, thirdImages, false),
        NewMockCollector(0 * time.Millisecond, firstImages, false),
        NewMockCollector(1 * time.Millisecond, secondImages, false),
    }
    mockDB := NewMockDB([]imagedata.Location{})
    region := imagedata.NewLocation(45, 66)
    populateImageDBWithCollectors(mockDB, collectorArray, region.Lat, region.Lng)
    if len(mockDB.Images) != len(firstImages) {
        t.Error("Expected", len(mockDB.Images), "to equal", len(firstImages))
    }
    // set regions so that images are equal
    for i := 0; i < len(firstImages); i++ {
        firstImages[i].Region = region
    }
    for i := 0; i < len(secondImages); i++ {
        secondImages[i].Region = region
    }
    for i := 0; i < len(thirdImages); i++ {
        thirdImages[i].Region = region
    }
    if !reflect.DeepEqual(mockDB.Images, firstImages) {
        t.Error("Expected", mockDB.Images, "to equal", firstImages)
    }
    time.Sleep(5 * time.Millisecond)
    allImages := firstImages
    allImages = append(allImages, secondImages...)
    allImages = append(allImages, thirdImages...)
    if !reflect.DeepEqual(mockDB.Images, allImages) {
        t.Error("Expected", mockDB.Images, "to equal", allImages)
    }
}