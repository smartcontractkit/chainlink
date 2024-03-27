pragma solidity 0.8.19;

import "./BaseTest.t.sol";
import {TrustedBlockhashStore} from "../dev/TrustedBlockhashStore.sol";
import {console} from "forge-std/console.sol";

contract TrustedBlockhashStoreTest is BaseTest {
  address internal constant LINK_WHALE = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;
  address internal constant LINK_WHALE_2 = 0xe9b2C5A6D9bA93dD354783a9De0a265da7551a20;
  TrustedBlockhashStore bhs;
  uint256 unreachableBlockNumber = 5;
  bytes32 unreachableBlockhash;

  function setUp() public override {
    BaseTest.setUp();

    // Get the blockhash for a block that later becomes unreachable in the EVM.
    vm.roll(10);
    unreachableBlockhash = blockhash(unreachableBlockNumber);

    // Fund our users.
    vm.roll(1000);
    vm.deal(LINK_WHALE, 10_000 ether);
    changePrank(LINK_WHALE);

    address[] memory whitelist = new address[](1);
    whitelist[0] = LINK_WHALE;
    bhs = new TrustedBlockhashStore(whitelist);
  }

  function testGenericBHSFunctions() public {
    // Should store.
    uint256 blockNumber = 999;
    bhs.store(blockNumber);
    assertEq(bhs.getBlockhash(blockNumber), blockhash(blockNumber));

    // Should store earliest.
    uint256 earliestBlockNumber = block.number - 256;
    bhs.storeEarliest();
    assertEq(bhs.getBlockhash(earliestBlockNumber), blockhash(earliestBlockNumber));
  }

  function testTrustedBHSFunctions() public {
    uint256 recentBlockNumber = 999;

    // Assume that the EVM cannot access the blockhash for block 5.
    uint256 unreachableBlock = 5;
    assertEq(blockhash(unreachableBlock), 0);

    // Store blockhash from whitelisted address;
    uint256[] memory invalidBlockNums = new uint256[](0);
    uint256[] memory blockNums = new uint256[](1);
    blockNums[0] = unreachableBlock;
    bytes32[] memory blockhashes = new bytes32[](1);
    blockhashes[0] = unreachableBlockhash;

    // Should not be able to store with invalid recent blockhash
    vm.expectRevert(TrustedBlockhashStore.InvalidRecentBlockhash.selector);
    bhs.storeTrusted(blockNums, blockhashes, recentBlockNumber, blockhash(998));

    // Should not be able to store or change whitelist for non-whitelisted address.
    changePrank(LINK_WHALE_2);
    vm.expectRevert(TrustedBlockhashStore.NotInWhitelist.selector);
    bhs.storeTrusted(blockNums, blockhashes, recentBlockNumber, blockhash(recentBlockNumber));
    vm.expectRevert("Only callable by owner");
    bhs.setWhitelist(new address[](0));

    // Should not store for a mismatched list of block numbers and hashes.
    changePrank(LINK_WHALE);
    vm.expectRevert(TrustedBlockhashStore.InvalidTrustedBlockhashes.selector);
    bhs.storeTrusted(invalidBlockNums, blockhashes, recentBlockNumber, blockhash(recentBlockNumber));

    // Should store unreachable blocks via whitelisted address.
    bhs.storeTrusted(blockNums, blockhashes, recentBlockNumber, blockhash(recentBlockNumber));
    assertEq(bhs.getBlockhash(unreachableBlock), unreachableBlockhash);

    // Change whitelist. Assert that the old whitelisted address can no longer store,
    // but the new one can.
    address[] memory newWhitelist = new address[](1);
    newWhitelist[0] = LINK_WHALE_2;
    bhs.setWhitelist(newWhitelist);

    vm.expectRevert(TrustedBlockhashStore.NotInWhitelist.selector);
    bhs.storeTrusted(blockNums, blockhashes, recentBlockNumber, blockhash(recentBlockNumber));

    changePrank(LINK_WHALE_2);
    bhs.storeTrusted(blockNums, blockhashes, recentBlockNumber, blockhash(recentBlockNumber));
  }
}
