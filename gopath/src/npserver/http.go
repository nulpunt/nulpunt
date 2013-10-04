package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

// Service is a combination of a ServiceHandlerFunc and options
// it is used by the rootServiceHandler that performs checks (depending on the options)
// when the rootServiceHandler is satitsfied, the function in this Service object is called
type Service struct {
	fn                ServiceHandlerFunc
	omitClientSession bool
}

// ServiceHandlerFunc defines the layout of a service handler func.. d'oh.
type ServiceHandlerFunc func(w http.ResponseWriter, r *http.Request, cs *ClientSession) (outData interface{}, err error)

// services is a list containing all registered Service instances
var services = map[string]Service{
	// NOTE: please keep this list sorted
	"/service/sessionInit":    Service{newSessionInitHandlerFunc(), true},
	"/service/sessionCheck":   Service{newSessionCheckHandlerFunc(), true},
	"/service/sessionDestroy": Service{newSessionDestroyHandlerFunc(), false},
}

// initHTTPServer sets up the http.FileServer and other http services.
func initHTTPServer() {
	// add handlers to http.DefaultServeMux
	// try to serve static files on root uri
	http.Handle("/", http.FileServer(http.Dir("./http-files/")))
	http.HandleFunc("/service/", rootServiceHandler)

	// run http server in goroutine
	go func() {
		//++ TODO: make configurable
		port := "8000"

		// inform user of startup
		log.Printf("starting http server on port %s\n", port)

		// listen and serve on given port
		// error is fatal
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if flags.UnixSocket {
		go func() {
			// socketClosign is used to omit error on socket read when closing down.
			var socketClosing bool

			//++ TODO: make configurable
			socketFilename := "./npserver.socket"

			// inform user of startup
			log.Printf("Starting http server on unix socket %s\n", socketFilename)

			// create and listen on this unix socket
			socket, err := net.ListenUnix("unix", &net.UnixAddr{
				Name: socketFilename,
				Net:  "unix",
			})
			if err != nil {
				log.Fatal(err)
			}

			// append a function on graceful shutdown to close the unix socket
			processEndFuncs = append(processEndFuncs, func() {
				socketClosing = true
				socket.Close()
			})

			// serve on the opened unix socket
			// an error (when not closing down) is fatal
			err = http.Serve(socket, nil)
			if !socketClosing && err != nil {
				log.Fatal(err)
			}
		}()
	}
}

// rootServiceHandler handles every service request in a generic way
// service requests are checked for autenticity etc
func rootServiceHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	// lookup service
	s, exists := services[r.RequestURI]
	if !exists {
		log.Printf("invalid request (404) for service on: %s\n", r.RequestURI)
		http.NotFound(w, r)
		return
	}

	//++ TODO: check origin
	fmt.Println("check origin")

	//++ TODO: check referer
	fmt.Println("check referer")

	// find ClientSession
	var cs *ClientSession
	if !s.omitClientSession {
		cs, err = getClientSession(r.Header.Get("X-Nulpunt-SessionKey"))
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

	// we're responding with json
	w.Header().Add("Content-Type", "application/json")

	// encode response data to json for client
	err = json.NewEncoder(w).Encode(outData)
	if err != nil {
		log.Printf("error encoding service outData to client for service %s: %s\n", r.RequestURI, err)
	}
}

// newSessionInitHandlerFunc creates a handler for new session initialization
func newSessionInitHandlerFunc() ServiceHandlerFunc {
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

// newSessionCheckHandlerFunc creates a handler for session checks
func newSessionCheckHandlerFunc() ServiceHandlerFunc {
	type inDataType struct {
		SessionKey string `json:"sessionKey"`
	}

	type outDataType struct {
		Result bool `json:"result"`
	}

	return func(w http.ResponseWriter, r *http.Request, cs *ClientSession) (interface{}, error) {
		// decode input data
		inData := &inDataType{}
		err := json.NewDecoder(r.Body).Decode(inData)
		if err != nil {
			return nil, err
		}

		// new outData
		outData := &outDataType{}

		// get ClientSession
		s, err := getClientSession(inData.SessionKey)
		if err != nil {
			log.Printf("Could not find CS for key %s\n", inData.SessionKey)
			return outData, nil
		}

		// already done
		s.done()

		// session is valid
		outData.Result = true

		// return data
		return outData, nil
	}
}

// newSessionDestroyHandlerFunc creates a handler for session destroy
func newSessionDestroyHandlerFunc() ServiceHandlerFunc {
	type outDataType struct{}

	return func(w http.ResponseWriter, r *http.Request, cs *ClientSession) (interface{}, error) {
		cs.ping <- false //++ TODO: race when this is called more then once
		return nil, nil
	}
}
