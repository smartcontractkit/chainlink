// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {IDestinationVerifier} from "../../interfaces/IDestinationVerifier.sol";
import {Common} from "../../../libraries/Common.sol";

contract ErroredVerifier is IDestinationVerifier {
  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == this.verify.selector;
  }

// fix all of this msissing interfaces


function verifyBulk(bytes[] memory, bytes memory, address) external payable returns (bytes[] memory) {
  revert("Failed to verify");
}

function getAccessController() external view returns (address) {
  revert("Failed to verify");
}

 function getFeeManager() external view returns (address) {
  revert("Failed to verify");
}

 function setAccessController(address _accessController) external{
  revert("Failed to verify");
}

  function setConfigActive(bytes24 DONConfigID, bool isActive) external{
  revert("Failed to verify");
}

 function setFeeManager(address _feeManager) external {
revert("Failed to verify");
}

  function verify(
    bytes memory,
    bytes memory,
    /**
     * signedReport*
     */
    address
  )
    external
    payable
    override
    returns (
      /**
       * sender*
       */
      bytes memory
    )
  {
    revert("Failed to verify");
  }

  function setConfig(
    address[] memory,
    uint8,
    Common.AddressAndWeight[] memory
  ) external pure override {
    revert("Failed to set config");
  }


  function activateConfig(bytes32, bytes32) external pure {
    revert("Failed to activate config");
  }

  function deactivateConfig(bytes32, bytes32) external pure {
    revert("Failed to deactivate config");
  }

  function activateFeed(bytes32) external pure {
    revert("Failed to activate feed");
  }

  function deactivateFeed(bytes32) external pure {
    revert("Failed to deactivate feed");
  }
}
