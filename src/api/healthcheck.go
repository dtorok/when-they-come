package api

import (
	"net/http"
	"fmt"
)

type Healthcheck struct {}


func (hc Healthcheck) AddHandlers() {
	http.HandleFunc("/_ah/health", hc.handler)
}

func (hc Healthcheck) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

