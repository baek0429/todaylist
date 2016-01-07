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

func SaveLocationIfNonExist(ctx context.Context, title string) (*datastore.Key, error) {
	if !IsLocationExist(ctx, title) {
		l := NewLocation()
		l.Title = title
		key, err := saveLocation(ctx, &l)
		return key, err
	}
	return nil, errors.New("there exists same location")
}

func IsLocationExist(ctx context.Context, title string) bool {
	var ls []Location
	datastore.NewQuery("Locaction").Filter("Title =", title).GetAll(ctx, &ls)
	if len(ls) != 0 {
		return true
	}
	return false
}
