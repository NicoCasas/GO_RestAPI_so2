#!/bin/bash

# Imprime test_pass si la cantidad de usuarios devueltos por el endpoint es igual a la cantidad
# de usuarios leidos en /etc/passwd

# Login para obtener el jwt
TOKEN=$(curl -s --request POST \
             --url http://dashboard.com/api/users/login \
                -u admin:admin \
                --header 'accept: application/json' \
                --header 'content-type: application/json' \
                --data '{"username": "nicoAdmin", "password": "nicoPass"}' \
                | jq .token)

[[ ${TOKEN} != "" ]] || (echo TEST_FAILED && exit 1)

# Sacamos las comillas
TOKEN=$(echo ${TOKEN:1:-1})

# Con el jwt, hacemos GET a listall
RESULTADO=$(curl -s --request GET \
              --url http://dashboard.com/api/users/listall \
                -u admin:admin \
                --header 'accept: application/json' \
                --header 'content-type: application/json' \
                --header "Authentification: ${TOKEN}")

# Comparamos contra los usuarios del sistemas operativo
N_USERS_OS=$(wc -l /etc/passwd | awk '{print $1}')
N_USERS_REQUEST=$(echo ${RESULTADO} | jq '.data | length')

# Imprimimos el resultado del test
[[ ${N_USERS_OS} == ${N_USERS_REQUEST} ]] && echo TEST_PASSED || (echo TEST_FAILED && exit 1)


