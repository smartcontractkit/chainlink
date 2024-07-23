// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

contract MockGasBoundCaller {
  function gasBoundCall(address, uint256, bytes calldata) external payable {
    uint256 pubdataGas = 500000;
    bytes memory returnData = abi.encode(address(0), uint256(500000));

    uint256 paddedReturndataLen = returnData.length + 96;
    if (paddedReturndataLen % 32 != 0) {
      paddedReturndataLen += 32 - (paddedReturndataLen % 32);
    }

    assembly {
      mstore(sub(returnData, 0x40), 0x40)
      mstore(sub(returnData, 0x20), pubdataGas)
      return(sub(returnData, 0x40), paddedReturndataLen)
    }
  }
}
