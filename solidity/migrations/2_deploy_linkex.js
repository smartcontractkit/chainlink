const LinkEx = artifacts.require('LinkEx.sol')

module.exports = async deployer => {
  // Deploy and update a LinkEx contract so it can be asserted upon in the tests
  // in a future block
  await deployer.deploy(LinkEx)
  const contract = await LinkEx.at(LinkEx.address)
  await contract.update('0xd321073d')
}
