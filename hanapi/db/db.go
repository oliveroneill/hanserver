package db

import (
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
)

// DatabaseInterface - a generic interface for database queries
type DatabaseInterface interface {
    GetRegions() []imagedata.ImageLocation
    AddRegion(lat float64, lng float64)
    AddImage(image imagedata.ImageData)
    GetImages(lat float64, lng float64) []imagedata.ImageData
    GetAllImages() []imagedata.ImageData
    SoftDelete(id string, reason string)
    Close()
}