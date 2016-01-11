package model

import (
	"time"
)

type Post struct {
	ID           int64
	Title        string
	Password     string
	Address      string
	Hour         string
	Description  string
	BlobKeys     []string
	Time         time.Time
	ImageSrc     []string
	CategoryKeys []int64
	ParentDSID   []int64
	ChildDSID    []int64
}

func (p Post) GetTitle() string {
	return p.Title
}

func (p Post) GetParentDSID() []int64 {
	return p.ParentDSID
}

func (p Post) SetParentDSID(inputs []int64) DataModel {
	for _, input := range inputs {
		p.ParentDSID = append(p.ParentDSID, input)
	}
	return p
}

func (p Post) SetChildDSID(inputs []int64) DataModel {
	for _, input := range inputs {
		p.ChildDSID = append(p.ChildDSID, input)
	}
	return p
}

func (p Post) GetChildDSID() []int64 {
	return p.ChildDSID
}
