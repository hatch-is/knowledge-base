package model

import "knowledge-base/store"

//TagsModel ...
type TagsModel struct{}

//All get all tags
func (tagModel *TagsModel) All() (result []string, err error) {
	tagDb := store.TagsCollectionConnect()
	tags, err := tagDb.All()

	result = make([]string, 0)
	if err != nil {
		return
	}
	for _, tag := range tags {
		result = append(result, tag.Name)
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
