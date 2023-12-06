package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type jsonResponse map[string]any

func postError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func postBodyResponse(w http.ResponseWriter, code int, content jsonResponse) {
	if content == nil {
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		return
	}

	js, err := json.Marshal(content)
	if err != nil {
		postError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}

func postOptionsResponse(w http.ResponseWriter, methods []string, content jsonResponse) {
	w.Header().Set("Allow", strings.Join(methods, ","))
	postBodyResponse(w, http.StatusOK, content)
}
