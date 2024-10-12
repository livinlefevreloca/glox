package main

import "fmt"

type Environment struct {
	name   string
	values map[string]any
	parent *Environment
}

func (e Environment) assign(name string, value any) error {
	// fmt.Println("assign ", value, " from ", e.name, " to ", e.values)
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return nil
	}

	if e.parent != nil {
		return e.parent.assign(name, value)
	}

	return fmt.Errorf("Undefined variable : %s", name)
}

func (e Environment) define(name string, value any) {
	e.values[name] = value
	// fmt.Println("define ", name, " from ", e.name, " to ", e.values)
}

func (e Environment) get(name string) (any, error) {
	// fmt.Println("get ", name, " from ", e.name, " with ", e.values)
	if value, ok := e.values[name]; ok {
		return value, nil
	}

	if e.parent != nil {
		return e.parent.get(name)
	}

	return nil, fmt.Errorf("Undefined variable : %s", name)
}
