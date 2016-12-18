package api

import (
	"net/http"
	"encoding/json"
	"remote"
	"strconv"
)

const baseUrl = "https://api.tfl.gov.uk/"

type http_handler func(w http.ResponseWriter, r *http.Request)
type internal_handler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

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
	http.HandleFunc("/api/v1/stops/", decorator(stopHandler))
}

func decorator(handler internal_handler) http_handler {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := handler(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		jsonString, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-type", "application/json")
		w.Write(jsonString)
	}

}

func stopHandler(w http.ResponseWriter, r *http.Request) (interface{}, error){
	stopId := r.URL.Path[len("/api/v1/stops/"):]

	if stopId == "" {
		return stopsByPosition(w, r)
	} else {
		return arrivalsByStop(w, r, stopId)
	}
}

func stopsByPosition(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil {
		return nil, err
	}

	api := remote.NewLondonTransportAPI()
	stops, err := api.ListStopPointsAround(lat, lon)

	if err != nil {
		return nil, err
	} else {
		return stops, nil
	}
}

func arrivalsByStop(w http.ResponseWriter, r *http.Request, stopId string) (interface{}, error) {
	api := remote.NewLondonTransportAPI()
	arrivals, err := api.ListArrivalsOf(stopId)

	if err != nil {
		return nil, err
	} else {
		return arrivals, nil
	}

}
