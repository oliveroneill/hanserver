package main

import (
	"fmt"
	"flag"
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/oliveroneill/hanserver/hanapi"
	"github.com/oliveroneill/hanserver/hanapi/dao"
	"github.com/oliveroneill/hanserver/hanapi/reporting"
	"github.com/oliveroneill/hanserver/hancollector/imagepopulation"
	"github.com/oliveroneill/hanserver/hanhttpserver/response"
)

// HanServer is a http server that also populates the database periodically
// This allows easy tracking of API usage
type HanServer struct {
	populator *imagepopulation.ImagePopulator
	db		  dao.DatabaseInterface
	logger    reporting.Logger
}

// NewHanServer will create a new http server and start population
// @param noCollection - set this to true if you don't want hancollector to
//                       start
// @param apiToken     - optional slack api token used for logging errors to
//                       Slack
func NewHanServer(noCollection bool, apiToken string) *HanServer {
	// this database session is kept onto over the lifetime of the server
	db := dao.NewMongoInterface()
	logger := reporting.NewSlackLogger(apiToken)
	populator := imagepopulation.NewImagePopulator(logger)
	if !noCollection {
		fmt.Println("Starting image collection")
		// populate image db in the background
		go populator.PopulateImageDB(db)
	}
	return &HanServer{populator: populator, db: db, logger: logger}
}

func (s *HanServer) imageSearchHandler(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
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
	if !hanapi.ContainsRegion(session, lat, lng) {
		hanapi.AddRegion(session, lat, lng)
		s.populator.PopulateImageDBWithLoc(session, lat, lng)
	}

	images := hanapi.GetImagesWithRange(session, lat, lng, start, end)
	response := new(response.ImageSearchResults)
	response.Images = images
	// return as a json response
	json.NewEncoder(w).Encode(response)
}

func (s *HanServer) reportImageHandler(w http.ResponseWriter, r *http.Request) {
	// for running locally with Javascript
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mongo := dao.NewMongoInterface()
	defer mongo.Close()
	// get the GET parameters
	params := r.URL.Query()
	// found strangeness passing in strings as parameters with mongo
	id := fmt.Sprintf("%s", params.Get("id"))
	reason := fmt.Sprintf("%s", params.Get("reason"))
	hanapi.ReportImage(mongo, id, reason, s.logger)
}

func getRegionHandler(w http.ResponseWriter, r *http.Request) {
	// for running locally with Javascript
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mongo := dao.NewMongoInterface()
	defer mongo.Close()
	// return regions as json
	regions := hanapi.GetRegions(mongo)
	json.NewEncoder(w).Encode(regions)
}

func main() {
	noCollectionPtr := flag.Bool("nocollection", false, "Use this argument to stop hancollector being started automatically")
	slackAPITokenPtr := flag.String("slacktoken", "", "Specify the API token for logging through Slack")
	flag.Parse()
	server := NewHanServer(*noCollectionPtr, *slackAPITokenPtr)
	http.HandleFunc("/api/image-search", server.imageSearchHandler)
	http.HandleFunc("/api/report-image", server.reportImageHandler)
	http.HandleFunc("/api/get-regions", getRegionHandler)
	http.ListenAndServe(":80", nil)
}
