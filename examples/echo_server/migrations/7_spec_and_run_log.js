const SpecAndRunLog = artifacts.require("./SpecAndRunLog.sol");
const LinkToken = artifacts.require("../../../solidity/contracts/lib/LinkToken.sol");
const Oracle = artifacts.require("../../../solidity/contracts/Oracle.sol");

module.exports = function(deployer) {
  deployer.deploy(SpecAndRunLog, LinkToken.address, Oracle.address).then(async function() {
    let linkInstance = await LinkToken.deployed();
    await linkInstance.transfer(SpecAndRunLog.address, web3.toWei(1000));
  });
};
