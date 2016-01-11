package model

import ()

type Location struct {
	Title      string
	ParentDSID []int64
	ChildDSID  []int64
}

func (c Location) GetTitle() string {
	return c.Title
}

func (c Location) GetParentDSID() []int64 {
	return c.ParentDSID
}

func (c Location) SetParentDSID(inputs []int64) DataModel {
	for _, input := range inputs {
		c.ParentDSID = append(c.ParentDSID, input)
	}
	return c
}

func (c Location) SetChildDSID(inputs []int64) DataModel {
	for _, input := range inputs {
		c.ChildDSID = append(c.ChildDSID, input)
	}
	return c
}

func (c Location) GetChildDSID() []int64 {
	return c.ChildDSID
}
