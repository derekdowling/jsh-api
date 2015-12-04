package jshapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/zenazn/goji/web"
)

// Resource is a handler object for dealing with CRUD API endpoints
type resource struct {
	// The singular name of the resource type ex) "user" or "post"
	Type string
	// An implemented jshapi.Storage interface
	Storage Storage
	// An implementation of Go's standard logger
	Logger *log.Logger
	// Prefix is set if the resource is not the top level of URI, "/prefix/resources
	Prefix string
	mux    *web.Mux
}

// Post => POST /resources
func (res *resource) Post(c web.C, w http.ResponseWriter, r *http.Request) {

	object, err := jsh.ParseObject(r)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	err = res.Storage.Save(object)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	res.sendAndLog(c, w, r, object)
}

// Get => GET /resources/:id
func (res *resource) Get(c web.C, w http.ResponseWriter, r *http.Request) {
	id, exists := c.URLParams["id"]
	if !exists {
		res.sendAndLog(c, w, r, jsh.ISE(fmt.Sprintf("Unable to parse resource ID from path: %s", r.URL.Path)))
		return
	}

	object, err := res.Storage.Get(id)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	res.sendAndLog(c, w, r, object)
}

// List => GET /resources
func (res *resource) List(c web.C, w http.ResponseWriter, r *http.Request) {
	list, err := res.Storage.List()
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	res.sendAndLog(c, w, r, list)
}

// Delete => DELETE /resources/:id
func (res *resource) Delete(c web.C, w http.ResponseWriter, r *http.Request) {
	id, exists := c.URLParams["id"]
	if !exists {
		res.sendAndLog(c, w, r, jsh.ISE(fmt.Sprintf("Unable to parse resource ID from path: %s", r.URL.Path)))
		return
	}

	err := res.Storage.Delete(id)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Patch => PATCH /resources/:id
func (res *resource) Patch(c web.C, w http.ResponseWriter, r *http.Request) {
	object, err := jsh.ParseObject(r)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	err = res.Storage.Patch(object)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	res.sendAndLog(c, w, r, object)
}

func (res *resource) sendAndLog(c web.C, w http.ResponseWriter, r *http.Request, sendable jsh.Sendable) {

	jshErr, isType := sendable.(*jsh.Error)
	if isType && jshErr.Status == http.StatusInternalServerError {
		res.Logger.Printf("JSH ISE: %s-%s", jshErr.Title, jshErr.Detail)
	}

	err := jsh.Send(w, r, sendable)
	if err != nil {
		res.Logger.Print(err.Error())
	}
}
