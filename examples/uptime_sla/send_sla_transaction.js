#!/usr/bin/env node

var Web3 = require('web3')
var contract = require('truffle-contract')
var path = require('path')
const UptimeSLAJSON = require(path.join(
  __dirname,
  'build/contracts/UptimeSLA.json',
))

var provider = new Web3.providers.HttpProvider('http://localhost:18545')
var UptimeSLA = contract(UptimeSLAJSON)
UptimeSLA.setProvider(provider)
var devnetAddress = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'

UptimeSLA.deployed()
  .then(function(instance) {
    return instance.updateUptime('0', { from: devnetAddress })
  })
  .then(console.log, console.log)
