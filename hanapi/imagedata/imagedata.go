package imagedata

// Location is a type which contains the lat and lng
type Location struct {
    Lat float64 `json:"lat" bson:"lat"`
    Lng float64 `json:"lng" bson:"lng"`
}

// User is a type used to keep track of the source of an image
type User struct {
    Username string `json:"username" bson:"username"`
    ProfilePictureURL string `json:"profile_picture" bson:"profile_picture"`
}

// ImageData is data that is stored and returned from `hanapi`
type ImageData struct {
    Caption       string `json:"caption" bson:"caption"`
    CreatedTime   int64 `json:"createdTime" bson:"createdTime"`
    ImageURL      string `json:"url" bson:"url"`
    Link          string `json:"link" bson:"link"`
    User          *User `json:"user" bson:"user"`
    ThumbnailURL  string `json:"thumbnail_url" bson:"thumbnail_url"`
    ID            string `json:"id" bson:"_id"`
    Location      *Location `json:"location" bson:"location"`
    // regions are specified imageops.go
    Region        *Location `json:"region" bson:"region"`
    // where the photo was taken
    Coordinates   []float64 `json:"coordinates" bson:"coordinates"`
    // will be set when querying using DatabaseInterface
    Distance      float64 `json:"distance" bson:"distance"`
    // the source of the image
    Source        string `json:"source" bson:"source"`
}

// NewLocation returns a new location
func NewLocation(lat float64, lng float64) *Location {
    loc := new(Location)
    loc.Lat = lat
    loc.Lng = lng
    return loc
}

// NewUser returns new user
func NewUser(username string, profileURL string) *User {
    user := new(User)
    user.Username = username
    user.ProfilePictureURL = profileURL
    return user
}

// NewImage returns a new image that's suitable for being added to the database.
// Note that region is not specified here, this is done before entry
// into the database. This is because the collectors have no real idea
// of which query belongs to which region
func NewImage(caption string, createdTime int64, imageURL string,
    thumbnailURL string, id string, lat float64, lng float64, link string,
    user string, profilePictureURL string, source string) *ImageData {
    i := new(ImageData)
    i.Caption = caption
    i.CreatedTime = createdTime
    i.ImageURL = imageURL
    i.ThumbnailURL = thumbnailURL
    i.ID = id
    i.Location = NewLocation(lat, lng)
    i.Coordinates = []float64{lng, lat}
    i.Link = link
    i.User = NewUser(user, profilePictureURL)
    i.Source = source
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
    i.Location = NewLocation(lat, lng)
    i.Coordinates = []float64{lng, lat}
    i.Distance = distance
    return i
}
