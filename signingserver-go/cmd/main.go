package main

import (
	"github.com/emicklei/go-restful"
	"signingserver"
)



func main() {
	signingserver.CAInit()
	signingserver.OIDCInit()
	ws := new(restful.WebService)
	ws.Route(ws.POST("/hashes").Consumes("application/json").To(signingserver.PostHashes))
	ws.Route(ws.POST("/sign").Consumes("application/json").To(signingserver.SignHashes))
}

