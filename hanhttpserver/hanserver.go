package main

import (
    "strconv"
    "encoding/json"
    "net/http"
    "github.com/oliveroneill/hanserver/hanapi"
    "github.com/oliveroneill/hanserver/hanapi/db"
    "github.com/oliveroneill/hanserver/hancollector/imagepopulation"
    "github.com/oliveroneill/hanserver/hanhttpserver/response"
)

func imageSearchHandler(w http.ResponseWriter, r *http.Request) {
    // for running locally with Javascript
    w.Header().Set("Access-Control-Allow-Origin", "*")
    mongo := db.NewMongoInterface()
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
    if !hanapi.ContainsRegion(mongo, lat, lng) {
        hanapi.AddRegion(mongo, lat, lng)
        imagepopulation.PopulateImageDBWithLoc(mongo, lat, lng)
    }

    images := hanapi.GetImagesWithRange(mongo, lat, lng, start, end)
    response := new(response.ImageSearchResults)
    response.Images = images
    // return as a json response
    json.NewEncoder(w).Encode(response)
}

func getRegionHandler(w http.ResponseWriter, r *http.Request) {
    // for running locally with Javascript
    w.Header().Set("Access-Control-Allow-Origin", "*")
    mongo := db.NewMongoInterface()
    // return regions as json
    regions := hanapi.GetRegions(mongo)
    json.NewEncoder(w).Encode(regions)
}

func main() {
    http.HandleFunc("/api/image-search", imageSearchHandler)
    http.HandleFunc("/api/get-regions", getRegionHandler)
    http.ListenAndServe(":80", nil)
}
