package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func NewResponseLogger() *ResponseLogger {
	return &ResponseLogger{}
}

type ResponseLogger struct {
	verb   string
	date   time.Time
	path   string
	bytes  int64
}

func (rl *ResponseLogger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before

		// call the next handler
		next.ServeHTTP(w, r)

		// after
		// log the response
		rl.LogResponse(w, r)
	})
}

// LogResponse logs the response
func (rl *ResponseLogger) LogResponse(w http.ResponseWriter, r *http.Request) {
	// - create the response log
	res := &ResponseLogger{
		verb:   r.Method,
		date:   time.Now(),
		path:   r.URL.Path,
		bytes:  r.ContentLength,
	}

	// - log the response
	fmt.Printf("method: %s; date: %s; path: %s; bytes: %d\n", res.verb, res.date, res.path, res.bytes)
}
