package imagepopulation

import (
	"reflect"
	"sync"
	"testing"
	"time"
	"errors"
	"github.com/oliveroneill/hanserver/hanapi"
	"github.com/oliveroneill/hanserver/hancollector/collectors"
	"github.com/oliveroneill/hanserver/hancollector/collectors/config"
)

type MockDB struct {
	Images  []hanapi.ImageData
	regions []hanapi.Location
	lock    sync.Mutex
}

func NewMockDB(regions []hanapi.Location) *MockDB {
	c := new(MockDB)
	c.regions = regions
	c.Images = []hanapi.ImageData{}
	return c
}

func (c *MockDB) GetRegions() []hanapi.Location {
	return c.regions
}

func (c *MockDB) AddRegion(lat float64, lng float64) {
}

func (c *MockDB) AddImage(image hanapi.ImageData) {
}

func (c *MockDB) AddBulkImagesToRegion(images []hanapi.ImageData,
	region *hanapi.Location) {
	c.lock.Lock()
	for _, image := range images {
		image.Region = region
		c.Images = append(c.Images, image)
	}
	c.lock.Unlock()
}

func (c *MockDB) GetImages(lat float64, lng float64, start int, end int) []hanapi.ImageData {
	return []hanapi.ImageData{}
}

func (c *MockDB) GetAllImages() []hanapi.ImageData {
	return []hanapi.ImageData{}
}

func (c *MockDB) DeleteOldImages(amount int) {}

func (c *MockDB) Size() int {
	return 0
}

func (c *MockDB) SoftDelete(id string, reason string) {}

func (c *MockDB) Copy() hanapi.DatabaseInterface {
	return c
}

func (c *MockDB) Close() {}

type MockCollector struct {
	collectors.ImageCollector
	images      []hanapi.ImageData
	sleepDelay  time.Duration
	shouldError bool
}

/**
 * Sleep delay, the amount that GetImages should delay for
 */
func NewMockCollector(sleepDelay time.Duration, images []hanapi.ImageData, shouldError bool) *MockCollector {
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

func (c *MockCollector) GetImages(lat float64, lng float64) ([]hanapi.ImageData, error) {
	if c.sleepDelay > 0 {
		time.Sleep(c.sleepDelay)
	}
	if c.shouldError {
		return nil, errors.New("Mock error")
	}
	return c.images, nil
}

func TestPopulateImageDB(t *testing.T) {
	firstImages := []hanapi.ImageData{
		*hanapi.NewImage("caption string", 10, "", "", "", 56, 33, "", "", "", ""),
		*hanapi.NewImage("dsgjsdk", 104, "", "", "", 5336, 3, "", "", "", ""),
	}
	secondImages := []hanapi.ImageData{
		*hanapi.NewImage("caption string2", 12, "", "", "", 532, 33, "", "", "", ""),
		*hanapi.NewImage("dsgjsdk2", 14, "", "", "", 5336, 3, "", "", "", ""),
	}
	thirdImages := []hanapi.ImageData{
		*hanapi.NewImage("caption string3", 10, "", "", "", 56, 233, "", "", "", ""),
		*hanapi.NewImage("dsgjsdk3", 104, "", "", "", 56, 32, "", "", "", ""),
	}
	collectorArray := []collectors.ImageCollector{
		NewMockCollector(1*time.Millisecond, []hanapi.ImageData{}, true),
		NewMockCollector(3*time.Millisecond, thirdImages, false),
		NewMockCollector(0*time.Millisecond, firstImages, false),
		NewMockCollector(1*time.Millisecond, secondImages, false),
	}
	mockDB := NewMockDB([]hanapi.Location{})
	region := hanapi.NewLocation(45, 66)
	populateImageDBWithCollectors(mockDB, collectorArray, region.Lat, region.Lng, nil)
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
	// TODO: shouldn't rely on timings
	time.Sleep(10 * time.Millisecond)
	allImages := firstImages
	allImages = append(allImages, secondImages...)
	allImages = append(allImages, thirdImages...)
	if !reflect.DeepEqual(mockDB.Images, allImages) {
		t.Error("Expected", mockDB.Images, "to equal", allImages)
	}
}
