package api

import (
	"net/http"
	"fmt"
)

func AddHandlers() {
	http.HandleFunc("/api/v1/stops/", stopsByPosition)
	http.HandleFunc("/api/v1/lines/", linesByStop)
	http.HandleFunc("/api/v1/vehicles/", vehiclesByLineAndStop)
}

func stopsByPosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, `[{"id": "stopid1", "coords": [123, 234], "name": "Stop1", "distance": 234}]`)
}

func linesByStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, `[{"id": "lineid1", "name": "Line1", "type": "tube"}]`)
}

func vehiclesByLineAndStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, `[{"id": "vehicle1", "name": "Vehicle1", "arrives_at": "2016-12-20 00:00:00"}]`)
}
