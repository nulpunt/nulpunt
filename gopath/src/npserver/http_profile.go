package main

import (
	"encoding/json"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

// type Profile struct is defined in profile.go

// Get a single profile, specified by Session.Account.Username.
// You can only get your own profile.
func getProfileHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("getProfile-request: %v\n", req)

	// assemble results into a json-object
	result := make(map[string]interface{})

	switch req.Method {
	case "GET":
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

		// get profile
		profile, err := getProfile(bson.M{"username": acc.Username})
		if err != nil {
			if err == mgo.ErrNotFound {
				// user does not have a profile. Create an empty one.
				profile = &Profile{
					ID:       bson.NewObjectId(),
					Username: acc.Username,
				}
			} else {
				log.Printf("Profile not found: error %#v\n", err)
				http.Error(rw, "Profile not found", http.StatusNotFound) // 404
				return
			}
		}
		result["profile"] = profile

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

	default: // request.Method
		http.Error(rw, "error", http.StatusMethodNotAllowed) // 405
	}
}

func updateProfileHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("\n\nupdateProfile-request: %v\n", req)

	switch req.Method {
	case "POST":
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("\n\nbody is %s\n", string(body))
		profile := &Profile{}
		err := json.Unmarshal(body, profile)
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

		// Set username to that of the session.
		// It means you can only update your own.
		// TODO: verify if usernames match and give a Forbidden.
		profile.Username = acc.Username

		log.Printf("\n\nProfile to update is: %#v\n", profile)
		err = upsertProfile(profile)
		if err != nil {
			log.Printf("Error inserting/updating  profile: error %#v\n", err)
			http.Error(rw, "error inserting/updating profile", http.StatusInternalServerError) // 500
			return
		}

		rw.WriteHeader(200)
		rw.Write([]byte(`OK, updated`))
		return
	default:
		http.Error(rw, "error", http.StatusMethodNotAllowed) // 405
	}
}
