pragma solidity ^0.8.0;

interface ChainlinkPoRAddressList {
    function getPoRAddressListLength() external view returns (uint);
    function getPoRAddressList(uint startIndex, uint endIndex) external view returns (string[] memory);
}
