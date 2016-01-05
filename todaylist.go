package todaylist

import (
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	err := t.Execute(w, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
