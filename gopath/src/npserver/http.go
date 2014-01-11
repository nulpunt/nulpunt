package main

import (
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

const headerKeySessionKey = `X-Nulpunt-SessionKey`

// initHTTPServer sets up the http.FileServer and other http services.
func initHTTPServer() {
	// create fileServer that servces that static files (html,css,js,etc.)
	fileServer := http.FileServer(http.Dir(flags.HTTPFiles))

	// rootRouter is directly linked to the http server.
	rootRouter := mux.NewRouter()

	// serve static files on / and several subdirs
	rootRouter.Path("/").Handler(fileServer)
	rootRouter.PathPrefix("/css/").Handler(fileServer)
	rootRouter.PathPrefix("/fonts/").Handler(fileServer)
	rootRouter.PathPrefix("/html/").Handler(fileServer)
	rootRouter.PathPrefix("/js/").Handler(fileServer)
	rootRouter.PathPrefix("/img/").Handler(fileServer)
	rootRouter.Path("/download-original/{documentIDHex}/{filename:.*}").HandlerFunc(downloadOriginalHandlerFunc)

	// serve document files on /docfiles/
	docfilesRouter := rootRouter.PathPrefix("/docfiles/").Subrouter()
	docfilesRouter.Path("/pages/{documentIDHex}/{pageNumber}.png").HandlerFunc(pageImageHandlerFunc)
	docfilesRouter.Path("/thumbnails/{documentIDHex}.png").HandlerFunc(thumbnailImageHandlerFunc)

	// create serviceRouter for /service/*
	serviceRouter := rootRouter.PathPrefix("/service/").Subrouter()
	serviceRouter.Path("/sessionInit").HandlerFunc(sessionInitHandlerFunc)
	serviceRouter.Path("/sessionCheck").HandlerFunc(sessionCheckHandlerFunc)

	// Document handlers
	serviceRouter.Path("/getDocument").HandlerFunc(getDocumentHandler)
	serviceRouter.Path("/getPage").HandlerFunc(getPageHandlerFunc)
	serviceRouter.Path("/getDocuments").HandlerFunc(getDocumentsHandler)
	serviceRouter.Path("/getDocumentList").HandlerFunc(getDocumentListHandler)

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
	sessionRouter.Path("/resume").HandlerFunc(sessionResumeHandlerFunc)
	sessionRouter.Path("/dataBlobSave").HandlerFunc(sessionDataBlobSave)
	sessionRouter.Path("/dataBlobLoad").HandlerFunc(sessionDataBlobLoad)

	sessionRouter.Path("/add-annotation").HandlerFunc(addAnnotationHandler)
	sessionRouter.Path("/add-comment").HandlerFunc(addCommentHandler)

	// Users can't make profiles yet. BLOCK em.
	//sessionRouter.Path("/get-profile").HandlerFunc(getProfileHandler)
	//sessionRouter.Path("/update-profile").HandlerFunc(updateProfileHandler)

	sessionRouter.Path("/get-documents-by-tags").HandlerFunc(getDocumentsByTagsHandler)

	// register /service/session/admin/* handlers
	adminRouter := sessionRouter.PathPrefix("/admin/").Subrouter()
	adminRouter.Path("/upload").HandlerFunc(adminUpload)
	adminRouter.Path("/getRawUploads").HandlerFunc(adminGetRawUploads)

	sessionRouter.Path("/get-tags").HandlerFunc(adminGetTags) //  /service/get-tags, ie for all
	adminRouter.Path("/add-tag").HandlerFunc(adminAddTag)     //  /service/add-tags, ie only for admins
	adminRouter.Path("/delete-tag").HandlerFunc(adminDeleteTag)

	adminRouter.Path("/updateDocument").HandlerFunc(updateDocumentHandler)
	adminRouter.Path("/insertDocument").HandlerFunc(insertDocumentHandler)
	adminRouter.Path("/deleteDocument").HandlerFunc(deleteDocumentHandler)

	adminRouter.PathPrefix("/").Handler(http.NotFoundHandler())

	// 404 when /service/session/* was not found
	sessionRouter.PathPrefix("/").Handler(http.NotFoundHandler())

	// when session auth failed, return 403 forbidden for /service/session/*
	sessionPathRouter.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden, invalid session key", http.StatusForbidden)
	})

	// run http server in goroutine
	go func() {
		// inform user of startup
		log.Printf("starting http server on http://localhost:%s\n", flags.HTTPPort)

		// listen and serve on given port
		// error is fatal
		err := http.ListenAndServe(":"+flags.HTTPPort, rootRouter)
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
			err = http.Serve(socket, rootRouter)
			if !socketClosing && err != nil {
				log.Fatalf("fatal error serving http on the unix socket: %s\n", err)
			}
		}()
	}
}
