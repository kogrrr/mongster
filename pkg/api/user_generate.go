// +build ignore

package main

import (
	"github.com/gargath/mongster/pkg/generator"
)

func main() {
	generator.GenerateEntityApi(generator.Options{
		EntityPackage: "github.com/gargath/mongster/pkg/entities",
	})

}
