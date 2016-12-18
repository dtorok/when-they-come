package api

import (
	"testing"
	"net/http/httptest"
	"remote"
	"strings"
	"errors"
)
func testStopByPosition(query string, trApi *remote.MockTransportAPI) ([]remote.Stop, error) {
	if trApi == nil {
		trApi = &remote.MockTransportAPI{}
	}

	bApi := NewBackendApi(trApi)

	r := httptest.NewRequest("GET", "/api/v1/stops/" + query, nil)
	w := httptest.NewRecorder()

	return bApi.stopsByPosition(w, r)
}

func checkError(t *testing.T, err error, errorWord string) {
	if err == nil {
		t.Error("err shouldn't have been nil")
	}

	if strings.Index(err.Error(), errorWord) == -1 {
		t.Errorf("The word \"%s\" should've appeard in the error message (%s)", errorWord, err.Error())
	}
}

func TestStopsByPositionFailsWithoutLon(t *testing.T) {
	// WHEN
	_, err := testStopByPosition("?lat=10", nil)

	// THEN
	checkError(t, err, "lon")
}

func TestStopsByPositionFailsWithoutLat(t *testing.T) {
	// WHEN
	_, err := testStopByPosition("?lon=10", nil)

	// THEN
	checkError(t, err, "lat")
}

func TestStopsByPositionFailsIfRemoteFails(t *testing.T) {
	// GIVEN
	trApi := remote.MockTransportAPI{
		ErrorListStopPointsAround: errors.New("Remote call failed"),
	}

	// WHEN
	_, err := testStopByPosition("?lon=10&lat=20", &trApi)

	// THEN
	checkError(t, err, "Remote call failed")
}

func TestStopsByPositionCallsRemote(t *testing.T) {
	// GIVEN
	returnValue := [...]remote.Stop{
		remote.Stop{Id: "1"},
		remote.Stop{Id: "2"},
		remote.Stop{Id: "3"},
	}

	trApi := remote.MockTransportAPI{
		ReturnValueListStopPointsAround: returnValue[:],
	}

	// WHEN
	stops, err := testStopByPosition("?lon=10&lat=20", &trApi)

	// THEN
	if err != nil {
		t.Errorf("Err shoul've been nil (%s)", err.Error())
	}

	if trApi.TrackListStopPointsAroundCalled == false {
		t.Error("Remote should've been called")
	}

	if trApi.TrackListStopPointsAroundLat != 20.0 {
		t.Errorf("Remote should've been called with lat=10 (%d)", trApi.TrackListStopPointsAroundLat)
	}

	if trApi.TrackListStopPointsAroundLon != 10.0 {
		t.Errorf("Remote should've been called with lon=20 (%d)", trApi.TrackListStopPointsAroundLon)
	}

	if len(stops) != len(returnValue) {
		t.Errorf("Different amount of data in result (%d) than expected (%d)", len(stops), len(returnValue))
	}

	for i, stop := range stops {
		expected := returnValue[i]
		if stop != expected {
			t.Errorf("Unxpected stop in result: %s != %s", stop, expected)
		}
	}
}
