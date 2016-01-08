package model

import "testing"

func TestCategory(t *testing.T) {

	cms:= []Category
	cms = append(cms, Category{})

	type Category struct {
	UUID       string
	Title      string
	ParentDSID string
}

}
