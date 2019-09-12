#!/usr/bin/env node

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
const devnetAddress = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'

UptimeSLA.deployed()
  .then(function(instance) {
    return instance.updateUptime('0', { from: devnetAddress })
  })
  .then(console.log, console.log)
