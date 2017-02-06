package feedsort

import (
    "time"
    "github.com/oliveroneill/hanserver/hanapi/imagedata"
)

/**
 * This is a sorting algorithm used to sort the images based on time and
 * distance
 */

// RecencyBias is a multiplier used to create a more equal weighting between
// recency and distance
const RecencyBias = 350

// BySum is a sorting heuristic based on distance and recency
type BySum []imagedata.ImageData

func (images BySum) Len() int {
    return len(images)
}
func (images BySum) Swap(i, j int) {
    images[i], images[j] = images[j], images[i]
}
func (images BySum) Less(i, j int) bool {
    t1 := time.Since(time.Unix(0, images[i].CreatedTime * 1000)).Seconds()
    t2 := time.Since(time.Unix(0, images[j].CreatedTime * 1000)).Seconds()

    return (images[i].Distance + t1*RecencyBias) <
        (images[j].Distance + t2*RecencyBias)
}
