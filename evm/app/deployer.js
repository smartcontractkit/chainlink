const TruffleContract = require('truffle-contract')
const ABI = require('ethereumjs-abi')
const compile = require('./compile.js')

module.exports = function Deployer(wallet, utils) {
  function contractify(abi, address) {
    const contract = TruffleContract({
      abi: abi,
      address: address
    })
    contract.setProvider(utils.provider)
    contract.defaults({
      from: wallet.address,
      gas: 3500000,
      gasPrice: 10000000000
    })
    return contract.at(address)
  }

  function getBytecode(contract) {
    return contract.evm.bytecode.object.toString()
  }

  function findConstructor(abi) {
    for (let method of abi) {
      if (method.type === 'constructor') return method
    }
  }

  function constructorInputTypes(abi) {
    const types = []
    for (let input of findConstructor(abi).inputs) {
      types.push(input.type)
    }
    return types
  }

  function encodeArgs(unencoded, abi) {
    if (unencoded.length === 0) {
      return ''
    }
    const buf = ABI.rawEncode(constructorInputTypes(abi), unencoded)
    return buf.toString('hex')
  }

  return {
    perform: async function perform(filename, ...contractArgs) {
      const compiled = compile(filename)
      const encodedArgs = encodeArgs(contractArgs, compiled.abi)

      const txHash = await wallet.send({
        gas: 4000000,
        gasPrice: 10000000000,
        from: wallet.address,
        data: `0x${getBytecode(compiled)}${encodedArgs}`
      })
      const receipt = await utils.getTxReceipt(txHash)
      const contract = await contractify(compiled.abi, receipt.contractAddress)
      contract.transactionHash = txHash
      return contract
    },
    load: async (filename, address) => {
      const compiled = compile(filename)
      const contract = await contractify(compiled.abi, address)
      return contract
    }
  }
}
