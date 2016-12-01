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
func (artModel *ArticlesModel) Read(qFilter string) (result []store.Article, total int, left int, err error) {
	artDb := store.ArticlesCollectionConnect()

	query, isSort, err := artModel.convertFilter(qFilter)
	if err != nil {
		return
	}
	skip, limit := getSkipLimit(qFilter)
	query["deleted"] = false
	if isSort == false {
		fields := bson.M{}
		result, total, left, err = artDb.Read(query, fields, skip, limit)
	} else {
		result, total, left, err = artDb.Search(query, skip, limit)
	}

	if err != nil {
		return
	}
	return result, total, left, nil
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

func getSkipLimit(filter string) (skip int, limit int) {
	if filter != "" {
		jParser, err := gabs.ParseJSON([]byte(filter))
		if err != nil {
			return
		}

		if s := jParser.Exists("skip"); s == true {
			sTmp, ok := jParser.Path("skip").Data().(string)
			if ok == true {
				skip, _ = strconv.Atoi(sTmp)
			}
		}

		if l := jParser.Exists("limit"); l == true {
			lTmp, ok := jParser.Path("limit").Data().(string)
			if ok == true {
				limit, _ = strconv.Atoi(lTmp)
			}
		}
	}
	return skip, limit
}
