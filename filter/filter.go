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
	id  = "id"
	_id = "_id"
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
	Skip   string
	Limit  string
	Where  map[string]interface{}
	Sort   map[string]string
	Search string
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
		f.limit, _ = strconv.Atoi(tmp.Limit)
		f.skip, _ = strconv.Atoi(tmp.Skip)
		f.search = tmp.Search

		for key, value := range tmp.Sort {
			if "desc" == strings.ToLower(value) {
				key = "-" + key
			}
			f.sort = append(f.sort, key)
		}
		f.query = tmp.parseWhere(t)
	}
	return f
}

func (f *BeforeFilter) parseWhere(t interface{}) bson.M {
	var q bson.M
	for key, values := range f.Where {
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
	if q == nil {
		q = bson.M{key: []bson.M{}}
	} else {
		q[key] = []bson.M{}
	}
	fieldName := convertName(key)
	var fieldType string
	for field, data := range values.(map[string]interface{}) {
		if val.FieldByName(fieldName).IsValid() == true {
			fieldType = val.FieldByName(fieldName).Type().String()
			if fieldType == "bson.ObjectId" {
				q = parseIDData(key, field, data, q)
			}
		}
	}
	return q
}

func parseIDData(key string, field string, data interface{}, q bson.M) bson.M {
	IDS := prepareInIds(data.(string))
	if len(IDS) > 0 {
		if field == "in" || field == "nin" || field == "all" {
			q[key] = bson.M{"$" + field: IDS}
		}
	}
	return q
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
		return "ID"
	}
	newName := []rune(name)
	newName[0] = unicode.ToUpper(newName[0])
	return string(newName)
}

func parseAndIDData(key string, field string, data interface{}, q bson.M) bson.M {
	if field == id {
		field = _id
	}
	slice := data.(map[string]interface{})
	for k, value := range slice {
		IDS := prepareInIds(value.(string))
		if len(IDS) > 0 {
			if k == "in" || k == "nin" || k == "all" {
				q[key] = append(q[key].([]bson.M), bson.M{field: bson.M{"$" + k: IDS}})
			}
		}
	}
	return q
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
