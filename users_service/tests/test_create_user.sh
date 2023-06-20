#!/bin/bash

# Imprime TEST_PASSED si logra crear el usuario en el sistema operativo a través de los endpoints login
# (para obtener el jwt) y createuser

# Flag para determinar si borrar el usuario creado al terminar el test
DEL_USER_FLAG=1

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

# Con el jwt, creamos el usuario
RESULTADO=$(curl -s --request POST \
              --url http://dashboard.com/api/users/createuser \
                -u admin:admin \
                --header 'accept: application/json' \
                --header 'content-type: application/json' \
                --header "Authentification: ${TOKEN}" \
                --data '{"username": "newUser", "password": "newUserPass"}' )



echo ${RESULTADO} | grep >/dev/null "username" && (cat /etc/passwd | grep >/dev/null newUser) && echo TEST_PASSED || echo TEST_FAILED

# Lo eliminamos
[[ DEL_USER_FLAG == 1 ]] && sudo deluser >/dev/null 2>&1 newUser 


