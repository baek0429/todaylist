package todaylist

import (
	"html/template"
	m "model"
	"net/http"

	// "golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	// "google.golang.org/appengine/datastore"
	// "google.golang.org/appengine/log"

	// for datastore
)

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
	if r.Method == "GET" {
		i := m.GetAddEmptyInput()
		//creating uploadUrl
		ctx := appengine.NewContext(r)
		uploadURL, err := blobstore.UploadURL(ctx, "/add", nil)
		if err != nil {
			return
		}
		i.BlobActionURL = uploadURL
		t.ExecuteTemplate(w, "base", i)
	}
	if r.Method == "POST" {
		// ctx := appengine.NewContext(r)
		r.ParseForm()
		blobs, _, err := blobstore.ParseUpload(r)
		if err != nil {
			panic(err)
			// serveError(ctx, w, err)
			return
		}
		file := blobs["img"]
		// title := others.Get("title")
		// password := others.Get("password")
		// description := others.Get("description")
		if len(file) == 0 {
			w.Write([]byte("no file was found"))
			return
		}
		http.Redirect(w, r, "/serve/?blobKey="+string(file[0].BlobKey), http.StatusFound)
		// http.Redirect(w, r, "/serve/?blobKey="+string(file[0].BlobKey)+"&title="+title+"&password="+password+"&description="+description, http.StatusFound)
	}
}

func addResultHandler(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
}
func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/serve/", addResultHandler)
}
