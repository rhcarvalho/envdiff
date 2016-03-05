// env-diff is a tool for comparing variables in two environments.
package main

import (
	"os"
	"strings"
)

const pathListSeparator = string(os.PathListSeparator)

// Var represents a single environment variable. A Var should have a non-empty
// Name, and zero or more values in Value. A single value should not contain
// os.PathListSeparator.
type Var struct {
	Name  string
	Value []string
}

// Env represents an environment, a list of environment variables.
type Env []Var

// Parse parses the string s with environment variables as printed by `env` or
// `printenv` and returns an Env.
func Parse(s string) Env {
	var env Env
	for _, v := range strings.Split(s, "\n") {
		pair := strings.SplitN(v, "=", 2)
		if len(pair) < 2 {
			pair = append(pair, "")
		}
		name, value := pair[0], pair[1]
		if name == "" {
			continue
		}
		env = append(env, Var{name, strings.Split(value, pathListSeparator)})
	}
	return env
}
