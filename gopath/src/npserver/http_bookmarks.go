package main

import (
	"encoding/json"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

// type Bookmark struct is defined in schema.go

// Get your Bookmarks, specified by Session.Account.Username.
// You can only get your own bookmarks.
func getBookmarksHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("getBookmarks-request: %v\n", req.URL)

	// get session
	cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	// get account
	acc := cs.account
	if acc == nil {
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}

	// assemble results into a json-object
	result := make(map[string]interface{})

	// get profile
	bookmarks, err := getBookmarks(acc.Username)
	if err != nil {
		if err == mgo.ErrNotFound {
			// user does not have bookmarks. Create an empty one.
			bookmarks = &Bookmark{
				ID:       bson.NewObjectId(),
				Username: acc.Username,
			}
		} else {
			log.Printf("Bookmarks not found: error %#v\n", err)
			http.Error(rw, "Bookmarks not found", http.StatusNotFound) // 404
			return
		}
	}
	documents, err := getDocuments(bson.M{"_id": bson.M{"$in": bookmarks.DocumentIDs}})
	if err != nil {
		log.Printf("error retrieving bookmarked documents %#v\n", err)
		http.Error(rw, "Bookmarked documents not found", http.StatusInternalServerError) // 500
		return
	}
	result["bookmarks"] = documents

	// marshal and write out.
	j, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling results: error %#v\n", err)
		http.Error(rw, "Marshaling error", http.StatusInternalServerError) // 500
		return
	}
	rw.WriteHeader(200)
	rw.Write(j)
	return
}

type BookmarkEntry struct {
	DocumentID bson.ObjectId
}

func addBookmarkHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("\n\naddBookmark-request: %v\n", req)

	body, _ := ioutil.ReadAll(req.Body)
	log.Printf("\n\nbody is %s\n", string(body))
	entry := &BookmarkEntry{}
	err := json.Unmarshal(body, entry)
	if err != nil {
		log.Printf("\n\nJSON unmarshal error %#v\n", err)
		http.Error(rw, "JSON unmarshal error", http.StatusBadRequest) // 400
		return
	}

	// get session
	cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	// get account
	acc := cs.account
	if acc == nil {
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}

	// Take username from the session.
	// It means you can only update your own bookmarks.
	err = addBookmark(acc.Username, entry.DocumentID)
	if err != nil {
		log.Printf("Error adding entry to bookmarks: error %#v\n", err)
		http.Error(rw, "error adding entry to bookmarks", http.StatusInternalServerError) // 500
		return
	}

	rw.WriteHeader(200)
	rw.Write([]byte(`OK, updated`))
	return
}

// getBookmarks gets one users' bookmarks from the DB
func getBookmarks(username string) (*Bookmark, error) {
	selection := bson.M{"username": username}
	bookmark := &Bookmark{}
	err := colBookmarks.Find(selection).One(bookmark)
	if err != nil {
		return nil, err
	}
	return bookmark, nil
}

func addBookmark(username string, docId bson.ObjectId) error {
	// fetch the users' bookmark record. If not there, create one,
	log.Printf("addBookmark:getBookmarks()")
	_, err := getBookmarks(username)
	if err == mgo.ErrNotFound {
		// user does not have bookmarks-entry, create one.
		bookmarks := &Bookmark{
			ID:          bson.NewObjectId(),
			Username:    username,
			DocumentIDs: []bson.ObjectId{docId},
		}
		log.Printf("inserting new bookmark")
		err := colBookmarks.Insert(bookmarks)
		if err != nil {
			log.Printf("unexpected error inserting bookmark-record: %#v\n", err)
			return err
		}
		// We've created an Bookmarks record for the user,
		return nil // signal correct insert
	}

	if err != nil {
		log.Printf("unexpected error looking up bookmarks: %#v\n", err)
		return err
	}

	log.Printf("Adding bookmark to existing record")
	// user has bookmark-record, push the new entry to it
	change := bson.M{"$push": bson.M{"documentIds": docId}}
	err = colBookmarks.Update(bson.M{"username": username}, change)
	if err != nil {
		return err
	}
	// all done
	return nil
}
