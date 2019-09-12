const LINK = artifacts.require('LinkToken')
const Oracle = artifacts.require('Oracle')

module.exports = function(deployer) {
  deployer.deploy(Oracle, LINK.address)
}
