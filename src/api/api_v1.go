package api

import (
	"net/http"
	"fmt"
)

func AddHandlers() {
	http.HandleFunc("/api/v1/stops/", stopHandler)
}


func stopHandler(w http.ResponseWriter, r *http.Request) {
	stopId := r.URL.Path[len("/api/v1/stops/"):]

	if stopId == "" {
		stopsByPosition(w, r)
	} else {
		arrivalsByStop(w, r, stopId)
	}
}

func stopsByPosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, `[{"id": "stopid1", "coords": [123, 234], "name": "Stop1", "distance": 234}]`)
}

func arrivalsByStop(w http.ResponseWriter, r *http.Request, stopId string) {
	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, `[{"id": "vehicle1", "lineName": "Bakerloo", "towards": "Almafa", "arrivesIn": 254, "arrivesAt": "2016-12-20 00:00:00"}]`)
}
