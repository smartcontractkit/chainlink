// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {AuthorizedCallers} from "../../../shared/access/AuthorizedCallers.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_updatePrices is FeeQuoterSetup {
  function test_OnlyTokenPrice() public {
    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](1),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    update.tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(
      update.tokenPriceUpdates[0].sourceToken, update.tokenPriceUpdates[0].usdPerToken, block.timestamp
    );

    s_feeQuoter.updatePrices(update);

    assertEq(s_feeQuoter.getTokenPrice(s_sourceTokens[0]).value, update.tokenPriceUpdates[0].usdPerToken);
  }

  function test_OnlyGasPrice() public {
    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](1)
    });
    update.gasPriceUpdates[0] =
      Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: 2000e18});

    vm.expectEmit();
    emit FeeQuoter.UsdPerUnitGasUpdated(
      update.gasPriceUpdates[0].destChainSelector, update.gasPriceUpdates[0].usdPerUnitGas, block.timestamp
    );

    s_feeQuoter.updatePrices(update);

    assertEq(
      s_feeQuoter.getDestinationChainGasPrice(DEST_CHAIN_SELECTOR).value, update.gasPriceUpdates[0].usdPerUnitGas
    );
  }

  function test_UpdateMultiplePrices() public {
    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](3);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});
    tokenPriceUpdates[1] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[1], usdPerToken: 1800e18});
    tokenPriceUpdates[2] = Internal.TokenPriceUpdate({sourceToken: address(12345), usdPerToken: 1e18});

    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](3);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: 2e6});
    gasPriceUpdates[1] = Internal.GasPriceUpdate({destChainSelector: SOURCE_CHAIN_SELECTOR, usdPerUnitGas: 2000e18});
    gasPriceUpdates[2] = Internal.GasPriceUpdate({destChainSelector: 12345, usdPerUnitGas: 1e18});

    Internal.PriceUpdates memory update =
      Internal.PriceUpdates({tokenPriceUpdates: tokenPriceUpdates, gasPriceUpdates: gasPriceUpdates});

    for (uint256 i = 0; i < tokenPriceUpdates.length; ++i) {
      vm.expectEmit();
      emit FeeQuoter.UsdPerTokenUpdated(
        update.tokenPriceUpdates[i].sourceToken, update.tokenPriceUpdates[i].usdPerToken, block.timestamp
      );
    }
    for (uint256 i = 0; i < gasPriceUpdates.length; ++i) {
      vm.expectEmit();
      emit FeeQuoter.UsdPerUnitGasUpdated(
        update.gasPriceUpdates[i].destChainSelector, update.gasPriceUpdates[i].usdPerUnitGas, block.timestamp
      );
    }

    s_feeQuoter.updatePrices(update);

    for (uint256 i = 0; i < tokenPriceUpdates.length; ++i) {
      assertEq(
        s_feeQuoter.getTokenPrice(update.tokenPriceUpdates[i].sourceToken).value, tokenPriceUpdates[i].usdPerToken
      );
    }
    for (uint256 i = 0; i < gasPriceUpdates.length; ++i) {
      assertEq(
        s_feeQuoter.getDestinationChainGasPrice(update.gasPriceUpdates[i].destChainSelector).value,
        gasPriceUpdates[i].usdPerUnitGas
      );
    }
  }

  function test_UpdatableByAuthorizedCaller() public {
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](1),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    priceUpdates.tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});

    // Revert when caller is not authorized
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_feeQuoter.updatePrices(priceUpdates);

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = STRANGER;
    vm.startPrank(OWNER);
    s_feeQuoter.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: priceUpdaters, removedCallers: new address[](0)})
    );

    // Stranger is now an authorized caller to update prices
    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(
      priceUpdates.tokenPriceUpdates[0].sourceToken, priceUpdates.tokenPriceUpdates[0].usdPerToken, block.timestamp
    );
    s_feeQuoter.updatePrices(priceUpdates);

    assertEq(s_feeQuoter.getTokenPrice(s_sourceTokens[0]).value, priceUpdates.tokenPriceUpdates[0].usdPerToken);

    vm.startPrank(OWNER);
    s_feeQuoter.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: new address[](0), removedCallers: priceUpdaters})
    );

    // Revert when authorized caller is removed
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_feeQuoter.updatePrices(priceUpdates);
  }

  // Reverts

  function test_RevertWhen_OnlyCallableByUpdater() public {
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });

    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_feeQuoter.updatePrices(priceUpdates);
  }
}
