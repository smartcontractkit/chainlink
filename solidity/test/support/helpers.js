process.env.SOLIDITY_INCLUDE = '../../solidity/contracts/:../../solidity/contracts/examples/:../../contracts/:../../contracts/lib/:../../node_modules/:../../node_modules/link_token/contracts:../../node_modules/openzeppelin-solidity/contracts/ownership/:../../node_modules/@ensdomains/ens/contracts/'

const PRIVATE_KEY = 'c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3'

const Wallet = require('../../app/wallet.js')
const Utils = require('../../app/utils.js')
const Deployer = require('../../app/deployer.js')

const abi = require('ethereumjs-abi')
const util = require('ethereumjs-util')
const BN = require('bn.js')

const HEX_BASE = 16

let consumer, defaultAccount, eth, oracleNode, stranger
let deployer, utils, wallet

(() => {
  eth = web3.eth

  before(async function () {
    // Default hard coded truffle accounts:
    // ==================
    // (0) 0x627306090abab3a6e1400e9345bc60c78a8bef57
    // (1) 0xf17f52151ebef6c7334fad080c5704d77216b732
    // (2) 0xc5fdf4076b8f3a5357c5e395ab970b5b54098fef
    // (3) 0x821aea9a577a9b44299b9c15c88cf3087f3b5544
    // (4) 0x0d1d4e623d10f9fba5db95830f7d3839406c6af2
    // (5) 0x2932b7a2355d6fecc4b5c0b6bd44cc31df247a2e
    // (6) 0x2191ef87e392377ec08e7c08eb105ef5448eced5
    // (7) 0x0f4f2ac550a1b4e2280d04c21cea7ebd822934b5
    // (8) 0x6330a553fc93768f612722bb8c2ec78ac90b3bbc
    // (9) 0x5aeda56215b167893e80b4fe645ba6d5bab767de

    // Private Keys
    // ==================
    // (0) c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3
    // (1) ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f
    // (2) 0dbbe8e4ae425a6d2687f1a7e3ba17bc98c673636790f1b8ad91193c05875ef1
    // (3) c88b703fb08cbea894b6aeff5a544fb92e78a18e19814cd85da83b71f772aa6c
    // (4) 388c684f0ba1ef5017716adb5d21a053ea8e90277d0868337519f97bede61418
    // (5) 659cbb0e2411a44db63778987b1e22153c086a95eb6b18bdf89de078917abc63
    // (6) 82d052c865f5763aad42add438569276c00d3d88a2d062d36b2bae914d58b8c8
    // (7) aa3680d5d48a8283413f7a108367c7299ca73f553735860a87b08f39395618b7
    // (8) 0f62d96d6675f32685bbdb8ac13cda7c23436f63efbb9d07700d8669ff12b7c4
    // (9) 8d5366123cb560bb606379f90a0bfd4769eecc0557f1b362dcae9012b548b1e5

    // HD Wallet
    // ==================
    // Mnemonic:      candy maple cake sugar pudding cream honey rich smooth crumble sweet treat
    // Base HD Path:  m/44'/60'/0'/0/{account_index}
    const accounts = eth.accounts

    defaultAccount = accounts[0]
    oracleNode = accounts[1]
    stranger = accounts[2]
    consumer = accounts[3]

    utils = Utils(web3.currentProvider)
    wallet = Wallet(PRIVATE_KEY, utils)
    deployer = Deployer(wallet, utils)
  })
})()

export {
  consumer,
  defaultAccount,
  eth,
  oracleNode,
  stranger
}

export const bigNum = number => web3.toBigNumber(number)

export const toWei = number => bigNum(web3.toWei(number))

export const hexToInt = string => web3.toBigNumber(string)

export const toHexWithoutPrefix = arg => {
  if (arg instanceof Buffer || arg instanceof BN) {
    return arg.toString('hex')
  } else if (arg instanceof Uint8Array) {
    return Array.prototype.reduce.call(arg, (a, v) => a + v.toString('16').padStart(2, '0'), '')
  } else {
    return Buffer.from(arg, 'ascii').toString('hex')
  }
}

export const toHex = value => {
  return `0x${toHexWithoutPrefix(value)}`
}

export const deploy = (filePath, ...args) => deployer.perform(filePath, ...args)

export const getEvents = contract => (
  new Promise(
    (resolve, reject) =>
      contract
        .allEvents()
        .get((error, events) => (error ? reject(error) : resolve(events)))
  )
)

export const getLatestEvent = async (contract) => {
  let events = await getEvents(contract)
  return events[events.length - 1]
}

export const requestDataFrom = (oc, link, amount, args) => link.transferAndCall(oc.address, amount, args)

export const functionSelector = signature => '0x' + web3.sha3(signature).slice(2).slice(0, 8)

