// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {Internal} from "../../libraries/Internal.sol";
import {TokenSetup} from "../TokenSetup.t.sol";
import {PriceRegistry} from "../../PriceRegistry.sol";

contract PriceRegistrySetup is TokenSetup {
  uint112 internal constant USD_PER_GAS = 1e6; // 0.001 gwei
  uint112 internal constant USD_PER_DATA_AVAILABILITY_GAS = 1e9; // 1 gwei

  // Encode L1 gas price and L2 gas price into a packed price.
  // L1 gas price is left-shifted to the higher-order bits.
  uint224 internal constant PACKED_USD_PER_GAS =
    (uint224(USD_PER_DATA_AVAILABILITY_GAS) << Internal.GAS_PRICE_BITS) + USD_PER_GAS;

  PriceRegistry internal s_priceRegistry;
  // Cheat to store the price updates in storage since struct arrays aren't supported.
  bytes internal s_encodedInitialPriceUpdates;
  address internal s_weth;

  address[] internal s_sourceFeeTokens;
  uint224[] internal s_sourceTokenPrices;
  address[] internal s_destFeeTokens;
  uint224[] internal s_destTokenPrices;

  function setUp() public virtual override {
    TokenSetup.setUp();

    s_weth = s_sourceRouter.getWrappedNative();

    address[] memory sourceFeeTokens = new address[](3);
    sourceFeeTokens[0] = s_sourceTokens[0];
    sourceFeeTokens[1] = s_sourceTokens[1];
    sourceFeeTokens[2] = s_sourceRouter.getWrappedNative();
    s_sourceFeeTokens = sourceFeeTokens;

    uint224[] memory sourceTokenPrices = new uint224[](3);
    sourceTokenPrices[0] = 5e18;
    sourceTokenPrices[1] = 2000e18;
    sourceTokenPrices[2] = 2000e18;
    s_sourceTokenPrices = sourceTokenPrices;

    address[] memory destFeeTokens = new address[](3);
    destFeeTokens[0] = s_destTokens[0];
    destFeeTokens[1] = s_destTokens[1];
    destFeeTokens[2] = s_destRouter.getWrappedNative();
    s_destFeeTokens = destFeeTokens;

    uint224[] memory destTokenPrices = new uint224[](3);
    destTokenPrices[0] = 5e18;
    destTokenPrices[1] = 2000e18;
    destTokenPrices[2] = 2000e18;
    s_destTokenPrices = destTokenPrices;

    uint256 sourceTokenCount = sourceFeeTokens.length;
    uint256 destTokenCount = destFeeTokens.length;
    address[] memory pricedTokens = new address[](sourceTokenCount + destTokenCount);
    uint224[] memory tokenPrices = new uint224[](sourceTokenCount + destTokenCount);
    for (uint256 i = 0; i < sourceTokenCount; ++i) {
      pricedTokens[i] = sourceFeeTokens[i];
      tokenPrices[i] = sourceTokenPrices[i];
    }
    for (uint256 i = 0; i < destTokenCount; ++i) {
      pricedTokens[i + sourceTokenCount] = destFeeTokens[i];
      tokenPrices[i + sourceTokenCount] = destTokenPrices[i];
    }

    Internal.PriceUpdates memory priceUpdates = getPriceUpdatesStruct(pricedTokens, tokenPrices);
    priceUpdates.gasPriceUpdates = getSingleGasPriceUpdateStruct(DEST_CHAIN_SELECTOR, PACKED_USD_PER_GAS)
      .gasPriceUpdates;

    s_encodedInitialPriceUpdates = abi.encode(priceUpdates);
    address[] memory priceUpdaters = new address[](0);
    address[] memory feeTokens = new address[](2);
    feeTokens[0] = s_sourceTokens[0];
    feeTokens[1] = s_weth;
    s_priceRegistry = new PriceRegistry(priceUpdaters, feeTokens, uint32(TWELVE_HOURS));
    s_priceRegistry.updatePrices(priceUpdates);
  }
}

