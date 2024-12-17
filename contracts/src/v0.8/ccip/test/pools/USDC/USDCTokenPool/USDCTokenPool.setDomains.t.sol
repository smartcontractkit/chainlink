// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../../shared/access/Ownable2Step.sol";
import {USDCTokenPool} from "../../../../pools/USDC/USDCTokenPool.sol";
import {USDCTokenPoolSetup} from "./USDCTokenPoolSetup.t.sol";

contract USDCTokenPool_setDomains is USDCTokenPoolSetup {
  mapping(uint64 destChainSelector => USDCTokenPool.Domain domain) private s_chainToDomain;

  // Setting lower fuzz run as 256 runs was causing differing gas results in snapshot.
  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function testFuzz_SetDomains_Success(
    bytes32[5] calldata allowedCallers,
    uint32[5] calldata domainIdentifiers,
    uint64[5] calldata destChainSelectors
  ) public {
    uint256 numberOfDomains = allowedCallers.length;
    USDCTokenPool.DomainUpdate[] memory domainUpdates = new USDCTokenPool.DomainUpdate[](numberOfDomains);
    for (uint256 i = 0; i < numberOfDomains; ++i) {
      vm.assume(allowedCallers[i] != bytes32(0) && domainIdentifiers[i] != 0 && destChainSelectors[i] != 0);

      domainUpdates[i] = USDCTokenPool.DomainUpdate({
        allowedCaller: allowedCallers[i],
        domainIdentifier: domainIdentifiers[i],
        destChainSelector: destChainSelectors[i],
        enabled: true
      });

      s_chainToDomain[destChainSelectors[i]] =
        USDCTokenPool.Domain({domainIdentifier: domainIdentifiers[i], allowedCaller: allowedCallers[i], enabled: true});
    }

    vm.expectEmit();
    emit USDCTokenPool.DomainsSet(domainUpdates);

    s_usdcTokenPool.setDomains(domainUpdates);

    for (uint256 i = 0; i < numberOfDomains; ++i) {
      USDCTokenPool.Domain memory expected = s_chainToDomain[destChainSelectors[i]];
      USDCTokenPool.Domain memory got = s_usdcTokenPool.getDomain(destChainSelectors[i]);
      assertEq(got.allowedCaller, expected.allowedCaller);
      assertEq(got.domainIdentifier, expected.domainIdentifier);
    }
  }

  // Reverts

  function test_RevertWhen_OnlyOwner() public {
    USDCTokenPool.DomainUpdate[] memory domainUpdates = new USDCTokenPool.DomainUpdate[](0);

    vm.startPrank(STRANGER);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);

    s_usdcTokenPool.setDomains(domainUpdates);
  }

  function test_RevertWhen_InvalidDomain() public {
    bytes32 validCaller = bytes32(uint256(25));
    // Ensure valid domain works
    USDCTokenPool.DomainUpdate[] memory domainUpdates = new USDCTokenPool.DomainUpdate[](1);
    domainUpdates[0] = USDCTokenPool.DomainUpdate({
      allowedCaller: validCaller,
      domainIdentifier: 0, // ensures 0 is valid, as this is eth mainnet
      destChainSelector: 45690,
      enabled: true
    });

    s_usdcTokenPool.setDomains(domainUpdates);

    // Make update invalid on allowedCaller
    domainUpdates[0].allowedCaller = bytes32(0);
    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.InvalidDomain.selector, domainUpdates[0]));

    s_usdcTokenPool.setDomains(domainUpdates);

    // Make valid again
    domainUpdates[0].allowedCaller = validCaller;

    // Make invalid on destChainSelector
    domainUpdates[0].destChainSelector = 0;
    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.InvalidDomain.selector, domainUpdates[0]));

    s_usdcTokenPool.setDomains(domainUpdates);
  }
}
