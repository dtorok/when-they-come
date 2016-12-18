package remote

import (
	"io/ioutil"
	"encoding/json"
	"log"
)

func httpJsonGET(client HttpClient, url string, res interface{}) error {
	log.Printf("outgoing_request %s", url)

	resp, err := client.Get(url)
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

func inBoundingBox(x, y, x1, y1, x2, y2 float64) bool {
	return x > x1 && x1 < x2 && y > y1 && y < y2
}