contract PriceRegistry_constructor is PriceRegistrySetup {
  function testSetupSuccess() public virtual {
    address[] memory priceUpdaters = new address[](2);
    priceUpdaters[0] = STRANGER;
    priceUpdaters[1] = OWNER;
    address[] memory feeTokens = new address[](2);
    feeTokens[0] = s_sourceTokens[0];
    feeTokens[1] = s_weth;

    s_priceRegistry = new PriceRegistry(priceUpdaters, feeTokens, uint32(TWELVE_HOURS));

    assertEq(feeTokens, s_priceRegistry.getFeeTokens());
    assertEq(uint32(TWELVE_HOURS), s_priceRegistry.getStalenessThreshold());
    assertEq(priceUpdaters, s_priceRegistry.getPriceUpdaters());
    assertEq(s_priceRegistry.typeAndVersion(), "PriceRegistry 1.2.0");
  }

  function testInvalidStalenessThresholdReverts() public {
    vm.expectRevert(PriceRegistry.InvalidStalenessThreshold.selector);
    s_priceRegistry = new PriceRegistry(new address[](0), new address[](0), 0);
  }
}

contract PriceRegistry_getTokenPrices is PriceRegistrySetup {
  function testGetTokenPricesSuccess() public {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));

    address[] memory tokens = new address[](3);
    tokens[0] = s_sourceTokens[0];
    tokens[1] = s_sourceTokens[1];
    tokens[2] = s_weth;

    Internal.TimestampedPackedUint224[] memory tokenPrices = s_priceRegistry.getTokenPrices(tokens);

    assertEq(tokenPrices.length, 3);
    assertEq(tokenPrices[0].value, priceUpdates.tokenPriceUpdates[0].usdPerToken);
    assertEq(tokenPrices[1].value, priceUpdates.tokenPriceUpdates[1].usdPerToken);
    assertEq(tokenPrices[2].value, priceUpdates.tokenPriceUpdates[2].usdPerToken);
  }
}

contract PriceRegistry_getValidatedTokenPrice is PriceRegistrySetup {
  function testGetValidatedTokenPriceSuccess() public {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));
    address token = priceUpdates.tokenPriceUpdates[0].sourceToken;

    uint224 tokenPrice = s_priceRegistry.getValidatedTokenPrice(token);

    assertEq(priceUpdates.tokenPriceUpdates[0].usdPerToken, tokenPrice);
  }

  function testStaleFeeTokenReverts() public {
    vm.warp(block.timestamp + TWELVE_HOURS + 1);

    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.StaleTokenPrice.selector, s_sourceTokens[0], TWELVE_HOURS, TWELVE_HOURS + 1)
    );
    s_priceRegistry.getValidatedTokenPrice(s_sourceTokens[0]);
  }

  function testTokenNotSupportedReverts() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_priceRegistry.getValidatedTokenPrice(DUMMY_CONTRACT_ADDRESS);
  }
}

contract PriceRegistry_applyPriceUpdatersUpdates is PriceRegistrySetup {
  event PriceUpdaterSet(address indexed priceUpdater);
  event PriceUpdaterRemoved(address indexed priceUpdater);

  function testApplyPriceUpdaterUpdatesSuccess() public {
    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = STRANGER;

    vm.expectEmit();
    emit PriceUpdaterSet(STRANGER);

    s_priceRegistry.applyPriceUpdatersUpdates(priceUpdaters, new address[](0));
    assertEq(s_priceRegistry.getPriceUpdaters().length, 1);
    assertEq(s_priceRegistry.getPriceUpdaters()[0], STRANGER);

    // add same priceUpdater is no-op
    s_priceRegistry.applyPriceUpdatersUpdates(priceUpdaters, new address[](0));
    assertEq(s_priceRegistry.getPriceUpdaters().length, 1);
    assertEq(s_priceRegistry.getPriceUpdaters()[0], STRANGER);

    vm.expectEmit();
    emit PriceUpdaterRemoved(STRANGER);

    s_priceRegistry.applyPriceUpdatersUpdates(new address[](0), priceUpdaters);
    assertEq(s_priceRegistry.getPriceUpdaters().length, 0);

    // removing already removed priceUpdater is no-op
    s_priceRegistry.applyPriceUpdatersUpdates(new address[](0), priceUpdaters);
    assertEq(s_priceRegistry.getPriceUpdaters().length, 0);
  }

  function testOnlyCallableByOwnerReverts() public {
    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = STRANGER;
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_priceRegistry.applyPriceUpdatersUpdates(priceUpdaters, new address[](0));
  }
}

