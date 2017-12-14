package hanapi

// DatabaseInterface - a generic interface for database queries
type DatabaseInterface interface {
	GetRegions() []Location
	AddRegion(lat float64, lng float64)
	AddImage(image ImageData)
	AddBulkImagesToRegion(images []ImageData, region *Location)
	GetImages(lat float64, lng float64, start int, end int) []ImageData
	GetAllImages() []ImageData
	SoftDelete(id string, reason string)
	DeleteOldImages(amount int)
	Size() int
	Copy() DatabaseInterface
	Close()
}
