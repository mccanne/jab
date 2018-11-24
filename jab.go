package jab

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type Matcher interface {
	match(map[string]interface{}) bool
}

type FieldMatcher struct {
	Field string
	Value string
}

type Binding struct {
	Parent *Template
	Matcher
	Struct   interface{}
	Children []Child
}

type Template struct {
	DefaultKey string
	Bindings   []Binding
}

func NewTemplate(key string) *Template {
	return &Template{key, nil}
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

func (t *Template) Match(value string, typ interface{}) *Binding {
	binding := &Binding{t, FieldMatcher{t.DefaultKey, value}, typ, nil}
	t.Bindings = append(t.Bindings, *binding)
	return binding
}

func (t *Template) match(node map[string]interface{}) *Binding {
	for _, binding := range t.Bindings {
		if binding.match(node) {
			return &binding
		}
	}
	return nil
}

func (b *Binding) AddChild(key string) *Template {
	t := NewTemplate(b.Parent.DefaultKey)
	b.Children = append(b.Children, Child{key, t})
	return t
}

func parse(node map[string]interface{}, t *Template) (interface{}, error) {
	match := t.match(node)
	if match == nil {
		return nil, nil
	}
	instance := reflect.New(reflect.TypeOf(match.Struct))
	for _, child := range match.Children {
		childNode, ok := node[child.Field]
		if !ok {
			err := fmt.Errorf("missing child: %s", child.Field)
			return nil, err
		}
		childMap, ok := childNode.(map[string]interface{})
		if !ok {
			err := fmt.Errorf("field %s is not an object", child.Field)
			return nil, err
		}
		childInstance, err := parse(childMap, child.Template)
		if err != nil {
			return nil, err
		}
		field := instance.FieldByName(child.Field)
		field.Set(reflect.ValueOf(childInstance))
	}
	return instance.Interface(), nil
}

func Parse(b []byte, t *Template) (interface{}, error) {
	var root interface{}
	err := json.Unmarshal(b, &root)
	if err != nil {
		return nil, err
	}
	object, ok := root.(map[string]interface{})
	if !ok {
		return nil, errors.New("input not a JSON object")
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
