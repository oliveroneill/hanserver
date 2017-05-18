package hanapi

import (
	"testing"
	"reflect"
	"github.com/kellydunn/golang-geo"
	"github.com/oliveroneill/hanserver/hanapi/db"
	"github.com/oliveroneill/hanserver/hanapi/imagedata"
)

type MockDB struct {
	regions []imagedata.Location
	images []imagedata.ImageData
}

func NewMockDB(regions []imagedata.Location, images []imagedata.ImageData) *MockDB {
	c := new(MockDB)
	c.regions = regions
	c.images = images
	return c
}

func (c MockDB) GetRegions() []imagedata.Location {
	return c.regions
}

func (c MockDB) AddRegion(lat float64, lng float64) {
}

func (c MockDB) AddImage(image imagedata.ImageData) {
}

func (c MockDB) AddBulkImagesToRegion(images []imagedata.ImageData,
									  region *imagedata.Location) {
}

func (c MockDB) GetImages(lat float64, lng float64, start int, end int) []imagedata.ImageData {
	if end > len(c.images) {
		end = len(c.images)
	}
	if start > len(c.images) {
		return []imagedata.ImageData{}
	}
	return c.images[start:end]
}

func (c MockDB) GetAllImages() []imagedata.ImageData {
	return c.images
}

func (c MockDB) SoftDelete(id string, reason string) {}

func (c MockDB) Copy() db.DatabaseInterface {
	return c
}

func (c MockDB) Close() {}

/**
 * Sets up DB with regions close to input argument but should
 * never match anything
 */
func setupNonMatchingDB(testRegion *imagedata.Location) (*MockDB) {
	p := geo.NewPoint(testRegion.Lat, testRegion.Lng)
	// a point just out of range
	newPoint := p.PointAtDistanceAndBearing(5.1, 0)
	// a point very out of range
	newPoint2 := p.PointAtDistanceAndBearing(10, 0)
	regions := []imagedata.Location{
		*imagedata.NewLocation(newPoint.Lat(), newPoint.Lng()),
		*imagedata.NewLocation(newPoint2.Lat(), newPoint2.Lng()),
	}
	nonMatchingDB := NewMockDB(regions, []imagedata.ImageData{})
	return nonMatchingDB
}

/**
 * Sets up DB with regions close to input argument and one region
 * that matches, this region is returned as the second return value
 */
func setupMatchingDB(testRegion *imagedata.Location) (*MockDB, *imagedata.Location) {
	p := geo.NewPoint(testRegion.Lat, testRegion.Lng)
	// a point just out of range
	newPoint := p.PointAtDistanceAndBearing(5.1, 0)
	// a point very out of range
	newPoint2 := p.PointAtDistanceAndBearing(10, 0)
	// a point in range
	newPoint3 := p.PointAtDistanceAndBearing(4.5, 0)
	expected := imagedata.NewLocation(newPoint3.Lat(), newPoint3.Lng())
	matchingRegions := []imagedata.Location{
		*imagedata.NewLocation(newPoint.Lat(), newPoint.Lng()),
		*imagedata.NewLocation(newPoint2.Lat(), newPoint2.Lng()),
		*expected,
	}
	matchDB := NewMockDB(matchingRegions, []imagedata.ImageData{})
	return matchDB, expected
}

func TestContainsRegion(t *testing.T) {
	testRegion := imagedata.NewLocation(-35.250327, 149.075300)
	// test that if there are no points within 5km then ContainsRegion is false
	db := setupNonMatchingDB(testRegion)
	if ContainsRegion(db, testRegion.Lat, testRegion.Lng) {
		t.Error("Expected no region match")
	}
	matchDB, _ := setupMatchingDB(testRegion)
	if !ContainsRegion(matchDB, testRegion.Lat, testRegion.Lng) {
		t.Error("Expected region match")
	}
}

func TestGetRegion(t *testing.T) {
	// test that if there are no points within 5km then ContainsRegion is false
	testRegion := imagedata.NewLocation(-35.250327, 149.075300)
	db := setupNonMatchingDB(testRegion)
	if GetRegion(db, testRegion.Lat, testRegion.Lng) != nil {
		t.Error("Expected no region match for GetRegion")
	}
	matchDB, expected := setupMatchingDB(testRegion)
	result := GetRegion(matchDB, testRegion.Lat, testRegion.Lng)
	if result.Lat != expected.Lat || result.Lng != expected.Lng {
		t.Error("Expected", expected, "region, got", result)
	}
}