contract PriceRegistry_applyFeeTokensUpdates is PriceRegistrySetup {
  event FeeTokenAdded(address indexed feeToken);
  event FeeTokenRemoved(address indexed feeToken);

  function testApplyFeeTokensUpdatesSuccess() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];

    vm.expectEmit();
    emit FeeTokenAdded(feeTokens[0]);

    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_priceRegistry.getFeeTokens().length, 3);
    assertEq(s_priceRegistry.getFeeTokens()[2], feeTokens[0]);

    // add same feeToken is no-op
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_priceRegistry.getFeeTokens().length, 3);
    assertEq(s_priceRegistry.getFeeTokens()[2], feeTokens[0]);

    vm.expectEmit();
    emit FeeTokenRemoved(feeTokens[0]);

    s_priceRegistry.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_priceRegistry.getFeeTokens().length, 2);

    // removing already removed feeToken is no-op
    s_priceRegistry.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_priceRegistry.getFeeTokens().length, 2);
  }

  function testOnlyCallableByOwnerReverts() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = STRANGER;
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));
  }
}

contract PriceRegistry_updatePrices is PriceRegistrySetup {
  event UsdPerTokenUpdated(address indexed token, uint256 value, uint256 timestamp);
  event UsdPerUnitGasUpdated(uint64 indexed destChain, uint256 value, uint256 timestamp);

  function testOnlyTokenPriceSuccess() public {
    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](1),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    update.tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});

    vm.expectEmit();
    emit UsdPerTokenUpdated(
      update.tokenPriceUpdates[0].sourceToken,
      update.tokenPriceUpdates[0].usdPerToken,
      block.timestamp
    );

    s_priceRegistry.updatePrices(update);

    assertEq(s_priceRegistry.getTokenPrice(s_sourceTokens[0]).value, update.tokenPriceUpdates[0].usdPerToken);
  }

  function testOnlyGasPriceSuccess() public {
    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](1)
    });
    update.gasPriceUpdates[0] = Internal.GasPriceUpdate({
      destChainSelector: DEST_CHAIN_SELECTOR,
      usdPerUnitGas: 2000e18
    });

    vm.expectEmit();
    emit UsdPerUnitGasUpdated(
      update.gasPriceUpdates[0].destChainSelector,
      update.gasPriceUpdates[0].usdPerUnitGas,
      block.timestamp
    );

    s_priceRegistry.updatePrices(update);

    assertEq(
      s_priceRegistry.getDestinationChainGasPrice(DEST_CHAIN_SELECTOR).value,
      update.gasPriceUpdates[0].usdPerUnitGas
    );
  }

  function testUpdateMultiplePricesSuccess() public {
    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](3);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});
    tokenPriceUpdates[1] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[1], usdPerToken: 1800e18});
    tokenPriceUpdates[2] = Internal.TokenPriceUpdate({sourceToken: address(12345), usdPerToken: 1e18});

    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](3);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: 2e6});
    gasPriceUpdates[1] = Internal.GasPriceUpdate({destChainSelector: SOURCE_CHAIN_SELECTOR, usdPerUnitGas: 2000e18});
    gasPriceUpdates[2] = Internal.GasPriceUpdate({destChainSelector: 12345, usdPerUnitGas: 1e18});

    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      gasPriceUpdates: gasPriceUpdates
    });

    for (uint256 i = 0; i < tokenPriceUpdates.length; ++i) {
      vm.expectEmit();
      emit UsdPerTokenUpdated(
        update.tokenPriceUpdates[i].sourceToken,
        update.tokenPriceUpdates[i].usdPerToken,
        block.timestamp
      );
    }
    for (uint256 i = 0; i < gasPriceUpdates.length; ++i) {
      vm.expectEmit();
      emit UsdPerUnitGasUpdated(
        update.gasPriceUpdates[i].destChainSelector,
        update.gasPriceUpdates[i].usdPerUnitGas,
        block.timestamp
      );
    }

    s_priceRegistry.updatePrices(update);

    for (uint256 i = 0; i < tokenPriceUpdates.length; ++i) {
      assertEq(
        s_priceRegistry.getTokenPrice(update.tokenPriceUpdates[i].sourceToken).value,
        tokenPriceUpdates[i].usdPerToken
      );
    }
    for (uint256 i = 0; i < gasPriceUpdates.length; ++i) {
      assertEq(
        s_priceRegistry.getDestinationChainGasPrice(update.gasPriceUpdates[i].destChainSelector).value,
        gasPriceUpdates[i].usdPerUnitGas
      );
    }
  }

  // Reverts

  function testOnlyCallableByUpdaterOrOwnerReverts() public {
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });

    changePrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.OnlyCallableByUpdaterOrOwner.selector));
    s_priceRegistry.updatePrices(priceUpdates);
  }
}

