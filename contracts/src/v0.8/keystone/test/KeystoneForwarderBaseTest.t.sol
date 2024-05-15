// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "forge-std/Test.sol";
import {Receiver} from "./mocks/Receiver.sol";
import {KeystoneForwarder} from "../KeystoneForwarder.sol";

contract BaseTest is Test {
  address internal ADMIN = address(1);
  address internal constant TRANSMITTER = address(50);
  uint256 internal constant MAX_ORACLES = 31;
  uint32 internal DON_ID = 0x01020304;
  uint8 internal F = 1;

  struct Signer {
    uint256 mockPrivateKey;
    address signerAddress;
  }

  Signer[MAX_ORACLES] internal s_signers;
  KeystoneForwarder internal s_forwarder;
  Receiver internal s_receiver;

  function setUp() public virtual {
    vm.startPrank(ADMIN);
    s_forwarder = new KeystoneForwarder();
    s_receiver = new Receiver();

    uint256 seed = 0;

    for (uint256 i; i < MAX_ORACLES; i++) {
      uint256 mockPK = seed + i + 1;
      s_signers[i].mockPrivateKey = mockPK;
      s_signers[i].signerAddress = vm.addr(mockPK);
    }
  }

  function _getSignerAddresses() internal view returns (address[] memory) {
    address[] memory signerAddrs = new address[](s_signers.length);
    for (uint256 i = 0; i < signerAddrs.length; i++) {
      signerAddrs[i] = s_signers[i].signerAddress;
    }
    return signerAddrs;
  }

  function _getSignerAddresses(uint256 limit) internal view returns (address[] memory) {
    address[] memory signerAddrs = new address[](limit);
    for (uint256 i = 0; i < limit; i++) {
      signerAddrs[i] = s_signers[i].signerAddress;
    }
    return signerAddrs;
  }
}
