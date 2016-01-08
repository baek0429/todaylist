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
	ParentDSID string
}

type LocationVM struct {
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
		l.ParentDSID = key1.StringID()
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