contract PriceRegistry_convertTokenAmount is PriceRegistrySetup {
  function testConvertTokenAmountSuccess() public {
    Internal.PriceUpdates memory initialPriceUpdates = abi.decode(
      s_encodedInitialPriceUpdates,
      (Internal.PriceUpdates)
    );
    uint256 amount = 3e16;
    uint256 conversionRate = (uint256(initialPriceUpdates.tokenPriceUpdates[2].usdPerToken) * 1e18) /
      uint256(initialPriceUpdates.tokenPriceUpdates[0].usdPerToken);
    uint256 expected = (amount * conversionRate) / 1e18;
    assertEq(s_priceRegistry.convertTokenAmount(s_weth, amount, s_sourceTokens[0]), expected);
  }

  function testFuzz_ConvertTokenAmountSuccess(
    uint256 feeTokenAmount,
    uint224 usdPerFeeToken,
    uint160 usdPerLinkToken,
    uint224 usdPerUnitGas
  ) public {
    vm.assume(usdPerFeeToken > 0);
    vm.assume(usdPerLinkToken > 0);
    // We bound the max fees to be at most uint96.max link.
    feeTokenAmount = bound(feeTokenAmount, 0, (uint256(type(uint96).max) * usdPerLinkToken) / usdPerFeeToken);

    address feeToken = address(1);
    address linkToken = address(2);
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = feeToken;
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](2);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: feeToken, usdPerToken: usdPerFeeToken});
    tokenPriceUpdates[1] = Internal.TokenPriceUpdate({sourceToken: linkToken, usdPerToken: usdPerLinkToken});

    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](1);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({
      destChainSelector: DEST_CHAIN_SELECTOR,
      usdPerUnitGas: usdPerUnitGas
    });

    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      gasPriceUpdates: gasPriceUpdates
    });

    s_priceRegistry.updatePrices(priceUpdates);

    uint256 linkFee = s_priceRegistry.convertTokenAmount(feeToken, feeTokenAmount, linkToken);
    assertEq(linkFee, (feeTokenAmount * usdPerFeeToken) / usdPerLinkToken);
  }

  // Reverts

  function testStaleFeeTokenReverts() public {
    vm.warp(block.timestamp + TWELVE_HOURS + 1);

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](1);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    s_priceRegistry.updatePrices(priceUpdates);

    vm.expectRevert(
      abi.encodeWithSelector(
        PriceRegistry.StaleTokenPrice.selector,
        s_weth,
        uint128(TWELVE_HOURS),
        uint128(TWELVE_HOURS + 1)
      )
    );
    s_priceRegistry.convertTokenAmount(s_weth, 3e16, s_sourceTokens[0]);
  }

  function testLinkTokenNotSupportedReverts() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_priceRegistry.convertTokenAmount(DUMMY_CONTRACT_ADDRESS, 3e16, s_sourceTokens[0]);

    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_priceRegistry.convertTokenAmount(s_sourceTokens[0], 3e16, DUMMY_CONTRACT_ADDRESS);
  }

  function testStaleLinkTokenReverts() public {
    vm.warp(block.timestamp + TWELVE_HOURS + 1);

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](1);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_weth, usdPerToken: 18e17});
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    s_priceRegistry.updatePrices(priceUpdates);

    vm.expectRevert(
      abi.encodeWithSelector(
        PriceRegistry.StaleTokenPrice.selector,
        s_sourceTokens[0],
        uint128(TWELVE_HOURS),
        uint128(TWELVE_HOURS + 1)
      )
    );
    s_priceRegistry.convertTokenAmount(s_weth, 3e16, s_sourceTokens[0]);
  }
}

