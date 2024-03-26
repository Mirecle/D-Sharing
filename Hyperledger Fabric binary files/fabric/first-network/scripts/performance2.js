/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process')
const { performance } = require('perf_hooks');

// Parameters for the test
const totalTransactions = 10; // Total number of transactions to be sent
const parallelClients = 2; // Number of parallel transactions to simulate concurrent clients

async function ShareInitiation(contract,carNumber) {
    // const makeArray = ['Toyota', 'Ford', 'Honda', 'Hyundai', 'Tesla'];
    // const modelArray = ['Corolla', 'Fiesta', 'Civic', 'Elantra', 'Model S'];
    // const colorArray = ['Red', 'Blue', 'Black', 'White', 'Green'];
    // const ownerArray = ['Tom', 'Jerry', 'Mickey', 'Donald', 'Goofy'];

    // const make = makeArray[Math.floor(Math.random() * makeArray.length)];
    // const model = modelArray[Math.floor(Math.random() * modelArray.length)];
    // const color = colorArray[Math.floor(Math.random() * colorArray.length)];
    // const owner = ownerArray[Math.floor(Math.random() * ownerArray.length)];

    // await contract.submitTransaction('createCar', `CAR${carNumber}`, make, model, color, owner);
    await contract.submitTransaction('ShareInitiation','yzao', `datahash0000${carNumber}`, 'rootkey', 'sharekey', 'token00001', 'address00001', 'loc00001');
}
async function ShareDissemination(contract,carNumber) {
    await contract.submitTransaction('ShareDissemination',`SU${carNumber}`, 'datahash00001', 'rootkey', 'token00001', 'address00001', 'loc00001' , 'yzao');
}
async function ShareUpdate(contract,carNumber) {
    await contract.submitTransaction('ShareUpdate','yzao', `datahash0000${carNumber}`);
}
async function ShareRevocation(contract,carNumber) {
    await contract.submitTransaction('ShareRevocation',`datahash0000${carNumber}`, 'yzao');
}

async function main() {
    var ccpPath = path.resolve(__dirname, 'connection.json');
    const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
    const ccp = JSON.parse(ccpJSON);
    var endorserAddresses = execSync("bash printEndorsers.sh").toString().replace("\n", "").split(" ")
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
        const start = performance.now();
        const promises = [];
        for (let i = 0; i < totalTransactions; i++) {
            if (i % parallelClients === 0) {
                await Promise.all(promises);
                promises.length = 0;
            }
            promises.push(ShareDissemination(contract,i));
        }
        await Promise.all(promises);
        const end = performance.now();
        const duration = (end - start) / 1000; // Duration in seconds
        const throughput = totalTransactions / duration;
        console.log(`Total transactions: ${totalTransactions}`);
        console.log(`Duration: ${duration.toFixed(2)} seconds`);
        console.log(`Throughput: ${throughput.toFixed(2)} transactions per second`);
        //这里结束测试代码		

        // Disconnect from the gateway.	
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