export const assertActionThrows = action => (
  Promise
    .resolve()
    .then(action)
    .catch(error => {
      assert(error, 'Expected an error to be raised')
      assert(error.message, 'Expected an error to be raised')
      return error.message
    })
    .then(errorMessage => {
      assert(errorMessage, 'Expected an error to be raised')
      const invalidOpcode = errorMessage.includes('invalid opcode')
      const reverted = errorMessage.includes('VM Exception while processing transaction: revert')
      assert.isTrue(invalidOpcode || reverted, 'expected error message to include "invalid JUMP" or "revert"')
      // see https://github.com/ethereumjs/testrpc/issues/39
      // for why the "invalid JUMP" is the throw related error when using TestRPC
    })
)

export const checkPublicABI = (contract, expectedPublic) => {
  let actualPublic = []
  for (const method of contract.abi) {
    if (method.type === 'function') actualPublic.push(method.name)
  };

  for (const method of actualPublic) {
    let index = expectedPublic.indexOf(method)
    assert.isAtLeast(index, 0, (`#${method} is NOT expected to be public`))
  }

  for (const method of expectedPublic) {
    let index = actualPublic.indexOf(method)
    assert.isAtLeast(index, 0, (`#${method} is expected to be public`))
  }
}

export const decodeRunABI = log => {
  let runABI = util.toBuffer(log.data)
  let types = ['bytes32', 'address', 'bytes4', 'bytes']
  return abi.rawDecode(types, runABI)
}

export const decodeRunRequest = log => {
  let runABI = util.toBuffer(log.data)
  let types = ['uint256', 'uint256', 'bytes']
  let [internalId, version, data] = abi.rawDecode(types, runABI)
  return [log.topics[1], log.topics[2], log.topics[3], toHex(internalId), version, data]
}

export const runRequestId = log => {
  var [_, _, _, internalId, _, _] = decodeRunRequest(log)
  return internalId
}

export const requestDataBytes = (specId, to, fHash, runId, data) => {
  let types = ['address', 'uint256', 'uint256', 'bytes32', 'address', 'bytes4', 'bytes32', 'bytes']
  let values = [0, 0, 1, specId, to, fHash, runId, data]
  let encoded = abi.rawEncode(types, values)
  let funcSelector = functionSelector('requestData(address,uint256,uint256,bytes32,address,bytes4,bytes32,bytes)')
  return funcSelector + encoded.toString('hex')
}

export const newUint8ArrayFromStr = (str) => {
  const codePoints = Array.prototype.map.call(str, c => c.charCodeAt(0))
  return Uint8Array.from(codePoints)
}

// newUint8Array returns a uint8array of count bytes from either a hex or
// decimal string, hex strings must begin with 0x
export const newUint8Array = (str, count) => {
  let result = new Uint8Array(count)

  if (str.startsWith('0x') || str.startsWith('0X')) {
    const hexStr = str.slice(2).padStart(count * 2, '0')
    for (let i = result.length; i >= 0; i--) {
      const offset = i * 2
      result[i] = parseInt(hexStr[offset] + hexStr[offset + 1], HEX_BASE)
    }
  } else {
    const num = bigNum(str)
    result = newHash('0x' + num.toString(HEX_BASE))
  }

  return result
}

// newSignature returns a 65 byte Uint8Array for representing a signature
export const newSignature = str => {
  return newUint8Array(str, 65)
}

// newHash returns a 65 byte Uint8Array for representing a hash
export const newHash = str => {
  return newUint8Array(str, 32)
}

// newAddress returns a 20 byte Uint8Array for representing an address
export const newAddress = str => {
  return newUint8Array(str, 20)
}

// lengthTypedArrays sums the length of all specified TypedArrays
export const lengthTypedArrays = (...arrays) => {
  return arrays.reduce((a, v) => a + v.length, 0)
}

export const toBuffer = uint8a => {
  return Buffer.from(uint8a)
}

// concatTypedArrays recursively concatenates TypedArrays into one big
// TypedArray
// TODO: Does not work recursively
export const concatTypedArrays = (...arrays) => {
  let size = lengthTypedArrays(...arrays)
  let result = new arrays[0].constructor(size)
  let offset = 0
  arrays.forEach((a) => {
    result.set(a, offset)
    offset += a.length
  })
  return result
}

export const splitRPCSignature = oracleSignature => {
  let v = oracleSignature[64]
  if (v < 27) {
    v += 27
  }
  return {
    v: v,
    r: oracleSignature.slice(0, 32),
    s: oracleSignature.slice(32, 64)
  }
}

export const increaseTime5Minutes = async () => {
  await web3.currentProvider.send({
    jsonrpc: "2.0", 
    method: "evm_increaseTime", 
    params: [300], id: 0
  });
};
