// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Create2} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/Create2.sol";
import {BurnMintTokenPool} from "../../../pools/BurnMintTokenPool.sol";
import {RegistryModuleOwnerCustom} from "../../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {FactoryBurnMintERC20} from "../../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";
import {TokenPoolFactory} from "../../../tokenAdminRegistry/TokenPoolFactory/TokenPoolFactory.sol";
import {TokenAdminRegistrySetup} from "../TokenAdminRegistry/TokenAdminRegistrySetup.t.sol";

contract TokenPoolFactorySetup is TokenAdminRegistrySetup {
  using Create2 for bytes32;

  TokenPoolFactory internal s_tokenPoolFactory;
  RegistryModuleOwnerCustom internal s_registryModuleOwnerCustom;

  bytes internal s_poolInitCode;
  bytes internal s_poolInitArgs;

  bytes32 internal constant FAKE_SALT = keccak256(abi.encode("FAKE_SALT"));

  address internal s_rmnProxy = address(0x1234);

  bytes internal s_tokenCreationParams;
  bytes internal s_tokenInitCode;

  uint256 public constant PREMINT_AMOUNT = 100 ether;

  function setUp() public virtual override {
    TokenAdminRegistrySetup.setUp();

    s_registryModuleOwnerCustom = new RegistryModuleOwnerCustom(address(s_tokenAdminRegistry));
    s_tokenAdminRegistry.addRegistryModule(address(s_registryModuleOwnerCustom));

    s_tokenPoolFactory =
      new TokenPoolFactory(s_tokenAdminRegistry, s_registryModuleOwnerCustom, s_rmnProxy, address(s_sourceRouter));

    // Create Init Code for BurnMintERC20 TestToken with 18 decimals and supply cap of max uint256 value
    s_tokenCreationParams = abi.encode("TestToken", "TT", 18, type(uint256).max, PREMINT_AMOUNT, OWNER);

    s_tokenInitCode = abi.encodePacked(type(FactoryBurnMintERC20).creationCode, s_tokenCreationParams);

    s_poolInitCode = type(BurnMintTokenPool).creationCode;

    // Create Init Args for BurnMintTokenPool with no allowlist minus the token address
    address[] memory allowlist = new address[](1);
    allowlist[0] = OWNER;
    s_poolInitArgs = abi.encode(allowlist, address(0x1234), s_sourceRouter);
  }
}
