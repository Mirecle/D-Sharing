#!/usr/bin/env bash
source base_parameters.sh

export CORE_PEER_MSPCONFIGPATH=./crypto-config/peerOrganizations/${PEER_DOMAIN}/users/Admin@${PEER_DOMAIN}/msp

peer=$1

export CORE_PEER_ADDRESS=${peer}:7051
#echo peer chaincode invoke -C ${CHANNEL} -n $CHAINCODE -c '{"Args":["transfer","account2", "account3", "20"]}'

#peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["transfer","account2", "account3", "20"]}'

#sleep 3


a="'{\"Args\":[\"querydata\",\"datahash000099\"]}'"
echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash

#a="'{\"Args\":[\"queryAllCars\"]}'"
#echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash



