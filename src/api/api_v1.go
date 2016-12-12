package api

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

const baseUrl = "https://api.tfl.gov.uk/"

type Stop struct {
	Id string    `json:"id"`
	Lat string   `json:"lat"`
	Lon string   `json:"lon"`
	Name string  `json:"name"`
	Distance int `json:"distance"`
}

type Arrival struct {
	LineName string  `json:"lineName"`
	Towards string   `json:"towards"`
	ArrivesIn int    `json:"arrivesIn"`
	ArrivesAt string `json:"arrivesAt"`
}

type LondonStopPointResult struct {
	StopPoints []LondonStopPoint `json:"stopPoints"`
}

type LondonStopPoint struct {
	Id         string  `json:"id"`
	CommonName string  `json:"commonName"`
	Distance   float32 `json:"distance"`
	Status     string  `json:"status"`
	Lat        string  `json:"lat"`
	Lon        string  `json:"lon"`
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
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	url := fmt.Sprintf("%s/StopPoint/?lat=%s&lon=%s&stopTypes=%s",
		baseUrl,
		lat, lon,
		"NaptanBusCoachStation,NaptanFerryPort,NaptanMetroStation,NaptanRailStation")

	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var res LondonStopPointResult

		_ = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.StopPoints[0].CommonName)

		//var sp LondonStopPoint
		var response []Stop = make([]Stop, len(res.StopPoints))
		for i, sp := range res.StopPoints {
			response[i] = Stop{
				sp.Id,
				sp.Lat, sp.Lon,
				sp.CommonName,
				int(sp.Distance),
			}
		}

		w.Header().Set("Content-type", "application/json")
		//s := []Stop{Stop{"stopid1", "123", "234", "Stop1", 234}}
		json, _ := json.Marshal(response)
		w.Write(json)
	}
}

func arrivalsByStop(w http.ResponseWriter, r *http.Request, stopId string) {
	w.Header().Set("Content-type", "application/json")
	arrivals := []Arrival{Arrival{"Bakerloo", "Almafa", 254, "2016-12-20 00:00:00"}}
	json, _ := json.Marshal(arrivals)
	w.Write(json)
}
