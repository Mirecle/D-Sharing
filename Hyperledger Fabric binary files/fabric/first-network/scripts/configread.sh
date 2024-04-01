#!/usr/bin/env bash
source base_parameters.sh
export CORE_PEER_MSPCONFIGPATH=./crypto-config/peerOrganizations/${PEER_DOMAIN}/users/Admin@${PEER_DOMAIN}/msp
peer channel fetch config config_block.pb -c channel123
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq .data.data[0].payload.data.config.channel_group.groups.Orderer.values.BatchSize.value.max_message_count config_block.json >batchsize.txt
jq .data.data[0].payload.data.config config_block.json > config.json
cp config.json modified_config.json
