package jshapi

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"goji.io/pat"

	"goji.io"

	"golang.org/x/net/context"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/jsh-api/store"
)

const (
	post   = "POST"
	get    = "GET"
	list   = "LIST"
	delete = "DELETE"
	patch  = "PATCH"
)

// Resource holds the necessary state for creating a REST API endpoint for a
// given resource type. Will be accessible via `/[prefix/]<type>s` where the
// proceeding `prefix/` is only precent if it is not empty.
//
// Using NewCRUDResource you can generate a generic CRUD handler for a
// JSON Specification Resource end point. If you wish to only implement a subset
// of these endpoints that is also available through NewResource() and manually
// registering storage handlers via .Post(), .Get(), .List(), .Patch(), and .Delete():
//
// You can add your own routes using the goji.Mux API:
//
//	func searchHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//		name := pat.Param(ctx, "name")
//		fmt.Fprintf(w, "Hello, %s!", name)
//	}
//
//	resource := jshapi.NewCRUDResource("user", userStorage)
//	// creates /users/search/:name
//	resource.HandleC(pat.New(resource.Matcher()+"/search/:name"), searchHandler)
//
// Or add a nested resources:
//
//	commentResource := resource.NewSubResource("post")
//	// creates /users/:id/posts* routes
//	resource.CRUD(postStorage)
type Resource struct {
	*goji.Mux
	// The singular name of the resource type("user", "post", etc)
	Type string
	// An implementation of Go's standard logger
	Logger *log.Logger
	// Prefix is set if the resource is not the top level of URI, "/prefix/resources
	prefix string
}

// NewResource is a resource constructor that makes no assumptions about routes
// that you'd like to implement, but still provides some basic utilities for
// managing routes and handling API calls.
func NewResource(resourceType string) *Resource {
	return &Resource{
		Mux:  goji.NewMux(),
		Type: resourceType,
	}
}

// NewCRUDResource generates a resource
func NewCRUDResource(resourceType string, storage store.CRUD) *Resource {
	resource := NewResource(resourceType)
	resource.CRUD(storage)
	return resource
}

// CRUD is syntactic sugar and a shortcut for registering all JSON API CRUD
// routes for a compatible storage implementation:
//
// Registers handlers for:
//	GET    /[prefix/]types
//	POST   /[prefix/]types
//	GET    /[prefix/]types/:id
//	DELETE /[prefix/]types/:id
//	PATCH  /[prefix/]types/:id
func (res *Resource) CRUD(storage store.CRUD) {
	// add all JSON API CRUD routes
	res.Get(storage)
	res.Patch(storage)
	res.Post(storage)
	res.List(storage)
	res.Delete(storage)
}

// Post registers a `POST /resources` handler with the resource
func (res *Resource) Post(storage store.Save) {
	res.HandleFuncC(
		pat.Post(res.Matcher()),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.postHandler(ctx, w, r, storage)
		},
	)
}

// Get registers a `GET /resources/:id` handler for the resource
func (res *Resource) Get(storage store.Get) {
	res.HandleFuncC(
		pat.Get(res.IDMatcher()),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.getHandler(ctx, w, r, storage)
		},
	)
}

// List registers a `GET /resources` handler for the resource
func (res *Resource) List(storage store.List) {
	res.HandleFuncC(
		pat.Get(res.Matcher()),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.listHandler(ctx, w, r, storage)
		},
	)
}

// Delete registers a `DELETE /resources/:id` handler for the resource
func (res *Resource) Delete(storage store.Delete) {
	res.HandleFuncC(
		pat.Delete(res.IDMatcher()),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.deleteHandler(ctx, w, r, storage)
		},
	)
}

// Patch registers a `PATCH /resources/:id` handler for the resource
func (res *Resource) Patch(storage store.Update) {
	res.HandleFuncC(
		pat.Patch(res.IDMatcher()),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.patchHandler(ctx, w, r, storage)
		},
	)
}

// POST /resources
func (res *Resource) postHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.Save) {
	object, err := jsh.ParseObject(r)
	if err != nil {
		res.SendAndLog(ctx, w, r, err)
		return
	}

	err = storage.Save(ctx, object)
	if err != nil {
		res.SendAndLog(ctx, w, r, err)
		return
	}

	res.SendAndLog(ctx, w, r, object)
}

// GET /resources/:id
func (res *Resource) getHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.Get) {
	id := pat.Param(ctx, "id")

	object, err := storage.Get(ctx, id)
	if err != nil {
		res.SendAndLog(ctx, w, r, err)
		return
	}

	res.SendAndLog(ctx, w, r, object)
}

// GET /resources
func (res *Resource) listHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.List) {
	list, err := storage.List(ctx)
	if err != nil {
		res.SendAndLog(ctx, w, r, err)
		return
	}

	res.SendAndLog(ctx, w, r, list)
}

// DELETE /resources/:id
func (res *Resource) deleteHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.Delete) {
	id := pat.Param(ctx, "id")

	err := storage.Delete(ctx, id)
	if err != nil {
		res.SendAndLog(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PATCH /resources/:id
func (res *Resource) patchHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.Update) {
	object, err := jsh.ParseObject(r)
	if err != nil {
		res.SendAndLog(ctx, w, r, err)
		return
	}

	err = storage.Update(ctx, object)
	if err != nil {
		res.SendAndLog(ctx, w, r, err)
		return
	}

	res.SendAndLog(ctx, w, r, object)
}

// SendAndLog is a jsh wrapper function that handles logging 500 errors and
// ensures that any errors that leak out of JSH are also captured
func (res *Resource) SendAndLog(ctx context.Context, w http.ResponseWriter, r *http.Request, sendable jsh.Sendable) {
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
	return res.Type + "s"
}

// IDMatcher returns a uri path matcher for the resource type
func (res *Resource) IDMatcher() string {
	return path.Join(res.Matcher(), ":id")
}

// Matcher returns the top level uri path matcher for the resource type
func (res *Resource) Matcher() string {
	if res.prefix == "" {
		return fmt.Sprintf("/%s", res.PluralType())
	}

	return fmt.Sprintf("%s", path.Join(res.prefix, res.PluralType()))
}

// CreateSubResource automatically builds a resource with the proper
// prefixes to ensure that it is accessible via /[prefix/]types/:id/subtypes
// and then returns it so you can register route. You can either manually add
// individual routes like normal, or make use of:
//	subResource.CRUD(storage)
// to register the equivalent of what NewCRUDResource() gives you.
func (res *Resource) CreateSubResource(resourceType string) *Resource {
	subResource := NewResource(resourceType)
	subResource.prefix = res.IDMatcher()

	res.HandleC(pat.New(subResource.Matcher()+"*"), subResource)
	return subResource
}
