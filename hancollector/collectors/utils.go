package collectors

import (
	"github.com/kellydunn/golang-geo"
	"github.com/oliveroneill/hanserver/hanapi"
)

// GetSurroundingPoints returns points surrounding the specified point by
// circling at query range
func GetSurroundingPoints(lat float64, lng float64, queryRange float64) []hanapi.Location {
	points := []hanapi.Location{}
	for degrees := float64(0); degrees < 360; degrees += 90 {
		// search 5 kilometers in each direction
		p := geo.NewPoint(lat, lng)
		// find another point that's at the edge of the previous query
		newPoint := p.PointAtDistanceAndBearing(queryRange/1000, degrees)
		points = append(points, *hanapi.NewLocation(newPoint.Lat(), newPoint.Lng()))
	}
	return points
}
