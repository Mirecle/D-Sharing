#  D-Sharing

---

This is the implementation of paper "D-Sharing: A Blockchain-based Dissemination Control Framework for Privacy-Preserving Data Sharing" by Ming Zhang, Hui Li, Zihao Yang, Chao Qu

Note that this is a work in progress, and the code here is not the final version of the original paper.

---

## Abstract

---

The growing awareness of privacy-preserving drives the development of anonymous sharing. By removing user identities, anonymous sharing offers greater privacy protection, yet also introduces the control challenge for data dissemination. In this paper, we propose D-Sharing, a blockchain-based dissemination control framework that provides both advanced anonymity and efficient dissemination control. Unlike traditional mechanisms compromising server-side anonymity, our framework eliminates reliance on trusted entities, ensuring full anonymity throughout the data dissemination lifecycle. In response to the authorization challenge for anonymous sharing, we introduce the Decentralized Multinomial Hash Key Assignment (DMHKA). DMHKA computes a unique shared key for each sharing instance, empowering previous sharers to derive subsequent keys locally and thus manage the authorization for arbitrary anonymous disseminators. With a carefully designed anonymous dissemination tree, D-Sharing tracks and maintains dissemination status on the blockchain, achieving real-time control for data providers. Furthermore, we design an efficient authentication protocol (zk-AnonyAuth) based on zero-knowledge proof (ZKP) to provide a verification-optimized solution for anonymous authentication. Experimental results from real-world datasets demonstrate D-Sharing's robust anonymity against identity reconstruction and relationship analysis (evidenced by high information entropy and adversarial graph-based deep learning experiments), good compatibility (minimum 0.0013 ms for key generation) for various computing-constrained clients, and high concurrency (over 1660 tps under 0.06 s latency) for extensive sharing demands.

---

## Platform

- HyperLedger Fabric
- Hyperledger Caliper
- Docker
- CentOS
- Ubuntu
## Codes

---

This project code consists of four parts:

1. Smart contracts
2. Zero knowledge proof for three operating systems
3. Hyperledger Fabric binary files
4. Performance Assessment

In this section, a general framework will be given, and more implementation details will be available after the acceptance of the paper.
## 1.Smart Contract
We implement the control framework for anonymous data sharing in Go language. The framework consists of four main parts: data upload, data authorization, permission update, and permission revocation.

**Preparation phase**
We use the examples provided by Fabric's Fabric samples to illustrate and modify the following code in the first network/scripts/script.sh script:
```shell
if [ "${NO_CHAINCODE}" != "false" ]; then
  # Install chaincode on peer0.org1 and peer0.org2
  echo "Installing chaincode on peer0.org1..."
  installChaincode 0 1       
  echo "Install chaincode on peer0.org2..."
  installChaincode 0 2       
  #...
  #...
fi
```
Before starting the fabric network, the following configuration files must be modified to use the compiled images modified in the project. For specific instructions on how to add images, please refer to section 3.
```shell
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

version: '2'

services:
  peer-base:
    image: yangzao/zpeer:1.4.0
    environment:
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      # the following setting starts chaincode containers on the same
      # bridge network as the peers
      # https://docs.docker.com/compose/networking/
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=${COMPOSE_PROJECT_NAME}_byfn
      - FABRIC_LOGGING_SPEC=INFO
      #- FABRIC_LOGGING_SPEC=DEBUG
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_ORGLEADER=false
      - CORE_PEER_PROFILE_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start

  orderer-base:
    image: hyperledger/fabric-orderer:$IMAGE_TAG
    environment:
      - FABRIC_LOGGING_SPEC=INFO
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      # enabled TLS
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
      - ORDERER_KAFKA_TOPIC_REPLICATIONFACTOR=1
      - ORDERER_KAFKA_VERBOSE=true
      - ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer

```
**Start network**
```shell
$ ./byfn.sh up
```
**Enter the CLI client container**
The CLI client defaults to connecting to the peer0.org1 node as Admin.org1:
```shell
$ docker exec -it cli bash
```
**Check the current node to join which channels**
```shell
# peer channel list
```
**Installation Chain Code**
```shell
peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/zeroproof/go/
```
**Set environment variables for channel names**
```shell
export CHANNEL_NAME=mychannel
echo $CHANNEL_NAME
```
**Set the environment variable for the certificate path of the orderer node**
```shell
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
echo $ORDERER_CA
```
**Instantiating Chain Code Using the Instantiate Command**
```shell
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -v 1.0 -c '{"Args":["init"," "]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
```
```shell
-o: Specify the Oderer service node address
--Tls: Enable TLS verification
--Cafile: specifies the root certificate path of the Orderer, used to verify TLS handshake
-n: Specify the chain code name to be instantiated, which must be the same as the chain code name specified during installation
-v: Specify the version number of the chain code to be instantiated, which must be the same as the chain code version number specified during installation
-C: Specify channel name
-c: Parameters specified when instantiating chain codes
-P: Specify endorsement strategy
```

## 2.Zero knowledge proof for three operating systems
### 2.1 Windows
The zero knowledge proof code on the Windows side is entirely written in Golang language and can run normally on operating systems with a go environment.
### 2.2 Android
The main method for developing Android apps in Go language is to use the Golang mobile binding tool library Gomobile. Gomobile compiles Go into a shared library compatible with Android modules, seamlessly integrating with Java and Kotlin code, and running on Android.
### 2.3 Raspberry Pi
Raspberry Pi generally uses Linux system, so the Go language environment needs to be installed on the Raspberry Pi end to run the code.

## 3. Performance Assessment
We used the HyperLedger caliper tool to evaluate the main functions of the system, including system throughput, maximum latency, and average latency.
The specific results can refer to the experiments in the paper or download test data.
