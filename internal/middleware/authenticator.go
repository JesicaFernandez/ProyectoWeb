package middleware

import (
	"fmt"
	"net/http"
)

func NewAuthenticator(token string) *Authenticator {
	return &Authenticator{
		token: token,
	}
}

type Authenticator struct {
	token string
}

func (a *Authenticator) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate the token
		if err := a.validateAuthToken(r); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// call the next handler
		next.ServeHTTP(w, r)
	})
}

// ValidateAuthToken validates that the request contains a valid auth token
func (a *Authenticator) validateAuthToken(r *http.Request) (err error) {
	// - get the auth token from the header
	authToken := r.Header.Get("Auth")
	expectedAuthToken := a.token
	if expectedAuthToken == "" {
		return fmt.Errorf("expected auth token is empty")
	}
	// - validate that the auth token is not empty
	if authToken == "" {
		return fmt.Errorf("auth token is empty")
	}
	// - validate that the auth token is valid
	if authToken != expectedAuthToken {
		return fmt.Errorf("auth token is invalid")
	}
	return
}
