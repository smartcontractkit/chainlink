#!/usr/bin/env node
/* eslint-disable @typescript-eslint/no-var-requires */

const Web3 = require('web3')
const contract = require('truffle-contract')
const path = require('path')
const UptimeSLAJSON = require(path.join(
  __dirname,
  'build/contracts/UptimeSLA.json',
))

const provider = new Web3.providers.HttpProvider('http://localhost:18545')
const UptimeSLA = contract(UptimeSLAJSON)
UptimeSLA.setProvider(provider)

UptimeSLA.deployed()
  .then(function(instance) {
    return instance.uptime.call()
  })
  .then(console.log, console.log)
