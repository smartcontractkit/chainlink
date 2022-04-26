pragma solidity ^0.8.0;

interface ChainlinkPoRAddressList {
  function getPoRAddressListLength() external view returns (uint256);

  function getPoRAddressList(uint256 startIndex, uint256 endIndex) external view returns (string[] memory);
}
