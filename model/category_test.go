package model

import "testing"

func TestCategory(t *testing.T) {

	cms:= []Category
	cms = append(cms, Category{Title:"test1",})

	type Category struct {
		UUID       string
		Title      string
		ParentDSID string
	}

}
