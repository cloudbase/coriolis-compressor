// Copyright 2019 Cloudbase Solutions Srl
// All Rights Reserved.

package routers

import (
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/cloudbase/coriolis-compressor/controllers"
)

// GetRouter returns a new HTTP router
func GetRouter() *mux.Router {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/").Subrouter()
	apiRouter.Handle("/", gorillaHandlers.LoggingHandler(os.Stdout, http.HandlerFunc(controllers.CompressorHandler))).Methods("POST")

	return router
}
