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
)

func handler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*")) // add sub-templates in /template
	t.ParseFiles("index.html")                           // parse main.html as main
	ctx := appengine.NewContext(r)
	posts, err := m.ParseAllPosts(ctx) // get all posts in db.
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	err = t.ExecuteTemplate(w, "base", posts) // exec templates
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func superHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*")) // add sub-templates in /template
	t.ParseFiles("super.html")                           //parse super

	c := make([]m.Category, 2)
	c[0] = m.NewCategory()
	c[0].Title = "test"
	c[0].Description = "testtesttesttesttesttesttesttest"
	c[1] = m.NewCategory()
	c[1].Title = "test"
	c[1].Description = "testtesttesttesttesttesttesttest"
	l := make([]m.Location, 2)
	l[0] = m.NewLocation()
	l[0].Title = "test"
	l[1] = m.NewLocation()
	l[1].Title = "test"
	err := t.ExecuteTemplate(w, "base",
		struct {
			Category []m.Category
			Location []m.Location
		}{c, l})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("template/*")) // add sub-templates in /template
	t.ParseFiles("main.html")                            //parse main
	r.ParseForm()
	ctx := appengine.NewContext(r)
	uuid := r.FormValue("uuid") // get post by uuid
	posts := m.ParsePostByUID(ctx, uuid)
	err := t.ExecuteTemplate(w, "base", posts)
	if err != nil {
		w.Write([]byte(err.Error()))
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

		post := m.NewPost() // creating post and fill the fields.
		post.Title = title
		post.Password = password
		post.Description = description
		post.BlobKeys = blobKeys
		post.Time = time.Now()
		post.ImageSrc = imgSrcs

		posts := []*m.Post{&post} // uploading posts
		c := make(chan int, 1)
		go func() { // it seems that I need to creat a term to make sure that the change appears in the db.
			er := m.SavePosts(ctx, &posts) // success?
			if er != nil {
				http.Error(w, er.Error(), http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Millisecond * 100)
			c <- 1
		}()
		if 1 == <-c {
			http.Redirect(w, r, "/main/?uuid="+posts[0].UUID, http.StatusFound) // redirect to /main/?uuid= with uuid
		}
		http.Redirect(w, r, "/error", http.StatusNotFound)
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

// serve image from blobstore
func serveImageHandler(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
}
func init() {
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/serve/", serveImageHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/loc", superHandler)
	http.HandleFunc("/main/", mainHandler)
	http.HandleFunc("/", handler)
}
