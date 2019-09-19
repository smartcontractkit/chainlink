const LinkToken = artifacts.require('LinkToken')

module.exports = function(deployer) {
  deployer.deploy(LinkToken)
}
