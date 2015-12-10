package jshapi

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"goji.io"

	"goji.io/pat"

	"golang.org/x/net/context"

	"github.com/derekdowling/go-json-spec-handler"
)

// Resource is a handler for a specific API endpoint type ("users", "posts", etc)
// which wraps Goji's(https://goji.io/) Mux and automatically populates a basic
// set of routes that represent a CRUD API:
//
// GET    /resource
// POST   /resource
// GET    /resource/:id
// DELETE /resource/:id
// PATCH  /resource/:id
//
// You can also add your own custom routes as well using Goji's Mux API:
//
//	func hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//		name := pat.Param(ctx, "name")
//		fmt.Fprintf(w, "Hello, %s!", name)
//	}
//
//	resource := jshapi.NewResource("api", "user", userStorage)
//	resource.Mux.HandleC(pat.New("/users/search/:name"), httpNameSearchHandler)
//
// Or add nested resources:
//
//	commentResource := jshapi.NewResource("/posts/:id/", "comment", commentStorage)
//	resource.Handle("/posts/:id/*", commentStorage)
//
type Resource struct {
	*goji.Mux
	// The singular name of the resource type("user", "post", etc)
	Name string
	// An implemented jshapi.Storage interface
	Storage Storage
	// An implementation of Go's standard logger
	Logger *log.Logger
	// Prefix is set if the resource is not the top level of URI, "/prefix/resources
	Prefix string
}

// NewResource is a resource constructor
func NewResource(prefix string, name string, storage Storage) *Resource {

	r := &Resource{
		Mux:     goji.NewMux(),
		Name:    name,
		Storage: storage,
		Prefix:  prefix,
	}

	// setup resource sub-router
	r.Mux.HandleFuncC(pat.Post(r.Matcher()), r.Post)
	r.Mux.HandleFuncC(pat.Get(r.IDMatcher()), r.Get)
	r.Mux.HandleFuncC(pat.Get(r.Matcher()), r.List)
	r.Mux.HandleFuncC(pat.Delete(r.IDMatcher()), r.Delete)
	r.Mux.HandleFuncC(pat.Patch(r.IDMatcher()), r.Patch)

	return r
}

// Post => POST /resources
func (res *Resource) Post(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	object, err := jsh.ParseObject(r)
	if err != nil {
		res.sendAndLog(ctx, w, r, err)
		return
	}

	err = res.Storage.Save(ctx, object)
	if err != nil {
		res.sendAndLog(ctx, w, r, err)
		return
	}

	res.sendAndLog(ctx, w, r, object)
}

// Get => GET /resources/:id
func (res *Resource) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := pat.Param(ctx, "id")

	object, err := res.Storage.Get(ctx, id)
	if err != nil {
		res.sendAndLog(ctx, w, r, err)
		return
	}

	res.sendAndLog(ctx, w, r, object)
}

// List => GET /resources
func (res *Resource) List(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	list, err := res.Storage.List(ctx)
	if err != nil {
		res.sendAndLog(ctx, w, r, err)
		return
	}

	res.sendAndLog(ctx, w, r, list)
}

// Delete => DELETE /resources/:id
func (res *Resource) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := pat.Param(ctx, "id")

	err := res.Storage.Delete(ctx, id)
	if err != nil {
		res.sendAndLog(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Patch => PATCH /resources/:id
func (res *Resource) Patch(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	object, err := jsh.ParseObject(r)
	if err != nil {
		res.sendAndLog(ctx, w, r, err)
		return
	}

	err = res.Storage.Patch(ctx, object)
	if err != nil {
		res.sendAndLog(ctx, w, r, err)
		return
	}

	res.sendAndLog(ctx, w, r, object)
}

func (res *Resource) sendAndLog(ctx context.Context, w http.ResponseWriter, r *http.Request, sendable jsh.Sendable) {
	response, err := sendable.Prepare(r, true)
	if err != nil && response.HTTPStatus == http.StatusInternalServerError {
		res.Logger.Printf("Error: %s", err.Internal())
	}

	sendErr := jsh.SendResponse(w, r, response)
	if sendErr != nil {
		res.Logger.Print(err.Error())
	}
}

// PluralType returns the resource's name, but pluralized
func (res *Resource) PluralType() string {
	return res.Name + "s"
}

// IDMatcher returns a uri path matcher for the resource type
func (res *Resource) IDMatcher() string {
	return path.Join(res.Matcher(), ":id")
}

// Matcher returns the top level uri path matcher for the resource type
func (res *Resource) Matcher() string {
	return fmt.Sprintf("/%s", path.Join(res.Prefix, res.PluralType()))
}
