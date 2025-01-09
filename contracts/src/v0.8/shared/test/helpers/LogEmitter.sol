// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// solhint-disable
contract LogEmitter {
  event Log1(uint256);
  event Log2(uint256 indexed);
  event Log3(string);
  event Log4(uint256 indexed, uint256 indexed);

  function EmitLog1(uint256[] memory v) public {
    for (uint256 i = 0; i < v.length; i++) {
      emit Log1(v[i]);
    }
  }

  function EmitLog2(uint256[] memory v) public {
    for (uint256 i = 0; i < v.length; i++) {
      emit Log2(v[i]);
    }
  }

  function EmitLog3(string[] memory v) public {
    for (uint256 i = 0; i < v.length; i++) {
      emit Log3(v[i]);
    }
  }

  function EmitLog4(uint256 v, uint256 w, uint256 c) public {
    for (uint256 i = 0; i < c; i++) {
      emit Log4(v, w);
    }
  }
}
