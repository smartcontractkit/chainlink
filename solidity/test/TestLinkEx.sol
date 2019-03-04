pragma solidity ^0.4.24 <0.6.0;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/LinkEx.sol";

contract TestLinkEx {
  function testCurrentRate() public {
    LinkEx ex = new LinkEx();

    Assert.equal(ex.currentRate(), 0, "Fresh contract should have 0 initial rate");
  }

  function testUpdate() public {
    LinkEx ex = new LinkEx();

    ex.update(8616460799);
    Assert.equal(ex.currentRate(), 0, "Contract should only return the rate set for future blocks");
  }
}
