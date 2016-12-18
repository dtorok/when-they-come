package api

import (
	"remote"
	"net/http/httptest"
	"testing"
	"errors"
	"strings"
)

const fake_stop_id = "fake-stop-id"

func testArrivalsByStop(trApi *remote.MockTransportAPI) ([]remote.Arrival, error) {
	if trApi == nil {
		trApi = &remote.MockTransportAPI{}
	}

	bApi := NewBackendApi(trApi)

	r := httptest.NewRequest("GET", "/api/v1/stops/" + fake_stop_id, nil)
	w := httptest.NewRecorder()

	return bApi.arrivalsByStop(w, r, fake_stop_id)
}

func TestArrivalByStopFailsIfRemoteFails(t *testing.T) {
	// GIVEN
	errMsg := "Remote call failure"

	trApi := remote.MockTransportAPI{
		ErrorListArrivalsOf: errors.New(errMsg),
	}

	// WHEN
	_, err := testArrivalsByStop(&trApi)

	// THEN
	if err == nil {
		t.Error("Error shouldn't have been nil")
	}

	if !trApi.TrackListArrivalsOfCalled {
		t.Error("Method `remote.ListArrivalsOf` wasn't called, although it should've been")
	}
	if strings.Index(err.Error(), errMsg) == -1 {
		t.Errorf("Error message should've passed the message \"%s\" (%s)", errMsg, err.Error())
	}
}

func TestArrivalByStopCallsRemote(t *testing.T) {
	// GIVEN
	returnValue := [...]remote.Arrival{
		remote.Arrival{ArrivesIn: 1},
		remote.Arrival{ArrivesIn: 2},
		remote.Arrival{ArrivesIn: 3},
	}

	trApi := remote.MockTransportAPI{
		ReturnValueListArrivalsOf: returnValue[:],
	}

	// WHEN
	arrivals, err := testArrivalsByStop(&trApi)

	// THEN
	if err != nil {
		t.Errorf("Error should've been nil (%s)", err.Error())
	}

	if !trApi.TrackListArrivalsOfCalled {
		t.Error("Method `remote.ListArrivalsOf` wasn't called, although it should've been")
	}

	if trApi.TrackListArrivalsOfStopPointId != fake_stop_id {
		t.Errorf("Method `remote.ListArrivalsOf` should've been called with %s (%s)", fake_stop_id, trApi.TrackListArrivalsOfStopPointId)
	}

	if len(arrivals) != len(returnValue) {
		t.Errorf("Different amount of data in result (%d) than expected (%d)", len(arrivals), len(returnValue))
	}

	for i, arrival := range arrivals {
		expectedArrival := returnValue[i]
		if arrival != expectedArrival {
			t.Errorf("Unxpected arrival in result: %s != %s", arrival, expectedArrival)
		}
	}
}