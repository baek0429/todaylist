package todaylist

import (
	"fmt"
	"html/template"
	m "model"
	"net/http"
	"time"

	//using blobstore to save image and save the blobkey to datastore.
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/mail"
	"strconv"
)

var initDBSet bool

func handler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*")) // add sub-templates in /template
	t.ParseFiles("index.html")                           // parse main.html as main
	ctx := appengine.NewContext(r)
	keys, err := m.ParseAll(ctx, m.Post{}) // get all posts in db.
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	dataModels, err := m.ParseEntitiesFromKeys(ctx, keys) // get DataModel Entities
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	var postModels []m.Post
	for i, model := range dataModels { // convert to ViewModel
		p := model.(m.Post)
		p.ID = keys[i].IntID()
		postModels = append(postModels, p)
	}
	err = t.ExecuteTemplate(w, "base", postModels) // exec templates
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func superHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*")) // add sub-templates in /template
	t.ParseFiles("super.html")                           //parse super
	ctx := appengine.NewContext(r)

	cKeys, err := m.ParseAll(ctx, m.Category{}) // get category keys
	if err != nil {
		panic(err)
	}
	cModels, err := m.ParseEntitiesFromKeys(ctx, cKeys) // get category model as datamodel interface
	if err != nil {
		panic(err)
	}
	var Cvms []m.Category // convert from interface to struct
	for _, model := range cModels {
		Cvms = append(Cvms, model.(m.Category))
	}

	lKeys, err := m.ParseAll(ctx, m.Location{}) // same procedure
	if err != nil {
		panic(err)
	}
	lModels, err := m.ParseEntitiesFromKeys(ctx, lKeys)
	if err != nil {
		panic(err)
	}
	var Lvms []m.Location
	for _, model := range lModels {
		Lvms = append(Lvms, model.(m.Location))
	}
	err = t.ExecuteTemplate(w, "base",
		struct {
			CVMS []m.Category
			LVMS []m.Location
		}{Cvms, Lvms})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*")) // add sub-templates in /template
	t.ParseFiles("main.html")                            //parse main
	r.ParseForm()
	ctx := appengine.NewContext(r)
	id := r.FormValue("id") // get post by id
	key := m.ParseKeyFromID(ctx, id, m.Post{})
	if key != nil {
		model, err := m.ParseEntityFromKey(ctx, key)
		err = t.ExecuteTemplate(w, "base", []m.Post{model.(m.Post)})
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	} else {
		w.Write([]byte("No items under the id" + id))
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*")) // add sub-templates in /template
	if r.Method == "GET" {                               // get method display form
		t.ParseFiles("add.html")                                // parse add.html if 'method' = 'get'
		ctx := appengine.NewContext(r)                          // get ctx
		uploadURL, err := blobstore.UploadURL(ctx, "/add", nil) // get uploadurl for blobstore
		if err != nil {
			return
		}
		actionURL := map[string]string{"BlobActionURL": uploadURL.String()} // provide blobuploadurl to action field of the form
		t.ExecuteTemplate(w, "base", actionURL)
	}
	if r.Method == "POST" {
		ctx := appengine.NewContext(r)
		r.ParseForm()
		blobs, others, err := blobstore.ParseUpload(r) // get upload blob info
		if err != nil {
			w.Write([]byte(err.Error())) // error
			return
		}
		files := blobs["img"]                    // name="img" in the html form
		title := others.Get("title")             // '' title
		password := others.Get("password")       // '' password
		description := others.Get("description") // '' description

		var blobKeys []string
		var imgSrcs []string
		for _, file := range files {
			imgSrc := "/serve/?blobKey=" + string(file.BlobKey) // create imgsrc url from blobkey
			imgSrcs = append(imgSrcs, imgSrc)                   //multiple images in singe post
			blobKeys = append(blobKeys, string(file.BlobKey))   // also save blobkey in case for use.
		}

		var post m.Post // creating post and fill the fields.
		post.Title = title
		post.Password = password
		post.Description = description
		post.BlobKeys = blobKeys
		post.Time = time.Now()
		post.ImageSrc = imgSrcs

		posts := []m.DataModel{post}            // uploading posts
		keys, _ := m.SaveDataModels(ctx, posts) // success?
		i := strconv.Itoa(int(keys[0].IntID()))
		http.Redirect(w, r, "/main/?id="+i, http.StatusFound) // redirect to /main/?uuid= with uuid
		// http.Redirect(w, r, "/error", http.StatusNotFound)
	}
}
func contactHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*")) // add sub-templates in /template
	if r.Method == "GET" {
		t.ParseFiles("contact.html")
		err := t.ExecuteTemplate(w, "base", nil)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}
	if r.Method == "POST" {
		r.ParseForm()
		ctx := appengine.NewContext(r)
		addr := r.FormValue("email")
		msg := &mail.Message{
			Sender:  "Example.com Support <" + addr + ">",
			To:      []string{"csbaek0429@gmail.com"},
			Subject: "Confirm your registration",
			Body:    fmt.Sprintf(r.FormValue("message"), nil),
		}
		if err := mail.Send(ctx, msg); err != nil {
			fmt.Println(err.Error())
		}
	}
}
func adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("admin.html")
		t.Execute(w, nil)
	}
	if r.Method == "POST" {
		ctx := appengine.NewContext(r)
		r.ParseForm()
		categoryTitle := r.FormValue("newCategory")
		categoryParentTitle := r.FormValue("categoryParent")
		locationTitle := r.FormValue("newLocation")
		locationParentTitle := r.FormValue("locationParent")

		// save category
		key, err := m.SaveIfTitleNoneExists(ctx, m.Category{Title: categoryParentTitle})
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		c := m.Category{Title: categoryTitle}
		if key != nil {
			c = c.SetChildDSID([]int64{key.IntID()}).(m.Category)
			m.SaveIfTitleNoneExists(ctx, c)
		}

		// save location
		key, err = m.SaveIfTitleNoneExists(ctx, m.Location{Title: locationParentTitle})
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		l := m.Location{Title: locationTitle}
		if key != nil {
			l = l.SetChildDSID([]int64{key.IntID()}).(m.Location)
			m.SaveIfTitleNoneExists(ctx, l)
		}
		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}

// serve image from blobstore
func serveImageHandler(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
}
func init() {
	initDBSet = false
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/serve/", serveImageHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/loc", superHandler)
	http.HandleFunc("/main/", mainHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/", handler)
}
