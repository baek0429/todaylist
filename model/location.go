package model

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Location struct {
	UUID       string
	Title      string
	ParentDSID int64
}

type LocationVM struct {
	Title    string
	Children []LocationVM
}

func SaveLocationWithTitles(ctx context.Context, title string, parentTitle string) error {
	if parentTitle == "" {
		_, err := SaveLocationIfNonExist(ctx, &Location{Title: title})
		return err
	} else {
		// try to get valid parent by title
		keys, _, err := GetLocationByTitle(ctx, parentTitle)
		if err != nil {
			return err
		}
		if len(keys) == 0 { // no parent, create it
			key, _ := SaveLocationIfNonExist(ctx, &Location{Title: parentTitle})
			keys = append(keys, key)
		}
		if len(keys) == 0 {
			return errors.New("Error in getting proper parent information")
		}
		for _, key := range keys { // the parents
			c := Location{Title: title, ParentDSID: key.IntID()}
			_, err := SaveLocationIfNonExist(ctx, &c)
			return err
		}
	}
	return nil
}

func GetLocationByTitle(ctx context.Context, title string) ([]*datastore.Key, []Location, error) {
	var ls []Location
	keys, err := datastore.NewQuery("Location").Filter("Title =", title).GetAll(ctx, &ls)
	return keys, ls, err
}

func GetLocationVM(ctx context.Context, ls *[]Location, title string) (*[]LocationVM, error) {
	lvms := make([]LocationVM, 0)
	if title == "" {
		for _, c := range *ls {
			if c.ParentDSID == 0 {
				lvm := LocationVM{
					Title: c.Title,
				}
				lvms = append(lvms, lvm)
			}
		}
	} else {
		for _, c := range *ls { // get lvm mathing with title
			if c.Title == title {
				lvm := LocationVM{
					Title: c.Title,
				}
				lvms = append(lvms, lvm)
			}
		}
	}
	for _, c := range *ls {
		if c.ParentDSID != 0 {
			location, err := getParentLocationFromID(ctx, c.ParentDSID)
			if err != nil {
				return nil, err
			}
			for i, _ := range lvms {
				if lvms[i].Title == location.Title {
					nlvms, err := GetLocationVM(ctx, ls, c.Title)
					if err != nil {
						return nil, err
					}
					for _, nlvm := range *nlvms {
						lvms[i].Children = append(lvms[i].Children, nlvm)
					}
				}
			}
		}
	}
	return &lvms, nil
}

func getParentLocationFromID(ctx context.Context, id int64) (Location, error) {
	key := datastore.NewKey(ctx, "Location", "", id, nil)
	var location Location
	err := datastore.Get(ctx, key, &location)
	return location, err
}

func NewLocation() Location {
	l := Location{
		UUID: uuid.New(),
	}
	return l
}

func locationInit(ctx context.Context) {
	locationSampling(ctx)
}

func locationSampling(ctx context.Context) {
	key1, _ := SaveLocationIfNonExist(ctx, &Location{Title: "Seoul", UUID: uuid.New()})
	l := NewLocation()
	if key1 != nil {
		l.Title = "Gangnam"
		l.ParentDSID = key1.IntID()
		saveLocation(ctx, &l)
	}
}

func ParseLocation(ctx context.Context) (*[]Location, error) {
	var ls []Location
	_, err := datastore.NewQuery("Location").GetAll(ctx, &ls)
	if err != nil {
		return &ls, err
	}
	return &ls, nil
}

func saveLocation(ctx context.Context, l *Location) (*datastore.Key, error) {
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Location", nil), l)
}

func SaveLocationIfNonExist(ctx context.Context, l *Location) (*datastore.Key, error) {
	if !IsLocationExist(ctx, l) {
		key, err := saveLocation(ctx, l)
		return key, err
	}
	return nil, errors.New("there exists same location")
}

func IsLocationExist(ctx context.Context, l *Location) bool {
	var ls []Location
	datastore.NewQuery("Location").Filter("Title =", l.Title).GetAll(ctx, &ls)
	if len(ls) != 0 {
		return true
	}
	return false
}
