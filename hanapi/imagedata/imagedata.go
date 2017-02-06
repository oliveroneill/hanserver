package imagedata

// ImageLocation is a type which contains the lat and lng
type ImageLocation struct {
    Lat float64 `json:"lat" bson:"lat"`
    Lng float64 `json:"lng" bson:"lng"`
}

// ImageData is data that is stored and returned from `hanapi`
type ImageData struct {
    Caption      string `json:"caption" bson:"caption"`
    CreatedTime  int64 `json:"createdTime" bson:"createdTime"`
    ImageURL     string `json:"url" bson:"url"`
    ThumbnailURL string `json:"thumbnail_url" bson:"thumbnail_url"`
    ID           string `json:"id" bson:"_id"`
    Location     *ImageLocation `json:"location" bson:"location"`
    // regions are specified imageops.go
    Region       *ImageLocation `json:"region" bson:"region"`
    // where the photo was taken
    Coordinates  []float64 `json:"coordinates" bson:"coordinates"`
    // will be set when querying using DatabaseInterface
    Distance     float64 `json:"distance" bson:"distance"`
}

// NewImageLocation returns a new location
func NewImageLocation(lat float64, lng float64) *ImageLocation {
    loc := new(ImageLocation)
    loc.Lat = lat
    loc.Lng = lng
    return loc
}

// NewImage returns a new image that's suitable for being added to the database.
// Note that region is not specified here, this is done before entry
// into the database. This is because the collectors have no real idea
// of which query belongs to which region
func NewImage(caption string, createdTime int64, imageURL string,
    thumbnailURL string, id string, lat float64, lng float64) *ImageData {
    i := new(ImageData)
    i.Caption = caption
    i.CreatedTime = createdTime
    i.ImageURL = imageURL
    i.ThumbnailURL = thumbnailURL
    i.ID = id
    i.Location = NewImageLocation(lat, lng)
    i.Coordinates = []float64{lng, lat}
    return i
}

// NewImageWithDistance returns a new image with distance specified
// Created purely for testing purposes, so that the distance can be specified
func NewImageWithDistance(caption string, createdTime int64, imageURL string,
    thumbnailURL string, id string, lat float64, lng float64, distance float64) *ImageData {
    i := new(ImageData)
    i.Caption = caption
    i.CreatedTime = createdTime
    i.ImageURL = imageURL
    i.ThumbnailURL = thumbnailURL
    i.ID = id
    i.Location = NewImageLocation(lat, lng)
    i.Coordinates = []float64{lng, lat}
    i.Distance = distance
    return i
}
