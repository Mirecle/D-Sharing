/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

'use strict';

module.exports.info  = 'opening accounts';

let bc, contx;
module.exports.init = function(blockchain, context, args) {
    bc = blockchain;
    contx = context;

    return Promise.resolve();
};



module.exports.run = function() {
    const hash = new Date().getTime().toString();
    let args = {
        chaincodeFunction:'ShareInitiation',
        chaincodeArguments: ['yzao', hash, 'rootkey','sharekey','token00001','address00001','loc00001'],
    };
    return bc.invokeSmartContract(contx, 'simple', 'v0', args, 100);
};

module.exports.end = function() {
    return Promise.resolve();
};

