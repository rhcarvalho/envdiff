// envdiff is a tool for comparing variables in two environments.
package main

import (
	"flag"
	"fmt"
	"os"
)

var output = flag.String("o", "", "set to 'docker' to print an ENV Dockerfile instruction")

func main() {
	flag.Parse()
	if flag.Arg(0) == "scl" {
		image := flag.Arg(1)
		if err := writeInstalledCollectionsEnvDiff(os.Stdout, image, *output == "docker"); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
		}
		return
	}
	old := flag.Arg(0)
	new := flag.Arg(1)
	fmt.Println(Diff(Parse(old), Parse(new)))
}
