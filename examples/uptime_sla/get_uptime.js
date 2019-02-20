#!/usr/bin/env node

var Web3 = require('web3'),
  contract = require('truffle-contract'),
  path = require('path')
UptimeSLAJSON = require(path.join(__dirname, 'build/contracts/UptimeSLA.json'))

var provider = new Web3.providers.HttpProvider('http://localhost:18545')
var UptimeSLA = contract(UptimeSLAJSON)
UptimeSLA.setProvider(provider)

UptimeSLA.deployed()
  .then(function(instance) {
    return instance.uptime.call()
  })
  .then(console.log, console.log)
