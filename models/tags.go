package models

import "knowledge-base/db"

//Tags ...
type Tags struct {
	ID   int64  `db:"id, primarykey, autoincrement"json:"id"`
	Name string `db:"name",json:"name"`
}

//TagsModel ...
type TagsModel struct{}

//Create add new tags to DB
func (tm TagsModel) Create(tags []string) error {
	var exising []Tags
	query := "SELECT name FROM tags"
	insertQuery := "INSERT INTO tags (name) VALUES ($1)"
	_, err := db.GetDB().Select(&exising, query)
	for _, tag := range tags {
		flag := false
		for _, eTag := range exising {
			if tag == eTag.Name {
				flag = true
				break
			}
		}
		if flag == false {
			_, err = db.GetDB().Exec(insertQuery, tag)
		}
	}
	return err
}

//All return tags as []string
func (tm TagsModel) All() []string {
	var exising []Tags
	var tags []string
	tags = make([]string, 0)
	query := "SELECT name FROM tags"
	_, err := db.GetDB().Select(&exising, query)
	if err != nil {
		return tags
	}

	for _, t := range exising {
		tags = append(tags, t.Name)
	}
	return tags
}
