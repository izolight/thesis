#!/bin/bash

#url="http://timestamp.comodoca.com/rfc3161"
url="http://tsa.swisssign.net"

openssl ts -query -data test.data -sha256 -cert -no_nonce -out request.tsq
cat request.tsq|curl -s -S -H'Content-Type: application/timestamp-query' --data-binary @- ${url} -o response.tsr
