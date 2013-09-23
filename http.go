package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Service is a combination of a ServiceHandlerFunc and options, used by the rootServiceHandler
type Service struct {
	fn                ServiceHandlerFunc
	omitClientSession bool
}

// ServiceHandlerFunc defines the layout of a service handler func.. d'oh.
type ServiceHandlerFunc func(w http.ResponseWriter, r *http.Request, cs *ClientSession) (outData interface{}, err error)

var services = map[string]Service{
	// list of services that don not require a clientSession.
	// please keep this list sorted
	"/service/init": Service{newInitHandlerFunc(), true},
}

// initHTTPServer sets up the http.FileServer and other http services.
func initHTTPServer() {
	// add handlers to http.DefaultServeMux
	// try to serve static files on root uri
	http.Handle("/", http.FileServer(http.Dir("./http-files/")))
	http.HandleFunc("/service/", rootServiceHandler)

	// run http server in goroutine
	go func() {
		port := "8000"
		log.Printf("starting http server on port %s\n", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func rootServiceHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	// lookup service
	s, exists := services[r.RequestURI]
	if !exists {
		log.Printf("invalid request (404) for service on: %s\n", r.RequestURI)
		http.NotFound(w, r)
		return
	}

	//++ check origin
	//++ check referer
	fmt.Println("check origin")
	fmt.Println("check referer")

	// find ClientSession
	var cs *ClientSession
	if !s.omitClientSession {
		cs, err = getClientSession(r.Header.Get("X-nulpunt-sessionKey"))
		if err != nil {
			http.Error(w, "forbidden without valid sessionKey", http.StatusForbidden)
			return
		}
		defer cs.done()
	}

	// call actual handler
	outData, err := s.fn(w, r, cs)
	if err != nil {
		log.Printf("error from service %s: %s\n", r.RequestURI, err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// encode response data to client
	err = json.NewEncoder(w).Encode(outData)
	if err != nil {
		log.Printf("error encoding service outData to client for service %s: %s\n", r.RequestURI, err)
	}
}

func newInitHandlerFunc() ServiceHandlerFunc {
	type inDataType struct{}

	type outDataType struct {
		SessionKey string `json:"sessionKey"`
	}

	return func(w http.ResponseWriter, r *http.Request, cs *ClientSession) (interface{}, error) {
		// create a new cs
		newCS := newClientSession()

		// return cs key
		out := &outDataType{
			SessionKey: newCS.key,
		}

		return out, nil
	}
}
