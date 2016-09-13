package store

import "gopkg.in/mgo.v2/bson"

//Tag ...
type Tag struct {
	ID   bson.ObjectId `bson:"_id,omitempty"json:"id"`
	Name string        `bson:"name"json:"name"`
}

//TagsCollection connection to DB
type TagsCollection struct {
	conn       *MongoConnection
	collection string
}

//TagsCollectionConnect return connect to collection Tags
func TagsCollectionConnect() *TagsCollection {
	t := &TagsCollection{
		conn:       db,
		collection: "Tags",
	}
	return t
}

//All return all tags in Tags collection
func (tagsCollection *TagsCollection) All() (result []Tag, err error) {
	session, collection, err := tagsCollection.conn.getSessionAndCollection(tagsCollection.collection)
	if err != nil {
		return
	}
	defer session.Close()
	result = make([]Tag, 0)
	err = collection.Find(bson.M{}).All(&result)

	if err != nil {
		return
	}

	return result, nil
}

//Create new tags
func (tagsCollection *TagsCollection) Create(entries []Tag) (err error) {
	session, collection, err := tagsCollection.conn.getSessionAndCollection(tagsCollection.collection)
	if err != nil {
		return
	}
	defer session.Close()

	for _, entre := range entries {
		err = collection.Insert(&entre)
	}
	return
}
