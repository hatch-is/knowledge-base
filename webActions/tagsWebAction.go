package webActions

import (
	"encoding/json"
	"knowledge-base/model"
	"net/http"
)

//TagsWebActions ...
type TagsWebActions struct {
	model model.TagsModel
}

//All return all ags in collection Tags
func (tag *TagsWebActions) All(w http.ResponseWriter, r *http.Request) {
	data, err := tag.model.All()

	if err != nil {
		ErrorWithJSON(w, err.Error(), 404)
	} else {
		result, _ := json.Marshal(data)
		ResponseWithJSON(w, result, 200)
	}
}
