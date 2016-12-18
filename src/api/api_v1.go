package api

import (
	"net/http"
	"encoding/json"
	"strconv"
	"remote"
	"strings"
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

type BackendApi struct {
	trApi *remote.LondonTransportAPI
}

func NewBackendApi(trApi *remote.LondonTransportAPI) BackendApi {
	return BackendApi{trApi}
}

func (api BackendApi) AddHandlers() {
	http.HandleFunc("/api/v1/stops/", api.decorator(api.handler))
}

func (api BackendApi) decorator(handler internal_handler) http_handler {
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

func (api BackendApi) handler(w http.ResponseWriter, r *http.Request) (interface{}, error){
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) <= 3 {
		return api.stopsByPosition(w, r)
	} else {
		stopId := parts[3]
		return api.arrivalsByStop(w, r, stopId)
	}
}

func (api BackendApi) stopsByPosition(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil {
		return nil, err
	}

	stops, err := api.trApi.ListStopPointsAround(lat, lon)

	if err != nil {
		return nil, err
	} else {
		return stops, nil
	}
}

func (api BackendApi) arrivalsByStop(w http.ResponseWriter, r *http.Request, stopId string) (interface{}, error) {
	arrivals, err := api.trApi.ListArrivalsOf(stopId)

	if err != nil {
		return nil, err
	} else {
		return arrivals, nil
	}

}
