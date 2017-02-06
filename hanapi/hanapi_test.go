package hanapi

import (
    "testing"
    "reflect"
    "github.com/kellydunn/golang-geo"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
)

type MockDB struct {
    regions []imagedata.ImageLocation
    images []imagedata.ImageData
}

func NewMockDB(regions []imagedata.ImageLocation, images []imagedata.ImageData) *MockDB {
    c := new(MockDB)
    c.regions = regions
    c.images = images
    return c
}

func (c MockDB) GetRegions() []imagedata.ImageLocation {
    return c.regions
}

func (c MockDB) AddRegion(lat float64, lng float64) {
}

func (c MockDB) AddImage(image imagedata.ImageData) {
}

func (c MockDB) GetImages(lat float64, lng float64) []imagedata.ImageData {
    return c.images
}

func (c MockDB) GetAllImages() []imagedata.ImageData {
    return c.images
}

func TestContainsRegion(t *testing.T) {
    // test that if there are no points within 5km then ContainsRegion is false
    testRegion := imagedata.NewImageLocation(-35.250327, 149.075300)
    p := geo.NewPoint(testRegion.Lat, testRegion.Lng)
    // a point just out of range
    newPoint := p.PointAtDistanceAndBearing(5.1, 0)
    newPoint2 := p.PointAtDistanceAndBearing(10, 0)
    regions := []imagedata.ImageLocation{
        *imagedata.NewImageLocation(newPoint.Lat(), newPoint.Lng()),
        *imagedata.NewImageLocation(newPoint2.Lat(), newPoint2.Lng()),
    }
    db := NewMockDB(regions, []imagedata.ImageData{})
    if ContainsRegion(db, testRegion.Lat, testRegion.Lng) {
        t.Error("Expected no region match")
    }
    // a point in range
    newPoint3 := p.PointAtDistanceAndBearing(4.5, 0)
    matchingRegions := []imagedata.ImageLocation{
        *imagedata.NewImageLocation(newPoint.Lat(), newPoint.Lng()),
        *imagedata.NewImageLocation(newPoint2.Lat(), newPoint2.Lng()),
        *imagedata.NewImageLocation(newPoint3.Lat(), newPoint3.Lng()),
    }
    matchDB := NewMockDB(matchingRegions, []imagedata.ImageData{})
    if !ContainsRegion(matchDB, testRegion.Lat, testRegion.Lng) {
        t.Error("Expected region match")
    }
}

func TestGetRegion(t *testing.T) {
    // test that if there are no points within 5km then ContainsRegion is false
    testRegion := imagedata.NewImageLocation(-35.250327, 149.075300)
    p := geo.NewPoint(testRegion.Lat, testRegion.Lng)
    // a point just out of range
    newPoint := p.PointAtDistanceAndBearing(5.1, 0)
    newPoint2 := p.PointAtDistanceAndBearing(10, 0)
    regions := []imagedata.ImageLocation{
        *imagedata.NewImageLocation(newPoint.Lat(), newPoint.Lng()),
        *imagedata.NewImageLocation(newPoint2.Lat(), newPoint2.Lng()),
    }
    db := NewMockDB(regions, []imagedata.ImageData{})
    if GetRegion(db, testRegion.Lat, testRegion.Lng) != nil {
        t.Error("Expected no region match for GetRegion")
    }
    // a point in range
    newPoint3 := p.PointAtDistanceAndBearing(4.5, 0)
    expected := imagedata.NewImageLocation(newPoint3.Lat(), newPoint3.Lng())
    matchingRegions := []imagedata.ImageLocation{
        *imagedata.NewImageLocation(newPoint.Lat(), newPoint.Lng()),
        *imagedata.NewImageLocation(newPoint2.Lat(), newPoint2.Lng()),
        *expected,
    }
    matchDB := NewMockDB(matchingRegions, []imagedata.ImageData{})
    result := GetRegion(matchDB, testRegion.Lat, testRegion.Lng)
    if result.Lat != expected.Lat || result.Lng != expected.Lng {
        t.Error("Expected", expected, "region, got", result)
    }
}

func TestGetImagesWithRange(t *testing.T) {
    testRegion := imagedata.NewImageLocation(-35.250327, 149.075300)
    // arbitrary images. ensure that the distance and created time only
    // increase, to avoid the sort reording
    images := []imagedata.ImageData{
        *imagedata.NewImageWithDistance("caption string", 10, "", "", "", testRegion.Lat, testRegion.Lng, 10),
        *imagedata.NewImageWithDistance("testCaption_2", 15, "", "", "", testRegion.Lat, testRegion.Lng, 15),
        *imagedata.NewImageWithDistance("dhfksdj", 100, "", "", "", testRegion.Lat, testRegion.Lng, 100),
        *imagedata.NewImageWithDistance("bla", 200, "", "", "", testRegion.Lat, testRegion.Lng, 200),
    }
    db := NewMockDB([]imagedata.ImageLocation{}, images)
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