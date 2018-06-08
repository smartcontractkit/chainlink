const TruffleContract = require('truffle-contract')
const ABI = require('ethereumjs-abi')
const compile = require('./compile.js')

const TruffleDefaults = {
  gas: 6721975,
  gasPrice: 100000000000
}

module.exports = function Deployer (wallet, utils) {
  this.perform = async function perform (filename, ...contractArgs) {
    const compiled = compile(filename)
    const encodedArgs = encodeArgs(contractArgs, compiled.abi)

    const txHash = await wallet.send({
      gas: 2500000,
      from: wallet.address,
      data: `0x${getBytecode(compiled)}${encodedArgs}`
    })
    const receipt = await utils.getTxReceipt(txHash)
    const contract = await contractify(compiled.abi, receipt.contractAddress)
    contract.transactionHash = txHash
    return contract
  }

  function contractify (abi, address) {
    const contract = TruffleContract({
      abi: abi,
      address: address
    })
    contract.setProvider(utils.provider)
    contract.defaults({
      from: wallet.address,
      gas: TruffleDefaults.gas,
      gasPrice: TruffleDefaults.gasPrice
    })
    return contract.at(address)
  }
}

function getBytecode (contract) {
  return contract.evm.bytecode.object.toString()
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
