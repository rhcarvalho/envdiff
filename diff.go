// envdiff is a tool for comparing variables in two environments.
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

// asMap return the environment e as a map.
func (e Env) asMap() map[string][]string {
	m := make(map[string][]string)
	for _, v := range e {
		m[v.Name] = v.Value
	}
	return m
}

// Parse parses the string s with environment variables as printed by `env` or
// `printenv` and returns an Env.
func Parse(s string) Env {
	var env Env
	sep := "\n"
	if strings.Contains(s, "\x00") {
		sep = "\x00"
	}
	for _, v := range strings.Split(s, sep) {
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

// Diff returns an Env with environment variables that were unset from old, and
// set or overwritten in new. The result is such that applying the diff to an
// old environment turns it into new.
func Diff(old, new Env) Env {
	var diff Env
	oldMap := old.asMap()
	newMap := new.asMap()
	// Look for variables that were unset.
	for _, v := range old {
		if _, ok := newMap[v.Name]; ok {
			continue
		}
		diff = append(diff, Var{Name: v.Name})
	}
	var (
		val []string
		ok  bool
	)
	// Look for variables that were kept or overwritten.
	for _, v := range new {
		// Ignore variables that are in a and b with the same value.
		if val, ok = oldMap[v.Name]; ok && stringSliceEqual(v.Value, val) {
			continue
		}
		// Look for common prefix or suffix.
		for i := 0; i < len(v.Value) && i < len(val); i++ {
			// Common suffix: PATH=/new/path:$PATH.
			if stringSliceEqual(v.Value[i+1:], val) {
				v.Value = append(v.Value[:i+1], "$"+v.Name)
				break
			}
			// Common prefix: PATH=$PATH:/new/path.
			if stringSliceEqual(v.Value[:i+1], val) {
				v.Value = append([]string{"$" + v.Name}, v.Value[i+1:]...)
				break
			}
		}
		diff = append(diff, v)
	}
	return diff
}

// stringSliceEqual returns true if s and t are equal.
func stringSliceEqual(s, t []string) bool {
	if (s == nil) != (t == nil) {
		return false
	}
	if len(s) != len(t) {
		return false
	}
	for i := range s {
		if s[i] != t[i] {
			return false
		}
	}
	return true
}
