// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistryWithFixture is WorkflowRegistrySetup {
  string internal s_workflowName1 = "Workflow1";
  bytes32 internal s_workflowID1 = keccak256("workflow1");
  string internal s_binaryURL1 = "https://example.com/binary1";
  string internal s_configURL1 = "https://example.com/config1";
  string internal s_secretsURL1 = "https://example.com/secrets1";

  string internal s_workflowName2 = "Workflow2";
  bytes32 internal s_workflowID2 = keccak256("workflow2");
  string internal s_binaryURL2 = "https://example.com/binary2";
  string internal s_configURL2 = "https://example.com/config2";
  string internal s_secretsURL2 = "https://example.com/secrets2";

  string internal s_workflowName3 = "Workflow3";
  bytes32 internal s_workflowID3 = keccak256("workflow3");
  string internal s_binaryURL3 = "https://example.com/binary3";
  string internal s_configURL3 = "https://example.com/config3";
  string internal s_secretsURL3 = "https://example.com/secrets3";

  function setUp() public override {
    super.setUp();

    // Register some workflows for s_authorizedAddress in s_allowedDonID
    string[] memory workflowNames = new string[](3);
    bytes32[] memory workflowIDs = new bytes32[](3);
    string[] memory binaryURLs = new string[](3);
    string[] memory configURLs = new string[](3);
    string[] memory secretsURLs = new string[](3);

    workflowNames[0] = s_workflowName1;
    workflowIDs[0] = s_workflowID1;
    binaryURLs[0] = s_binaryURL1;
    configURLs[0] = s_configURL1;
    secretsURLs[0] = s_secretsURL1;

    workflowNames[1] = s_workflowName2;
    workflowIDs[1] = s_workflowID2;
    binaryURLs[1] = s_binaryURL2;
    configURLs[1] = s_configURL2;
    secretsURLs[1] = s_secretsURL2;

    workflowNames[2] = s_workflowName3;
    workflowIDs[2] = s_workflowID3;
    binaryURLs[2] = s_binaryURL3;
    configURLs[2] = s_configURL3;
    secretsURLs[2] = s_secretsURL3;

    vm.startPrank(s_authorizedAddress);
    for (uint256 i = 0; i < workflowNames.length; ++i) {
      s_registry.registerWorkflow(
        workflowNames[i],
        workflowIDs[i],
        s_allowedDonID,
        WorkflowRegistry.WorkflowStatus.ACTIVE,
        binaryURLs[i],
        configURLs[i],
        secretsURLs[i]
      );
    }
    vm.stopPrank();
  }
}
