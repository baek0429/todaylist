package model

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Category struct {
	UUID       string
	Title      string
	ParentDSID string
}

type CategoryViewModel struct {
	Single   Category
	Children []Category
}

func (cs *[]Category) Sort(title string) *[]CategoryViewModel {
	var cvms *[]CategoryViewModel
	for _, c := range *cs {
		if c.Title == title { // if title same to title
			cvm := CategoryViewModel{
				Single: c,
			}
			cvms = append(cvms, cvm)
		} else { // not same
			scvms := cs.Sort(c.Title)
			c.Children = scvms
		}
	}
	return cvms
}

func InitiateSamples(ctx context.Context) {
	categoryInit(ctx)
	locationInit(ctx)
}

func categoryInit(ctx context.Context) {
	categorySampling(ctx)
}

func categorySampling(ctx context.Context) {
	key1, _ := SaveCategoryIfNonExist(ctx, &Category{Title: "Apartment", UUID: uuid.New()})
	c2 := NewCategory()
	if key1 != nil {
		c2.Title = "Gangnam"
		c2.ParentDSID = key1.StringID()
		saveCategory(ctx, &c2)
	}
}

func SaveCategoryIfNonExist(ctx context.Context, c *Category) (*datastore.Key, error) {
	if !IsCategoryExist(ctx, c) {
		key, err := saveCategory(ctx, c)
		return key, err
	}
	return nil, errors.New("there exists same category name")
}

func ParseCategory(ctx context.Context) (*[]Category, error) {
	var cs []Category
	_, err := datastore.NewQuery("Category").GetAll(ctx, &cs)
	if err != nil {
		return &cs, err
	}
	return &cs, nil
}

func IsCategoryExist(ctx context.Context, c *Category) bool {
	var cs []Category
	datastore.NewQuery("Category").Filter("Title =", c.Title).GetAll(ctx, &cs)
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
