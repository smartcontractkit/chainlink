const Eth = require('ethjs')
const TruffleContract = require('truffle-contract')
const ABI = require('ethereumjs-abi')

let compile = require('./compile.js')
let wallet = require('./wallet.js')

function getBytecode (contract) {
  return contract.evm.bytecode.object.toString()
}

function contractify (abi, address) {
  let contract = TruffleContract({abi: abi, address: address})
  contract.setProvider(clUtils.provider)
  return contract.at(address)
}

function findConstructor(abi) {
  for (let method of abi) {
    if (method.type === 'constructor') return method
  }
}

function constructorInputTypes(abi) {
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
    to: wallet.address,
    value: clUtils.toWei(1)
  })
  await clUtils.getTxReceipt(fundingTx)
  let txHash = await wallet.send({
    gas: 2000000,
    data: `0x${getBytecode(compiled)}${encodedArgs}`
  })
  let receipt = await clUtils.getTxReceipt(txHash)
  return contractify(compiled.abi, receipt.contractAddress)
}
