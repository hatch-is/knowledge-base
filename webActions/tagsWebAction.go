package webActions

import (
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
		ErrorWithJSON(w, r, err.Error(), 404)
	} else {
		ResponseWithJSON(w, r, data, 200)
	}
}
