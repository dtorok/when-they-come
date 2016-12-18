package remote

import (
	"encoding/json"
	"net/http"
	"log"
	"io/ioutil"
)


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
	return nil, nil
}
