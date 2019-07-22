let MyContract = artifacts.require('MyContract')
let LinkToken = artifacts.require('LinkToken')

/*
  This script is meant to assist with funding the requesting
  contract with LINK. It will send 1 LINK to the requesting
  contract for ease-of-use. Any extra LINK present on the contract
  can be retrieved by calling the withdrawLink() function.
*/

const payment = process.env.TRUFFLE_CL_BOX_PAYMENT || '1000000000000000000'

module.exports = async (callback) => {
  let mc = await MyContract.deployed()
  let tokenAddress = await mc.getChainlinkToken()
  let token = await LinkToken.at(tokenAddress)
  console.log('Funding contract:', mc.address)
  let tx = await token.transfer(mc.address, payment)
  callback(tx.tx)
}
