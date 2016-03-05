package main

import (
	"reflect"
	"testing"
)

var parseTests = []struct {
	in   string
	want Env
}{
	{
		in:   "",
		want: nil,
	},
	{
		in:   "\n\n\n\n",
		want: nil,
	},
	{
		in:   "TERM=",
		want: Env{{"TERM", []string{""}}},
	},
	{
		in:   "TERM=xterm",
		want: Env{{"TERM", []string{"xterm"}}},
	},
	{
		in:   "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		want: Env{{"PATH", []string{"/usr/local/sbin", "/usr/local/bin", "/usr/sbin", "/usr/bin", "/sbin", "/bin"}}},
	},
	{
		in:   "LS_COLORS=rs=0:di=01;34:ln=01;36",
		want: Env{{"LS_COLORS", []string{"rs=0", "di=01;34", "ln=01;36"}}},
	},
	{
		in: `HOSTNAME=9c584fc5c0d4
TERM=xterm
LS_COLORS=...:*.spx=01;36:*.xspf=01;36:
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
PWD=/
SHLVL=1
HOME=/root
_=/usr/bin/printenv`,
		want: Env{
			{"HOSTNAME", []string{"9c584fc5c0d4"}},
			{"TERM", []string{"xterm"}},
			{"LS_COLORS", []string{"...", "*.spx=01;36", "*.xspf=01;36", ""}},
			{"PATH", []string{"/usr/local/sbin", "/usr/local/bin", "/usr/sbin", "/usr/bin", "/sbin", "/bin"}},
			{"PWD", []string{"/"}},
			{"SHLVL", []string{"1"}},
			{"HOME", []string{"/root"}},
			{"_", []string{"/usr/bin/printenv"}},
		},
	},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		if got := Parse(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Parse(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

var diffTests = []struct {
	desc     string
	old, new Env
	diff     Env
}{
	{},
	{
		desc: "set variable",
		old:  nil,
		new:  Env{{"TERM", []string{"xterm"}}},
		diff: Env{{"TERM", []string{"xterm"}}},
	},
	{
		desc: "unset variable",
		old:  Env{{"TERM", []string{"xterm"}}},
		new:  nil,
		diff: Env{{"TERM", nil}},
	},
	{
		desc: "unset and set different variables",
		old:  Env{{"TERM", []string{"xterm"}}},
		new:  Env{{"PWD", []string{"/"}}},
		diff: Env{
			{"TERM", nil},
			{"PWD", []string{"/"}},
		},
	},
	{
		desc: "unset, set and keep existing (sorted)",
		old: Env{
			{"HOME", []string{"/root"}},
			{"TERM", []string{"xterm"}},
		},
		new: Env{
			{"HOME", []string{"/root"}},
			{"PWD", []string{"/"}},
		},
		diff: Env{
			{"TERM", nil},
			{"PWD", []string{"/"}},
		},
	},
	{
		desc: "unset, set and keep existing (unsorted)",
		old: Env{
			{"TERM", []string{"xterm"}},
			{"HOME", []string{"/root"}},
		},
		new: Env{
			{"HOME", []string{"/root"}},
			{"PWD", []string{"/"}},
		},
		diff: Env{
			{"TERM", nil},
			{"PWD", []string{"/"}},
		},
	},
	{
		desc: "unset, overwrite and set",
		old: Env{
			{"TERM", []string{"xterm"}},
			{"HOME", []string{"/root"}},
		},
		new: Env{
			{"HOME", []string{"/home/user"}},
			{"PWD", []string{"/"}},
		},
		diff: Env{
			{"TERM", nil},
			{"HOME", []string{"/home/user"}},
			{"PWD", []string{"/"}},
		},
	},
	{
		desc: "unset, overwrite and keep existing (including empty)",
		old: Env{
			{"TERM", []string{"xterm"}},
			{"HOME", []string{"/root"}},
		},
		new: Env{
			{"HOME", []string{"/home/user"}},
			{"PWD", []string{"/"}},
			{"FOO", []string{}},
		},
		diff: Env{
			{"TERM", nil},
			{"HOME", []string{"/home/user"}},
			{"PWD", []string{"/"}},
			{"FOO", []string{}},
		},
	},
	{
		desc: "prepend value",
		old:  Env{{"PATH", []string{"/bin"}}},
		new:  Env{{"PATH", []string{"/sbin", "/bin"}}},
		diff: Env{{"PATH", []string{"/sbin", "$PATH"}}},
	},
	{
		desc: "prepend values",
		old:  Env{{"PATH", []string{"/sbin", "/bin"}}},
		new:  Env{{"PATH", []string{"/usr/sbin", "/usr/bin", "/sbin", "/bin"}}},
		diff: Env{{"PATH", []string{"/usr/sbin", "/usr/bin", "$PATH"}}},
	},
	{
		desc: "append value",
		old:  Env{{"PATH", []string{"/sbin"}}},
		new:  Env{{"PATH", []string{"/sbin", "/bin"}}},
		diff: Env{{"PATH", []string{"$PATH", "/bin"}}},
	},
	{
		desc: "append values",
		old:  Env{{"PATH", []string{"/usr/sbin", "/usr/bin"}}},
		new:  Env{{"PATH", []string{"/usr/sbin", "/usr/bin", "/sbin", "/bin"}}},
		diff: Env{{"PATH", []string{"$PATH", "/sbin", "/bin"}}},
	},
}

func TestDiff(t *testing.T) {
	for _, tt := range diffTests {
		if got := Diff(tt.old, tt.new); !envEqual(got, tt.diff) {
			t.Errorf("%s:\nDiff(%v, %v) = %v, want %v", tt.desc, tt.old, tt.new, got, tt.diff)
		}
	}
}

func envEqual(a, b Env) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name || !stringSliceEqual(a[i].Value, b[i].Value) {
			return false
		}
	}
	return true
}
