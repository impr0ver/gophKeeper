#!/bin/sh
# delete pem files
rm *.pem 

# Create CA private key and self-signed certificate
# adding -nodes to not encrypt the private key
openssl req -x509 -newkey rsa:4096 -nodes -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=RU/ST=EUROPE/L=MSK/O=DEV/OU=YandexPraktikum/CN=*.test.dev/emailAddress=ggg@test.com"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text 

# Create Web Server private key and CSR
# adding -nodes to not encrypt the private key
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=RU/ST=EUROPE/L=MSK/O=DEV/OU=YandexPraktikum/CN=*.test.dev/emailAddress=ggg@test.com"

# Sign the Web Server Certificate Request (CSR)
# !!!!WARNING!!!
#Remember that when we develop on localhost, it’s important to add (via -extfile server-ext.conf) the IP:0.0.0.0 or/and 127.0.0.1 as an Subject Alternative Name (SAN) extension to the certificate.
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.conf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text

# Verify certificate
echo "Verifying certificate"
openssl verify -CAfile ca-cert.pem server-cert.pem

# Generate client key for mutual TLS realisation, if need
#
# Generate client's private key and certificate signing request (CSR)
#openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/C=RU/ST=EUROPE/L=MSK/O=DEV/OU=CLIENT/CN=*.someclient.com/emailAddress=someclient@gmail.com"

#  Sign the Client Certificate Request (CSR)
#openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.conf

#echo "Client's signed certificate"
#openssl x509 -in client-cert.pem -noout -text




