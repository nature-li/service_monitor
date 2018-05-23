package handler

import "net/http"

func MainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("monitor is running"))
}
