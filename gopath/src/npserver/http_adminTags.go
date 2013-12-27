package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func adminGetTags(rw http.ResponseWriter, req *http.Request) {
	// Ignore req.method. We only return data, who cares about the Method.
	log.Printf("tag-request: %v\n", req)

	// // get session
	// cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
	// if err != nil {
	// 	http.Error(rw, "error", http.StatusInternalServerError)
	// 	return
	// }
	// defer cs.done()

	// // get accuont
	// acc := cs.account
	// if acc == nil {
	// 	http.Error(rw, "forbidden", http.StatusForbidden)
	// 	return
	// }

	// Get the tags and send them out.
	getEm(rw, req)
	return
}

// getEm gets the tags and sends them out.
func getEm(rw http.ResponseWriter, req *http.Request) {
	tags, err := getTags()
	j, err := json.Marshal(map[string]interface{}{"tags": tags})
	if err != nil {
		http.Error(rw, "error", http.StatusInternalServerError) // 500
		return
	}
	rw.WriteHeader(200)
	rw.Write(j)
	return
}

func adminAddTag(rw http.ResponseWriter, req *http.Request) {
	log.Printf("tag-request: %v\n", req)

	// get session
	cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	// get account
	acc := cs.account
	if acc == nil || acc.Admin == false {
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}

	switch req.Method {
	case "POST":
		body, _ := ioutil.ReadAll(req.Body)
		tag := &Tag{}
		err := json.Unmarshal(body, tag)

		if err != nil {
			log.Printf("Tag is empty.\n")
			http.Error(rw, "error", http.StatusBadRequest) // 400
			return
		}

		if tag.Tag == "" {
			log.Printf("Tag is empty.\n")
			http.Error(rw, "error", http.StatusBadRequest) // 400
			return
		}
		// Todo UPDATE complete tag when tag.ID != nil.
		// Now, we just want the string value, to insert.
		err = insertTag(newTag(tag.Tag))
		if err != nil {
			log.Printf("insertTag error: %v\n", err)
			http.Error(rw, "We already have that tag.", http.StatusInternalServerError) // 500
			return
		}

		// Get the tags and send them out.
		getEm(rw, req)
		return

	default:
		http.Error(rw, "error", http.StatusMethodNotAllowed) // 405
	}
}

func adminDeleteTag(rw http.ResponseWriter, req *http.Request) {
	log.Printf("tag-request: %v\n", req)

	// get session
	cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	// get account
	acc := cs.account
	if acc == nil || acc.Admin == false {
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}

	switch req.Method {
	case "POST":
		body, _ := ioutil.ReadAll(req.Body)
		tag := &Tag{}
		err := json.Unmarshal(body, tag)

		if err != nil {
			log.Printf("Tag is empty.\n")
			http.Error(rw, "error", http.StatusBadRequest) // 400
			return
		}

		if tag.Tag == "" {
			log.Printf("Tag is empty.\n")
			http.Error(rw, "error", http.StatusBadRequest) // 400
			return
		}

		err = removeTag(tag)
		if err != nil {
			log.Printf("removeTag error: %v\n", err)
			http.Error(rw, "Tag wasn't there", http.StatusInternalServerError) // 500
			return
		}

		// Get the tags and send them out.
		getEm(rw, req)
		return
		// All done, get the results and show the updated list
	}
}
