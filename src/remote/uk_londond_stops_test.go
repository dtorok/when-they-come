package remote

import (
	"testing"
	"errors"
	"strings"
	"net/http"
)

var londonStopsExpectedUrl = "https://api.tfl.gov.uk/StopPoint/?lat=10.000000&lon=20.000000&radius=1000&stopTypes=NaptanBusCoachStation%2CNaptanFerryPort%2CNaptanMetroStation%2CNaptanRailStation"
var londonStopsExpectedStopIds = [...]string{"stopPointId1", "stopPointId2"}
var londonStopsReturnValue = `
		{
			"stopPoints" : [
				{
					"status" : true,
					"id" : "stopPointId1",
					"lat" : 1.1,
					"lon" : 1.2,
					"distance" : 1.3,
					"commonName" : "stopPointCommonName"
				},
				{
					"status" : true,
					"id" : "stopPointId2",
					"lat" : 2.1,
					"lon" : 2.2,
					"distance" : 2.3,
					"commonName" : "stopPointCommonName"
				}
				]
		}`


func TestLondondListStopPointsAroundFailsIfRemoteCallFails(t *testing.T) {
	errMsg := "External call failed"

	// GIVEN
	client := MockHttpClient{
		ErrorGet: errors.New(errMsg),
	}
	trApi := NewLondonTransportAPI(&client)

	// WHEN
	_, err := trApi.ListStopPointsAround(10, 20)

	// THEN
	checkExpectedHttpCall(t, client, londonStopsExpectedUrl)

	if err == nil {
		t.Error("Method `ListStopPointsAround` should've failed, but it didn't")
	}

	if strings.Index(err.Error(), errMsg) == -1 {
		t.Errorf("Error message should've contained \"%s\" (%s)", errMsg, err.Error())
	}
}

func TestLondonListStopPointsAroundCallsExternalAPI(t *testing.T) {
	// GIVEN
	client := MockHttpClient{
			ReturnValueGet: &http.Response{
			StatusCode: 200,
			Body: fakeReadCloser{strings.NewReader(londonStopsReturnValue)},
		},
	}
	trApi := NewLondonTransportAPI(&client)

	// WHEN
	stops, err := trApi.ListStopPointsAround(10, 20)

	// THEN
	checkExpectedHttpCall(t, client, londonStopsExpectedUrl)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if len(stops) != len(londonStopsExpectedStopIds) {
		t.Errorf("Unexpected number of ids received (%d != %d)", len(stops), len(londonStopsExpectedStopIds))
	}

	for i, stop := range stops {
		if stop.Id != londonStopsExpectedStopIds[i] {
			t.Errorf("Unexpected stop id received (%s != %s)", stop.Id, londonStopsExpectedStopIds[i])
		}
	}
}