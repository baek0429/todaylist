package model

import "testing"

func TestCategory(t *testing.T) {
	c := CategoryVM{Title: "Hello", Children: []CategoryVM{CategoryVM{Title: "child"}}}
	c.SetTitle("title").PrintTitle()
}
