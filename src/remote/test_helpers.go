package remote

import (
	"net/http"
	"io"
	"testing"
)


// MockTransportAPI
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


// MockHttpClient
type MockHttpClient struct {
	ReturnValueGet *http.Response
	ErrorGet error

	TrackGetCalled bool
	TrackGetUrl string
}

func (client *MockHttpClient) Get(url string) (resp *http.Response, err error) {
	client.TrackGetCalled = true
	client.TrackGetUrl = url

	return client.ReturnValueGet, client.ErrorGet
}


// fakeReadCloser
type fakeReadCloser struct {
	io.Reader
}

func (rc fakeReadCloser) Close() error { return nil }


// checkExpectedHttpCall
func checkExpectedHttpCall(t *testing.T, client MockHttpClient, expectedUrl string) {
	if !client.TrackGetCalled {
		t.Error("Method `TrackGetCalled` should've been called")
	}

	if client.TrackGetUrl != expectedUrl {
		t.Errorf("Unexpected URL called (%s != %s)", client.TrackGetUrl, expectedUrl)
	}
}
