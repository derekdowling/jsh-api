// Package store is a collection of composable interfaces that are can be implemented
// in order to build a storage driver
package store

import (
	"github.com/derekdowling/go-json-spec-handler"
	"golang.org/x/net/context"
)

// CRUD implements all sub-storage functions
type CRUD interface {
	Save(ctx context.Context, object *jsh.Object) jsh.SendableError
	Get(ctx context.Context, id string) (*jsh.Object, jsh.SendableError)
	List(ctx context.Context) (jsh.List, jsh.SendableError)
	Update(ctx context.Context, object *jsh.Object) jsh.SendableError
	Delete(ctx context.Context, id string) jsh.SendableError
}

// Save a new resource to storage
type Save func(ctx context.Context, object *jsh.Object) jsh.SendableError

// Get a specific instance of a resource by id from storage
type Get func(ctx context.Context, id string) (*jsh.Object, jsh.SendableError)

// List all instances of a resource from storage
type List func(ctx context.Context) (jsh.List, jsh.SendableError)

// Update an existing object in storage
type Update func(ctx context.Context, object *jsh.Object) jsh.SendableError

// Delete an object from storage by id
type Delete func(ctx context.Context, id string) jsh.SendableError
