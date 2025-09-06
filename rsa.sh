#!/bin/bash

mkdir -p cert
openssl genrsa -out cert/ca.key 4096
openssl req -new -x509 -days 365 -key cert/ca.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out cert/ca.crt 
openssl req -newkey rsa:4096 -nodes -keyout cert/server.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=localhost" -out cert/server.csr
openssl x509 -req -extfile <(printf "subjectAltName=DNS:localhost") -days 365 -in cert/server.csr -CA cert/ca.crt -CAkey cert/ca.key -CAcreateserial -out cert/server.crt