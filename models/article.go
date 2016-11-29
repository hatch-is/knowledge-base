package models

import (
	"encoding/json"
	"knowledge-base/db"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
)

//Article ...
type Article struct {
	ID           int64     `db:"id, primarykey, autoincrement" json:"id"`
	Deleted      bool      `db:"deleted" json:"-" binding:"omitempty"`
	Subject      string    `db:"subject" json:"subject" binding:"required,max=160"`
	Body         string    `db:"body" json:"body" binding:"required"`
	Summary      string    `db:"summary" json:"summary" binding:"omitempty,max=300"`
	CreatedDate  time.Time `db:"created_date" json:"createdDate"`
	ModifiedDate time.Time `db:"modified_date" json:"modifiedDate"`
	Tags         *JSONRaw  `db:"tags" json:"tags" binding:"omitempty"`
	PCC          *JSONRaw  `db:"pcc" json:"pcc" binding:"omitempty"`
	CreatedBy    string    `db:"created_by" json:"createdBy" binding:"omitempty"`
	ModifiedBy   string    `db:"modified_by" json:"modifiedBy" binding:"omitempty"`
	Attachments  *JSONRaw  `db:"attachments" json:"attachments"`
	CustomFields *JSONRaw  `db:"custom_fields" json:"customFields"`
}

//ArticleModel ...
type ArticleModel struct{}

//Create add new article
func (a *Article) Create() error {
	getDB := db.GetDB()
	var ID int64
	t := time.Now()

	query := `INSERT INTO articles (
							subject, body, summary, created_by, modified_by, created_date, modified_date
						) VALUES (
							$1, $2, $3, $4, $5, $6, $7
						) RETURNING id;`
	err := getDB.QueryRow(query, a.Subject, a.Body, a.Summary, a.CreatedBy, a.ModifiedBy, t, t).Scan(&ID)
	if err != nil {
		return err
	}

	if a.Tags != nil {
		//add tags
		go func() {
			var t interface{}
			data, _ := a.Tags.MarshalJSON()
			json.Unmarshal(data, &t)
			var tags []string
			for _, a := range t.([]interface{}) {
				tags = append(tags, a.(string))
			}
			var tagsModel TagsModel
			updateQuery := `UPDATE articles SET tags_string=$1 WHERE id = $2;`
			getDB.Exec(updateQuery, strings.Join(tags, ", "), ID)
			err = tagsModel.Create(tags)
			log.Println(err)
		}()
		updateQuery := `UPDATE articles SET tags=$1 WHERE id = $2;`
		getDB.Exec(updateQuery, a.Tags, ID)
	}
	if a.PCC != nil {
		updateQuery := `UPDATE articles SET pcc=$1 WHERE id = $2;`
		getDB.Exec(updateQuery, a.PCC, ID)
	}
	if a.Attachments != nil {
		updateQuery := `UPDATE articles SET attachments=$1 WHERE id = $2;`
		getDB.Exec(updateQuery, a.Attachments, ID)
	}
	if a.CustomFields != nil {
		updateQuery := `UPDATE articles SET custom_fields=$1 WHERE id = $2;`
		getDB.Exec(updateQuery, a.CustomFields, ID)
	}

	selectQuery := `SELECT
										id, subject, body, summary, created_by, modified_by,
										created_date, modified_date, tags, attachments, custom_fields
									FROM articles WHERE id = $1`

	err = getDB.SelectOne(&a, selectQuery, ID)
	if err != nil {
		return err
	}

	return nil
}

