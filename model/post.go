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
