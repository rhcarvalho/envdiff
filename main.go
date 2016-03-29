// envdiff is a tool for comparing variables in two environments.
package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	old := flag.Arg(0)
	new := flag.Arg(1)
	fmt.Println(Diff(Parse(old), Parse(new)))
}
