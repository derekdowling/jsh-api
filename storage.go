package jshapi

import "github.com/derekdowling/go-json-spec-handler"

// Storage is an interface that allows you to manage your resource objects and is
// the only required implementation to use jshapi.
//
// See the MockStorage(https://github.com/derekdowling/jsh-api/blob/master/test_util.go#L13)
// object for a very basic sample implementation that is used for testing jshapi.
type Storage interface {
	// Save a new resource to storage
	Save(object *jsh.Object) jsh.SendableError
	// Get a specific instance of a resource from storage
	Get(id string) (*jsh.Object, jsh.SendableError)
	// List all instances of a resource from storage
	List() (jsh.List, jsh.SendableError)
	// Save an object to storage
	Patch(object *jsh.Object) jsh.SendableError
	// Delete from storage
	Delete(id string) jsh.SendableError
}
