package model

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Category struct {
	UUID        string
	Title       string
	Description string
	ParentDSID  string
}

func SaveCategoryIfNonExist(ctx context.Context, title string, dcrp string) (*datastore.Key, error) {
	if !IsCategoryExist(ctx, title) {
		c := NewCategory()
		c.Title = title
		c.Description = dcrp
		key, err := saveCategory(ctx, &c)
		return key, err
	}
	return nil, errors.New("there exists same category name")
}

func IsCategoryExist(ctx context.Context, title string) bool {
	var cs []Category
	datastore.NewQuery("Category").Filter("Title =", title).GetAll(ctx, &cs)
	if len(cs) != 0 {
		return true
	}
	return false
}

func saveCategory(ctx context.Context, c *Category) (*datastore.Key, error) {
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Category", nil), c)
}

func NewCategory() Category {
	c := Category{
		UUID: uuid.New(),
	}
	return c
}
