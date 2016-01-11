package model

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"reflect"
	"strconv"
	"strings"
)

type DataModel interface {
	GetTitle() string
	GetParentDSID() []int64
	SetParentDSID([]int64) DataModel
	GetChildDSID() []int64
	SetChildDSID([]int64) DataModel
}

func getType(m DataModel) string {
	t := reflect.TypeOf(m).String()
	ty := strings.Split(t, ".")[1]
	return ty
}

func ParseKeyFromID(ctx context.Context, id string, m DataModel) *datastore.Key {
	ty := getType(m)
	parentKey := datastore.NewKey(ctx, "Model", "Model", 0, nil)
	i, _ := strconv.Atoi(id)
	key := datastore.NewKey(ctx, ty, "", int64(i), parentKey)
	return key
}

func ParseAll(ctx context.Context, m DataModel) ([]*datastore.Key, error) {
	ty := getType(m)
	keys, err := datastore.NewQuery(ty).KeysOnly().GetAll(ctx, nil)
	return keys, err
}

func ParseEntitiesFromKeys(ctx context.Context, keys []*datastore.Key) ([]DataModel, error) {
	var dms []DataModel
	for _, key := range keys {
		dm, err := ParseEntityFromKey(ctx, key)
		if err != nil {
			return nil, err
		}
		dms = append(dms, dm)
	}
	return dms, nil
}

func ParseEntityFromKey(ctx context.Context, key *datastore.Key) (DataModel, error) {
	ty := key.Kind()
	var err error
	switch ty {
	case "Category":
		var m Category
		err = datastore.Get(ctx, key, &m)
		return m, nil
	case "Location":
		var m Location
		err = datastore.Get(ctx, key, &m)
		return m, nil
	case "Post":
		var p Post
		err = datastore.Get(ctx, key, &p)
		return p, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func SaveIfTitleNoneExists(ctx context.Context, m DataModel) (*datastore.Key, error) {
	ty := getType(m)
	keys, err := datastore.NewQuery(ty).Filter("Title =", m.GetTitle()).KeysOnly().GetAll(ctx, nil)
	if len(keys) == 0 {
		key, err := SaveDataModel(ctx, m)
		return key, err
	}
	return keys[0], err
}

func SaveDataModels(ctx context.Context, models []DataModel) ([]*datastore.Key, error) {
	var keys []*datastore.Key
	for _, model := range models {
		key, err := SaveDataModel(ctx, model)
		if err != nil {
			return keys, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func SaveDataModel(ctx context.Context, m DataModel) (*datastore.Key, error) {
	ty := getType(m)
	parentKey := datastore.NewKey(ctx, "Model", "Model", 0, nil)
	key := datastore.NewIncompleteKey(ctx, ty, parentKey) // create unique id all kinds of the datamodel
	switch ty {
	case "Category":
		c := m.(Category)
		return datastore.Put(ctx, key, &c)
	case "Location":
		l := m.(Location)
		return datastore.Put(ctx, key, &l)
	case "Post":
		p := m.(Post)
		return datastore.Put(ctx, key, &p)
	}
	return nil, errors.New("errors in saving model")
}
