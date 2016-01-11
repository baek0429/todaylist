package model

import ()

/* define Category and implement DataModel */
type Category struct {
	Title      string
	ParentDSID []int64
	ChildDSID  []int64
}

func (c Category) GetTitle() string {
	return c.Title
}

func (c Category) GetParentDSID() []int64 {
	return c.ParentDSID
}

func (c Category) SetParentDSID(inputs []int64) DataModel {
	for _, input := range inputs {
		c.ParentDSID = append(c.ParentDSID, input)
	}
	return c
}

func (c Category) SetChildDSID(inputs []int64) DataModel {
	for _, input := range inputs {
		c.ChildDSID = append(c.ChildDSID, input)
	}
	return c
}

func (c Category) GetChildDSID() []int64 {
	return c.ChildDSID
}
