#!/usr/bin/env node

const Web3 = require('web3')
const path = require('path')
const RunLogJSON = require(path.join(__dirname, 'build/contracts/RunLog.json'))

const provider = new Web3.providers.HttpProvider('http://localhost:18545')
const devnetAddress = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'
const web3 = new Web3(provider)
const networkId = Object.keys(RunLogJSON.networks)[0]
const wc3 = new web3.eth.Contract(
  RunLogJSON.abi,
  RunLogJSON.networks[networkId].address
)

wc3.methods.request().send(
  {
    from: devnetAddress,
    gas: 200000
  },
  (error, transactionHash) => {
    if (error) {
      console.error('encountered error: ', error, transactionHash)
    } else {
      console.log('sent tx: ', transactionHash)
    }
  }
)
