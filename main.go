// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample helloworld is a basic App Engine flexible app.
package main

import (
	"log"
	"net/http"
	"healthcheck"
	"api"
	"web"
	"remote"
)

func main() {
	http.HandleFunc("/_ah/health", healthcheck.HealthCheckHandler)

	web.AddHandlers()

	httpClient := http.Client{}
	//trApi := remote.NewLondonTransportAPI(&httpClient)
	trApi := remote.NewBudapestTransportAPI(&httpClient)
	backend := api.NewBackendApi(&trApi)
	backend.AddHandlers()

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
