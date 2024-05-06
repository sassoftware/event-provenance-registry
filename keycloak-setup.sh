#!/bin/bash

REALM=junk

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

# need a client-id
# need a client secret

#curl -H "Authorization: bearer $access_token" -H "content-type: application/json" http://localhost:8083/admin/realms/junk/users

#docker run --user=1000 --env=KC_HTTP_PORT=8083 --env=KEYCLOAK_ADMIN=admin --env=KEYCLOAK_ADMIN_PASSWORD=admin --env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin --env=LANG=en_US.UTF-8 --env=KC_RUN_IN_CONTAINER=true --env=KC_LOG_LEVEL=DEBUG -p 8083:8083 --restart=no --runtime=runc quay.io/keycloak/keycloak:24.0.1 start-dev