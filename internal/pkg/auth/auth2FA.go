package auth

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ServiceWeaver/weaver"
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

// GenerateAuthKey2FA generates a 2FA authentication key.
// The generated key will be hashed and stored via localStorage in JavaScript
// and should be removed from localStorage cache once verified.

type Auth2FA interface {
	GenerateAuthKey2FA(context.Context) (string, error)
	VerifyAuthKey2FA(context.Context, string, string) (bool, error)
}

type auth2FA struct {
	weaver.Implements[Auth2FA]
}

func (t *auth2FA) GenerateAuthKey2FA(ctx context.Context) (string, error) {
	log := t.Logger(ctx)

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	source := rand.New(rand.NewSource(time.Now().Unix()))
	var authKey strings.Builder

	for i := 0; i < authKeySize; i++ {
		sample := source.Intn(len(charset))
		authKey.WriteString(string(charset[sample]))
	}

	//switch receiver.Receiver {
	//case Email:
	//	receiver.EmailSendAuthKey2FA(authKey.String())
	//case Mobile:
	//	receiver.MobileSendAuthKey2FA(authKey.String())
	//}

	fmt.Println(authKey.String())
	authKeyHash, error := HashPassword(authKey.String())

	log.Debug("vivian: STATUS!", "authentication key generated", http.StatusOK)
	return authKeyHash, error
}

func (t *auth2FA) VerifyAuthKey2FA(ctx context.Context, authkey_hash, input string) (bool, error) {
	log := t.Logger(ctx)
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	log.Debug("vivian: STATUS!", "action", "verifying input...")
	if SanitizeCheck(input) {
		status := bcrypt.CompareHashAndPassword([]byte(authkey_hash), []byte(input))
		return status == nil, status
	}

	return false, nil
}
