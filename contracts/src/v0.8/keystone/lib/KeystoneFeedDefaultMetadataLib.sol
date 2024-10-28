// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library KeystoneFeedDefaultMetadataLib {
  /**
   * Metadata Layout:
   *
   * +-------------------------------+--------------------+---------------------+---------------+
   * | 32 bytes (length prefix)      | 32 bytes           | 10 bytes            | 20 bytes      | 2 bytes        |
   * | (Not used in function)        | workflow_cid       | workflow_name       | workflow_owner| report_name    |
   * +-------------------------------+--------------------+---------------------+---------------+----------------+
   * |                               |                    |                     |               |                |
   * |          (Offset 0)           |     (Offset 32)    |     (Offset 64)     |  (Offset 74)  |  (Offset 94)   |
   * +-------------------------------+--------------------+---------------------+---------------+----------------+
   * @dev used to slice metadata bytes into workflowName, workflowOwner and report name
   */
  function _extractMetadataInfo(
    bytes memory metadata
  ) internal pure returns (bytes10 workflowName, address workflowOwner, bytes2 reportName) {
    // (first 32 bytes contain length of the byte array)
    // workflow_cid             // offset 32, size 32
    // workflow_name            // offset 64, size 10
    // workflow_owner           // offset 74, size 20
    // report_name              // offset 94, size  2
    assembly {
      // no shifting needed for bytes10 type
      workflowName := mload(add(metadata, 64))
      // shift right by 12 bytes to get the actual value
      workflowOwner := shr(mul(12, 8), mload(add(metadata, 74)))
      // no shifting needed for bytes2 type
      reportName := mload(add(metadata, 94))
    }
    return (workflowName, workflowOwner, reportName);
  }
}
