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

UptimeSLA.deployed()
  .then(function(instance) {
    return instance.uptime.call()
  })
  .then(console.log, console.log)
