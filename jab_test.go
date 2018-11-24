package jab_test

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
	"github.com/mccanne/jab"
)

type Foo struct {
	Op string
	A  int
	B  string
}

type Bar struct {
	Op string
	A  int
	B  int
}

type Baz struct {
	Op string
	A  Foo
	B  Bar
}

func gen() *jab.Template {
	template := jab.NewTemplate("Op")
	foo := template.Match("foo", Foo{})
	bar := template.Match("bar", Bar{})
	baz := template.Match("baz", Baz{})
	baz.AddChild("A").Match("foo", foo)
	baz.AddChild("B").Match("bar", bar)
	return template
}

func test(s string, template *jab.Template) {
	out, err := jab.Parse([]byte(s), template)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		pretty.Println(out)
	}
}

func Test_Jab(t *testing.T) {
	s1 := `{
		"Op": "foo",
		"A": 12,
		"B": "bar"
	}`
	s2 := `{
		"Op": "bar",
		"A": 12,
		"B": 17
	}`
	s3 := `{
		"Op": "baz",
		"A": {
			"Op": "foo",
			"A": 77,
			"B": "hello"
		},
		"B": {
			"Op": "bar",
			"A": 9809,
			"B": 31416
		}
	}`
	template := gen()
	test(s1, template)
	test(s2, template)
	test(s3, template)
}
