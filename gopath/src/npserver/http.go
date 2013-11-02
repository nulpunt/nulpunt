package main

import (
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"strings"
)

const headerKeySessionKey = `X-Nulpunt-SessionKey`

// initHTTPServer sets up the http.FileServer and other http services.
func initHTTPServer() {
	// create fileServer that servces that static files (html,css,js,etc.)
	fileServer := http.FileServer(http.Dir(flags.HTTPFiles))

	// normally, rootRouter would be directly linked to the http server.
	// during closed alpha, the alphaRouter takes over, it checks for closed-alpha credentials.
	// when everything is ok, the rootRouter is allowed to handle the requests.
	alphaRouter := mux.NewRouter()

	// Chrome doesn't apply basic auth to requests for map files
	// therefore we must work arround this and provide the map files without basic auth
	// TODO FIX WORKARROUND NASTY YOLO ALPHA
	alphaRouter.PathPrefix("/js/").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return strings.HasSuffix(r.RequestURI, ".map")
	}).Handler(fileServer)

	// proceed to the rootRouter when basic auth is satisfied
	rootRouter := alphaRouter.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		if flags.DisableAlphaAuth {
			return true
		}
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
	rootRouter.Path("/").Handler(fileServer)
	rootRouter.PathPrefix("/css/").Handler(fileServer)
	rootRouter.PathPrefix("/fonts/").Handler(fileServer)
	rootRouter.PathPrefix("/html/").Handler(fileServer)
	rootRouter.PathPrefix("/js/").Handler(fileServer)

	// create serviceRouter for /service/*
	serviceRouter := rootRouter.PathPrefix("/service/").Subrouter()
	serviceRouter.Path("/sessionInit").HandlerFunc(sessionInitHandlerFunc)
	serviceRouter.Path("/sessionCheck").HandlerFunc(sessionCheckHandlerFunc)

	// create sessionPathRouter for /service/session/*
	sessionPathRouter := rootRouter.PathPrefix("/service/session/").Subrouter()

	// sessionRouter handles valid authenticated requests for /service/session
	sessionRouter := sessionPathRouter.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		sessionKey := r.Header.Get(headerKeySessionKey)
		return isValidClientSession(sessionKey)
	}).Subrouter()

	// register /service/session/* handlers
	sessionRouter.Path("/ping").HandlerFunc(sessionPingHandlerFunc)
	sessionRouter.Path("/destroy").HandlerFunc(sessionDestroyHandlerFunc)
	sessionRouter.Path("/registerAccount").HandlerFunc(registerAccountHandlerFunc)
	sessionRouter.Path("/authenticateAccount").HandlerFunc(sessionAuthenticateAccountHandlerFunc)
	sessionRouter.Path("/dataBlobSave").HandlerFunc(sessionDataBlobSave)
	sessionRouter.Path("/dataBlobLoad").HandlerFunc(sessionDataBlobLoad)

	// 404 when /service/session/* was not found
	sessionRouter.PathPrefix("/").Handler(http.NotFoundHandler())

	// when session auth failed, return 403 forbidden for /service/session/*
	sessionPathRouter.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			log.Fatalf("fatal error listening/serving http on tcp: %s\n", err)
		}
	}()

	if len(flags.UnixSocket) > 0 {
		go func() {
			// socketClosing is used to omit error on socket read when closing down.
			var socketClosing bool

			// inform user of startup
			log.Printf("Starting http server on unix socket %s\n", flags.UnixSocket)

			// create and listen on this unix socket
			socket, err := net.ListenUnix("unix", &net.UnixAddr{
				Name: flags.UnixSocket,
				Net:  "unix",
			})
			if err != nil {
				log.Fatalf("fatal error on listening on unix socket: %s\n", err)
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
				log.Fatalf("fatal error serving http on the unix socket: %s\n", err)
			}
		}()
	}
}
