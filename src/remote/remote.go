package remote

import "net/http"

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

type TransportAPI interface {
	ListStopPointsAround(lat, lon float64) ([]Stop, error)
	ListArrivalsOf(stopPointId string) ([]Arrival, error)
}

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}