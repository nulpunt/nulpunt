package main

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var errTagNotUnique = errors.New("We already have that tag.")

// newTag returns a new Tag struct ready to be inserted into the DB.
func newTag(tag string) *Tag {
	return &Tag{
		ID:  bson.NewObjectId(),
		Tag: tag,
	}
}

// getTags gets all tags from the DB
func getTags() ([]Tag, error) {
	tags := []Tag{}
	err := colTags.Find(nil).All(&tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// insertTag inserts a new tag in the DB.
// Tag must be created with newTag
func insertTag(tag *Tag) error {
	err := colTags.Insert(tag)
	if err != nil {
		mgoErr := err.(*mgo.LastError)
		if mgoErr.Code == 11000 {
			return errTagNotUnique
		}
		return err
	}
	// all done
	return nil
}

// removeTag removes a tag from the DB
func removeTag(tag *Tag) error {
	err := colTags.Remove(bson.M{"tag": tag.Tag})
	return err
}
