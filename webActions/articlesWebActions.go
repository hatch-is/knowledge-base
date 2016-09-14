package webActions

import (
	"encoding/json"
	"knowledge-base/model"
	"net/http"

	"github.com/gorilla/mux"
)

//ArticlesWebActions ...
type ArticlesWebActions struct {
	model model.ArticlesModel
}

//Read get all articles
func (art *ArticlesWebActions) Read(w http.ResponseWriter, r *http.Request) {
	skip := r.URL.Query().Get("skip")
	limit := r.URL.Query().Get("limit")
	filter := r.URL.Query().Get("filter")
	q := r.URL.Query().Get("q")

	data, err := art.model.Read(filter, skip, limit, q)

	if err != nil {
		ErrorWithJSON(w, err.Error(), 404)
	} else {
		result, _ := json.Marshal(data)
		ResponseWithJSON(w, result, 200)
	}
}

//ReadOne get article by ID
func (art *ArticlesWebActions) ReadOne(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	data, err := art.model.ReadOne(ID)

	if err != nil {
		ErrorWithJSON(w, err.Error(), 404)
	} else {
		result, _ := json.Marshal(data)
		ResponseWithJSON(w, result, 200)
	}
}

//Create add new article
func (art *ArticlesWebActions) Create(w http.ResponseWriter, r *http.Request) {
	body := r.Body

	data, err := art.model.Create(body)

	if err != nil {
		ErrorWithJSON(w, err.Error(), 404)
	} else {
		result, _ := json.Marshal(data)
		ResponseWithJSON(w, result, 200)
	}
}

//Update entry by ID
func (art *ArticlesWebActions) Update(w http.ResponseWriter, r *http.Request) {

	ID := mux.Vars(r)["id"]
	body := r.Body

	data, err := art.model.Update(ID, body)

	if err != nil {
		ErrorWithJSON(w, err.Error(), 404)
	} else {
		result, _ := json.Marshal(data)
		ResponseWithJSON(w, result, 200)
	}
}

//Delete entry by ID
func (art *ArticlesWebActions) Delete(w http.ResponseWriter, r *http.Request) {

	ID := mux.Vars(r)["id"]

	err := art.model.Delete(ID)

	if err != nil {
		ErrorWithJSON(w, err.Error(), 404)
	} else {
		result, _ := json.Marshal("")
		ResponseWithJSON(w, result, 200)
	}
}
