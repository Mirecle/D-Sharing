#!/usr/bin/env bash


export PEER_DOMAIN="org1.example.com"           # can be anything if running on localhost
export ORDERER_DOMAIN="example.com"        # can be anything if running on localhost

# fill in the addresses without domain suffix and without ports
export FAST_PEER_ADDRESS="localhost"
export ORDERER_ADDRESS="localhost"

# leave endorser address and storage address blank if you want to run on a single server
export ENDORSER_ADDRESS=()      # can define multiple addresses in the form ( "addr1" "addr2" ... )
export STORAGE_ADDRESS=""

export CHANNEL="channel123"               # the name of the created channel of the network
export CHAINCODE="chaincode"             # the name of the chaincode used on the channel
