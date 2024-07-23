// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

contract MockGasBoundCaller {
  function gasBoundCall(address target, uint256 gasAmount, bytes memory data) external payable {
    //    uint256 gasUsed = gasleft();
    bool success;
    assembly {
      success := call(gasAmount, target, 0, add(data, 0x20), mload(data), 0, 0)
    }
    //    gasUsed = gasUsed - gasleft();

    uint256 pubdataGas = 500000;
    bytes memory returnData = abi.encode(address(0), pubdataGas);

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

  //  function gasBoundCall(address, uint256, bytes memory) external payable {
  //    uint256 pubdataGas = 500000;
  //    bytes memory returnData = abi.encode(address(0), uint256(500000));
  //
  //    uint256 paddedReturndataLen = returnData.length + 96;
  //    if (paddedReturndataLen % 32 != 0) {
  //      paddedReturndataLen += 32 - (paddedReturndataLen % 32);
  //    }
  //
  //    assembly {
  //      mstore(sub(returnData, 0x40), 0x40)
  //      mstore(sub(returnData, 0x20), pubdataGas)
  //      return(sub(returnData, 0x40), paddedReturndataLen)
  //    }
  //  }
}
