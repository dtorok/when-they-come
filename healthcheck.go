package main

import (
	"net/http"
	"fmt"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

