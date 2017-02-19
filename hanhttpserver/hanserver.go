package main

import (
    "os"
    "fmt"
    "flag"
    "strconv"
    "encoding/json"
    "net/url"
    "net/http"
    "net/http/httputil"
    "github.com/oliveroneill/hanserver/hanapi"
    "github.com/oliveroneill/hanserver/hanapi/db"
    "github.com/oliveroneill/hanserver/hancollector/imagepopulation"
    "github.com/oliveroneill/hanserver/hanhttpserver/response"
)

// HanServer is a http server that also populates the database periodically
// This allows easy tracking of API usage
type HanServer struct {
    populator *imagepopulation.ImagePopulator
    db        db.DatabaseInterface
}

// NewHanServer will create a new http server and start population
// @param noCollection - set this to true if you don't want hancollector to start
func NewHanServer(noCollection bool) *HanServer {
    // this database session is kept onto over the lifetime of the server
    db := db.NewMongoInterface()
    populator := imagepopulation.NewImagePopulator()
    if !noCollection {
        fmt.Println("Starting image collection")
        // populate image db in the background
        go populator.PopulateImageDB(db)
    }
    return &HanServer{populator: populator, db: db}
}

func (s *HanServer) imageSearchHandler(w http.ResponseWriter, r *http.Request) {
    // for running locally with Javascript
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // get the GET parameters
    params := r.URL.Query()
    lat, err := strconv.ParseFloat(params.Get("lat"), 64)
    if err != nil {
        http.Error(w, "Invalid latitude", 400)
        return
    }
    lng, err := strconv.ParseFloat(params.Get("lng"), 64)
    if err != nil {
        http.Error(w, "Invalid longitude", 400)
        return
    }
    // optional range values
    start, err := strconv.Atoi(params.Get("start"))
    if err != nil {
        start = -1
    }
    end, err := strconv.Atoi(params.Get("end"))
    if err != nil {
        end = -1
    }
    // if the region does not exist then we create it and populate it with
    // images
    if !hanapi.ContainsRegion(s.db, lat, lng) {
        hanapi.AddRegion(s.db, lat, lng)
        s.populator.PopulateImageDBWithLoc(s.db, lat, lng)
    }

    images := hanapi.GetImagesWithRange(s.db, lat, lng, start, end)
    response := new(response.ImageSearchResults)
    response.Images = images
    // return as a json response
    json.NewEncoder(w).Encode(response)
}

func reportImageHandler(w http.ResponseWriter, r *http.Request) {
    // for running locally with Javascript
    w.Header().Set("Access-Control-Allow-Origin", "*")
    mongo := db.NewMongoInterface()
    // get the GET parameters
    params := r.URL.Query()
    // found strangeness passing in strings as parameters with mongo
    id := fmt.Sprintf("%s", params.Get("id"))
    reason := fmt.Sprintf("%s", params.Get("reason"))
    hanapi.ReportImage(mongo, id, reason)
    mongo.Close()
}

func getRegionHandler(w http.ResponseWriter, r *http.Request) {
    // for running locally with Javascript
    w.Header().Set("Access-Control-Allow-Origin", "*")
    mongo := db.NewMongoInterface()
    // return regions as json
    regions := hanapi.GetRegions(mongo)
    json.NewEncoder(w).Encode(regions)
    mongo.Close()
}

// PlacesProxy is used to avoid revealing the
// Google Maps API key and also avoid storing each image
// ourselves
type PlacesProxy struct {
    proxy  *httputil.ReverseProxy
    host   string
    apiKey string
}

// NewProxy will create a new Places Proxy that will add your API key
// to all requests
func NewProxy(apiKey string) *PlacesProxy {
    host := "maps.googleapis.com"
    url := &url.URL{
        Scheme: "https",
        Host: host,
    }
    return &PlacesProxy{proxy: httputil.NewSingleHostReverseProxy(url), host:host, apiKey:apiKey}
}

func (p *PlacesProxy) handle(w http.ResponseWriter, r *http.Request) {
    // for running locally with Javascript
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // add API key as parameter
    q := r.URL.Query()
    q.Set("key", p.apiKey)
    r.URL.RawQuery = q.Encode()
    // reset the host
    r.Host = p.host
    p.proxy.ServeHTTP(w, r)
}

func main() {
    noCollectionPtr := flag.Bool("nocollection", false, "use this argument to stop hancollector being started automatically")
    flag.Parse()
    server := NewHanServer(*noCollectionPtr)
    http.HandleFunc("/api/image-search", server.imageSearchHandler)
    http.HandleFunc("/api/report-image", reportImageHandler)
    http.HandleFunc("/api/get-regions", getRegionHandler)
    // used to retrieve Google Places photos
    apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
    if len(apiKey) > 0 {
        proxy := NewProxy(apiKey)
        // we use a matching path to make it easier for the reverse proxy
        http.HandleFunc("/maps/api/place/photo", proxy.handle)
    }
    http.ListenAndServe(":80", nil)
}
