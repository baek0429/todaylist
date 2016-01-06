package todaylist

import (
	"encoding/json"
	"html/template"
	m "model"
	"net/http"
)

var t = template.Must(template.ParseGlob("template/*"))

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/add", addHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t.ParseFiles("index.html")
	i := m.GetMainInput()
	err := t.ExecuteTemplate(w, "index", i)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	t.ParseFiles("main.html")
	t, _ := template.ParseFiles("main.html")
	err := t.Execute(w, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if "" == r.URL.Query().Get("key") {
		t.ParseFiles("add.html")
		i := m.GetAddEmptyInput()
		err := t.ExecuteTemplate(w, "add", i)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	} else {
		decoder := json.NewDecoder(r.Body)
		var t test_struct
		err := decoder.Decode(&t)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}
	// do something with t
}

type test_struct struct {
}
