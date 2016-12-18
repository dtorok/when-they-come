package remote

import (
	"fmt"
	"net/url"
	"log"
)

const londonBaseScheme = "https"
const londonBaseHost = "api.tfl.gov.uk"


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
	client HttpClient
}

func NewLondonTransportAPI(client HttpClient) LondonTransportAPI {
	return LondonTransportAPI{client}
}

func (api LondonTransportAPI) ListStopPointsAround(lat, lon float64) ([]Stop, error) {
	log.Printf("uk_london list_stop_points_around %f %f", lat, lon)

	query := url.Values{}
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("radius", "1000")
	query.Set("stopTypes", "NaptanBusCoachStation,NaptanFerryPort,NaptanMetroStation,NaptanRailStation")

	trUrl := url.URL{
		Scheme: londonBaseScheme,
		Host: londonBaseHost,
		Path: "/StopPoint/",
		RawQuery: query.Encode(),
	}

	var res LondonStopPointResult

	err := httpJsonGET(api.client, trUrl.String(), &res)

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
	log.Printf("uk_london list_arrivals_of %s", stopPointId)

	trUrl := url.URL{
		Scheme: londonBaseScheme,
		Host: londonBaseHost,
		Path: "/StopPoint/" + stopPointId + "/Arrivals/",
	}

	var res []LondonArrival

	err := httpJsonGET(api.client, trUrl.String(), &res)

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