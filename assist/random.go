package assist

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// +-------------------------------------------------------------------------------------+
// 									EXPOSED FUNCTIONS
// +-------------------------------------------------------------------------------------+

// GetRandom returns a copy of a random element from the array
func GetRandom(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

// RandomSliceIndex returns one of the indices of
// the given slice at random.
func RandomSliceIndex(slice []string) int {
	return rand.Intn(len(slice))
}

// GetRandomIndex returns a random index from the passed in strin
func GetRandomIndex(str string) int {
	return rand.Intn(len(str))
}

// GetRandomSliceIndex returns a random index from the passed in slice
func GetRandomSliceIndex(slice []string) int {
	return rand.Intn(len(slice))
}

// RandomScale returns either -1 or 1.
func RandomScale() int {
	if rand.Intn(2) == 1 {
		return 1
	}

	return -1
}

// RandomDelay returns the base, either plus or minus a
// random number up to the value of the variance.
func RandomDelay(base int, variance int) int {
	return base + rand.Intn(variance)*RandomScale()
}

// RandomBetween generates a number that is at least
// min and no more than max.
func RandomBetween(min int, max int) int {
	return rand.Intn(max-(min)-1) + min + 1
}

// SmallRandomWait halts the current goroutine for a small
// period of time. Used to simulate the way a user usually
// doesn't perform consecutive actions instantly.
func SmallRandomWait() {
	Wait(RandomBetween(30, 100))
}
