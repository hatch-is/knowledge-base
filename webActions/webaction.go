package webActions

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/raven-go"
	sm "github.com/phemmer/sawmill"
)

//ResponseWithJSON return Response
func ResponseWithJSON(w http.ResponseWriter, r *http.Request, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	sm.Info(getUlrForLog(r), data)
	result, _ := json.Marshal(data)
	w.Write(result)
}

//ErrorWithJSON return Error Response
func ErrorWithJSON(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	sm.Error(getUlrForLog(r), []string{message})
	raven.CaptureMessage(message, nil, nil)
}

//URLRoot return hello message
func URLRoot(w http.ResponseWriter, r *http.Request) {
	result := "Hello and welcome to the Hatch Knowledge Base"
	ResponseWithJSON(w, r, []string{result}, 200)
}

func getUlrForLog(r *http.Request) string {
	return "\"" + r.Method + " " + r.URL.RequestURI() + "\""
}
