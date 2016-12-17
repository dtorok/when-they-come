package web

import (
	"net/http"
	"text/template"
)

func AddHandlers() {
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("web/static")))

	http.Handle("/static/", staticHandler)
	http.HandleFunc("/test", test)
	http.HandleFunc("/", main)
}

func main(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("main.html").Delims("{{{", "}}}").ParseFiles("src/web/templates/main.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, "")
}
