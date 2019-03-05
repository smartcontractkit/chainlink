pragma solidity ^0.4.24 <0.6.0;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/LinkEx.sol";

contract TestLinkEx {
  function testCurrentRate() public {
    LinkEx ex = new LinkEx();

    Assert.equal(ex.currentRate(), 0, "Fresh contract should have 0 initial rate");
  }

  function testCurrentRateNewTransaction() public {
    LinkEx ex = LinkEx(DeployedAddresses.LinkEx());

    Assert.equal(ex.currentRate(), 3542157117, "Already deployed contract should have non zero rate");
  }

  function testUpdateSameTransaction() public {
    LinkEx ex = new LinkEx();

    ex.update(8616460799);
    Assert.equal(ex.currentRate(), 0, "Contract should only return the rate set for future blocks");
  }

  function testUpdateNewTransaction() public {
    LinkEx ex = LinkEx(DeployedAddresses.LinkEx());

    ex.update(8616460799);
    Assert.equal(ex.currentRate(), 3542157117, "Contract should only return the rate set for future blocks");

    ex.update(8616460799);
    Assert.equal(ex.currentRate(), 3542157117, "Update rate shouldn't overwrite historic rate");
  }
}
