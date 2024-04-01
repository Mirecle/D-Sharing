#!/usr/bin/env bash
source base_parameters.sh

export CORE_PEER_MSPCONFIGPATH=./crypto-config/peerOrganizations/${PEER_DOMAIN}/users/Admin@${PEER_DOMAIN}/msp

peer=$1

export CORE_PEER_ADDRESS=${peer}:7051
#echo peer chaincode invoke -C ${CHANNEL} -n $CHAINCODE -c '{"Args":["transfer","account2", "account3", "20"]}'

#peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["transfer","account2", "account3", "20"]}'

#sleep 3

#a="'{\"Args\":[\"query\",\"mprooft\"]}'"
#echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash
#a="'{\"Args\":[\"query\",\"mprooff\"]}'"
#echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash
#a="'{\"Args\":[\"query\",\"mproof\"]}'"
#echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash

#a="'{\"Args\":[\"query\",\"Initiation#datahash00001Owner#yzao\"]}'"
#echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash
#a="'{\"Args\":[\"query\",\"Initiation#datahash00001Owner#SU2\"]}'"
#echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash
a="'{\"Args\":[\"query\",\"account19999\"]}'"
echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash
