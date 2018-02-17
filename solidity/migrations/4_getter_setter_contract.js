var GetterSetter = artifacts.require("./GetterSetter.sol");

module.exports = function(deployer) {
  deployer.deploy(GetterSetter);
};
