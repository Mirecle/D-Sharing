/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process')

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms))
}
// async function ShareInitiation(contract, carNumber) {
//     await contract.submitTransaction('ShareInitiation', 'yzao', `datahash0000${carNumber}`, 'rootkey', 'sharekey', 'token00001', 'address00001', 'loc00001');
// }
// async function ShareDissemination(contract, carNumber) {
//     await contract.submitTransaction('ShareDissemination', `SU${carNumber}`, `datahash0000${carNumber}`, 'rootkey', 'token00001', 'address00001', 'loc00001', 'yzao');
// }
// async function ShareUpdate(contract, carNumber) {
//     await contract.submitTransaction('ShareUpdate', 'yzao', `datahash0000${carNumber}`);
// }
// async function ShareRevocation(contract, carNumber) {
//     await contract.submitTransaction('ShareRevocation', `datahash0000${carNumber}`, 'yzao');
// }

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
        const totalTransactions = 10000; // Total number of transactions to be sent
        var resp = new Array()
        console.log(`test begin!`)
        //开始测试记录开始时间节点
        var inittime = Date.now()
        for (var r = 0; r < totalTransactions; r++) {
            //var t = Promise.resolve(contract.submitTransaction('ShareInitiation', 'yzao', `datahash0000${r}`, 'rootkey', 'sharekey', 'token00001', 'address00001', 'loc00001'))
            //var t = Promise.resolve(contract.submitTransaction('ShareDissemination', `SU${r}`, `datahash0000${r}`, 'rootkey', 'token00001', 'address00001', 'loc00001', 'yzao'))
            //var t = Promise.resolve(contract.submitTransaction('ShareUpdate', 'yzao', `datahash0000${r}`))
            //var t = Promise.resolve(contract.submitTransaction('ShareRevocation', `datahash0000${r}`, 'yzao'))
              var t = Promise.resolve(contract.submitTransaction('transfer', `account${2*r}`, `account${2*r+1}`, '1'))
 			await sleep(5)	
            //替换其他三个方法
            //此处并发执行100或更多交易
            resp.push(t)
			console.log(`transaction${r} has been submitted`)
        }
        var i = totalTransactions - 1
	 //for (var i = 0; i < totalTransactions; i++) {
        try {
            await resp[i].then(function () { console.log(` transactions${i} end`); })
                .catch(function (error) { console.log(` transactions${i} end`); })
        	}
        catch (error) { }
	//}
        //等待最后一笔交易完成并打时间节点
        var endtime = Date.now()
        console.log(`The total cost is ${endtime - inittime} ms`)

        fs.appendFile('./transfertps.txt', "\n" + "The totalTransactions = " + totalTransactions, (err) => {
            if (err) {
                console.log(err);
            }
        })
        fs.appendFile('./transfertps.txt', " " + "The tps = " + totalTransactions / (endtime - inittime) * 1000, (err) => {
            if (err) {
                console.log(err);
            }
        })
        //这里结束测试代码		

        // Disconnect from the gateway.
        await sleep(5000)	
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
