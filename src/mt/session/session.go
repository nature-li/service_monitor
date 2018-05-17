package session

import "net/http"

type Session interface {
	SessionId() string
	Set(string, string) error
	Get(string) string
	Del(string) error
}

type Manager interface {
	SessionStart(w http.ResponseWriter, r *http.Request) Session
	SessionDestroy(w http.ResponseWriter, r *http.Request) error
}