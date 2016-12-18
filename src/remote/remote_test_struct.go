package remote

type MockTransportAPI struct {
	ReturnValueListStopPointsAround []Stop
	ErrorListStopPointsAround error
	TrackListStopPointsAroundCalled bool
	TrackListStopPointsAroundLat float64
	TrackListStopPointsAroundLon float64

	ReturnValueListArrivalsOf []Arrival
	ErrorListArrivalsOf error
	TrackListArrivalsOfCalled bool
	TrackListArrivalsOfStopPointId string
}

func (api *MockTransportAPI) ListStopPointsAround(lat, lon float64) ([]Stop, error) {
	api.TrackListStopPointsAroundCalled = true
	api.TrackListStopPointsAroundLat = lat
	api.TrackListStopPointsAroundLon = lon

	return api.ReturnValueListStopPointsAround, api.ErrorListStopPointsAround
}

func (api *MockTransportAPI) ListArrivalsOf(stopPointId string) ([]Arrival, error) {
	api.TrackListArrivalsOfCalled = true
	api.TrackListArrivalsOfStopPointId = stopPointId

	return api.ReturnValueListArrivalsOf, api.ErrorListArrivalsOf
}

