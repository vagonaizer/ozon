package logger

import (
	"log"
	"net/http"
	"time"
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) LogIncomingRequest(r *http.Request) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[INFO] %s | LOG INCOMING REQUEST\n", timestamp)
	log.Printf("Method: %s, Endpoint: %s\n", r.Method, r.URL.Path)
	log.Printf("Headers:\n")
	for key, values := range r.Header {
		log.Printf("  %s: %v\n", key, values)
	}
	log.Printf("\n")
}

func (l *Logger) LogOutgoingResponse(statusCode int, description string, executionTime time.Duration) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[INFO] %s | LOG OUTGOING RESPONSE\n", timestamp)
	log.Printf("Status Code: %d\n", statusCode)
	log.Printf("Description Status Code: %s\n", description)
	log.Printf("Execution Time: %v\n", executionTime)
	log.Printf("\n")
}
