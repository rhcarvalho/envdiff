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
