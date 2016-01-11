package model

import "testing"

func TestCategory(t *testing.T) {
	c := Category{Title: "hello"}
	a := c.SetTitle("a").SetParentDSID(1)
	t.Log(a)

	vm := CategoryVM{Title: "hello", Children: []CategoryVM{{Title: "hello2"}}}
	t.Log(vm.GetChildren())
}
