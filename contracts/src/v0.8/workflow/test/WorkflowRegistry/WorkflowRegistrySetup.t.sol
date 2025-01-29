// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {Test} from "forge-std/Test.sol";

contract WorkflowRegistrySetup is Test {
  WorkflowRegistry internal s_registry;
  address internal s_owner;
  address internal s_stranger;
  address internal s_authorizedAddress;
  address internal s_unauthorizedAddress;
  uint32 internal s_allowedDonID;
  uint32 internal s_disallowedDonID;
  bytes32 internal s_validWorkflowID;
  string internal s_validWorkflowName;
  string internal s_validBinaryURL;
  string internal s_validConfigURL;
  string internal s_validSecretsURL;
  string internal s_invalidWorkflowName;
  string internal s_invalidURL;
  bytes32 internal s_validWorkflowKey;

  function setUp() public virtual {
    s_owner = makeAddr("owner");
    s_stranger = makeAddr("nonOwner");
    s_authorizedAddress = makeAddr("authorizedAddress");
    s_unauthorizedAddress = makeAddr("unauthorizedAddress");
    s_allowedDonID = 1;
    s_disallowedDonID = 99;
    s_validWorkflowID = keccak256("validWorkflow");
    s_validWorkflowName = "ValidWorkflow";
    s_validBinaryURL = "https://example.com/valid-binary";
    s_validConfigURL = "https://example.com/valid-config";
    s_validSecretsURL = "https://example.com/valid-secrets";
    s_invalidWorkflowName =
      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcd";
    s_invalidURL =
      "https://www.example.com/this/is/a/very/long/url/that/keeps/going/on/and/on/to/ensure/that/it/exceeds/two/hundred/and/one/characters/in/length/for/testing/purposes/and/it/should/be/sufficiently/long/to/meet/your/requirements/for/this/test";

    uint32[] memory allowedDONs = new uint32[](1);
    allowedDONs[0] = s_allowedDonID;
    address[] memory authorizedAddresses = new address[](1);
    authorizedAddresses[0] = s_authorizedAddress;

    // Deploy the WorkflowRegistry contract
    vm.startPrank(s_owner);
    s_registry = new WorkflowRegistry();

    s_validWorkflowKey = s_registry.computeHashKey(s_authorizedAddress, s_validWorkflowName);

    // Perform initial setup as the owner
    s_registry.updateAllowedDONs(allowedDONs, true);
    s_registry.updateAuthorizedAddresses(authorizedAddresses, true);
    vm.stopPrank();
  }

  // Helper function to register a valid workflow
  function _registerValidWorkflow() internal {
    vm.prank(s_authorizedAddress);
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );
  }

  // Helper function to remove an address from the authorized addresses list
  function _removeAddressFromAuthorizedAddresses(
    address addressToRemove
  ) internal {
    address[] memory addressesToRemove = new address[](1);
    addressesToRemove[0] = addressToRemove;
    vm.prank(s_owner);
    s_registry.updateAuthorizedAddresses(addressesToRemove, false);
  }

  // Helper function to remove a DON from the allowed DONs list
  function _removeDONFromAllowedDONs(
    uint32 donIDToRemove
  ) internal {
    uint32[] memory donIDsToRemove = new uint32[](1);
    donIDsToRemove[0] = donIDToRemove;
    vm.prank(s_owner);
    s_registry.updateAllowedDONs(donIDsToRemove, false);
  }

  // Helper function to add an address to the authorized addresses list
  function _addAddressToAuthorizedAddresses(
    address addressToAdd
  ) internal {
    address[] memory addressesToAdd = new address[](1);
    addressesToAdd[0] = addressToAdd;
    vm.prank(s_owner);
    s_registry.updateAuthorizedAddresses(addressesToAdd, true);
  }
}
