package model

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type DataModel interface {
}

type Category struct {
	UUID       string
	Title      string
	ParentDSID int64
}

type CategoryVM struct {
	Title    string
	Children []CategoryVM
}

func SaveCategoryWithTitles(ctx context.Context, title string, parentTitle string) error {
	if parentTitle == "" {
		_, err := SaveCategoryIfNonExist(ctx, &Category{Title: title})
		return err
	} else {
		// try to get valid parent by title
		keys, _, err := GetCategoryByTitle(ctx, parentTitle)
		if err != nil {
			return err
		}
		if len(keys) == 0 { // no parent, create it
			key, _ := SaveCategoryIfNonExist(ctx, &Category{Title: parentTitle})
			keys = append(keys, key)
		}
		if len(keys) == 0 {
			return errors.New("Error in getting proper parent information")
		}
		for _, key := range keys { // the parents
			c := Category{Title: title, ParentDSID: key.IntID()}
			_, err := SaveCategoryIfNonExist(ctx, &c)
			return err
		}
	}
	return nil
}

func GetCategoryByTitle(ctx context.Context, title string) ([]*datastore.Key, []Category, error) {
	var cs []Category
	keys, err := datastore.NewQuery("Category").Filter("Title =", title).GetAll(ctx, &cs)
	return keys, cs, err
}

func GetCategoryVM(ctx context.Context, cs *[]Category, title string) (*[]CategoryVM, error) {
	cvms := make([]CategoryVM, 0)
	if title == "" {
		for _, c := range *cs {
			if c.ParentDSID == 0 {
				cvm := CategoryVM{
					Title: c.Title,
				}
				cvms = append(cvms, cvm)
			}
		}
	} else {
		for _, c := range *cs { // get cvm mathing with title
			if c.Title == title {
				cvm := CategoryVM{
					Title: c.Title,
				}
				cvms = append(cvms, cvm)
			}
		}
	}
	for _, c := range *cs {
		if c.ParentDSID != 0 {
			category, err := getParentCategoryFromID(ctx, c.ParentDSID)
			if err != nil {
				return nil, err
			}
			for i, _ := range cvms {
				if cvms[i].Title == category.Title {
					ncvms, err := GetCategoryVM(ctx, cs, c.Title)
					if err != nil {
						return nil, err
					}
					for _, ncvm := range *ncvms {
						cvms[i].Children = append(cvms[i].Children, ncvm)
					}
				}
			}
		}
	}
	return &cvms, nil
}

func getParentCategoryFromID(ctx context.Context, id int64) (Category, error) {
	key := datastore.NewKey(ctx, "Category", "", id, nil)
	var category Category
	err := datastore.Get(ctx, key, &category)
	return category, err
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
		c2.ParentDSID = key1.IntID()
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
