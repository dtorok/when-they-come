package web

import (
	"net/http"
	"text/template"
)

func AddHandlers() {
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("src/web/static")))

	http.Handle("/static/", staticHandler)
	http.HandleFunc("/test", test)
	http.HandleFunc("/", main)
}

func test(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("test.html").Delims("{{{", "}}}").ParseFiles("src/web/templates/test.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, "")
}

func main(w http.ResponseWriter, r *http.Request) {
	data := struct {
		GoogleAPIKey string
	}{
		GoogleAPIKey: "SECRET",
	}
	t, err := template.New("main.html").ParseFiles("src/web/templates/main.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
