package webActions

import (
	"knowledge-base/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"knowledge-base/conf"
)

//ArticlesWebActions ...
type ArticlesWebActions struct {
	model model.ArticlesModel
}

//Read get all articles
func (art *ArticlesWebActions) Read(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	lg := r.Header.Get("x-location-group")
	conf.Config.LocationGroup = lg

	data, total, left, err := art.model.Read(filter)
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(left))
	if err != nil {
		ErrorWithJSON(w, r, err.Error(), 404)
	} else {
		ResponseWithJSON(w, r, data, 200)
	}
}

//ReadOne get article by ID
func (art *ArticlesWebActions) ReadOne(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]
	lg := r.Header.Get("x-location-group")
	conf.Config.LocationGroup = lg
	data, err := art.model.ReadOne(ID)

	if err != nil {
		ErrorWithJSON(w, r, err.Error(), 404)
	} else {
		ResponseWithJSON(w, r, data, 200)
	}
}

//Create add new article
func (art *ArticlesWebActions) Create(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	lg := r.Header.Get("x-location-group")
	conf.Config.LocationGroup = lg
	data, err := art.model.Create(body)

	if err != nil {
		ErrorWithJSON(w, r, err.Error(), 404)
	} else {
		ResponseWithJSON(w, r, data, 200)
	}
}

//Update entry by ID
func (art *ArticlesWebActions) Update(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]
	body := r.Body
	lg := r.Header.Get("x-location-group")
	conf.Config.LocationGroup = lg

	data, err := art.model.Update(ID, body)

	if err != nil {
		ErrorWithJSON(w, r, err.Error(), 404)
	} else {
		ResponseWithJSON(w, r, data, 200)
	}
}

//Delete entry by ID
func (art *ArticlesWebActions) Delete(w http.ResponseWriter, r *http.Request) {
	lg := r.Header.Get("x-location-group")
	conf.Config.LocationGroup = lg
	ID := mux.Vars(r)["id"]

	err := art.model.Delete(ID)

	if err != nil {
		ErrorWithJSON(w, r, err.Error(), 404)
	} else {
		ResponseWithJSON(w, r, []string{""}, 200)
	}
}
