// +build ignore

package main

import (
	"github.com/shurcooL/vfsgen"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier/static"
	"log"
)

//go:generate go run -tags=dev assets_generate.go
func main() {
	err := vfsgen.Generate(static.Assets, vfsgen.Options{
		PackageName: "static",
		BuildTags: "!dev",
		VariableName: "Assets",
	});
	if err != nil {
		log.Fatalln(err)
	}
}