//Update article by ID
func (a *Article) Update(ID int64) error {
	getDB := db.GetDB()
	t := time.Now()

	var article Article
	err := getDB.SelectOne(&article, "SELECT * FROM articles WHERE id = $1 AND deleted = FALSE", ID)

	if err != nil {
		return err
	}
	query := `UPDATE articles SET
						subject = $1, body = $2, summary = $3, created_by = $4,
						modified_by = $5, modified_date = $6
						WHERE id = $7`
	_, err = getDB.Exec(query, a.Subject, a.Body, a.Summary, a.CreatedBy, a.ModifiedBy, t, ID)
	if err != nil {
		return err
	}
	if a.Tags != nil {

		//add tags
		go func() {
			var t interface{}
			data, _ := a.Tags.MarshalJSON()
			json.Unmarshal(data, &t)
			var tags []string
			for _, a := range t.([]interface{}) {
				tags = append(tags, a.(string))
			}
			var tagsModel TagsModel
			updateQuery := `UPDATE articles SET tags_string=$1 WHERE id = $2;`
			getDB.Exec(updateQuery, strings.Join(tags, ", "), ID)
			err = tagsModel.Create(tags)
			log.Println(err)
		}()
		updateQuery := `UPDATE articles SET tags=$1 WHERE id = $2;`
		getDB.Exec(updateQuery, a.Tags, ID)

	}
	if a.PCC != nil {
		updateQuery := `UPDATE articles SET pcc=$1 WHERE id = $2;`
		getDB.Exec(updateQuery, a.PCC, ID)
	}
	if a.Attachments != nil {
		updateQuery := `UPDATE articles SET attachments=$1 WHERE id = $2;`
		getDB.Exec(updateQuery, a.Attachments, ID)
	}
	if a.CustomFields != nil {
		updateQuery := `UPDATE articles SET custom_fields=$1 WHERE id = $2;`
		getDB.Exec(updateQuery, a.CustomFields, ID)
	}
	selectQuery := `SELECT
										id, subject, body, summary, created_by, modified_by,
										created_date, modified_date, tags, attachments, custom_fields
									FROM articles WHERE id = $1`

	err = getDB.SelectOne(&a, selectQuery, ID)
	if err != nil {
		return err
	}

	return nil
}

//Read get articles
func (am ArticleModel) Read(filter string, limit string, offset string) ([]Article, error, int, int) {
	getDB := db.GetDB()
	var articles []Article
	var allArticles []Article

	where, search, err := am.convertFilter(filter)
	query := `SELECT
		id, subject, body, summary, created_by, modified_by,
		created_date, modified_date, tags, attachments, custom_fields
		FROM articles
		WHERE deleted = false`
	if search != "" {
		query = `
			SELECT
			id, subject, body, summary, created_by, modified_by,
			created_date, modified_date, tags, attachments, custom_fields,
			ts_rank_cd(to_tsvector(subject) || to_tsvector(body) || to_tsvector(tags_string), to_tsquery('$1')) as score
			where ts_rank_cd(to_tsvector(subject) || to_tsvector(body) || to_tsvector(tags_string), to_tsquery('$2')) > 0 AND deleted = false
		`

		if where != "" {
			query += " AND " + where
		}
		query += " ORDER BY score DESC"
		_, err := getDB.Select(&allArticles, query, search, search)
		total := len(allArticles)
		var left int
		query += " OFFSET = " + offset + " LIMIT = " + limit
		_, err = getDB.Select(&articles, query)
		if err != nil {
			return nil, err, 0, 0
		}
		count := len(articles)
		if offset != "NULL" {
			l, _ := strconv.Atoi(offset)
			left = int(total) - (l + count)
		} else {
			left = int(total) - count
		}
		return articles, nil, int(total), left
	} else if where != "" {
		query += " AND " + where
	}

	query += " ORDER BY created_date DESC"
	_, err = getDB.Select(&allArticles, query)
	total := len(allArticles)
	var left int
	query += " OFFSET " + offset + " LIMIT " + limit
	_, err = getDB.Select(&articles, query)
	if err != nil {
		return nil, err, 0, 0
	}
	if err != nil {
		return nil, err, 0, 0
	}
	count := len(articles)
	if offset != "NULL" {
		l, _ := strconv.Atoi(offset)
		left = int(total) - (l + count)
	} else {
		left = int(total) - count
	}
	return articles, nil, int(total), left
}

//ReadOne get article by ID
func (am ArticleModel) ReadOne(ID int64) (Article, error) {
	query := `SELECT
							id, subject, body, summary, created_by, modified_by,
							created_date, modified_date, tags, attachments, custom_fields
						FROM articles WHERE deleted = false AND id = $1`
	var article Article
	err := db.GetDB().SelectOne(&article, query, ID)

	return article, err
}

//Delete soft delete
func (am ArticleModel) Delete(ID int64) error {
	t := time.Now()
	query := `UPDATE articles SET deleted=true, modified_date=$1 WHERE id = $2 AND deleted = false;`
	_, err := db.GetDB().Exec(query, t, ID)
	return err
}

func (am ArticleModel) convertFilter(filter string) (where string, search string, err error) {
	if filter != "" {
		jParsed, err := gabs.ParseJSON([]byte(filter))
		if err == nil {
			IDs, err := jParsed.S("where", "id", "in").Children()
			if err == nil {
				ids := make([]string, 0)
				for _, id := range IDs {
					ids = append(ids, id.Data().(string))
				}
				if len(ids) > 0 {
					where += " id IN (" + strings.Join(ids, ", ") + ") "
				}
			}
			search, _ := jParsed.S("where", "search").Data().(string)
			return where, search, nil
		}
	}
	return "", "", nil
}
