#!/bin/bash

# Imprime TEST_PASSED si la respuesta a la request no contiene el token

TOKEN=$(curl -s --request POST \
             --url http://dashboard.com/api/users/login \
                -u admin:admin \
                --header 'accept: application/json' \
                --header 'content-type: application/json' \
                --data '{"username": "nicoAdmin", "password": "invalidPass"}' \
                | jq >/dev/null .token)

[[ !$? ]] && echo TEST_PASSED || echo TEST_FAILED
