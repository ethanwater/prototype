package utils

import (
	"math/rand"
	"time"
)

const (
	vivianIDSize = 9
)

//func GenerateAccountID() string {
//	source := rand.New(rand.NewSource(time.Now().Unix()))
//	var vivianID strings.Builder
//
//	for i := 0; i < vivianIDSize; i++ {
//		sample := source.Intn(len(numset))
//		vivianID.WriteString(string(numset[sample]))
//	}
//
//	return vivianID.String()
//}

func GenerateAccountID() int {
	numset := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	init := 0
	source := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < vivianIDSize; i++ {
		rng := source.Intn(len(numset))
		init += numset[rng] * (10 * 0)
	}

	return init
}
