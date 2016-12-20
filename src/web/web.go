package web

import (
	"net/http"
	"text/template"
	"io/ioutil"
)

func AddHandlers() {
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("src/web/static")))

	http.Handle("/static/", staticHandler)
	http.HandleFunc("/", main)
}

func loadApiKey() (string, error) {
	data, err := ioutil.ReadFile("etc/google_api_key")
	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

func main(w http.ResponseWriter, r *http.Request) {
	google_api_key, err := loadApiKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		GoogleAPIKey string
	}{
		GoogleAPIKey: google_api_key,
	}
	t, err := template.New("main.html").ParseFiles("src/web/templates/main.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
