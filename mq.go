package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var (
	workers int
	peers   int
	root    string
	port    string
	address string
)

type FrontHandler struct {
	Store     *Store
	Router    *Router
	Endpoints map[string]func(*Session)
}

type Session struct {
	Store    *Store
	Match    *RouteMatch
	Request  *http.Request
	Response http.ResponseWriter
}

func (handler *FrontHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	match := handler.Router.Match(request.Method, request.URL.Path)

	if match != nil {
		session := &Session{
			Store:    handler.Store,
			Match:    match,
			Request:  request,
			Response: response,
		}

		handler.Endpoints[match.Name](session)
		return
	}

	response.WriteHeader(http.StatusNotFound)
}

func init() {
	flag.IntVar(&workers, "workers", 8, "Number of workers")
	flag.IntVar(&peers, "peers", 0, "Number of peers")
	flag.StringVar(&root, "root", "/tmp/mq", "File system storage path")
	flag.StringVar(&port, "port", "8080", "Port to listen on")
	flag.StringVar(&address, "address", "0.0.0.0", "Address to listen on")
	flag.Parse()
}

func main() {
	store := NewStore(workers, peers, root)

	// Our storage mechanism needs to make sure our folders
	// and workers are standing up.
	store.PrepareFolders()
	store.PrepareWorkers()

	router := &Router{}
	router.AddRoute("GetQueue", "GET", "^/(?P<queue>[a-z]+)$")
	router.AddRoute("CreateQueue", "PUT", "^/(?P<queue>[a-z]+)$")
	router.AddRoute("DeleteQueue", "DELETE", "^/(?P<queue>[a-z]+)$")
	router.AddRoute("CreateMessage", "POST", "^/(?P<queue>[a-z]+)/messages$")
	router.AddRoute("GetMessage", "GET", "^/(?P<queue>[a-z]+)/messages$")
	router.AddRoute("DeleteMessage", "DELETE", "^/(?P<queue>[a-z]+)/messages/(?P<message>[a-z0-9-]+)$")

	handler := &FrontHandler{
		Store:     store,
		Router:    router,
		Endpoints: make(map[string]func(*Session)),
	}

	// Our handler functions by name. This can easily be looked up by the name
	// our RouteMatch contains.
	handler.Endpoints["GetQueue"] = GetQueue
	handler.Endpoints["CreateQueue"] = CreateQueue
	handler.Endpoints["DeleteQueue"] = DeleteQueue
	handler.Endpoints["CreateMessage"] = CreateMessage
	handler.Endpoints["GetMessage"] = GetMessage
	handler.Endpoints["DeleteMessage"] = DeleteMessage

	server := &http.Server{
		Addr:           address + ":" + port,
		Handler:        handler,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println(server.ListenAndServe())
}
