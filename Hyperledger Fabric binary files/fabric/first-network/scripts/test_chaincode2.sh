#!/usr/bin/env bash
source base_parameters.sh

export CORE_PEER_MSPCONFIGPATH=./crypto-config/peerOrganizations/${PEER_DOMAIN}/users/Admin@${PEER_DOMAIN}/msp

peer=$1

export CORE_PEER_ADDRESS=${peer}:7051
peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["ShareInitiation","yzao", "datahash00001", "rootkey","sharekey","token00001","address00001","loc00001"]}'

sleep 5

peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["ShareDissemination","SU2", "datahash00001", "rootkey","token00002","address00002","loc00002","yzao"]}'
peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["ShareDissemination","SU3", "datahash00001", "rootkey","token00002","address00002","loc00002","yzao"]}'
sleep 5

peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["ShareUpdate","yzao", "datahash00001"]}'

sleep 5

peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["ShareRevocation", "datahash00001", "yzao"]}'

sleep 5
