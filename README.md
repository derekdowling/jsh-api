# JSH-API

[![GoDoc](https://godoc.org/github.com/derekdowling/go-json-spec-handler?status.png)](https://godoc.org/github.com/derekdowling/jsh-api)

A [JSON Specification](http://jsonapi.org) API Build created on top of
[jsh](http://github.com/derekdowling/go-json-spec-handler). Bring your own
router, bring your own storage, focus on functionality, and let jsh-api do the
rest.

## Setup

```go
import github.com/derekdowling/jsh-api

// can specify an api prefix(/prefix/resources) or leave blank
api := jshapi.New("")

// implement the jshapi.Storage interface, then:
userStorage := &UserStorage{}
api.AddResource("user", userStorage)

yourRouter.Handle(api.Handle)
```

## What It Handles

All of the dirty work for parsing all supported JSON API request endpoints for
each resource:

```
POST /resources
GET /resources
GET /resources/:id
DELETE /resources/:id
PATCH /resources/:id
```

## Implementing a Storage Driver with jsh

Below is a simple example of how you might implement your Storage handle using
[jsh](https://godoc.org/github.com/derekdowling/go-json-spec-handler) for Save
and Patch. This should give you a pretty good idea of how easy it is to
implement the Storage driver with jsh.


```go
type User struct {
    ID string
    Name string `json:"name"`
}

func Save(object *jsh.Object) jsh.SendableError {
    user := &User{}
    err := object.Unmarshal("user", user)
    if err != nil {
        return err
    }

    // generate your id, however you choose
    user.ID = "1234"

    // do save logic
    return
}

func Patch(object *jsh.Object) jsh.SendableError {
    user := &User{}
    err := object.Unmarshal("user", user)
    if err != nil {
        return err
    }

    // object has the lookup ID
    id := object.ID

    // perform patch
    return
}
```
