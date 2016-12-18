package remote

import (
	"encoding/json"
	"net/http"
	"log"
	"io/ioutil"
	"strings"
	"fmt"
	"time"
	"net/url"
)

const budapestBaseScheme = "http"
const budapestBaseHost = "futar.bkk.hu"


type BudapestTransportAPI struct {
	client *http.Client
	stops []BudapestStop
}

type BudapestStop struct {
	Id string   `json:"id"`
	Name string `json:"name"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type BudapestArrivaleStopTime struct {
	ArrivalTime int64 `json:"arrivalTime"`
	TripId string `json:"tripId"`

}

type BudapestArrivalTrip struct {
	RouteId string `json:"routeId"`
	TripHeadsign string `json:"tripHeadsign"`
}

type BudapestArrivalRoute struct {
	ShortName string `json:"shortName"`
	Description string `json:"description"`
}

type BudapestArrivalReferences struct {
	Trips map[string]BudapestArrivalTrip `json:"trips"`
	Routes map[string]BudapestArrivalRoute `json:"routes"`
}

type BudapestArrivaleEntry struct {
	StopTimes []BudapestArrivaleStopTime `json:"stopTimes"`
}

type BudapestArrivalData struct {
	Entry BudapestArrivaleEntry `json:"entry"`
	References BudapestArrivalReferences `json:"references"`
}

type BudapestArrivalResponse struct {
	Data BudapestArrivalData
}

func (api BudapestTransportAPI) getCall(url string, res interface{}) error {
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

func NewBudapestTransportAPI(client *http.Client) BudapestTransportAPI {
	data, err := ioutil.ReadFile("src/remote/hu_budapest_stops.json")
	if err != nil {
		log.Fatal("Couldn't open src/remote/hu_budapest_stops.json file")
	}

	var stops []BudapestStop

	err = json.Unmarshal(data, &stops)
	if err != nil {
		log.Fatal("Couldn't parse src/remote/hu_budapest_stops.json file")
	}

	return BudapestTransportAPI{client, stops}
}

func (api BudapestTransportAPI) ListStopPointsAround(lat, lon float64) ([]Stop, error) {
	latFrom := lat - 0.01
	latTo := lat + 0.01
	lonFrom := lon - 0.01
	lonTo := lon + 0.01

	var result []Stop = make([]Stop, 200)
	var cnt = 0
	for _, stop := range api.stops {
		if stop.Lat > latFrom && stop.Lat < latTo && stop.Lon > lonFrom && stop.Lon < lonTo {
			var id = stop.Id
			if strings.Index(id, "F") == 0 {
				id = id[1:]
			}

			result[cnt] = Stop{
				stop.Id,
				stop.Lat,
				stop.Lon,
				stop.Name,
				0,
			}
			cnt += 1
		}
	}

	return result[:cnt], nil
}

func (api BudapestTransportAPI) ListArrivalsOf(stopPointId string) ([]Arrival, error) {
	query := url.Values{}
	query.Set("stopId", "BKK_" + stopPointId)
	query.Set("onlyDepartures", "1")
	query.Set("minutesBefore", "0")
	query.Set("minutesAfter", "40")

	trUrl := url.URL{
		Scheme: budapestBaseScheme,
		Host: budapestBaseHost,
		Path: "/bkk-utvonaltervezo-api/ws/otp/api/where/arrivals-and-departures-for-stop.json",
		RawQuery: query.Encode(),
	}

	var res BudapestArrivalResponse
	var now = time.Now().Unix()

	err := api.getCall(trUrl.String(), &res)
	if err != nil {
		return nil, err
	}

	var arrivals []Arrival = make([]Arrival, len(res.Data.Entry.StopTimes))

	cnt := 0
	for _, stopTime := range res.Data.Entry.StopTimes {
		trip, found := res.Data.References.Trips[stopTime.TripId]
		if !found {
			continue
		}

		route, found := res.Data.References.Routes[trip.RouteId]
		if !found {
			continue
		}

		arrivals[cnt] = Arrival{
			route.ShortName,
			trip.TripHeadsign,
			int(stopTime.ArrivalTime - now),
			"",
		}

		cnt += 1
	}

	return arrivals[:cnt], nil
}
