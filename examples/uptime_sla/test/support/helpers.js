import {
  eth,
  functionSelector,
  getEvents
} from '../../../../solidity/test/support/helpers'

const BigNumber = require('bignumber.js')
const abi = require('ethereumjs-abi')
const cbor = require('cbor')
const moment = require('moment')
const util = require('ethereumjs-util')

let Eth,
  emptyAddress,
  sealBlock,
  sendTransaction,
  getBalance,
  bigNum,
  toWei,
  tokens,
  intToHex,
  hexToInt,
  hexToAddress,
  unixTime,
  seconds,
  minutes,
  hours,
  days,
  keccak256,
  logTopic,
  getLatestBlock,
  getLatestTimestamp,
  fastForwardTo,
  eventsOfType,
  getEventsOfType,
  encodeUint256,
  encodeAddress,
  encodeBytes,
  checkPublicABI,
  randomHex,
  newAddress

(() => {
  Eth = function sendEth (method, params) {
    params = params || []

    return new Promise((resolve, reject) => {
      web3.currentProvider.sendAsync({
        jsonrpc: '2.0',
        method: method,
        params: params || [],
        id: new Date().getTime()
      }, function sendEthResponse (error, response) {
        if (error) {
          reject(error)
        } else {
          resolve(response.result)
        };
      }, () => {}, () => {})
    })
  }

  emptyAddress = '0x0000000000000000000000000000000000000000'

  sealBlock = async function () {
    return Eth('evm_mine')
  }

  sendTransaction = async function (params) {
    return await eth.sendTransaction(params)
  }

  getBalance = async function (account) {
    return bigNum(await eth.getBalance(account))
  }

  bigNum = function (number) {
    return new BigNumber(number)
  }

  toWei = function (number) {
    return web3.toWei(number)
  }

  tokens = function (number) {
    return bigNum(number * 10 ** 18)
  }

  intToHex = function (number) {
    return '0x' + bigNum(number).toString(16)
  }

  hexToInt = function (string) {
    return web3.toBigNumber(string)
  }

  hexToAddress = function (string) {
    return '0x' + string.slice(string.length - 40)
  }

  unixTime = function (time) {
    return moment(time).unix()
  }

  seconds = function (number) {
    return number
  }

  minutes = function (number) {
    return number * 60
  }

  hours = function (number) {
    return number * minutes(60)
  }

  days = function (number) {
    return number * hours(24)
  }

  keccak256 = function (string) {
    return web3.sha3(string)
  }

  logTopic = function (string) {
    let hash = keccak256(string)
    return '0x' + hash.slice(26)
  }

  getLatestBlock = async function () {
    return await eth.getBlock('latest', false)
  }

  getLatestTimestamp = async function () {
    let latestBlock = await getLatestBlock()
    return web3.toDecimal(latestBlock.timestamp)
  }

  fastForwardTo = async function (target) {
    let now = await getLatestTimestamp()
    assert.isAbove(target, now, 'Cannot fast forward to the past')
    let difference = target - now
    await Eth('evm_increaseTime', [difference])
    await sealBlock()
  }

  eventsOfType = function (events, type) {
    let filteredEvents = []
    for (event of events) {
      if (event.event === type) filteredEvents.push(event)
    }
    return filteredEvents
  }

  getEventsOfType = async function (contract, type) {
    return eventsOfType(await getEvents(contract), type)
  }

  encodeUint256 = function (int) {
    let zeros = '0000000000000000000000000000000000000000000000000000000000000000'
    let payload = int.toString(16)
    return (zeros + payload).slice(payload.length)
  }

  encodeAddress = function (address) {
    return '000000000000000000000000' + address.slice(2)
  }

  encodeBytes = function (bytes) {
    let zeros = '0000000000000000000000000000000000000000000000000000000000000000'
    let padded = bytes.padEnd(64, 0)
    let length = encodeUint256(bytes.length / 2)
    return length + padded
  }

  checkPublicABI = function (contract, expectedPublic) {
    let actualPublic = []
    for (method of contract.abi) {
      if (method.type == 'function') actualPublic.push(method.name)
    };

    for (method of actualPublic) {
      let index = expectedPublic.indexOf(method)
      assert.isAtLeast(index, 0, (`#${method} is NOT expected to be public`))
    }

    for (method of expectedPublic) {
      let index = actualPublic.indexOf(method)
      assert.isAtLeast(index, 0, (`#${method} is expected to be public`))
    }
  }

  // https://codepen.io/code_monk/pen/FvpfI
  randomHex = function (len) {
    var maxlen = 8
    var min = Math.pow(16, Math.min(len, maxlen) - 1)
    var max = Math.pow(16, Math.min(len, maxlen)) - 1
    var n = Math.floor(Math.random() * (max - min + 1)) + min
    var r = n.toString(16)
    while (r.length < len) {
      r = r + randomHex(len - maxlen)
    }
    return r
  }

  newAddress = () => {
    return '0x' + randomHex(40)
  }
})()

export {
  Eth,
  bigNum,
  checkPublicABI,
  days,
  emptyAddress,
  encodeAddress,
  encodeBytes,
  encodeUint256,
  eventsOfType,
  fastForwardTo,
  getBalance,
  getEventsOfType,
  getLatestBlock,
  getLatestTimestamp,
  hexToAddress,
  hexToInt,
  hours,
  intToHex,
  keccak256,
  logTopic,
  minutes,
  newAddress,
  randomHex,
  sealBlock,
  seconds,
  sendTransaction,
  toWei,
  tokens,
  unixTime
}
