package jab_test

import (
	"encoding/json"
	"fmt"
	"reflect"
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

type Baz struct {
	Op string
	A  Foo
	B  Bar
}

type Bug struct {
	Op string
	A  Foo
}

func test(s string, template *jab.Template) {
	out, err := jab.Parse([]byte(s), template)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		pretty.Println(out)
	}
}

func test0() {
	//	var x interface{} = Bug{}
	x := Bar{}
	t := reflect.TypeOf(x)
	pretty.Println(t)
	instance := reflect.New(t).Elem()
	pretty.Println(instance.Kind().String())
	pretty.Println(instance.FieldByName("A"))

	s3 := `{
		"op": "bug",
		"A": {
			"op": "foo",
			"value": 13
		}
	}`

	template := jab.NewTemplate()
	template.MatchKey("op", "foo", Foo{})
	bug := template.MatchKey("op", "bug", Bug{})
	bug.ChildObject("A", template)
	test(s3, template)
}

func test1() {
	s1 := `{
		"op": "foo",
		"value": 12
	}`
	s2 := `{
		"op": "bar",
		"value": "hello"
	}`
	s3 := `{
		"op": "baz",
		"A": {
			"op": "foo",
			"value": 13
		},
		"B": {
			"op": "bar",
			"value": "there"
		}
	}`
	/*
		s4 := `{
			"op": "baz",
			"A": {
				"op": "foo",
				"value": 13
			},
			"B": {
				"op": "baz",
				"A": {
					"op": "foo",
					"value": 13
				},
				"B": {
					"op": "baz",
					"value": "oops"
				}
			}
		}`
	*/
	template := jab.NewTemplate()
	template.MatchKey("op", "foo", Foo{})
	template.MatchKey("op", "bar", Bar{})
	baz := template.MatchKey("op", "baz", Baz{})
	baz.ChildObject("A", template)
	baz.ChildObject("B", template)
	test(s1, template)
	test(s2, template)
	test(s3, template)
	//(s4, template)
}

type Node interface {
	nodeType()
}

type InternalNode struct {
	Type  string `json:"type"`
	Left  Node   `json:"left"`
	Right Node   `json:"right"`
}

func (p *InternalNode) nodeType() {}
func (p *LeafNode) nodeType()     {}

type LeafNode struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

const tree = `{
	"type": "internal",
	"left": {
		"type": "internal",
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
		"type": "leaf",
		"value": 7
	}
}`

const bug2 = `{
	"type": "internal",
	"left": {
		"type": "leaf",
		"value": 5
	},
	"right": {
		"type": "leaf",
		"value": 7
	}
}`

func test3() {
	node := &InternalNode{
		Type:  "",
		Left:  &LeafNode{},
		Right: &LeafNode{}}
	pretty.Println(node)
	var v interface{} = node
	err := json.Unmarshal([]byte(bug2), &v)
	if err != nil {
		pretty.Println(err.Error())
	}
	pretty.Println(v)
}

func test2() {
	template := jab.NewTemplate()
	template.MatchKey("type", "leaf", LeafNode{})
	internal := template.MatchKey("type", "internal", InternalNode{})
	internal.ChildObject("left", template)
	internal.ChildObject("right", template)
	//test(tree, template)
	test(tree, template)
}

func test4() {
	template := jab.NewTemplate()
	template.MatchKey("type", "leaf", LeafNode{})
	internal := template.MatchKey("type", "internal", InternalNode{})
	internal.ChildObject("left", template)
	internal.ChildObject("right", template)

	out, err := jab.Parsex([]byte(bug2), template)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		pretty.Println(out)
	}
	pretty.Println("FINAL UNMARSHALL")
	pretty.Println(out)
	err = json.Unmarshal([]byte(bug2), &out)
	if err != nil {
		pretty.Println(err.Error())
	}
	pretty.Println(out)
}

func Test_Jab(t *testing.T) {
	//test0()
	//test1()
	test2()
	//test3()
	//test4()
}
