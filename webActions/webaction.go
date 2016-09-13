package webActions

import (
	"fmt"
	"net/http"
)

//ResponseWithJSON return Response
func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

//ErrorWithJSON return Error Response
func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"result":"","message":%q}`, message)
}

//URLRoot return hello message
func URLRoot(w http.ResponseWriter, r *http.Request) {
	result := "Hello and welcome to the Hatch Knowledge Base"
	ErrorWithJSON(w, result, 200)
}