func TestGetImagesWithRange(t *testing.T) {
	testRegion := imagedata.NewLocation(-35.250327, 149.075300)
	// arbitrary images. ensure that the distance and created time only
	// increase, to avoid the sort reording
	images := []imagedata.ImageData{
		*imagedata.NewImageWithDistance("caption string", 10, "", "", "", testRegion.Lat, testRegion.Lng, 10),
		*imagedata.NewImageWithDistance("testCaption_2", 15, "", "", "", testRegion.Lat, testRegion.Lng, 15),
		*imagedata.NewImageWithDistance("dhfksdj", 100, "", "", "", testRegion.Lat, testRegion.Lng, 100),
		*imagedata.NewImageWithDistance("bla", 200, "", "", "", testRegion.Lat, testRegion.Lng, 200),
	}
	db := NewMockDB([]imagedata.Location{}, images)
	result := GetImagesWithRange(db, testRegion.Lat, testRegion.Lng, 1, 3)
	if len(result) != 2 {
		t.Error("Expected length of result to be 2")
	}
	for i := 0; i < len(result); i++ {
		if !reflect.DeepEqual(result[i], images[i + 1]) {
			t.Error("Expected", result[i], "to equal", images[i + 1])
		}
	}

	// check start specified only
	start := 1
	result = GetImagesWithRange(db, testRegion.Lat, testRegion.Lng, start, -1)
	if len(result) != len(images) - start {
		t.Error("Expected length of result to be", (len(images) - start))
	}
	for i := 0; i < len(result); i++ {
		if !reflect.DeepEqual(result[i], images[i + start]) {
			t.Error("Expected", result[i], "to equal", images[i + start])
		}
	}

	// check end specified only
	end := 2
	result = GetImagesWithRange(db, testRegion.Lat, testRegion.Lng, -1, end)
	if len(result) != end {
		t.Error("Expected length of result to be", end)
	}
	for i := 0; i < end; i++ {
		if !reflect.DeepEqual(result[i], images[i]) {
			t.Error("Expected", result[i], "to equal", images[i])
		}
	}

	// check that it handles the end being greater than the number of images
	end = len(images) + 1
	result = GetImagesWithRange(db, testRegion.Lat, testRegion.Lng, 0, end)
	if len(result) != len(images) {
		t.Error("Expected length of result to be", end)
	}
	for i := 0; i < len(images); i++ {
		if !reflect.DeepEqual(result[i], images[i]) {
			t.Error("Expected", result[i], "to equal", images[i])
		}
	}
}

func TestGetRange(t *testing.T) {
	start := 10
	end := 50
	sampleSize := 100
	r1, r2 := getRange(sampleSize, start, end)
	if r1 != 0 || r2 != sampleSize {
		t.Error("Failed with range", r1, r2)
	}

	start = 10
	end = 150
	sampleSize = 100
	r1, r2 = getRange(sampleSize, start, end)
	if r1 != 0 || r2 != sampleSize * 2 {
		t.Error("Failed with range", r1, r2)
	}

	start = 110
	end = 150
	sampleSize = 100
	r1, r2 = getRange(sampleSize, start, end)
	if r1 != sampleSize || r2 != sampleSize * 2 {
		t.Error("Failed with range", r1, r2)
	}

	start = 110
	end = 250
	sampleSize = 100
	r1, r2 = getRange(sampleSize, start, end)
	if r1 != sampleSize || r2 != sampleSize * 3 {
		t.Error("Failed with range", r1, r2)
	}

	start = 0
	end = 100
	sampleSize = 50
	r1, r2 = getRange(sampleSize, start, end)
	if r1 != 0 || r2 != end {
		t.Error("Failed with range", r1, r2)
	}

}

// test that images are sorted in sample sizes
func TestGetImagesWithRangeAndSampleSize(t *testing.T) {
	sampleSize := 2
	testRegion := imagedata.NewLocation(-35.250327, 149.075300)
	// arbitrary images. ensure that the distance and created time only
	// increase, to avoid the sort reording
	images := []imagedata.ImageData{
		*imagedata.NewImageWithDistance("caption string", 10, "", "", "", testRegion.Lat, testRegion.Lng, 10),
		*imagedata.NewImageWithDistance("testCaption_2", 5, "", "", "", testRegion.Lat, testRegion.Lng, 5),
		*imagedata.NewImageWithDistance("dhfksdj", 1, "", "", "", testRegion.Lat, testRegion.Lng, 1),
		*imagedata.NewImageWithDistance("bla", 200, "", "", "", testRegion.Lat, testRegion.Lng, 200),
	}
	// sorted images, where sample size is 2 -- so the first two images are
	// sorted separately to the second two
	sorted := []imagedata.ImageData{images[1], images[0], images[2], images[3]}
	db := NewMockDB([]imagedata.Location{}, images)
	result := getImagesWithRangeAndSampleSize(db, testRegion.Lat, testRegion.Lng, 0, len(images), sampleSize)
	if len(result) != len(images) {
		t.Error("Expected length of result to be", len(images), "but was", len(result))
	}
	for i := 0; i < len(result); i++ {
		if !reflect.DeepEqual(result[i], sorted[i]) {
			t.Error("Expected", i, ":", result[i], "to equal", sorted[i])
		}
	}
}

// Test that the correct amount of images is returned when the range
// doesn't land on the border of sampleSize
func TestGetImagesWithRangeAndSampleSizeNotOnBorder(t *testing.T) {
	sampleSize := 2
	testRegion := imagedata.NewLocation(-35.250327, 149.075300)
	// arbitrary images. ensure that the distance and created time only
	// increase, to avoid the sort reording
	images := []imagedata.ImageData{
		*imagedata.NewImageWithDistance("caption string", 10, "", "", "", testRegion.Lat, testRegion.Lng, 10),
		*imagedata.NewImageWithDistance("testCaption_2", 5, "", "", "", testRegion.Lat, testRegion.Lng, 5),
		*imagedata.NewImageWithDistance("dhfksdj", 1, "", "", "", testRegion.Lat, testRegion.Lng, 1),
		*imagedata.NewImageWithDistance("bla", 200, "", "", "", testRegion.Lat, testRegion.Lng, 200),
	}
	// sorted images, where sample size is 2 -- so the first two images are
	// sorted separately to the second two
	sorted := []imagedata.ImageData{images[1], images[0], images[2], images[3]}
	db := NewMockDB([]imagedata.Location{}, images)
	result := getImagesWithRangeAndSampleSize(db, testRegion.Lat, testRegion.Lng, 0, 3, sampleSize)
	if len(result) != 3 {
		t.Error("Expected length of result to be", len(images), "but was", len(result))
	}
	for i := 0; i < len(result); i++ {
		if !reflect.DeepEqual(result[i], sorted[i]) {
			t.Error("Expected", i, ":", result[i], "to equal", sorted[i])
		}
	}
}