package db

import (
    "log"
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
)

// MongoInterface - a mongodb implementation of `DatabaseInterface`
type MongoInterface struct {
    DatabaseInterface
    session *mgo.Session
}

// NewMongoInterface - use to create a new mongo connection
func NewMongoInterface() DatabaseInterface {
    c := new(MongoInterface)
    // use for Docker:
    session, err := mgo.Dial("mongodb")
    // use locally:
    // session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }
    c.session = session
    // if geospatial index hasn't been set up this will create it
    getImageCollection(session).EnsureIndex(mgo.Index{Key: []string{"coordinates:2dsphere"}})
    return c
}

func getHanDB(session *mgo.Session) *mgo.Database {
    return session.DB("han")
}

func getRegionCollection(session *mgo.Session) *mgo.Collection {
    return getHanDB(session).C("regions")
}

func getImageCollection(session *mgo.Session) *mgo.Collection {
    return getHanDB(session).C("images")
}

// GetRegions returns the watched locations that are stored in the database
// These locations are queried to populate the database with images
func (c MongoInterface) GetRegions() []imagedata.Location {
    collection := getRegionCollection(c.session)
    var regions []imagedata.Location
    collection.Find(map[string]interface{}{}).All(&regions)
    return regions
}

// AddRegion adds this new location as a place to query images on
func (c MongoInterface) AddRegion(lat float64, lng float64) {
    collection := getRegionCollection(c.session)
    collection.Insert(map[string]interface{}{ "lat": lat, "lng": lng })
}

// AddImage adds new image data for the feed
func (c MongoInterface) AddImage(image imagedata.ImageData) {
    collection := getImageCollection(c.session)
    // insert if it's not already there
    _, err := collection.Upsert(bson.M{ "_id": image.ID }, bson.M{"$set":image})
    if err != nil {
        log.Fatal(err)
    }
}

// GetImages returns images closest to the specified location
func (c MongoInterface) GetImages(lat float64, lng float64, start int, end int) []imagedata.ImageData {
    if start == -1 {
        start = 0
    }
    // if end is unspecified then we'll only return 100 images
    if end == -1 {
        end = start + 100
    }
    // TODO: may need some limit here to avoid queries taking long amounts of
    // time, however mongo seems quite fast at this
    // convert to response data
    var response []imagedata.ImageData
    collection := getImageCollection(c.session)
    // Mongo allows us to aggregate based on distance from the query
    agg := []bson.M{
        bson.M{
            "$geoNear": bson.M{
                "spherical": true,
                "near": bson.M{
                    "type": "Point",
                    "coordinates": []float64{lng, lat},
                },
                "distanceField": "distance",
                // ensure that deleted images aren't in here
                "query": map[string]interface{}{"deleted": nil},
                "num": end,
            },
        },
    }
    iter := collection.Pipe(agg).Iter()
    for i := 0; i < start; i++ {
        // we throw away these values but passing in nil seems to break the
        // iter
        unused := map[string]interface{}{}
        iter.Next(&unused)
    }
    for i := start; i < end; i++ {
        image := imagedata.ImageData{}
        iter.Next(&image)
        response = append(response, image)
    }
    return response
}

// GetAllImages returns all images stored
func (c MongoInterface) GetAllImages() []imagedata.ImageData {
    var response []imagedata.ImageData
    collection := getImageCollection(c.session)
    err := collection.Find(nil).All(&response)
    if err != nil {
        panic(err)
    }
    return response
}

// SoftDelete will add a delete field to image so it's no longer visible in
// feed
func (c MongoInterface) SoftDelete(id string, reason string) {
    collection := getImageCollection(c.session)
    // update the image with a "deleted" field
    err := collection.UpdateId(
        id,
        bson.M{"$set": bson.M{"deleted": true, "deleted_reason": reason}},
    )
    if err != nil {
        fmt.Println(err)
    }
}

// Close will close the current mongo connection
func (c MongoInterface) Close() {
    c.session.Close()
}