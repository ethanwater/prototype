package utils

import (
	"math/rand"
	"strings"
	"time"
)

const (
	numset       = "0123456789"
	vivianIDSize = 9
)

func GenerateVivianID() string {
	source := rand.New(rand.NewSource(time.Now().Unix()))
	var vivianID strings.Builder

	for i := 0; i < vivianIDSize; i++ {
		sample := source.Intn(len(numset))
		vivianID.WriteString(string(numset[sample]))
	}

	return vivianID.String()
}
