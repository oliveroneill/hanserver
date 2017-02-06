package response

import (
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
)

// ImageSearchResults is a list of images used for responses from `hanhttpserver`
type ImageSearchResults struct {
    Images []imagedata.ImageData `json:"images" bson:"images"`
}