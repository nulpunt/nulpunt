package main

import (
	//"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

// getTrending gets some trending based upon the specified selection.
// Limits to at most limit trends
func getTrending(selection interface{}, limit int) ([]Trending, error) {
	docs := []Trending{}
	err := colTrending.Find(selection).Sort("-score").Limit(limit).All(&docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

// calculateScore creates a score based solely on what's given.
// The score be calulated independent of other documents, is always guaranteed to be correct with respect to the other documents.
// Scores are based upon the timestamps of create/publish of the document and the annotations and comments.
// Scores are increasing. Leaving documents that don't see annotation/comment left behind.
// The score is updated in place.
func calculateScore(trend *Trending) {
	// Score == UploadDate
	trend.Score = float32(trend.Document.UploadDate.Unix()) / 1.0 // make it float.
	// Ignore annotations and comments for now.
}

// updateTrend updates or inserts a trend record
// The key is Trend.Document.ID
// It replaces any old data.
func updateTrend(trend *Trending) error {
	calculateScore(trend) //  in place
	oldTrend := &Trending{}
	err := colTrending.Find(bson.M{"document._id": trend.Document.ID}).One(oldTrend)
	if err == mgo.ErrNotFound {
		trend.ID = bson.NewObjectId()
		err = insertTrend(trend)
		if err != nil {
			log.Printf("error insert trend: %#v\n", err)
		}
		return nil
	}
	if err != nil {
		log.Printf("Unexpected error: %#v\n", err)
		return err
	}

	trend.ID = oldTrend.ID
	_, err = colTrending.UpsertId(trend.ID, trend)
	//log.Printf("Upsert:change info: %#v\n", chinfo)
	if err != nil {
		return err
	}
	// all done
	return nil
}

// insertDocument inserts a new document in the DB. Or updates an existing one.
// Document must have a valid ID, eg from newDocument
func insertTrend(trend *Trending) error {
	err := colTrending.Insert(trend)
	if err != nil {
		return err
	}
	// all done
	return nil
}
