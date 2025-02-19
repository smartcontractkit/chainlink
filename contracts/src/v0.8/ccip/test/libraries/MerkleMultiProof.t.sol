// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MerkleMultiProof} from "../../libraries/MerkleMultiProof.sol";
import {MerkleHelper} from "../helpers/MerkleHelper.sol";
import {Test} from "forge-std/Test.sol";

contract MerkleMultiProofTest is Test {
  // This must match the spec
  function test_SpecSync_gas() public pure {
    bytes32 expectedRoot = 0xd4f0f3c40a4d583d98c17d89e550b1143fe4d3d759f25ccc63131c90b183928e;

    bytes32[] memory leaves = new bytes32[](10);
    leaves[0] = 0xa20c0244af79697a4ef4e2378c9d5d14cbd49ddab3427b12594c7cfa67a7f240;
    leaves[1] = 0x3de96afb24ce2ac45a5595aa13d1a5163ae0b3c94cef6b2dc306b5966f32dfa5;
    leaves[2] = 0xacadf7b4d13cd57c5d25f1d27be39b656347fe8f8e0de8db9c76d979dff57736;
    leaves[3] = 0xc21c26a709802fe1ae52a9cd8ad94d15bf142ded26314339cd87a13e5b468165;
    leaves[4] = 0x55f6df03562738c9a6437cd9ad221c52b76906a175ae96188cff60e0a2a59933;
    leaves[5] = 0x2dbbe66452e43fec839dc65d5945aad6433d410c65863eaf1d876e1e0b06343c;
    leaves[6] = 0x8beab00297b94bf079fcd5893b0a33ebf6b0ce862cd06be07c87d3c63e1c4acf;
    leaves[7] = 0xcabdd3ad25daeb1e0541042f2ea4cd177f54e67aa4a2c697acd4bb682e94de59;
    leaves[8] = 0x7e01d497203685e99e34df33d55465c66b2253fa1630ee2fe5c4997968e4a6fa;
    leaves[9] = 0x1a03d013f1e2fa9cc04f89c7528ac3216e3e096a1185d7247304e97c59f9661f;

    bytes32[] memory proofs = new bytes32[](33);
    proofs[0] = 0xde96f24fcf9ddd20c803dc9c5fba7c478a5598a08a0faa5f032c65823b8e26a3;
    proofs[1] = 0xe1303cffc3958a6b93e2dc04caf21f200ff5aa5be090c5013f37804b91488bc2;
    proofs[2] = 0x90d80c76bccb44a91f4e16604976163aaa39e9a1588b0b24b33a61f1d4ba7bb5;
    proofs[3] = 0x012a299b25539d513c8677ecf37968774e9e4b045e79737f48defd350224cdfd;
    proofs[4] = 0x420a36c5a73f87d8fb98e70c48d0d6f9dd83f50b7b91416a6f5f91fac4db800f;
    proofs[5] = 0x5857d8d1b56abcd7f863cedd3c3f8677256f54d675be61f05efa45d6495fc30a;
    proofs[6] = 0xbf176d20166fdeb72593ff97efec1ce6244af41ca46cf0bc902d19d50c446f7b;
    proofs[7] = 0xa9221608e4380250a1815fb308632bce99f611a673d2e17fc617123fdc6afcd2;
    proofs[8] = 0xbd14f3366c73186314f182027217d0f70eba55817561de9e9a1f2c78bf5cbead;
    proofs[9] = 0x2f9aa48c0c9f82aaac65d7a9374a52d9dc138ed100a5809ede57e70697f48b56;
    proofs[10] = 0x2ae60afa54271cb421c12e4441c2dac0a25f25c9433a6d07cb32419e993fe344;
    proofs[11] = 0xc765c091680f0434b74c44507b932e5c80f6e995a975a275e5b130af1de1064c;
    proofs[12] = 0x59d2d6e0c4a5d07b169dbcdfa39dad7aea7b7783a814399f4f44c4a36b6336d3;
    proofs[13] = 0xdd14d1387d10740187d71ad9500475399559c0922dbe2576882e61f1edd84692;
    proofs[14] = 0x5412b8395509935406811ab3da43ab80be7acd8ffb5f398ab70f056ff3740f46;
    proofs[15] = 0xeadab258ae7d779ce5f10fbb1bb0273116b8eccbf738ed878db570de78bed1e4;
    proofs[16] = 0x6133aa40e6db75373b7cfc79e6f8b8ce80e441e6c1f98b85a593464dda3cf9c0;
    proofs[17] = 0x5418948467112660639b932af9b1b212e40d71b24326b4606679d168a765af4f;
    proofs[18] = 0x44f618505355c7e4e7c0f81d6bb15d2ec9cf9b366f9e1dc37db52745486e6b0f;
    proofs[19] = 0xa410ee174a66a4d64f3c000b93efe15b5b1f3e39e962af2580fcd30bce07d039;
    proofs[20] = 0x09c3eb05ac9552022a45c00d01a47cd56f95f94afdd4402299dba1291a17f976;
    proofs[21] = 0x0e780f6acd081b07320a55208fa3e1d884e2e95cb13d1c98c74b7e853372c813;
    proofs[22] = 0x2b60e8c21f78ef22fa4297f28f1d8c747181edfc465121b39c16be97d4fb8a04;
    proofs[23] = 0xf24da95060a8598c06e9dfb3926e1a8c8bd8ec2c65be10e69323442840724888;
    proofs[24] = 0x7e220fc095bcd2b0f5ef134d9620d89f6d7a1e8719ce8893bb9aff15e847578f;
    proofs[25] = 0xcfe9e475c4bd32f1e36b2cc65a959c403c59979ff914fb629a64385b0c680a71;
    proofs[26] = 0x25237fb8d1bfdc01ca5363ec3166a2b40789e38d5adcc8627801da683d2e1d76;
    proofs[27] = 0x42647949fed0250139c01212d739d8c83d2852589ebc892d3490ae52e411432c;
    proofs[28] = 0x34397a30930e6dd4fb5af48084afc5cfbe02c18dd9544b3faff4e2e90bf00cb9;
    proofs[29] = 0xa028f33226adc3d1cb72b19eb6808dab9190b25066a45cacb5dfe5d640e57cf2;
    proofs[30] = 0x7cff66ba47a05f932d06d168c294266dcb0d3943a4f2a4a75c860b9fd6e53092;
    proofs[31] = 0x5ca1b32f1dbfadd83205882be5eb76f34c49e834726f5239905a0e70d0a5e0eb;
    proofs[32] = 0x1b4b087a89e4eca6cdd237210932559dc8fd167d5f4f2d9acb13264e1e305479;

    uint256 flagsUint256 = 0x2f3c0000000;

    bytes32 root = MerkleMultiProof._merkleRoot(leaves, proofs, flagsUint256);

    assertEq(expectedRoot, root);
  }

  function testFuzz_MerkleRoot2(bytes32 left, bytes32 right) public pure {
    bytes32[] memory leaves = new bytes32[](2);
    leaves[0] = left;
    leaves[1] = right;
    bytes32[] memory proofs = new bytes32[](0);

    bytes32 expectedRoot = MerkleHelper.hashPair(left, right);

    bytes32 root = MerkleMultiProof._merkleRoot(leaves, proofs, 2 ** 2 - 1);

    assertEq(root, expectedRoot);
  }

  function test_MerkleRoot256() public pure {
    bytes32[] memory leaves = new bytes32[](256);
    for (uint256 i = 0; i < leaves.length; ++i) {
      leaves[i] = keccak256("a");
    }
    bytes32[] memory proofs = new bytes32[](0);

    bytes32 expectedRoot = MerkleHelper.getMerkleRoot(leaves);

    bytes32 root = MerkleMultiProof._merkleRoot(leaves, proofs, 2 ** 256 - 1);

    assertEq(root, expectedRoot);
  }

  function testFuzz_MerkleMulti1of4(bytes32 leaf1, bytes32 proof1, bytes32 proof2) public pure {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = leaf1;
    bytes32[] memory proofs = new bytes32[](2);
    proofs[0] = proof1;
    proofs[1] = proof2;

    // Proof flag = false
    bytes32 result = MerkleHelper.hashPair(leaves[0], proofs[0]);
    // Proof flag = false
    result = MerkleHelper.hashPair(result, proofs[1]);

    assertEq(MerkleMultiProof._merkleRoot(leaves, proofs, 0), result);
  }

  function testFuzz_MerkleMulti2of4(bytes32 leaf1, bytes32 leaf2, bytes32 proof1, bytes32 proof2) public pure {
    bytes32[] memory leaves = new bytes32[](2);
    leaves[0] = leaf1;
    leaves[1] = leaf2;
    bytes32[] memory proofs = new bytes32[](2);
    proofs[0] = proof1;
    proofs[1] = proof2;

    // Proof flag = false
    bytes32 result1 = MerkleHelper.hashPair(leaves[0], proofs[0]);
    // Proof flag = false
    bytes32 result2 = MerkleHelper.hashPair(leaves[1], proofs[1]);
    // Proof flag = true
    bytes32 finalResult = MerkleHelper.hashPair(result1, result2);

    assertEq(MerkleMultiProof._merkleRoot(leaves, proofs, 4), finalResult);
  }

  function testFuzz_MerkleMulti3of4(bytes32 leaf1, bytes32 leaf2, bytes32 leaf3, bytes32 proof) public pure {
    bytes32[] memory leaves = new bytes32[](3);
    leaves[0] = leaf1;
    leaves[1] = leaf2;
    leaves[2] = leaf3;
    bytes32[] memory proofs = new bytes32[](1);
    proofs[0] = proof;

    // Proof flag = true
    bytes32 result1 = MerkleHelper.hashPair(leaves[0], leaves[1]);
    // Proof flag = false
    bytes32 result2 = MerkleHelper.hashPair(leaves[2], proofs[0]);
    // Proof flag = true
    bytes32 finalResult = MerkleHelper.hashPair(result1, result2);

    assertEq(MerkleMultiProof._merkleRoot(leaves, proofs, 5), finalResult);
  }

  function testFuzz_MerkleMulti4of4(bytes32 leaf1, bytes32 leaf2, bytes32 leaf3, bytes32 leaf4) public pure {
    bytes32[] memory leaves = new bytes32[](4);
    leaves[0] = leaf1;
    leaves[1] = leaf2;
    leaves[2] = leaf3;
    leaves[3] = leaf4;
    bytes32[] memory proofs = new bytes32[](0);

    // Proof flag = true
    bytes32 result1 = MerkleHelper.hashPair(leaves[0], leaves[1]);
    // Proof flag = true
    bytes32 result2 = MerkleHelper.hashPair(leaves[2], leaves[3]);
    // Proof flag = true
    bytes32 finalResult = MerkleHelper.hashPair(result1, result2);

    assertEq(MerkleMultiProof._merkleRoot(leaves, proofs, 7), finalResult);
  }

  function test_MerkleRootSingleLeaf() public pure {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";
    bytes32[] memory proofs = new bytes32[](0);
    assertEq(MerkleMultiProof._merkleRoot(leaves, proofs, 0), leaves[0]);
  }

  /// forge-config: default.allow_internal_expect_revert = true
  function test_RevertWhen_EmptyLeaf() public {
    bytes32[] memory leaves = new bytes32[](0);
    bytes32[] memory proofs = new bytes32[](0);

    vm.expectRevert(abi.encodeWithSelector(MerkleMultiProof.LeavesCannotBeEmpty.selector));
    MerkleMultiProof._merkleRoot(leaves, proofs, 0);
  }

  /// forge-config: default.allow_internal_expect_revert = true
  function test_CVE_2023_34459() public {
    bytes32[] memory leaves = new bytes32[](2);
    // leaves[0] stays uninitialized, i.e., 0x000...0
    leaves[1] = "leaf";

    bytes32[] memory proof = new bytes32[](2);
    proof[0] = leaves[1];
    proof[1] = "will never be used";

    bytes32[] memory malicious = new bytes32[](2);
    malicious[0] = "malicious leaf";
    malicious[1] = "another malicious leaf";

    vm.expectRevert(abi.encodeWithSelector(MerkleMultiProof.InvalidProof.selector));
    MerkleMultiProof._merkleRoot(malicious, proof, 3);
    // Note, that without the revert the above computed root
    // would equal MerkleHelper.hashPair(leaves[0], leaves[1]).
  }
}
