// Package store is a collection of composable interfaces that are can be implemented
// in order to build a storage driver
package store

import (
	"github.com/derekdowling/go-json-spec-handler"
	"golang.org/x/net/context"
)

// CRUD implements all sub-storage interfaces
type CRUD interface {
	Save
	Get
	List
	Update
	Delete
}

// Save a new resource to storage
type Save interface {
	Save(ctx context.Context, object *jsh.Object) jsh.SendableError
}

// Get a specific instance of a resource by id from storage
type Get interface {
	Get(ctx context.Context, id string) (*jsh.Object, jsh.SendableError)
}

// List all instances of a resource from storage
type List interface {
	List(ctx context.Context) (jsh.List, jsh.SendableError)
}

// Update an existing object in storage
type Update interface {
	Update(ctx context.Context, object *jsh.Object) jsh.SendableError
}

// Delete an object from storage by id
type Delete interface {
	Delete(ctx context.Context, id string) jsh.SendableError
}
