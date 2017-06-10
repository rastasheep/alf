package main

import (
	"gopkg.in/matryer/respond.v1"
	"net/http"
)

func RespondOptions() *respond.Options {
	return &respond.Options{
		Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
			if err, ok := data.(error); ok {
				return status, map[string]interface{}{"error": err.Error()}
			}
			return status, data
		},
	}
}

func JSONResponse(options *respond.Options) Adapter {
	return options.Handler
}
