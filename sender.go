package jshapi

import (
	"net/http"

	"github.com/derekdowling/go-json-spec-handler"
	"golang.org/x/net/context"
)

/*
Sender is a function type definition that allows consumers to customize how they
send and log API responses.
*/
type Sender func(context.Context, http.ResponseWriter, *http.Request, jsh.Sendable)

/*
SendHandler is what JSHAPI uses to send all API responses. This can be overidden to
provide custom send or logging functionality.
*/
var SendHandler Sender = DefaultSender

/*
DefaultSender is the default sender that will log 5XX errors that it encounters
in the process of sending a response.
*/
func DefaultSender(ctx context.Context, w http.ResponseWriter, r *http.Request, sendable jsh.Sendable) {

	sendableError, isType := sendable.(jsh.ErrorType)
	if isType && sendableError.StatusCode() >= 500 {
		Logger.Printf("Returning ISE: %s\n", sendableError.Error())
	}

	sendError := jsh.Send(w, r, sendable)
	if sendError != nil && sendError.Status >= 500 {
		Logger.Printf("Error sending response: %s\n", sendError.Error())
	}
}
