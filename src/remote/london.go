package remote

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const baseUrl = "https://api.tfl.gov.uk/"

type LondonStopPointResult struct {
	StopPoints []LondonStopPoint `json:"stopPoints"`
}

type LondonStopPoint struct {
	Id         string  `json:"id"`
	CommonName string  `json:"commonName"`
	Distance   float32 `json:"distance"`
	Status     bool    `json:"status"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
}

type LondonArrival struct {
	Towards         string `json:"towards"`
	LineName        string `json:"lineName"`
	TimeToStation   int    `json:"timeToStation"`
	ExpectedArrival string `json:"expectedArrival"`
	DestinationName string `json:"destinationName"`
}

func getCall(url string, res interface{}) error {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	} else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body[:]))
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &res)
		//fmt.Println(res)
		//fmt.Println(err)
		if err != nil {
			return err
		}

		return nil
	}
}

func LondonListStopPoints(lat, lon float64) ([]LondonStopPoint, error) {
	url := fmt.Sprintf("%s/StopPoint/?lat=%f&lon=%f&radius=%d&stopTypes=%s",
		baseUrl,
		lat, lon,
		1000,
		"NaptanBusCoachStation,NaptanFerryPort,NaptanMetroStation,NaptanRailStation")

	var res LondonStopPointResult

	err := getCall(url, &res)

	if err != nil {
		return nil, err
	} else {
		return res.StopPoints, nil
	}
}

func LondonArrivals(stopPointId string) ([]LondonArrival, error) {
	url := fmt.Sprintf("%s/StopPoint/%s/Arrivals/",
		baseUrl,
		stopPointId)

	var res []LondonArrival

	err := getCall(url, &res)

	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}