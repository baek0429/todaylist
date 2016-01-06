package todaylist

import (
	"html/template"
	m "model"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/add", addHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*"))
	t.ParseFiles("index.html")
	i := m.GetMainInput()
	err := t.ExecuteTemplate(w, "base", i)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*"))
	t.ParseFiles("main.html")
	err := t.Execute(w, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*"))
	t.ParseFiles("add.html")
	i := m.GetAddEmptyInput()
	t.ExecuteTemplate(w, "base", i)
	// do something with t
}
