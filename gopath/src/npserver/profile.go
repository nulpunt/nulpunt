package main

import (
	"errors"
	"labix.org/v2/mgo/bson"
	"log"
)

var errProfileNotUnique = errors.New("Profile not unique")

// type Profile holds the profile properties
type Profile struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string
	Tags     []string // contains tag.Tag
}

// newProfile returns a new empty one.
func newProfile() *Profile {
	return &Profile{
		ID: bson.NewObjectId(),
	}
}

// getProfile gets a single profile based upon the specified selection.
// Selection must lead to a unique profile. Otherwise, results are undefined
func getProfile(selection interface{}) (*Profile, error) {
	doc := &Profile{}
	err := colProfiles.Find(selection).One(doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// upsertProfile inserts a new profile in the DB. Or updates an existing one.
// Profile must have a valid ID, eg from newProfile
func upsertProfile(doc *Profile) error {
	chinfo, err := colProfiles.UpsertId(doc.ID, doc)
	log.Printf("Upsert:change info: %#v\n", chinfo)
	if err != nil {
		// mgoErr := err.(*mgo.LastError)
		// if mgoErr.Code == 11000 {
		// 	return errProfileNotUnique
		// }
		return err
	}
	// all done
	return nil
}

// removeProfile removes a profile from the DB
//func removeTag(prf *Profile) error {
//	err := colProfiles.RemoveId(profile.ID)
//	return err
//}
