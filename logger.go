package jshapi

import (
	"log"
	"net/http"
	"time"

	"github.com/zenazn/goji/web/mutil"

	"golang.org/x/net/context"

	"goji.io"
)

// Logger describes a logger interface that is compatible with the standard
// log.Logger but also logrus and others. As not to limit which loggers can and
// can't be used with the API.
//
// This interface is from https://godoc.org/github.com/Sirupsen/logrus#StdLogger
type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

// logMiddleware logs inbound requests and outbound responses from the API
func (a *API) logMiddleware(next goji.Handler) goji.Handler {
	logger := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

		// print basic method info
		log.Printf("Serving %s %q from %s", r.Method, r.URL.String(), r.RemoteAddr)
		startTime := time.Now()

		// use WrapWriter so we can peek at response
		lw := mutil.WrapWriter(w)

		next.ServeHTTPC(ctx, lw, r)

		// no status
		if lw.Status() == 0 {
			lw.WriteHeader(http.StatusOK)
		}

		stopTime := time.Now()
		deltaTime := stopTime.Sub(startTime)

		log.Printf("Returning HTTP %03d after %s", lw.Status(), deltaTime.String())
	}

	return goji.HandlerFunc(logger)
}
