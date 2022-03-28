// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../KeeperCompatible.sol";
import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/KeeperRegistryInterface.sol";
import "../ConfirmedOwner.sol";

contract UpkeepAutoFunder is KeeperCompatible, ConfirmedOwner {
  bool public s_isEligible;
  bool public s_shouldCancel;
  uint256 public s_upkeepId;
  uint96 public s_autoFundLink;
  LinkTokenInterface public immutable LINK;
  KeeperRegistryBaseInterface public immutable s_keeperRegistry;

  constructor(address linkAddress, address registryAddress) ConfirmedOwner(msg.sender) {
    LINK = LinkTokenInterface(linkAddress);
    s_keeperRegistry = KeeperRegistryBaseInterface(registryAddress);

    s_isEligible = false;
    s_shouldCancel = false;
    s_upkeepId = 0;
    s_autoFundLink = 0;
  }

  function setShouldCancel(bool value) external onlyOwner {
    s_shouldCancel = value;
  }

  function setIsEligible(bool value) external onlyOwner {
    s_isEligible = value;
  }

  function setAutoFundLink(uint96 value) external onlyOwner {
    s_autoFundLink = value;
  }

  function setUpkeepId(uint256 value) external onlyOwner {
    s_upkeepId = value;
  }

  function checkUpkeep(bytes calldata data)
    external
    override
    cannotExecute
    returns (bool callable, bytes calldata executedata)
  {
    return (s_isEligible, data);
  }

  function performUpkeep(bytes calldata data) external override {
    require(s_isEligible, "Upkeep should be eligible");
    s_isEligible = false; // Allow upkeep only once until it is set again

    // Topup upkeep so it can be called again
    LINK.transferAndCall(address(s_keeperRegistry), s_autoFundLink, abi.encode(s_upkeepId));

    if (s_shouldCancel) {
      s_keeperRegistry.cancelUpkeep(s_upkeepId);
    }
  }
}
