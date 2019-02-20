let LINK = artifacts.require('LinkToken')
let Oracle = artifacts.require('Oracle')

module.exports = function(deployer) {
  deployer.deploy(Oracle, LINK.address)
}
