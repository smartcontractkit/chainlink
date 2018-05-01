#!/usr/bin/env node

var Web3            = require('web3'),
    contract        = require("truffle-contract"),
    path            = require('path')
    RunLogJSON      = require(path.join(__dirname, 'build/contracts/RunLog.json'));

var provider = new Web3.providers.HttpProvider("http://localhost:18545");
var RunLog = contract(RunLogJSON);
RunLog.setProvider(provider);
var devnetAddress = "0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f";

RunLog.deployed().then(function(instance) {
  return instance.request({
    from: devnetAddress,
    gas: 200000
  });
}).then(console.log, console.log);
