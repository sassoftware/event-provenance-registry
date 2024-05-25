#!/bin/bash
# This script must be run somewhere with access to the kcadm.sh script that invokes the Keycloak Admin CLI.
REALM=junk

kcadm.sh config credentials --server http://localhost:8083 --realm master --user admin --password admin
kcadm.sh create realms -s realm=$REALM -s enabled=true -o
kcadm.sh create users -r $REALM -s username=bob -s email=bob@example.com -s enabled=true -o --fields id,username
kcadm.sh create clients -r $REALM -s clientId=epr-client -o
kcadm.sh create client-scopes -r $REALM -s name=epr-client-scope #c7683fc0-8e21-4c22-b1ee-9ac27f3ac135 TODO: need more than this, otherwise we hit an error trying to view the scope in the UI.
kcadm.sh create client-scopes/c7683fc0-8e21-4c22-b1ee-9ac27f3ac135/mappers -r $REALM -o
# TODO: we can read from json files with the CLI which may make more sense long term.

access_token=$( curl -d "client_id=admin-cli" -d "username=admin" -d "password=admin" -d "grant_type=password" "http://localhost:8083/realms/master/protocol/openid-connect/token" | sed -n 's|.*"access_token":"\([^"]*\)".*|\1|p')

# Post a new realm. The docs are broke, so I had to wing it by posting everything.

curl -H "Authorization: bearer $access_token" -H "content-type: application/json" -X POST \
  -d '@realm.json' \
  http://localhost:8083/admin/realms

# need a user

curl -H "Authorization: bearer $access_token" -H "content-type: application/json" -X POST \
  -d '{
  "username": "bob",
  "email": "bob@example.com",
  "emailVerified": true,
  "enabled": true,
  "firstName": "Bob",
  "lastName": "Testy",
  "credentials": [
    {
      "temporary": false,
      "type": "password",
      "value": "abc123"
    }
  ]
}' http://localhost:8083/admin/realms/junk/users

# need a client

curl -H "Authorization: bearer $access_token" -H "content-type: application/json" \
  http://localhost:8083/admin/realms/junk/clients

curl -H "Authorization: bearer $access_token" -H "content-type: application/json" -X POST \
  -d '{
  "clientId": "dummy-client",
  "name": "my-dummy-client",
  "directAccessGrantsEnabled": true
}' http://localhost:8083/admin/realms/junk/clients

# TODO: need the audience mapper and audience resolve

curl -H "Authorization: bearer $access_token" -H "content-type: application/json" -X POST \
  -d '{
  "name": "epr-client-id-dedicated"
}' http://localhost:8083/admin/realms/junk/client-scopes

# need a client secret
curl -H "Authorization: bearer $access_token" -H "content-type: application/json" -X POST \
  -d '{
  "name": "epr-client-id-dedicated"
}' http://localhost:8083/admin/realms/junk/client-scopes/9940fe9c-005b-4d52-b25f-6865bbd9d60e/scope-mappings/clients/{asdf}


#curl -H "Authorization: bearer $access_token" -H "content-type: application/json" http://localhost:8083/admin/realms/junk/users

#docker run --user=1000 --env=KC_HTTP_PORT=8083 --env=KEYCLOAK_ADMIN=admin --env=KEYCLOAK_ADMIN_PASSWORD=admin --env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin --env=LANG=en_US.UTF-8 --env=KC_RUN_IN_CONTAINER=true --env=KC_LOG_LEVEL=DEBUG -p 8083:8083 --restart=no --runtime=runc quay.io/keycloak/keycloak:24.0.1 start-dev