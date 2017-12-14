package hanapi

import (
	"fmt"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	getImageCollection(session).EnsureIndex(mgo.Index{Key: []string{"$2dsphere:coordinates"}})
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
func (c *MongoInterface) GetRegions() []Location {
	collection := getRegionCollection(c.session)
	var regions []Location
	collection.Find(map[string]interface{}{}).All(&regions)
	return regions
}

// AddRegion adds this new location as a place to query images on
func (c *MongoInterface) AddRegion(lat float64, lng float64) {
	collection := getRegionCollection(c.session)
	collection.Insert(map[string]interface{}{"lat": lat, "lng": lng})
}

// AddImage adds new image data for the feed
func (c *MongoInterface) AddImage(image ImageData) {
	collection := getImageCollection(c.session)
	// insert if it's not already there
	_, err := collection.Upsert(bson.M{"_id": image.ID}, bson.M{"$set": image})
	if err != nil {
		log.Fatal(err)
	}
}

// AddBulkImagesToRegion adds new images in bulk, also setting the region
func (c *MongoInterface) AddBulkImagesToRegion(images []ImageData,
	region *Location) {
	collection := getImageCollection(c.session)
	bulk := collection.Bulk()
	for _, img := range images {
		img.Region = region
		// insert if it's not already there
		bulk.Upsert(bson.M{"_id": img.ID}, bson.M{"$set": img})
	}
	_, err := bulk.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// GetImages returns images closest to the specified location
func (c *MongoInterface) GetImages(lat float64, lng float64, start int, end int) []ImageData {
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
	response := make([]ImageData, 0)
	collection := getImageCollection(c.session)
	// Mongo allows us to aggregate based on distance from the query
	agg := []bson.M{
		bson.M{
			"$geoNear": bson.M{
				"spherical": true,
				"near": bson.M{
					"type":        "Point",
					"coordinates": []float64{lng, lat},
				},
				"distanceField": "distance",
				// ensure that deleted images aren't in here
				"query": map[string]interface{}{"deleted": nil},
				"num":   end,
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
		image := ImageData{}
		success := iter.Next(&image)
		if !success {
			err := iter.Err()
			// if there is no error then we've reached the end of the results
			if err == nil {
				break
			} else {
				// otherwise something went wrong so we skip this image
				continue
			}
		}
		response = append(response, image)
	}
	return response
}

// GetAllImages returns all images stored
func (c *MongoInterface) GetAllImages() []ImageData {
	var response []ImageData
	collection := getImageCollection(c.session)
	err := collection.Find(nil).All(&response)
	if err != nil {
		panic(err)
	}
	return response
}

// SoftDelete will add a delete field to image so it's no longer visible in
// feed
func (c *MongoInterface) SoftDelete(id string, reason string) {
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

// DeleteOldImages will clear `amount` worth of images starting at the oldest
func (c *MongoInterface) DeleteOldImages(amount int) {
	collection := getImageCollection(c.session)
	change := mgo.Change{
		Remove: true,
	}
	query := collection.Find(nil).Sort("createdTime")
	// sort by the oldest images and remove those first
	for i := 0; i < amount && c.Size() > 0; i++ {
		_, err := query.Apply(change, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Size will return the amount of images in the database
func (c *MongoInterface) Size() int {
	collection := getImageCollection(c.session)
	// TODO: check error
	count, _ := collection.Count()
	return count
}

// Copy the interface for added concurrency
func (c *MongoInterface) Copy() DatabaseInterface {
	i := new(MongoInterface)
	i.session = c.session.Copy()
	return i
}

// Close will close the current mongo connection
func (c *MongoInterface) Close() {
	c.session.Close()
}
