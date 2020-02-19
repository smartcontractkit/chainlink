const MyContract = artifacts.require('MyContract')

/*
  This script makes it easy to read the data variable
  of the requesting contract.
*/

module.exports = async callback => {
  const mc = await MyContract.deployed()
  const data = await mc.data.call()
  callback(data)
}
