let Web3 = require('web3');
let web3 = new Web3();
let x = module.exports = {}

x.fHash = function(signature) { return web3.utils.sha3(signature).slice(2,10) }
