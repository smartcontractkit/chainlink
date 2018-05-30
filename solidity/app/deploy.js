const Web3 = require('web3')

let compile = require('./compile.js')
let personal = require('./personal.js')
let utils = require('./utils.js')
let wallet = require('./wallet.js')
let web3 = new Web3('http://localhost:18545')

function getBytecode (contract) {
  return contract.evm.bytecode.object.toString()
}

function objectify (compiled, address) {
  let contract = new web3.eth.Contract(compiled.abi, address)
  contract.address = address
  return contract
}

module.exports = async function deploy (filename) {
  let compiled = compile(filename)

  await personal.send({to: wallet.address, value: utils.toWei(1)})
  let data = '0x' +
    getBytecode(compiled) +
    '0000000000000000000000004b274dfcd56656742A55ad54549b3770c392aA87'
  let txHash = await wallet.send({
    gas: 2000000,
    data: data
  })
  let receipt = await utils.getTxReceipt(txHash)
  return objectify(compiled, receipt.contractAddress)
}
