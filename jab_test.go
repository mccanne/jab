package jab_test

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
	"github.com/mccanne/jab"
)

type Foo struct {
	Op    string
	Value int
}

type Bar struct {
	Op    string
	Value string
}

//type Node interface{}

//type Leaf struct {
//	Value
//}

//func genSimple() *jab.Template {
//	template := jab.NewTemplate()
//	foo := template.Match("foo", Foo{})
//	bar := template.Match("bar", Bar{})
//	baz := template.Match("baz", Baz{})
//	baz.AddChild("A").Match("foo", foo)
//	baz.AddChild("B").Match("bar", bar)
//	return template

func test1() {
	s := `{
		"Op": "foo",
		"value": 12
	}`
	template := jab.NewTemplate()
	template.MatchKey("Op", "foo", Foo{})
	test(s, template)
}

func test2() {
	s := `[
			{
				"Op": "foo",
				"value": 12
			},
		{
			"Op": "bar",
			"value": "hello"
		}

	]`
	template := jab.NewTemplate()
	template.MatchKey("Op", "foo", Foo{})
	test(s, template)
}

func test(s string, template *jab.Template) {
	out, err := jab.Parse([]byte(s), template)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		pretty.Println(out)
	}
}

const tree = `{
	"type": "interal",
	"left": {
		"left": {
			"type": "leaf",
			"value": 1
		},
		"right": {
			"type": "internal",
			"left": {
				"type": "leaf",
				"value": 4
			},
			"right": {
				"type": "leaf",
				"value": 5
			}
		}

	},
	"right": {
		"type": "leaf"
		"value": 7
	}
}`

func Test_Jab(t *testing.T) {
	test1()
	test2()
}
