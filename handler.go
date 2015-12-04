package jshapi

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zenazn/goji/web"
)

// Handler is used to direct HTTP requests to resources
type Handler struct {
	resources []*resource
	Prefix    string
	Logger    *log.Logger
	Mux       *web.Mux
}

// NewHandler initializes a Handler object
func NewHandler(prefix string) *Handler {
	return &Handler{
		resources: []*resource{},
		Prefix:    prefix,
		Logger:    log.New(os.Stdout, "jshapi: ", log.Ldate|log.Ltime|log.Lshortfile),
		Mux:       web.New(),
	}
}

// Handle is an http.Handle compatible function that should be used to passed the
// jshapi Handler into a router
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.Mux.ServeHTTP(w, r)
}

// AddResource adds a new resource of type "name" to the handler's router
func (h *Handler) AddResource(name string, storage Storage) {

	resource := &resource{
		Type:    name,
		Storage: storage,
		Logger:  h.Logger,
		Prefix:  h.Prefix,
		mux:     web.New(),
	}

	plural := fmt.Sprintf("/%ss", name)
	pluralSelecter := fmt.Sprintf("%s/:id", name)

	// setup resource sub-router
	resource.mux.Post(plural, resource.Post)
	resource.mux.Delete(plural, resource.Delete)
	resource.mux.Get(plural, resource.List)
	resource.mux.Get(pluralSelecter, resource.Get)
	resource.mux.Patch(pluralSelecter, resource.Patch)

	h.resources = append(h.resources, resource)

	// Add subrouter to main Handler mux
	h.Mux.Handle(plural+"/*", resource.mux)
}
