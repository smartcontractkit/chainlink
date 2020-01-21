pragma solidity 0.5.0;

contract VRFRequestIDBase {
  /**
   * @notice Returns the id for this request
   * @param _keyHash The serviceAgreement ID to be used for this request
   * @param _seed The seed to be used in generating this randomness.
   */
  function makeRequestId(
    bytes32 _keyHash, uint256 _seed) public pure returns (bytes32) {
    return keccak256(abi.encodePacked(_keyHash, _seed));
  }
}
