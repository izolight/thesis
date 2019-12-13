#!/bin/bash

ocsp_uri=$(openssl x509 -noout -ocsp_uri -in "$1")

openssl ocsp -issuer "$2" -cert "$1" -url "$ocsp_uri" -respout "${1}.ocsp" -noverify -text
 
