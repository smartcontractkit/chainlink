let LinkToken = artifacts.require("../node_modules/smartcontractkit/chainlink/solidity/contracts/LinkToken.sol");
let RunLog = artifacts.require("./RunLog.sol");
let devnetAddress = "0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f";

module.exports = async function(deployer) {
  await LinkToken.deployed().then(async function(linkInstance) {
    await RunLog.deployed().then(async function(runLogInstance) {
      await linkInstance.transfer(runLogInstance.address, web3.toWei(1000));
      await linkInstance.transfer(devnetAddress, web3.toWei(1000));
    });
  });
};
