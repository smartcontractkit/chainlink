const Eth = require('ethjs')
const TruffleContract = require('truffle-contract')

let compile = require('./compile.js')
let personal = require('./personal.js')
let wallet = require('./wallet.js')

function getBytecode (contract) {
  return contract.evm.bytecode.object.toString()
}

async function contractify (compiled, address) {
  let contract = TruffleContract({abi: compiled.abi, address: address})
  contract.setProvider(clUtils.provider)
  return contract.at(address)
}

module.exports = async function deploy (filename) {
  let compiled = compile(filename)

  await personal.send({to: wallet.address, value: clUtils.toWei(1)})
  let data = '0x' +
    getBytecode(compiled) +
    '0000000000000000000000004b274dfcd56656742A55ad54549b3770c392aA87'
  let txHash = await wallet.send({
    gas: 2000000,
    data: data
  })
  let receipt = await clUtils.getTxReceipt(txHash)
  return contractify(compiled, receipt.contractAddress)
}
