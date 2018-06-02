const Eth = require('ethjs')
const TruffleContract = require('truffle-contract')
const ABI = require('ethereumjs-abi')

require('./cl_utils.js')
require('./cl_wallet.js')
const compile = require('./compile.js')

function getBytecode (contract) {
  return contract.evm.bytecode.object.toString()
}

function contractify (abi, address) {
  let contract = TruffleContract({
    abi: abi,
    address: address
  })
  contract.setProvider(clUtils.provider)
  contract.defaults({
    from: clWallet.address,
    gas: 6721975, // Truffle default
    gasPrice: 100000000000 // Truffle default
  })
  return contract.at(address)
}

function findConstructor (abi) {
  for (let method of abi) {
    if (method.type === 'constructor') return method
  }
}

function constructorInputTypes (abi) {
  let types = []
  for (let input of findConstructor(abi).inputs) {
    types.push(input.type)
  }
  return types
}

function encodeArgs (unencoded, abi) {
  if (unencoded.length === 0) {
    return ''
  }
  const buf = ABI.rawEncode(constructorInputTypes(abi), unencoded)
  return buf.toString('hex')
}

module.exports = async function deploy (filename) {
  const unencodedArgs = Array.prototype.slice.call(arguments).slice(1)
  const compiled = compile(filename)
  const encodedArgs = encodeArgs(unencodedArgs, compiled.abi)

  const fundingTx = await clUtils.send({
    to: clWallet.address,
    value: clUtils.toWei(1)
  })
  await clUtils.getTxReceipt(fundingTx)
  let txHash = await clWallet.send({
    gas: 2500000,
    data: `0x${getBytecode(compiled)}${encodedArgs}`
  })
  const receipt = await clUtils.getTxReceipt(txHash)
  let contract = await contractify(compiled.abi, receipt.contractAddress)
  contract.transactionHash = txHash
  return contract
}
