package remote

import (
	"testing"
	"errors"
	"strings"
	"net/http"
)

var londonArrivalsFakeStopId = "fake-stop-id"
var londonArrivalsExpectedUrl = "https://api.tfl.gov.uk/StopPoint/fake-stop-id/Arrivals/"
var londonArrivalsLineNames = [...]string{"lineName1", "lineName2"}
var londonArrivalsReturnValue = `
	[
		{
			"lineName" : "lineName1",
			"timeToStation" : 1,
			"destinationName" : "destinationName1",
			"towards" : "towards1"
		},
		{
			"lineName" : "lineName2",
			"timeToStation" : 2,
			"destinationName" : "destinationName2",
			"towards" : "towards2"
		}
	]
`

func TestLondondListArrivalsOfFailsIfRemoteFails(t *testing.T) {
	errMsg := "External call failed"

	// GIVEN
	client := MockHttpClient{
		ErrorGet: errors.New(errMsg),
	}
	trApi := NewLondonTransportAPI(&client)

	// WHEN
	_, err := trApi.ListArrivalsOf(londonArrivalsFakeStopId)

	// THEN
	checkExpectedHttpCall(t, client, londonArrivalsExpectedUrl)

	if err == nil {
		t.Error("Method `ListStopPointsAround` should've failed, but it didn't")
	}

	if strings.Index(err.Error(), errMsg) == -1 {
		t.Errorf("Error message should've contained \"%s\" (%s)", errMsg, err.Error())
	}
}

func TestLondondListArrivalsOfCallsRemote(t *testing.T) {
	// GIVEN
	client := MockHttpClient{
		ReturnValueGet: &http.Response{
			StatusCode: 200,
			Body: fakeReadCloser{strings.NewReader(londonArrivalsReturnValue)},
		},
	}
	trApi := NewLondonTransportAPI(&client)

	// WHEN
	arrivals, err := trApi.ListArrivalsOf(londonArrivalsFakeStopId)

	// THEN
	checkExpectedHttpCall(t, client, londonArrivalsExpectedUrl)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if len(arrivals) != len(londonArrivalsLineNames) {
		t.Errorf("Unexpected number of arrivals received (%d != %d)", len(arrivals), len(londonArrivalsLineNames))
	}

	for i, arrival := range arrivals {
		if arrival.LineName != londonArrivalsLineNames[i] {
			t.Errorf("Unexpected arrival lineName received (%s != %s)", arrival.LineName, londonArrivalsLineNames[i])
		}
	}
}