contract PriceRegistry_getTokenAndGasPrices is PriceRegistrySetup {
  function testGetFeeTokenAndGasPricesSuccess() public {
    (uint224 feeTokenPrice, uint224 gasPrice) = s_priceRegistry.getTokenAndGasPrices(
      s_sourceFeeToken,
      DEST_CHAIN_SELECTOR
    );

    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));

    assertEq(feeTokenPrice, s_sourceTokenPrices[0]);
    assertEq(gasPrice, priceUpdates.gasPriceUpdates[0].usdPerUnitGas);
  }

  function testZeroGasPriceSuccess() public {
    uint64 zeroGasDestChainSelector = 345678;
    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](1);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: zeroGasDestChainSelector, usdPerUnitGas: 0});

    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: gasPriceUpdates
    });
    s_priceRegistry.updatePrices(priceUpdates);

    (, uint224 gasPrice) = s_priceRegistry.getTokenAndGasPrices(s_sourceFeeToken, zeroGasDestChainSelector);

    assertEq(gasPrice, priceUpdates.gasPriceUpdates[0].usdPerUnitGas);
  }

  function testUnsupportedChainReverts() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.ChainNotSupported.selector, DEST_CHAIN_SELECTOR + 1));
    s_priceRegistry.getTokenAndGasPrices(s_sourceTokens[0], DEST_CHAIN_SELECTOR + 1);
  }

  function testStaleGasPriceReverts() public {
    uint256 diff = TWELVE_HOURS + 1;
    vm.warp(block.timestamp + diff);
    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.StaleGasPrice.selector, DEST_CHAIN_SELECTOR, TWELVE_HOURS, diff)
    );
    s_priceRegistry.getTokenAndGasPrices(s_sourceTokens[0], DEST_CHAIN_SELECTOR);
  }

  function testStaleTokenPriceReverts() public {
    uint256 diff = TWELVE_HOURS + 1;
    vm.warp(block.timestamp + diff);

    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](1);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({
      destChainSelector: DEST_CHAIN_SELECTOR,
      usdPerUnitGas: PACKED_USD_PER_GAS
    });

    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: gasPriceUpdates
    });
    s_priceRegistry.updatePrices(priceUpdates);

    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.StaleTokenPrice.selector, s_sourceTokens[0], TWELVE_HOURS, diff)
    );
    s_priceRegistry.getTokenAndGasPrices(s_sourceTokens[0], DEST_CHAIN_SELECTOR);
  }
}
