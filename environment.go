package main

import "fmt"

type Environment struct {
	values map[string]any
}

func (e Environment) assign(name string, value any) error {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return nil
	}

	return fmt.Errorf("Undefined variable : %s", name)
}

func (e Environment) define(name string, value any) {
	e.values[name] = value
}

func (e Environment) get(name string) (any, error) {
	if value, ok := e.values[name]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("Undefined variable : %s", name)
}
