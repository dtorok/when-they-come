package api

import (
	"net/http"
	"encoding/json"
	"remote"
	"strconv"
)

const baseUrl = "https://api.tfl.gov.uk/"

type Stop struct {
	Id string    `json:"id"`
	Lat float64  `json:"lat"`
	Lon float64  `json:"lon"`
	Name string  `json:"name"`
	Distance int `json:"distance"`
}

type Arrival struct {
	LineName string  `json:"lineName"`
	Towards string   `json:"towards"`
	ArrivesIn int    `json:"arrivesIn"`
	ArrivesAt string `json:"arrivesAt"`
}


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
	err := _stopsByPosition(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func _stopsByPosition(w http.ResponseWriter, r *http.Request) error {
	lat, err := strconv.ParseFloat(
		r.URL.Query().Get("lat"), 64)
	if err != nil {
		return err
	}

	lon, err := strconv.ParseFloat(
		r.URL.Query().Get("lon"), 64)
	if err != nil {
		return err
	}

	api := remote.NewLondonTransportAPI()
	stops, err := api.ListStopPointsAround(lat, lon)

	if err != nil {
		return err
	}

	jsonstring, err := json.Marshal(stops)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(jsonstring)

	return nil
}

func arrivalsByStop(w http.ResponseWriter, r *http.Request, stopId string) {
	api := remote.NewLondonTransportAPI()
	arrivals, err := api.ListArrivalsOf(stopId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonstring, err := json.Marshal(arrivals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(jsonstring)
}
