package healthcheck

import (
	"net/http"
	"fmt"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

