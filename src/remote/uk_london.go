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

type LondonTransportAPI struct {
	client *http.Client
}

func NewLondonTransportAPI(client *http.Client) LondonTransportAPI {
	return LondonTransportAPI{client}
}

func (api LondonTransportAPI) getCall(url string, res interface{}) error {
	fmt.Println(url)
	resp, err := api.client.Get(url)
	if err != nil {
		return err
	} else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &res)
		if err != nil {
			return err
		}

		return nil
	}
}

func (api LondonTransportAPI) ListStopPointsAround(lat, lon float64) ([]Stop, error) {
	url := fmt.Sprintf("%s/StopPoint/?lat=%f&lon=%f&radius=%d&stopTypes=%s",
		baseUrl,
		lat, lon,
		1000,
		"NaptanBusCoachStation,NaptanFerryPort,NaptanMetroStation,NaptanRailStation")

	var res LondonStopPointResult

	err := api.getCall(url, &res)

	if err != nil {
		return nil, err
	}

	stops := res.StopPoints

	var response = make([]Stop, len(stops))
	for i, sp := range stops {
		response[i] = Stop{
			sp.Id,
			sp.Lat, sp.Lon,
			sp.CommonName,
			int(sp.Distance),
		}
	}

	return response, nil
}

func (api LondonTransportAPI) ListArrivalsOf(stopPointId string) ([]Arrival, error) {
	url := fmt.Sprintf("%s/StopPoint/%s/Arrivals/",
		baseUrl,
		stopPointId)

	var res []LondonArrival

	err := api.getCall(url, &res)

	if err != nil {
		return nil, err
	}

	var response []Arrival = make([]Arrival, len(res))
	for i, arr := range res {
		var towards = arr.Towards
		if towards == "" {
			towards = arr.DestinationName
		}

		response[i] = Arrival{
			arr.LineName,
			towards,
			arr.TimeToStation,
			arr.ExpectedArrival,
		}
	}

	return response, nil
}