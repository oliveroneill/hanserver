package main

import (
	"github.com/oliveroneill/hanserver/hanapi"
)

// ImageSearchResults is a list of images used for responses from `hanhttpserver`
type ImageSearchResults struct {
	Images []hanapi.ImageData `json:"images" bson:"images"`
}
