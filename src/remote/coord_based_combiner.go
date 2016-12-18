package remote

import "strings"

type CoordBasedCombinedTransportAPI struct {
	ukLondonPrefix string
	ukLondonApi LondonTransportAPI
	huBudapestPrefix string
	huBudapestApi BudapestTransportAPI
}

func NewCoordBasedCombinedTransportAPI(client HttpClient) CoordBasedCombinedTransportAPI {
	return CoordBasedCombinedTransportAPI{
		"LON", NewLondonTransportAPI(client),
		"BUD", NewBudapestTransportAPI(client),
	}
}

func (api CoordBasedCombinedTransportAPI) ListStopPointsAround(lat, lon float64) ([]Stop, error) {
	var prefix string
	var trApi TransportAPI

	if inBoundingBox(lat, lon, 49.9, -11.05, 58.7, 1.78) {
		prefix = api.ukLondonPrefix
		trApi = api.ukLondonApi
	} else {
		prefix = api.huBudapestPrefix
		trApi = api.huBudapestApi
	}

	stops, err := trApi.ListStopPointsAround(lat, lon)

	if err != nil {
		return nil, err
	}

	for i, st := range stops {
		stops[i].Id = prefix + st.Id
	}

	return stops, nil
}

func (api CoordBasedCombinedTransportAPI) ListArrivalsOf(stopPointId string) ([]Arrival, error) {
	var prefix string
	var trApi TransportAPI

	if strings.Index(stopPointId, api.ukLondonPrefix) == 0 {
		prefix = api.ukLondonPrefix
		trApi = api.ukLondonApi
	} else {
		prefix = api.huBudapestPrefix
		trApi = api.huBudapestApi
	}

	stopPointId = stopPointId[len(prefix):]

	return trApi.ListArrivalsOf(stopPointId)
}
