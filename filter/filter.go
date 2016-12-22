package filter

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/pquerna/ffjson/ffjson"
	"gopkg.in/mgo.v2/bson"
)

const (
	id   = "id"
	_id  = "_id"
	desc = "desc"
	//ID ...
	ID = "ID"
)

//Filter describe fiter query
type Filter struct {
	skip   int
	limit  int
	query  bson.M
	sort   []string
	search string
}

//BeforeFilter ...
type BeforeFilter struct {
	Skip   interface{}
	Limit  interface{}
	Where  map[string]interface{}
	Sort   map[string]string
	Search string
}

func (bf *BeforeFilter) getSort() []string {
	var sort []string
	sort = make([]string, 0)
	for key, value := range bf.Sort {
		if desc == strings.ToLower(value) {
			key = "-" + key
		}
		sort = append(sort, key)
	}
	return sort
}

func (bf *BeforeFilter) getLimit() int {
	return convertParam(bf.Limit)
}

func (bf *BeforeFilter) getSkip() int {
	return convertParam(bf.Skip)
}

//convertParam to convert skip/limit to int
//by default return 0
func convertParam(param interface{}) int {
	var result int
	switch param.(type) {
	default:
		return int(0)
	case float64:
		result = int(param.(float64))
	case string:
		t, err := strconv.Atoi(param.(string))
		if err != nil {
			result = int(0)
		}
		result = t
	}
	if result < 0 {
		return int(0)
	}
	return result
}

//GetSkip ...
func (f *Filter) GetSkip() int {
	return f.skip
}

//GetLimit ...
func (f *Filter) GetLimit() int {
	return f.limit
}

//GetSearch ...
func (f *Filter) GetSearch() string {
	return f.search
}

//GetSort ...
func (f *Filter) GetSort() []string {
	return f.sort
}

//GetQuery ...
func (f *Filter) GetQuery() bson.M {
	return f.query
}

//AddQuery ...
func (f *Filter) AddQuery(q bson.M) {
	if f.query == nil {
		f.query = q
	} else {
		f.query["$and"] = append(f.query["$and"].([]bson.M), q)
	}
}

//GetFilterData ...
func GetFilterData(filter string, t interface{}) Filter {
	var f Filter
	if filter != "" {
		var tmp BeforeFilter
		ffjson.Unmarshal([]byte(filter), &tmp)
		f.limit = tmp.getLimit()
		f.skip = tmp.getSkip()
		f.sort = tmp.getSort()
		f.query = tmp.parseWhere(t)
	}
	return f
}

func (bf *BeforeFilter) parseWhere(t interface{}) bson.M {
	var q bson.M
	for key, values := range bf.Where {
		if key == "and" || key == "or" {
			key = "$" + key
			q = parseAnd(key, values.([]interface{}), q, t)
		} else {
			q = parse(key, values, q, t)
		}
	}
	return q
}

func parse(key string, values interface{}, q bson.M, t interface{}) bson.M {
	if key == id {
		key = _id
	}
	val := reflect.ValueOf(t)
	var fieldType string
	if val.FieldByName(convertName(key)).IsValid() == true {
		fieldType = val.FieldByName(convertName(key)).Type().String()
		if fieldType == "*time.Time" || fieldType == "time.Time" {
			if q == nil {
				q = bson.M{"$and": []bson.M{}}
			} else {
				if q["$and"] == nil {
					q["$and"] = []bson.M{}
				}
			}
			q = parseFilterDate("$and", key, values, q)
		}
		if fieldType == "bson.ObjectId" {
			q = parseID(key, values, q)
		}
	}
	return q
}

func addKeyToQuery(key string, q bson.M) bson.M {
	if q == nil {
		q = bson.M{key: []bson.M{}}
	} else {
		q[key] = []bson.M{}
	}
	return q
}

func parseID(key string, values interface{}, q bson.M) bson.M {
	q = addKeyToQuery(key, q)
	for field, data := range values.(map[string]interface{}) {
		ids := convertStringIDtoObjectID(data)
		if len(ids) > 0 {
			if field == "in" || field == "nin" || field == "all" {
				q[key] = bson.M{"$" + field: ids}
			}
		}
	}
	return q
}

func parseAndIDData(key string, field string, data interface{}, q bson.M) bson.M {
	if field == id {
		field = _id
	}
	slice := data.(map[string]interface{})
	for k, value := range slice {
		IDS := convertStringIDtoObjectID(value)
		if len(IDS) > 0 {
			if k == "in" || k == "nin" || k == "all" {
				q[key] = append(q[key].([]bson.M), bson.M{field: bson.M{"$" + k: IDS}})
			}
		}
	}
	return q
}

func convertStringIDtoObjectID(data interface{}) []bson.ObjectId {
	var ids []bson.ObjectId
	ids = make([]bson.ObjectId, 0)
	switch data.(type) {
	case []interface{}:
		ids = prepareIDSArray(data.([]interface{}))
	case string:
		ids = prepareInIds(data.(string))
	}
	return ids
}

func prepareIDSArray(values []interface{}) []bson.ObjectId {
	var ids []bson.ObjectId
	ids = make([]bson.ObjectId, 0)
	for _, id := range values {
		ids = append(ids, bson.ObjectIdHex(id.(string)))
	}
	return ids
}

func parseAnd(key string, values []interface{}, q bson.M, t interface{}) bson.M {
	val := reflect.ValueOf(t)
	if q == nil {
		q = bson.M{key: []bson.M{}}
	} else {
		q[key] = []bson.M{}
	}
	for _, value := range values {
		for field, data := range value.(map[string]interface{}) {
			fieldName := convertName(field)
			var fieldType string
			if val.FieldByName(fieldName).IsValid() == true {
				fieldType = val.FieldByName(fieldName).Type().String()
				if fieldType == "*time.Time" || fieldType == "time.Time" {
					q = parseFilterDate(key, field, data, q)
				}
				if fieldType == "bson.ObjectId" {
					q = parseAndIDData(key, field, data, q)
				}
			}
		}
	}
	return q
}

func convertName(name string) string {
	if name == id || name == _id {
		return ID
	}
	newName := []rune(name)
	newName[0] = unicode.ToUpper(newName[0])
	return string(newName)
}

func parseFilterDate(key string, field string, data interface{}, q bson.M) bson.M {
	for k, v := range data.(map[string]interface{}) {
		if k == "lte" || k == "lt" || k == "gte" || k == "gt" {
			t, err := time.Parse(time.RFC3339, v.(string))
			if err == nil {
				q[key] = append(q[key].([]bson.M), bson.M{field: bson.M{"$" + k: t}})
			}
		}
	}
	return q
}

func prepareInIds(sID string) []bson.ObjectId {
	var ids []bson.ObjectId
	ids = make([]bson.ObjectId, 0)
	IDs := strings.Split(sID, ",")
	for _, id := range IDs {
		ids = append(ids, bson.ObjectIdHex(id))
	}
	return ids
}
