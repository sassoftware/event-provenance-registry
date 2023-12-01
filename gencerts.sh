#!/bin/bash
set -x

SUBJECT='/C=US/ST=North Carolina/L=Cary/O=Event Provenance Registry/CN=localhost'

mkdir -p certs
pushd certs || exit;
    openssl req -newkey rsa:4096 -new -nodes -x509 -days 3650 -out cert.pem -keyout key.pem -subj "${SUBJECT}"
popd || exit