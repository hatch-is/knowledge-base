package model

import "knowledge-base/store"

//TagsModel ...
type TagsModel struct{}

//All get all tags
func (tagModel *TagsModel) All() (result []store.Tag, err error) {
	tagDb := store.TagsCollectionConnect()
	result, err = tagDb.All()

	if err != nil {
		return
	}
	return result, nil
}

//CampareAndCreate create tags if they not exists
func (tagModel *TagsModel) CampareAndCreate(sTags []string) (err error) {
	if len(sTags) > 0 {
		tagDb := store.TagsCollectionConnect()
		tags, _ := tagDb.All()
		var newTags []store.Tag
		for _, sTag := range sTags {
			bExists := false
			for _, tag := range tags {
				if tag.Name == sTag {
					bExists = true
					break
				}
			}
			if bExists == false {
				newTags = append(newTags, store.Tag{Name: sTag})
			}
		}
		if len(newTags) > 0 {
			err = tagDb.Create(newTags)
		}
	}
	return
}
