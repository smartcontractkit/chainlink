let MyContract = artifacts.require('MyContract')

/*
  This script makes it easy to read the data variable
  of the requesting contract.
*/

module.exports = async callback => {
  let mc = await MyContract.deployed()
  let data = await mc.data.call()
  callback(data)
}
