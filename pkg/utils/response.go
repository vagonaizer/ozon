package utils

import "net/http"

func SuccessResponse(w http.ResponseWriter) {
	w.Write([]byte("{}")) // TODO исправить
}
