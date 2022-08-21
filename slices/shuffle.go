package slices

import (
	"math/rand"
	"time"
)

// Shuffle slice
func Shuffle[T any](slice []T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(slice); i++ {
		rand.Shuffle(len(slice), func(i, j int) {
			slice[i], slice[j] = slice[j], slice[i]
		})
	}
}
