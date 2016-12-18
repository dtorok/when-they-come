package api

import (
	"net/http"
	"encoding/json"
	"strconv"
	"remote"
	"strings"
	"fmt"
	"errors"
)

const baseUrl = "https://api.tfl.gov.uk/"

type http_handler func(w http.ResponseWriter, r *http.Request)
type internal_handler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

type BackendApi struct {
	trApi remote.TransportAPI
}

func NewBackendApi(trApi remote.TransportAPI) BackendApi {
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

func (api BackendApi) parseFloatFromQuery(r *http.Request, name string) (float64, error) {
	sVal := r.URL.Query().Get(name)
	if sVal == "" {
		return 0, errors.New(fmt.Sprintf("Parameter `%s` not found or empty", name))
	}

	val, err := strconv.ParseFloat(sVal, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Parameter `%s` is not a proper float (%s)", name, err))
	}

	return val, nil
}

func (api *BackendApi) stopsByPosition(w http.ResponseWriter, r *http.Request) ([]remote.Stop, error) {
	lat, err := api.parseFloatFromQuery(r, "lat")
	if err != nil {
		return nil, err
	}

	lon, err := api.parseFloatFromQuery(r, "lon")
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

func (api BackendApi) arrivalsByStop(w http.ResponseWriter, r *http.Request, stopId string) ([]remote.Arrival, error) {
	arrivals, err := api.trApi.ListArrivalsOf(stopId)

	if err != nil {
		return nil, err
	} else {
		return arrivals, nil
	}

}
