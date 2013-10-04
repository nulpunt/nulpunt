package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net"
	"net/http"
)

const headerKeySessionKey = `X-Nulpunt-SessionKey`

// initHTTPServer sets up the http.FileServer and other http services.
func initHTTPServer() {
	// normally, rootRouter would be directly linked to the http server.
	// during closed alpha, the alphaRouter takes over, it checks for closed-alpha credentials.
	// when everything is ok, the rootRouter is allowed to handle the requests.
	alphaRouter := mux.NewRouter()

	// proceed to the rootRouter when basic auth is satisfied
	rootRouter := alphaRouter.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return alphaCheckBasicAuth(r)
	}).Subrouter()

	// otherwise present request for basic auth
	alphaRouter.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return !alphaCheckBasicAuth(r)
	}).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Nulpunt alpha access"`)
		http.Error(w, "Please enter valid Nulpunt alpha credentials", http.StatusUnauthorized)
	})

	// serve static files on / and several subdirs
	fileServer := http.FileServer(http.Dir("./http-files/"))
	rootRouter.Path("/").Handler(fileServer)
	rootRouter.PathPrefix("/css/").Handler(fileServer)
	rootRouter.PathPrefix("/fonts/").Handler(fileServer)
	rootRouter.PathPrefix("/html/").Handler(fileServer)
	rootRouter.PathPrefix("/js/").Handler(fileServer)

	// create serviceRouter for everything beneath /service/
	serviceRouter := rootRouter.PathPrefix("/service/").Subrouter()
	serviceRouter.Path("/service/sessionInit").HandlerFunc(sessionInitHandlerFunc)
	// serviceRouter.Path("/service/").Handler(http.NotFoundHandler())

	// create sessionRouter for everything beneath /service/session/
	sessionRouter := rootRouter.PathPrefix("/service/session/").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		//++ TODO: should be simple `return checkClientSessionValid(key string)`
		// 			locking is unnecicary
		cs, err := getClientSession(r.Header.Get(headerKeySessionKey))
		if err != nil {
			return false
		}
		cs.done()
		return true
	}).Subrouter()
	sessionRouter.Path("/service/session/destroy").HandlerFunc(sessionDestroyHandlerFunc)
	serviceRouter.Path("/service/session/check").HandlerFunc(sessionCheckHandlerFunc)
	// sessionRouter.Path("/service/session/").Handler(http.NotFoundHandler())

	// send 403 forbidden to requests on /service/session/* that don't have a valid session
	serviceRouter.PathPrefix("/service/session/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden, invalid session key", http.StatusForbidden)
	})

	// run http server in goroutine
	go func() {
		//++ TODO: make configurable
		port := "8000"

		// inform user of startup
		log.Printf("starting http server on port %s\n", port)

		// listen and serve on given port
		// error is fatal
		err := http.ListenAndServe(":"+port, alphaRouter)
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
			err = http.Serve(socket, alphaRouter)
			if !socketClosing && err != nil {
				log.Fatal(err)
			}
		}()
	}
}

// sessionInitHandlerFunc creates a new session
func sessionInitHandlerFunc(w http.ResponseWriter, r *http.Request) {
	type outDataType struct {
		SessionKey string `json:"sessionKey"`
	}

	// create a new cs
	newCS := newClientSession()

	// return cs key
	out := &outDataType{
		SessionKey: newCS.key,
	}
	err := json.NewEncoder(w).Encode(out)
	if err != nil {
		log.Println("Could not encode output data for service sessionInitHandlerFunc. %s\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

// sessionCheckHandlerFunc checks if session is ok
func sessionCheckHandlerFunc(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "session ok")
}

// sessionDestroyHandlerFunc destroys a session
func sessionDestroyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	cs, err := getClientSession(r.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	cs.ping <- false //++ TODO: probable race when this is called more then once
	cs.done()
}
