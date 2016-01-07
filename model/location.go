package model

import (
	"code.google.com/p/go-uuid/uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

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

func saveLocation(ctx context.Context, l *Location) (*datastore.Key, error) {
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Location", nil), l)
}

func SaveLocationIfNonExist(ctx context.Context, title string, dcrp string) (*datastore.Key, error) {
	if !IsLocationExist(ctx, title) {
		l := NewLocation()
		l.Title = title
		l.Description = dcrp
		key, err := saveLocation(ctx, &l)
		return key, err
	}
	return nil, error.Error("There exists same name category")
}

func IsLocactionExist(ctx context.Context, string title) bool {
	var ls []Location
	datastore.NewQuery("Locaction").Filter("Title =", title).GetAll(ctx, &ls)
	if len(ls) != 0 {
		return true
	}
	return false
}
