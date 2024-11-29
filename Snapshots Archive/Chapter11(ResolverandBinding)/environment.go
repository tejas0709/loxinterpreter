package main

import "fmt"

// Environment represents a variable scope and stores variable bindings.
type Environment struct {
	values map[string]interface{}
	parent *Environment
}

// NewEnvironment creates a new environment.
func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]interface{})}
}

// NewEnclosedEnvironment creates a new environment with a parent.
func NewEnclosedEnvironment(parent *Environment) *Environment {
	return &Environment{values: make(map[string]interface{}), parent: parent}
}

// Define adds a new variable to the environment.
func (env *Environment) Define(name string, value interface{}) {
	env.values[name] = value
}

// Get retrieves the value of a variable, checking parent environments if necessary.
func (env *Environment) Get(name Token) interface{} {
	// First check the current environment
	if value, found := env.values[name.Lexeme]; found {
		return value
	}

	// Then check parent environments
	if env.parent != nil {
		return env.parent.Get(name)
	}

	// If not found in any environment, raise an undefined variable error
	panic(fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}

// Assign updates the value of an existing variable, checking parent environments if necessary.
func (env *Environment) Assign(name Token, value interface{}) {
	// First check the current environment
	if _, found := env.values[name.Lexeme]; found {
		env.values[name.Lexeme] = value
		return
	}

	// Then check and update parent environments
	if env.parent != nil {
		env.parent.Assign(name, value)
		return
	}

	// If not found in any environment, raise an undefined variable error
	panic(fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}

// Get a variable value at a specific depth.
func (env *Environment) getAt(distance int, name string) interface{} {
	current := env
	for i := 0; i < distance; i++ {
		current = current.parent
	}
	return current.values[name]
}

// Assign a value to a variable at a specific depth.
func (env *Environment) assignAt(distance int, name Token, value interface{}) {
	current := env
	for i := 0; i < distance; i++ {
		current = current.parent
	}
	current.values[name.Lexeme] = value
}
