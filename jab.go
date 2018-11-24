package jab

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Rule interface {
	match(map[string]interface{}) bool
}

type FieldRule struct {
	Field string
	Value string
}

type Binding struct {
	Rule
	Struct   interface{}
	Children []Child
}

type Template []Binding

type Child struct {
	Field    string
	Template Template
}

func (rule FieldRule) match(node map[string]interface{}) bool {
	v, ok := node[rule.Field]
	if !ok {
		return false
	}
	return v == rule.Value
}

func (t Template) match(node map[string]interface{}) *Binding {
	for _, binding := range t {
		if binding.match(node) {
			return &binding
		}
	}
	return nil
}

var nilValue reflect.Value

func parse(node map[string]interface{}, t Template) (interface{}, error) {
	match := t.match(node)
	if match == nil {
		return nilValue, nil
	}
	instance := reflect.New(reflect.TypeOf(match.Struct))
	for _, child := range match.Children {
		childNode, ok := node[child.Field]
		if !ok {
			//XXX better error message
			err := fmt.Errorf("missing child: %s", child.Field)
			return nilValue, err
		}
		childMap, ok := childNode.(map[string]interface{})
		if !ok {
			err := fmt.Errorf("field %s is not an object", child.Field)
			return nilValue, err
		}
		childInstance, err := parse(childMap, child.Template)
		if err != nil {
			return nilValue, err
		}
		field := instance.FieldByName(child.Field)
		field.Set(reflect.ValueOf(childInstance))
	}
	return instance, nil
}

func Parse(b []byte, t Template) (interface{}, error) {
	var object map[string]interface{}
	err := json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}
	out, err := parse(object, t)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
