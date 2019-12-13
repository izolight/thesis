module gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier

go 1.13

require (
	github.com/coreos/go-oidc v2.1.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/mux v1.7.3
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0 // indirect
	go.mozilla.org/pkcs7 v0.0.0-20180702141046-401d3877331b
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/square/go-jose.v2 v2.4.0
)

replace go.mozilla.org/pkcs7 v0.0.0-20180702141046-401d3877331b => github.com/izolight/pkcs7 v0.0.0-20191208070232-70a1b50704d6
