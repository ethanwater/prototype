package utils

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Authentication Key Generation Config
const (
	charset     string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	authKeySize int    = 5
)

// Receiver Config
type ReceiverType int

const (
	Email ReceiverType = iota + 1
	Mobile
)

type Receiver struct {
	Receiver ReceiverType
}

func (r *Receiver) EmailSendAuthKey2FA(authKey string) error {
	return nil
}

func (r *Receiver) MobileSendAuthKey2FA(authKey string) error {
	return nil
}

// GenerateAuthKey2FA generates a 2FA authentication key.
// The generated key will be hashed and stored via localStorage in JavaScript
// and should be removed from localStorage cache once verified.
func GenerateAuthKey2FA(receiver Receiver) (string, error) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	source := rand.New(rand.NewSource(time.Now().Unix()))
	var authKey strings.Builder

	for i := 0; i < authKeySize; i++ {
		sample := source.Intn(len(charset))
		authKey.WriteString(string(charset[sample]))
	}

	switch receiver.Receiver {
	case Email:
		receiver.EmailSendAuthKey2FA(authKey.String())
	case Mobile:
		receiver.MobileSendAuthKey2FA(authKey.String())
	}

	authKeyHash, error := HashPassword(authKey.String())
	return string(authKeyHash), error
}

func VerifyAuthKey2FA(authkey_hash, input string) bool {
	if SanitizeCheck(input) {
		status := bcrypt.CompareHashAndPassword([]byte(authkey_hash), []byte(input))
		return status == nil
	}

	return false
}
