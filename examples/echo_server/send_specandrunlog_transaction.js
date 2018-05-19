#!/usr/bin/env node

var Web3            = require('web3'),
    contract        = require("truffle-contract"),
    path            = require('path')
    SpecAndRunJSON  = require(path.join(__dirname, 'build/contracts/SpecAndRunLog.json'));

var devnetAddress = "0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f";
var provider      = new Web3.providers.HttpProvider("http://localhost:18545");
var SpecAndRunLog = contract(SpecAndRunJSON);
SpecAndRunLog.setProvider(provider);

SpecAndRunLog.deployed().then(function(instance) {
  return instance.request({
    from: devnetAddress,
    gas: 200000
  });
}).then(console.log, console.log);
