package main

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"sort"
)

var errAnnotationNotUnique = errors.New("We already have an annotation with that ID.")

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
	sort.Sort(AnByDate(annotations))
	return annotations, nil
}

func getLatestAnnotation(selection interface{}) (*Annotation, error) {
	ann := &Annotation{}
	err := colAnnotations.Find(selection).Sort("-createDate").One(ann)
	if err != nil {
		return nil, err
	}
	return ann, nil
}

// insertAnnotation inserts a new annotation in the DB.
// Annotation must be created with newAnnotation
func insertAnnotation(annotation *Annotation) error {
	err := colAnnotations.Insert(annotation)
	if err != nil {
		//mgoErr := err.(*mgo.LastError)
		//if mgoErr.Code == 11000 {
		//	return errAnnotationNotUnique
		//}
		return err
	}
	// all done
	return nil
}

func updateAnnotationID(annotationID bson.ObjectId, change interface{}) error {
	err := colAnnotations.UpdateId(annotationID, change)
	if err != nil {
		return err
	}
	// all done
	return nil
}

// removeAnnotations remove  annotations from the DB
func removeAnnotations(selection interface{}) error {
	err := colAnnotations.Remove(selection)
	// removing an annotation that does not exist is NOT an error. It's just not there.
	if err == mgo.ErrNotFound {
		return nil
	}
	return err
}

// Sorting

type AnByDate []Annotation

func (as AnByDate) Len() int           { return len(as) }
func (as AnByDate) Swap(i, j int)      { as[i], as[j] = as[j], as[i] }
func (as AnByDate) Less(i, j int) bool { return as[i].CreateDate.Before(as[j].CreateDate) }
