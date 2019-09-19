#!/usr/bin/env node
/* eslint-disable @typescript-eslint/no-var-requires */

const Web3 = require('web3')
const path = require('path')
const EthLogJSON = require(path.join(__dirname, 'build/contracts/EthLog.json'))

const provider = new Web3.providers.HttpProvider('http://localhost:18545')
const devnetAddress = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'

const web3 = new Web3(provider)
const networkId = Object.keys(EthLogJSON.networks)[0]
const wc3 = new web3.eth.Contract(
  EthLogJSON.abi,
  EthLogJSON.networks[networkId].address,
)

wc3.methods.logEvent().send(
  {
    from: devnetAddress,
  },
  (error, transactionHash) => {
    if (error) {
      console.error('encountered error: ', error, transactionHash)
    } else {
      console.log('sent tx: ', transactionHash)
    }
  },
)
