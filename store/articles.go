package store

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Article store article
type Article struct {
	ID           bson.ObjectId  `bson:"_id,omitempty"json:"id"`
	Deleted      bool           `bson:"deleted"json:"deleted"validate:"omitempty"`
	Subject      string         `bson:"subject"json:"subject"validate:"required,min=2,max=160"`
	Body         string         `bson:"body"json:"body"validate:"required"`
	Summary      string         `bson:"summary"json:"summary"validate:"omitempty,min=2,max=300"`
	CreatedDate  time.Time      `bson:"createdDate"json:"createdDate"`
	ModifiedDate time.Time      `bson:"modifiedDate"json:"modifiedDate"`
	Tags         []string       `bson:"tags"json:"tags"validate:"omitempty"`
	Attachments  []Attachment   `bson:"attachments"json:"attachments"validate:"omitempty"`
	CustomFields []CustomFields `bson:"customFields"json:"customFields"validate:"omitempty"`
	PCC          []string       `bson:"pcc"json:"pcc"validate:"omitempty"`
	CreatedBy    string         `bson:"createdBy"json:"createdBy"validate:"omitempty"`
	ModifiedBy   string         `bson:"modifiedBy"json:"modifiedBy"validate:"omitempty"`
}

//Attachment store link and file name
type Attachment struct {
	URL      string `bson:"url"json:"url"validate:"omitempty"`
	FileName string `bson:"fileName"json:"fileName"validate:"omitempty"`
}

//CustomFields store custome fields for article
type CustomFields struct {
	Key   string `bson:"key"json:"key"validate:"omitempty"`
	Value string `bson:"value"json:"value"validate:"omitempty"`
}

//ArticlesCollection connection to DB
type ArticlesCollection struct {
	conn       *MongoConnection
	collection string
}

//ArticlesCollectionConnect return connect to collection Articles
func ArticlesCollectionConnect() *ArticlesCollection {
	t := &ArticlesCollection{
		conn:       db,
		collection: "Articles",
	}
	return t
}

//Read return entries of Articles collections
func (art *ArticlesCollection) Read(query bson.M, fields bson.M, skip int, limit int) (result []Article, total int, left int, err error) {
	session, artCollection, err := art.conn.getSessionAndCollection(art.collection)
	if err != nil {
		return
	}
	defer session.Close()
	query = checkUpdated(query)
	result = make([]Article, 0)
	err = artCollection.Find(query).Select(fields).Skip(skip).Limit(limit).All(&result)
	if err != nil {
		return
	}
	count := len(result)
	total, err = artCollection.Find(query).Count()
	if err != nil {
		return
	}
	left = total - (skip + count)
	if left < 0 {
		left = 0
	}
	return result, total, left, nil
}

//ReadOne return one enrty of Articles collection by query
func (art *ArticlesCollection) ReadOne(query bson.M, fields bson.M) (result Article, err error) {
	session, artCollection, err := art.conn.getSessionAndCollection(art.collection)
	if err != nil {
		return
	}
	defer session.Close()
	err = artCollection.Find(query).Select(fields).One(&result)

	if err != nil {
		return
	}

	return result, nil
}

//Create add new article
func (art *ArticlesCollection) Create(entry Article) (result Article, err error) {
	session, artCollection, err := art.conn.getSessionAndCollection(art.collection)
	if err != nil {
		return
	}
	defer session.Close()

	entry.ID = bson.NewObjectId()
	entry.CreatedDate = time.Now().UTC()
	entry.ModifiedDate = time.Now().UTC()
	entry.Deleted = false

	err = artCollection.Insert(entry)

	if err != nil {
		return
	}
	return entry, nil
}

//Update existing entry
func (art *ArticlesCollection) Update(ID bson.ObjectId, entry Article) (err error) {
	session, artCollection, err := art.conn.getSessionAndCollection(art.collection)
	if err != nil {
		return
	}
	defer session.Close()
	entry.ModifiedDate = time.Now().UTC()
	query := bson.M{
		"_id": ID,
	}
	change := bson.M{
		"$set": bson.M{
			"deleted":      entry.Deleted,
			"subject":      entry.Subject,
			"body":         entry.Body,
			"summary":      entry.Summary,
			"modifiedDate": entry.ModifiedDate,
			"tags":         entry.Tags,
			"attachments":  entry.Attachments,
			"customFields": entry.CustomFields,
			"pcc":          entry.PCC,
			"createdBy":    entry.CreatedBy,
			"modifiedBy":   entry.ModifiedBy,
		},
	}

	err = artCollection.Update(query, change)
	if err != nil {
		return
	}

	return nil
}

//Delete entry from Articles collection (soft delete)
func (art *ArticlesCollection) Delete(ID bson.ObjectId) (err error) {
	session, artCollection, err := art.conn.getSessionAndCollection(art.collection)
	if err != nil {
		return
	}
	defer session.Close()
	query := bson.M{
		"_id": ID,
	}
	change := bson.M{
		"$set": bson.M{
			"deleted":      true,
			"modifiedDate": time.Now().UTC(),
		},
	}

	err = artCollection.Update(query, change)
	return
}

//Search implement full-text search
func (art *ArticlesCollection) Search(q bson.M, skip int, limit int) (result []Article, total int, left int, err error) {
	session, artCollection, err := art.conn.getSessionAndCollection(art.collection)
	if err != nil {
		return
	}
	defer session.Close()
	q = checkUpdated(q)
	fields := bson.M{
		"score": bson.M{
			"$meta": "textScore",
		},
	}
	sort := "$textScore:score"
	result = make([]Article, 0)
	err = artCollection.Find(q).Select(fields).Sort(sort).Skip(skip).Limit(limit).All(&result)
	if err != nil {
		return
	}
	count := len(result)
	total, err = artCollection.Find(q).Select(fields).Sort(sort).Count()
	if err != nil {
		return
	}
	left = total - (skip + count)
	if left < 0 {
		left = 0
	}
	if err != nil {
		return
	}

	return result, total, left, nil
}
