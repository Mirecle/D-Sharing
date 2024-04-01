#  Transparent-Kites

---

This is the implementation of paper "Transparent-Kites: a blockchain-based dissemination control framework for anonymous data sharing" by Ming Zhang, Hui Li*, Zihao Yang, Chao Qu
*: This author is the corresponding author.
Note that this is a work in progress, and the code here is not the final version of the original paper.

---

## Abstract

---

The growing awareness of privacy-preserving drives the development of anonymous sharing. By striping the user’s identity from the sharing, anonymous sharing provides greater privacy protection, yet it also raises the challenge of ownership conservatory in wide dissemination. In this paper, we propose D-Sharing, a blockchain-based dissemination control framework that provides strong ownership protection for anonymity sharing. In contrast to conventional mechanisms that offer dissemination control at the cost of server-side anonymity, our framework achieves complete anonymity where any sensitive information related to user identity will not be handed over throughout the whole sharing process. To overcome the challenge of maintaining ownership across multi-tier peer-to-peer anonymous sharing, we propose a Fast Multinomial Hash Key Assignment (FMHKA). FMHKA computes a unique shared key for each sharing, em-powering prior sharers with higher-level keys to efficiently derive subsequent keys and manage ensuing dissemination. Utilizing a carefully designed anonymous dissemination tree, D-Sharing tracks and maintains complete data dissemination status on the blockchain, providing real-time control of subsequent data dissemination for data owners. To achieve anonymous authen tication of data sharing, we also designed a zero-knowledge proof-based ownership authentication mechanism to verify users without identity disclosure. We conducted extensive real-world simulations, and our experimental results demonstrate the ca-pability and effectiveness of the framework across a number of performance metrics.  

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

We will go through the details in the following sections.
## 1.Smart Contract
In this section, we provide a detailed implementation of the control framework for anonymous data sharing in Go language. The framework consists of four main parts: data upload, data authorization, permission update, and permission revocation.

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
**Execution of commands**
The client initiates a transaction to make changes to the ledger data, and the endorsed transaction needs to be sent to the sorting node on the chain. Therefore, it is necessary to enable TLS verification and specify the corresponding orderer certificate path.

1. ShareInitiation

Require: Call method name,owner,data_Hash,root_Key,share_Key ,token,address ,location
```shell
peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["ShareInitiation","Acc1", "datahash00001", "rkeyV1","skeyV1","token00001","address00001","loc00001"]}'
```
```shell
-o: Specify the orderer node address
--Tls: Enable TLS verification
--Cafile: specifies the root certificate path of the Orderer, used to verify TLS handshake
-n: Specify chain code name
-C: Specify channel name
-c: Specify the required parameters for calling the chain code
```

2. ShareDissemination

Require: Call method name,owner,data_Hash,root_Key,token,address,loc,Original_owner
```shell
peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["ShareDissemination","SU2", "datahash00001", "rootkey","token00002","address00002","loc00002","Acc1"]}'
```

3. ShareUpdate

Require: Call method name,Original_owner,data_Hash
```shell
peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["ShareUpdate","Acc1", "datahash00001"]}'
```

4. ShareRevocation

Require: Call method name,Original_owner,data_Hash
```shell
peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["ShareRevocation", "datahash00001", "Acc1"]}'
```
## 2.Zero knowledge proof for three operating systems
### 2.1 Windows
The zero knowledge proof code on the Windows side is entirely written in Golang language and can run normally on operating systems with a go environment.
### 2.2 Android
The main method for developing Android apps in Go language is to use the Golang mobile binding tool library Gomobile. Gomobile compiles Go into a shared library compatible with Android modules, seamlessly integrating with Java and Kotlin code, and running on Android.
### 2.3 Raspberry Pi
Raspberry Pi generally uses Linux system, so the Go language environment needs to be installed on the Raspberry Pi end to run the code.
## 3.Hyperledger Fabric binary files
We will integrate the zero knowledge proof function into the underlying source code of the fabric to extend its peer binary file functionality. There are two ways to use code modification: 1. Download the source code and compile binary files; 2. Download Docker images
### 3.1 Download the source code and compile binary files
**Compile binary files using the following command in the fabric directory**
```shell
make release
```
**After compilation, the binary files can be viewed in the following directory：**
```shell
/opt/gopath/src/github.com/hyperledger/fabric1/release/linux-amd64/bin
```
**Add binary files to Docker images**
```shell
docker pull hyperledger/fabric-peer:1.4.3
docker run -it hyperledger/fabric-peer:1.4.3 /bin/bash
```
**Start a new terminal and copy the peer to the image**
```shell
docker cp ./bin/peer ${container_name}:/usr/local/bin
```
**Submit image**
```shell
docker commit -a "${author_name}" -m "${explain}" ${containerID} ${REPOSITORY} :${TAG}
```
```shell
-a: Create Author
-M: Explanation
ContainerID: Depends on which container to create a new image (generated by the Docker Run command for Dockers PS viewing)
REPOSITORY: Mirror name
TAG: Image label, Docker Compose will find the image based on the label
```
**Usage**
Usually, the image name and version used in the image parameter settings in Docker Compose.yaml are set
### 3.2 Download Docker images
Similarly, we will also upload the compiled binary files to the Docker official website, which can be pulled using the following command
```shell
docker pull yangzao/zpeer:1.4.0
```
Afterwards, using the same method as above, set the image to be used in the .yaml file

