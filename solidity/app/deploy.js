const Eth = require('ethjs')
const TruffleContract = require('truffle-contract')

let compile = require('./compile.js')
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

  let fundingTx = await clUtils.send({
    to: wallet.address,
    value: clUtils.toWei(1)
  })
  await clUtils.getTxReceipt(fundingTx)
  let txHash = await wallet.send({
    gas: 2000000,
    data: '0x' + getBytecode(compiled)
  })
  let receipt = await clUtils.getTxReceipt(txHash)
  return contractify(compiled, receipt.contractAddress)
}
