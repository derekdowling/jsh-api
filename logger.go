package jshapi

import (
	"log"
	"os"

	"github.com/derekdowling/go-stdlogger"
)

// Logger can be overridden with your own logger to utilize any custom features
// it might have. Interface defined here: https://github.com/derekdowling/go-stdlogger/blob/master/logger.go
var Logger std.Logger = log.New(os.Stderr, "jshapi: ", log.LstdFlags)
