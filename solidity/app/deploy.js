const Eth = require('ethjs')
const TruffleContract = require('truffle-contract')
const ABI = require('ethereumjs-abi')

const clUtils = require('./cl_utils.js')
const clWallet = require('./cl_wallet.js')
const compile = require('./compile.js')

const TruffleDefaults = {
  gas: 6721975,
  gasPrice: 100000000000
}

function getBytecode (contract) {
  return contract.evm.bytecode.object.toString()
}

function contractify (abi, address) {
  const contract = TruffleContract({
    abi: abi,
    address: address
  })
  contract.setProvider(clUtils.provider)
  contract.defaults({
    from: clWallet.address,
    gas: TruffleDefaults.gas,
    gasPrice: TruffleDefaults.gasPrice
  })
  return contract.at(address)
}

function findConstructor (abi) {
  for (let method of abi) {
    if (method.type === 'constructor') return method
  }
}

function constructorInputTypes (abi) {
  const types = []
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

module.exports = async function deploy (filename, ...contractArgs) {
  const compiled = compile(filename)
  const encodedArgs = encodeArgs(contractArgs, compiled.abi)

  const fundingTx = await clUtils.send({
    to: clWallet.address,
    value: clUtils.toWei(1)
  })
  await clUtils.getTxReceipt(fundingTx)
  const txHash = await clWallet.send({
    gas: 2500000,
    data: `0x${getBytecode(compiled)}${encodedArgs}`
  })
  const receipt = await clUtils.getTxReceipt(txHash)
  const contract = await contractify(compiled.abi, receipt.contractAddress)
  contract.transactionHash = txHash
  return contract
}
