package model

import (
	"encoding/json"
	"io"
	"knowledge-base/store"
	"strconv"
	"time"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

//ArticlesModel ...
type ArticlesModel struct{}

//Read get articles
func (artModel *ArticlesModel) Read(qFilter string, qSkip string, qLimit string) (result []store.Article, err error) {
	query := artModel.convertFilter(qFilter)
	query["deleted"] = false

	fields := bson.M{}

	skip, err := strconv.Atoi(qSkip)
	if err != nil {
		skip = 0
	}
	limit, err := strconv.Atoi(qLimit)
	if err != nil {
		limit = 0
	}

	artDb := store.ArticlesCollectionConnect()
	result, err = artDb.Read(query, fields, skip, limit)

	if err != nil {
		return
	}
	return result, nil
}

//ReadOne get article by ID
func (artModel *ArticlesModel) ReadOne(qID string) (result store.Article, err error) {
	ID := bson.ObjectIdHex(qID)

	query := bson.M{
		"_id":     ID,
		"deleted": false,
	}

	fields := bson.M{}
	artDb := store.ArticlesCollectionConnect()
	result, err = artDb.ReadOne(query, fields)

	if err != nil {
		return
	}
	return result, nil
}

func (artModel *ArticlesModel) convertFilter(filter string) bson.M {
	var m map[string]map[string]map[string]interface{}

	json.Unmarshal([]byte(filter), &m)
	q := bson.M{"$and": []bson.M{}}
	for key, value := range m {
		if key == "where" {
			for field, data := range value {
				for k, v := range data {
					if k == "lte" || k == "lt" || k == "gte" || k == "gt" {
						t, err := time.Parse(time.RFC3339, v.(string))
						if err == nil {
							q["$and"] = append(q["$and"].([]bson.M), bson.M{field: bson.M{"$" + k: t}})
						}
					}
				}
			}
		}
	}
	if cap(q["$and"].([]bson.M)) < 1 {
		return bson.M{}
	}
	return q
}

//Create add new article
func (artModel *ArticlesModel) Create(data io.ReadCloser) (result store.Article, err error) {

	var validate *validator.Validate
	validate = validator.New()

	decoder := json.NewDecoder(data)
	var article store.Article
	err = decoder.Decode(&article)
	if err != nil {
		return
	}

	err = validate.Struct(article)

	if err != nil {
		return
	}

	artDb := store.ArticlesCollectionConnect()
	result, err = artDb.Create(article)
	tagModel := TagsModel{}
	go tagModel.CampareAndCreate(result.Tags)
	if err != nil {
		return
	}
	return result, nil
}

//Update article by ID
func (artModel *ArticlesModel) Update(qID string, data io.ReadCloser) (result store.Article, err error) {

	ID := bson.ObjectIdHex(qID)
	var validate *validator.Validate
	validate = validator.New()

	decoder := json.NewDecoder(data)
	var article store.Article
	err = decoder.Decode(&article)
	if err != nil {
		return
	}

	err = validate.Struct(article)

	if err != nil {
		return
	}

	artDb := store.ArticlesCollectionConnect()
	result, err = artDb.Update(ID, article)

	if err != nil {
		return
	}
	return result, nil
}
