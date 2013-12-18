package main

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

var errAnnotationNotUnique = errors.New("We already have an annotation with that ID.")

// type Annotation hold the document annotations
type Annotation struct {
	ID          bson.ObjectId `bson:"_id"`
	AnnotatorID bson.ObjectId
	CreateDate  time.Time
	Annotation  string
	Location    []Location
	Comments    []Comment
}

type Location struct {
	Page int
	X1   int
	Y1   int
	X2   int
	Y2   int
}

type Comment struct {
	ID         bson.ObjectId `bson:"_id"`
	Commenter  string
	CreateDate time.Time
	Comment    string
	Comments   []Comment
}

// newAnnotation returns a new Annotation struct ready to be inserted into the DB.
func newAnnotation() *Annotation {
	return &Annotation{
		ID: bson.NewObjectId(),
	}
}

// getAnnotation gets one annotation from the DB
func getAnnotation(selection interface{}) (*Annotation, error) {
	annotation := &Annotation{}
	err := colAnnotations.Find(selection).One(annotation)
	if err != nil {
		return nil, err
	}
	return annotation, nil
}

// getAnnotations gets all annotations from the DB
func getAnnotations(selection interface{}) ([]Annotation, error) {
	annotations := []Annotation{}
	err := colAnnotations.Find(selection).All(&annotations)
	if err != nil {
		return nil, err
	}
	return annotations, nil
}

// insertAnnotation inserts a new annotation in the DB.
// Annotation must be created with newAnnotation
func insertAnnotation(annotation *Annotation) error {
	err := colAnnotations.Insert(annotation)
	if err != nil {
		mgoErr := err.(*mgo.LastError)
		if mgoErr.Code == 11000 {
			return errAnnotationNotUnique
		}
		return err
	}
	// all done
	return nil
}

// removeAnnotation removes a annotation from the DB
// func removeAnnotation(annotation *Annotation) error {
// 	err := colAnnotations.Remove(bson.M{"annotation": annotation.Annotation})
// 	return err
// }
