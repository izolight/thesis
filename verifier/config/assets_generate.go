// +build ignore

package main

import (
	"github.com/shurcooL/vfsgen"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier/config"
	"log"
)

//go:generate go run -tags=dev assets_generate.go
func main() {
	err := vfsgen.Generate(config.Assets, vfsgen.Options{
		PackageName: "config",
		BuildTags: "!dev",
		VariableName: "Assets",
	});
	if err != nil {
		log.Fatalln(err)
	}
}