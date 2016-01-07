package model

import (
	"code.google.com/p/go-uuid/uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

type Post struct {
	UUID         string
	Title        string
	Password     string
	Address      string
	Hour         string
	Description  string
	BlobKeys     []string
	Time         time.Time
	ImageSrc     []string
	CategoryKeys []string
}

type Category struct {
	UUID        string
	Title       string
	Description string
	ParentDSID  string
}

type Location struct {
	UUID       string
	Title      string
	PArentDSID string
}

func NewLocation() Location {
	l := Location{
		UUID: uuid.New(),
	}
	return l
}

func saveCategory(ctx context.Context, p *Category) (*datastore.Key, error) {
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Category", nil), p)
}

func NewCategory() Category {
	c := Category{
		UUID: uuid.New(),
	}
	return c
}

func NewPost() Post {
	post := Post{
		UUID: uuid.New(),
	}
	return post
}

func ParsePostByUID(ctx context.Context, uuid string) []Post {
	var posts []Post
	q := datastore.NewQuery("Post").Filter("UUID =", uuid)
	q.GetAll(ctx, &posts)
	return posts
}

func ParseAllPosts(ctx context.Context) ([]Post, error) {
	var posts []Post
	_, err := datastore.NewQuery("Post").Order("Time").GetAll(ctx, &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func savePost(ctx context.Context, p *Post) (*datastore.Key, error) {
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Post", nil), p)
}

func SavePosts(ctx context.Context, ps *([]*Post)) error {
	for _, p := range *ps {
		_, err := savePost(ctx, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetSamplePosts() []Post {
	posts := make([]Post, 10)
	for i, post := range posts {
		post.Title = "title"
		post.Address = "address"
		post.Hour = "hour"
		post.Description = "description"
		post.Password = "1"
		posts[i] = post
	}
	return posts
}

func GetMainPosts() []Post {
	return GetSamplePosts()
}

func GetPostByCategory() []Post {
	return GetSamplePosts()
}
