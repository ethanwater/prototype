package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

const cost int = 13

func HashPassword(_ context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}

// when called in register, should be processed in parallel. go func()
//func HashPassword(password string) (string, error) {
//	resultCh := make(chan struct {
//		hash string
//		err  error
//	})
//
//	go func() {
//		hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
//		resultCh <- struct {
//			hash string
//			err  error
//		}{string(hash), err}
//	}()
//
//	result := <-resultCh
//	return result.hash, result.err
//}

func VerfiyHashPassword(hash, password string) bool {
	status := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return status == nil
}
