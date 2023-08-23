// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import {BaseTest} from "./BaseTest.t.sol";
import {Forwarder} from "../Forwarder.sol";
import {IForwarder} from "../IForwarder.sol";
import {BankERC20} from "../BankERC20.sol";
import "forge-std/console.sol";

contract ForwarderTest is BaseTest {
  string public constant TOKEN_NAME = "BankToken";
  string public constant TOKEN_VERSION = "1";
  uint256 public constant TOTAL_SUPPLY = 1e27;
  uint64 public constant CCIP_CHAIN_ID = 51;
  uint256 FROM_PKEY = uint256(bytes32(hex"e9abcd10b456381f69182ba2b563fe8a4223e172d6f7864c1ee1e01e722054d0"));
  address FROM = 0xE13CD8937CADa693048dC358734064C6a848f8F8;
  uint256 FROM_PKEY2 = uint256(bytes32(hex"e9abcd10b456381f69182ba2b563fe8a4223e172d6f7864c1ee1e01e722054ac"));
  address FROM2 = 0x920e8Dc2085341B8A8E8b1684bC7BE88Df373487;

  Forwarder private s_forwarder;
  BankERC20 private s_bankERC20;
  address private s_ccipRouterMock;
  address private s_feeProvider;
  uint256 private s_validUntilTime;
  address private s_toAddress = makeAddr("to");
  uint256 private s_nonce;
  uint256 private s_amount;
  bytes32 private s_domainSeparatorHash;
  bytes32 private s_requestTypeHash;

  function setUp() public override {
    s_forwarder = new Forwarder();
    s_ccipRouterMock = makeAddr("ccipRouter");
    s_feeProvider = makeAddr("feeProvider");
    s_bankERC20 = new BankERC20(
      TOKEN_NAME,
      TOKEN_VERSION,
      TOTAL_SUPPLY,
      address(s_forwarder),
      s_ccipRouterMock,
      s_feeProvider,
      CCIP_CHAIN_ID
    );
    s_validUntilTime = block.timestamp + 60 * 60; // valid for 1 hour
    s_nonce = uint256(keccak256("nonce"));
    s_amount = 1e18;
    s_forwarder.registerDomainSeparator(TOKEN_NAME, TOKEN_VERSION);
    bytes memory domainSeparator = s_forwarder.getDomainSeparator(TOKEN_NAME, TOKEN_VERSION);
    s_domainSeparatorHash = keccak256(domainSeparator);
    bytes memory requestType = abi.encodePacked(
      "ForwardRequest(address from,address target,uint256 nonce,bytes data,uint256 validUntilTime)"
    );
    s_requestTypeHash = keccak256(requestType);
  }

  function testExecuteRevertsOnUnavailableNonce() public {
    (IForwarder.ForwardRequest memory req, bytes memory signature) = _generateRequestAndSig(FROM_PKEY, FROM);

    (bool success, ) = s_forwarder.execute(req, s_domainSeparatorHash, s_requestTypeHash, "", signature);
    assertTrue(success);

    vm.expectRevert(abi.encodeWithSelector(Forwarder.NonceAlreadyUsed.selector, s_nonce));
    s_forwarder.execute(req, s_domainSeparatorHash, s_requestTypeHash, "", signature);
  }

  function testExecuteSameNonceMultipleFromAddresses() public {
    (IForwarder.ForwardRequest memory req, bytes memory signature) = _generateRequestAndSig(FROM_PKEY, FROM);
    (bool success, bytes memory ret) = s_forwarder.execute(
      req,
      s_domainSeparatorHash,
      s_requestTypeHash,
      "",
      signature
    );
    assertTrue(success);

    (req, signature) = _generateRequestAndSig(FROM_PKEY2, FROM2);
    (success, ret) = s_forwarder.execute(req, s_domainSeparatorHash, s_requestTypeHash, "", signature);
    assertTrue(success);
  }

  function _generateRequestAndSig(
    uint256 fromPK,
    address from
  ) internal returns (IForwarder.ForwardRequest memory req, bytes memory signature) {
    bytes memory encodedCalldata = abi.encodeWithSignature(
      "metaTransfer(address,uint256,uint64)",
      s_toAddress,
      s_amount,
      CCIP_CHAIN_ID
    );

    bytes32 calldataHash = keccak256(encodedCalldata);

    IForwarder.ForwardRequest memory forwardRequest = IForwarder.ForwardRequest({
      from: from,
      target: address(s_bankERC20),
      nonce: s_nonce,
      data: encodedCalldata,
      validUntilTime: s_validUntilTime
    });

    bytes memory forwardData = abi.encodePacked(
      "\x19\x01",
      s_domainSeparatorHash,
      keccak256(
        abi.encodePacked(
          s_requestTypeHash,
          uint256(uint160(from)),
          uint256(uint160(address(s_bankERC20))),
          s_nonce,
          calldataHash,
          s_validUntilTime,
          ""
        )
      )
    );

    (uint8 v, bytes32 r, bytes32 s) = vm.sign(fromPK, keccak256(forwardData));
    signature = abi.encodePacked(r, s, v);

    assertTrue(s_bankERC20.transfer(from, s_amount));

    return (forwardRequest, signature);
  }
}
