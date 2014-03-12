package main

import (
	//"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
	//"log"
)

// getTrending gets some trending based upon the specified selection.
// Limits to at most limit documents
func getTrendingDocs(limit int) ([]Document, error) {
	docs := []Document{}
	selection := bson.M{"published": true}
	err := colDocuments.Find(selection).Sort("-score").Limit(limit).All(&docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

// calculateTrendingScore creates a score based solely on what's given.
// The score is calulated independent of other documents.
// It is always guaranteed to be correct with respect to the other documents.
// Scores are based upon the timestamps of create/publish of the document and the annotations and comments.
// Scores are increasing at every document/annotation-update.
// Leaving documents behind that don't see updates.
// The score is returned
func calculateTrendingScore(document *Document, annotation *Annotation) float32 {
	latest := document.UploadDate
	if annotation != nil && annotation.CreateDate.After(latest) {
		latest = annotation.CreateDate
	}
	score := dateToTrending(latest)
	document.Score = score
	return score
}

func dateToTrending(date time.Time) float32 {
	return float32(date.Unix())
}

func updateDocumentScore(docId bson.ObjectId, score float32) error {
	err := colDocuments.UpdateId(docId, bson.M{"$set": bson.M{"score": score}})
	return err
}