Before using the caliper for testing, the following configuration files must be modified to use an image file that contains zero knowledge proof.
```shell
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

version: '2'

services:
  ca.org1.example.com:
    image: hyperledger/fabric-ca:${FABRIC_VERSION}
    environment:
      - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME=ca.org1.example.com
    ports:
      - "7054:7054"
    command: sh -c 'fabric-ca-server start --ca.certfile /etc/hyperledger/fabric-ca-server-config/ca.org1.example.com-cert.pem --ca.keyfile /etc/hyperledger/fabric-ca-server-config/key.pem -b admin:adminpw -d'
    volumes:
      - ../../config_solo/crypto-config/peerOrganizations/org1.example.com/ca/:/etc/hyperledger/fabric-ca-server-config
    container_name: ca.org1.example.com

  ca.org2.example.com:
    image: hyperledger/fabric-ca:${FABRIC_VERSION}
    environment:
      - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME=ca.org2.example.com
    ports:
      - "8054:7054"
    command: sh -c 'fabric-ca-server start --ca.certfile /etc/hyperledger/fabric-ca-server-config/ca.org2.example.com-cert.pem --ca.keyfile /etc/hyperledger/fabric-ca-server-config/key.pem -b admin:adminpw -d'
    volumes:
      - ../../config_solo/crypto-config/peerOrganizations/org2.example.com/ca/:/etc/hyperledger/fabric-ca-server-config
    container_name: ca.org2.example.com

  orderer.example.com:
    container_name: orderer.example.com
    image: hyperledger/fabric-orderer:${FABRIC_VERSION}
    environment:
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/etc/hyperledger/configtx/orgs.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/etc/hyperledger/msp/orderer/msp
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    ports:
      - 7050:7050
    volumes:
        - ../../config_solo/:/etc/hyperledger/configtx
        - ../../config_solo/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp:/etc/hyperledger/msp/orderer/msp

  peer0.org1.example.com:
    container_name: peer0.org1.example.com
    image: yangzao/zpeer:1.4.0
    environment:
      - CORE_LOGGING_PEER=debug
      - CORE_CHAINCODE_LOGGING_LEVEL=DEBUG
      - CORE_CHAINCODE_EXECUTETIMEOUT=999999
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_PEER_ID=peer0.org1.example.com
      - CORE_PEER_ENDORSER_ENABLED=true
      - CORE_PEER_ADDRESS=peer0.org1.example.com:7051
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=2org1peergoleveldb_default
      - CORE_PEER_LOCALMSPID=Org1MSP
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/msp
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_ORGLEADER=false
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org1.example.com:7051
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: peer node start
    ports:
      - 7051:7051
      - 7053:7053
    volumes:
        - /var/run/:/host/var/run/
        - ../../config_solo/mychannel.tx:/etc/hyperledger/configtx/mychannel.tx
        - ../../config_solo/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/etc/hyperledger/peer/msp
        - ../../config_solo/crypto-config/peerOrganizations/org1.example.com/users:/etc/hyperledger/msp/users
    depends_on:
      - orderer.example.com

  peer0.org2.example.com:
    container_name: peer0.org2.example.com
    image: yangzao/zpeer:1.4.0
    environment:
      - CORE_LOGGING_PEER=debug
      - CORE_CHAINCODE_LOGGING_LEVEL=DEBUG
      - CORE_CHAINCODE_EXECUTETIMEOUT=999999
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_PEER_ID=peer0.org2.example.com
      - CORE_PEER_ENDORSER_ENABLED=true
      - CORE_PEER_ADDRESS=peer0.org2.example.com:7051
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=2org1peergoleveldb_default
      - CORE_PEER_LOCALMSPID=Org2MSP
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/msp
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_ORGLEADER=false
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org2.example.com:7051
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: peer node start
    ports:
      - 8051:7051
      - 8053:7053
    volumes:
        - /var/run/:/host/var/run/
        - ../../config_solo/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/msp:/etc/hyperledger/peer/msp
        - ../../config_solo/crypto-config/peerOrganizations/org2.example.com/users:/etc/hyperledger/msp/users
    depends_on:
      - orderer.example.com
```
## 4. Performance Assessment
We used the HyperLedger caliper tool to evaluate the main functions of the system, including system throughput, maximum latency, and average latency.
The specific results can refer to the experiments in the paper or download test data.
