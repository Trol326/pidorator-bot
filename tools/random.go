package tools

import (
	"math/rand"
	"time"
)

// returns random positive int [0,n)
func GetRandomInt32(n int) int32 {
	random := rand.New(rand.NewSource(time.Now().Unix()))
	return random.Int31n(int32(n))
}
