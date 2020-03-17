// +build ignore

package main

import (
	"github.com/shurcooL/vfsgen"
	"log"

	"github.com/gargath/mongoose/pkg/static"
)

func main() {
	err := vfsgen.Generate(static.Assets, vfsgen.Options{
		PackageName:  "static",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
