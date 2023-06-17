#!/bin/bash

N1_POST=150
N2_POST=100

N_POST=$(( ${N1_POST} + ${N2_POST} ))

BUILD_PATH="../build"

SERV_IP="sensors.com"
SERV_PORT=80

# Obtenemos el contador al momento
CONTADOR_BASE=$(curl -s --request GET --url ${SERV_IP}:${SERV_PORT}/api/processing/summary \
					-u admin:admin \
					| grep counter | tr -d -c 0-9)

# Realizamos los primeros N1 posts
for i in $(seq ${N1_POST}); do
    curl -X POST ${SERV_IP}:${SERV_PORT}/api/processing/submit -u admin:admin > /dev/null 2>&1
done

# Realizamos el get y obtenemos un primer contador
CONTADOR_1=$(curl -s --request GET --url ${SERV_IP}:${SERV_PORT}/api/processing/summary \
					-u admin:admin \
					| grep counter | tr -d -c 0-9)

# Realizamos los N2 posts
for i in $(seq ${N2_POST}); do
    curl -X POST ${SERV_IP}:${SERV_PORT}/api/processing/submit -u admin:admin  > /dev/null 2>&1
done

# Realizamos el get y obtenemos el total
CONTADOR_2=$(curl -s --request GET --url ${SERV_IP}:${SERV_PORT}/api/processing/summary \
					-u admin:admin \
					| grep counter | tr -d -c 0-9)

#echo "CONTADOR_BASE: ${CONTADOR_BASE}"
#echo "C1: ${CONTADOR_1}"
#echo "C2: ${CONTADOR_2}"

CONTADOR_1=$(expr ${CONTADOR_1} - ${CONTADOR_BASE} )
CONTADOR_2=$(expr ${CONTADOR_2} - ${CONTADOR_BASE} )

# Comparamos e imprimos el resultado
[[ ${CONTADOR_1} == ${N1_POST} ]] && [[ ${CONTADOR_2} == ${N_POST} ]] && echo "TEST PASSED" || echo "TEST FAILED"


