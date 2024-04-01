/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } =require( 'fabric-network');
const fs =require('fs');
const path = require('path');
const {execSync} = require('child_process')


async function main() {
    var endorserIdx = 0;
    var repetitions = 1;
    var transferCount = 10;
    var iteration = 0;


    var ccpPath = path.resolve(__dirname, 'connection.json');
    const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
    const ccp = JSON.parse(ccpJSON);
    var endorserAddresses =execSync("bash printEndorsers.sh").toString().replace("\n","").split(" ")
	process.on('unhandledRejection', error => {
 
});
process.on('PromiseRejectionHandledWarning', warn => {
 console.log('PromiseRejectionHandledWarning', warn.message);
});
    ccp.name = ccp.name.replace("CHANNEL", 'channel123')
    ccp.organizations.Org1.signedCert.path = ccp.organizations.Org1.signedCert.path.split("DOMAIN").join('org1.example.com')
    ccp.orderers.address.url = ccp.orderers.address.url.replace("ADDRESS", 'localhost')
    //ccp.peers.address.url = ccp.peers.address.url.replace("ADDRESS", endorserAddresses[endorserIdx])
    ccp.peers.address.url = ccp.peers.address.url.replace("ADDRESS", 'localhost')
    const user = "Admin@" + 'org1.example.com'
    console.log(`Thread no.: ${iteration},\tTx per process: ${transferCount / 2},\trepetitions: ${repetitions}`);
    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(__dirname, './wallet');
    const wallet = new FileSystemWallet(walletPath);
    try {

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists(user);
        if (!userExists) {
            console.log(`An identity for the user "${user}" does not exist in the wallet`);
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        var gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: user, discovery: { enabled: true } });

        // Get the network (channel) our contract is deployed to.
        var network = await gateway.getNetwork(String('channel123'));
        var contract = network.getContract(String('chaincode'))


//这里开始替换测试代码
        //await contract.submitTransaction('ShareInitiation','yzao', 'datahash00002', 'rootkey', 'sharekey', 'token00001', 'address00001', 'loc00001');
        //console.log('Transaction has been submitted');
        //const result = await contract.evaluateTransaction('querydata', 'datahash00002');
        //console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
		//await contract.submitTransaction('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom');
		//console.log('Transaction has been submitted');
 const result = await contract.evaluateTransaction('queryAllCars');
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
//这里结束测试代码		

        // Disconnect from the gateway.	
        await gateway.disconnect();
        console.log(`Thread ${iteration} is done!`);
    } catch (error) {
        console.error(`Thread ${iteration}: Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
