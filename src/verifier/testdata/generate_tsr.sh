#!/bin/bash

#url="http://timestamp.comodoca.com/rfc3161"
url="http://tsa.swisssign.net"

openssl ts -query -data $1 -sha256 -cert -no_nonce -out "$1_request.tsq"
cat "$1_request.tsq"|curl -s -S -H'Content-Type: application/timestamp-query' --data-binary @- ${url} -o "$1_response.tsr"
rm "$1_request.tsq"