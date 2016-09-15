package model

import (
	"encoding/json"
	"io"
	"knowledge-base/store"
	"strconv"

	"github.com/Jeffail/gabs"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

//ArticlesModel ...
type ArticlesModel struct{}

//Read get articles
func (artModel *ArticlesModel) Read(qFilter string, qSkip string, qLimit string) (result []store.Article, err error) {
	artDb := store.ArticlesCollectionConnect()

	query, isSort, err := artModel.convertFilter(qFilter)
	if err != nil {
		return
	}
	query["deleted"] = false
	if isSort == false {

		fields := bson.M{}

		skip, err := strconv.Atoi(qSkip)
		if err != nil {
			skip = 0
		}
		limit, err := strconv.Atoi(qLimit)
		if err != nil {
			limit = 0
		}
		result, err = artDb.Read(query, fields, skip, limit)
	} else {
		result, err = artDb.Search(query)
	}

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

func (artModel *ArticlesModel) convertFilter(filter string) (q bson.M, isSearch bool, err error) {
	if filter != "" {
		jParsed, err := gabs.ParseJSON([]byte(filter))
		isSearch = false
		if err == nil {
			q = bson.M{"$and": []bson.M{}}
			IDs, err := jParsed.S("where", "id", "in").Children()
			if err == nil {
				var ids []bson.ObjectId
				ids = make([]bson.ObjectId, 0)
				for _, id := range IDs {
					ids = append(ids, bson.ObjectIdHex(id.Data().(string)))
				}
				if len(ids) > 0 {
					q["$and"] = append(q["$and"].([]bson.M), bson.M{"_id": bson.M{"$in": ids}})
				}
			}
			search, ok := jParsed.S("where", "search").Data().(string)
			if ok == true {
				isSearch = true
				q["$and"] = append(q["$and"].([]bson.M), bson.M{"$text": bson.M{"$search": search}})
			}
			if cap(q["$and"].([]bson.M)) < 1 {
				return bson.M{}, false, nil
			}
			return q, isSearch, nil
		}
	}
	return bson.M{}, false, nil
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

	if err != nil {
		return
	}
	tagModel := TagsModel{}
	go tagModel.CampareAndCreate(result.Tags)
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
	err = artDb.Update(ID, article)

	if err != nil {
		return
	}
	tagModel := TagsModel{}
	go tagModel.CampareAndCreate(result.Tags)
	query := bson.M{
		"_id": ID,
	}
	result, err = artDb.ReadOne(query, bson.M{})
	if err != nil {
		return
	}
	return result, nil
}

//Delete article by ID
func (artModel *ArticlesModel) Delete(qID string) (err error) {
	ID := bson.ObjectIdHex(qID)
	artDb := store.ArticlesCollectionConnect()
	err = artDb.Delete(ID)

	return
}
