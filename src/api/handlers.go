package api

import (
	"net/http"
	"remote"
)

func AddHandlers() {
	// healthcheck
	healthCheck := Healthcheck{}
	healthCheck.AddHandlers()

	// backend api
	httpClient := http.Client{}
	//trApi := remote.NewLondonTransportAPI(&httpClient)
	trApi := remote.NewBudapestTransportAPI(&httpClient)

	backend := NewBackendApi(&trApi)
	backend.AddHandlers()

}