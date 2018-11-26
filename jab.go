package jab

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/kr/pretty"
)

type Matcher interface {
	match(map[string]interface{}) bool
}

type FieldMatcher struct {
	Field string
	Value string
}

type Binding struct {
	Matcher
	Struct   interface{}
	Children []Child
}

type Template struct {
	Bindings []*Binding
}

func NewTemplate() *Template {
	return &Template{}
}

type Child struct {
	Field    string
	Template *Template
}

func (matcher FieldMatcher) match(node map[string]interface{}) bool {
	v, ok := node[matcher.Field]
	if !ok {
		return false
	}
	return v == matcher.Value
}

func (t *Template) MatchKey(key, value string, typ interface{}) *Binding {
	binding := &Binding{FieldMatcher{key, value}, typ, nil}
	t.Bindings = append(t.Bindings, binding)
	return binding
}

func (t *Template) match(node map[string]interface{}) *Binding {
	for _, binding := range t.Bindings {
		if binding.match(node) {
			return binding
		}
	}
	return nil
}

func (b *Binding) ChildObject(key string, t *Template) *Binding {
	b.Children = append(b.Children, Child{key, t})
	return b
}

func old(node map[string]interface{}, t *Template) (interface{}, error) {

	pretty.Println("PARSE")
	pretty.Println(node)
	pretty.Println(t)

	match := t.match(node)

	if match == nil {
		pretty.Println("NO MATCH")
		pretty.Println(node)
		pretty.Println(t)
		return nil, nil
	}
	pretty.Println("MATCH & MAKE")
	pretty.Println(match.Struct)
	instance := reflect.New(reflect.TypeOf(match.Struct))
	parent := instance.Elem()

	for _, child := range match.Children {
		childNode, ok := node[child.Field]
		if !ok {
			err := fmt.Errorf("missing child %s for struct %T", child.Field, match.Struct)
			return nil, err
		}
		childMap, ok := childNode.(map[string]interface{})
		if !ok {
			err := fmt.Errorf("field %s is not an object", child.Field)
			return nil, err
		}
		pretty.Println("PARSE CHILD")
		childInstance, err := parse(childMap, child.Template)
		if err != nil {
			return nil, err
		}

		fname := child.Field
		field := parent.FieldByName(fname)
		if !field.IsValid() {
			fname = strings.Title(fname)
			field = parent.FieldByName(fname)
		}

		pretty.Println("SETTING")
		pretty.Println("===")
		pretty.Println(fname)
		pretty.Println(match.Struct)
		pretty.Println(childInstance)
		pretty.Println(reflect.ValueOf(childInstance).Elem())
		field.Set(reflect.ValueOf(childInstance).Elem())
		pretty.Println("===")
	}
	return instance.Interface(), nil
}

var nilValue reflect.Value

func parse(node map[string]interface{}, t *Template) (reflect.Value, error) {

	pretty.Println("PARSE")
	pretty.Println(node)
	pretty.Println(t)

	match := t.match(node)

	if match == nil {
		pretty.Println("NO MATCH")
		pretty.Println(node)
		pretty.Println(t)
		return nilValue, nil
	}
	pretty.Println("MATCH & MAKE")
	pretty.Println(match.Struct)
	parentPtr := reflect.New(reflect.TypeOf(match.Struct))
	parent := parentPtr.Elem()

	for _, child := range match.Children {
		childNode, ok := node[child.Field]
		if !ok {
			err := fmt.Errorf("missing child %s for struct %T", child.Field, match.Struct)
			return nilValue, err
		}
		childMap, ok := childNode.(map[string]interface{})
		if !ok {
			err := fmt.Errorf("field %s is not an object", child.Field)
			return nilValue, err
		}
		pretty.Println("PARSE CHILD")
		childVal, err := parse(childMap, child.Template)
		if err != nil {
			return nilValue, err
		}
		if childVal == nilValue {
			continue
		}

		fname := child.Field
		field := parent.FieldByName(fname)
		if !field.IsValid() {
			fname = strings.Title(fname)
			field = parent.FieldByName(fname)
		}

		pretty.Println("SETTING")
		pretty.Println("===")
		pretty.Println(fname)
		pretty.Println(childVal.Interface())
		pretty.Println("PARENT-PRE-SET")
		pretty.Println(parent.Interface())
		pretty.Println(field.Type().Kind().String())
		pretty.Println(childVal.Type().Kind().String())
		if field.Type().Kind().String() == "interface" {
			field.Set(childVal.Addr())
		} else {
			field.Set(childVal)
		}
		pretty.Println("PARENT-POST-SET")
		pretty.Println(parent.Interface())
		pretty.Println("===")
	}
	//return parentPtr, nil
	return parent, nil
}

/*
func parseArray(a []interface{}, t *Template) ([]interface{}, error) {
	for k, v := range a {
		o, ok := v.(map[string]interface{})
		if ok {
			out, err := parse(o, t)
			if err != nil {
				return nil, err
			}
			a[k] = out
		}
	}
	return a, nil
}
*/

func Parse(b []byte, t *Template) (interface{}, error) {
	var root interface{}
	err := json.Unmarshal(b, &root)
	if err != nil {
		return nil, err
	}
	object, ok := root.(map[string]interface{})
	var out reflect.Value
	if ok {
		out, err = parse(object, t)
		if err != nil {
			return nil, err
		}
	} else {
		/*
			a, ok := root.([]interface{})
			if !ok {
				return nil, errors.New("input is either a JSON object or a JSON array of objects")
			}
			out, err = parseArray(a, t)
			if err != nil {
				return nil, err
			}*/
	}
	v := out.Addr().Interface()
	pretty.Println("FINAL UNMARSHALL")
	pretty.Println(string(b))
	pretty.Println(v)
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	pretty.Println(v)
	return v, nil
}

func Parsex(b []byte, t *Template) (interface{}, error) {
	var root interface{}
	err := json.Unmarshal(b, &root)
	if err != nil {
		return nil, err
	}
	object, ok := root.(map[string]interface{})
	var out reflect.Value
	if ok {
		out, err = parse(object, t)
		if err != nil {
			return nil, err
		}
	} else {
		/*
			a, ok := root.([]interface{})
			if !ok {
				return nil, errors.New("input is either a JSON object or a JSON array of objects")
			}
			out, err = parseArray(a, t)
			if err != nil {
				return nil, err
			}*/
	}
	return out.Interface(), nil
}
