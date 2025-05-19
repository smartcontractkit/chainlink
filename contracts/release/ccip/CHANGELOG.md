# @chainlink/contracts-ccip

## 1.6.0

### Minor Changes

- [#15804](https://github.com/smartcontractkit/chainlink/pull/15804) [`46ef625`](https://github.com/smartcontractkit/chainlink/commit/46ef62537ea0389a86de03465253a8629766c2c9) - #feature Add two new pool types: Siloed-LockRelease and BurnToAddress and fix bug in HybridUSDCTokenPool for transferLiqudity #bugfix

  PR issue: CCIP-4723

  Solidity Review issue: CCIP-3966

- [#15811](https://github.com/smartcontractkit/chainlink/pull/15811) [`4e28497`](https://github.com/smartcontractkit/chainlink/commit/4e284976ea8ca7c0e355efce6336742d70918ac2) - Update FeeQuoter to support Solana chain families #feature

  PR issue: CCIP-4687

  Solidity Review issue: CCIP-3966

- [#14924](https://github.com/smartcontractkit/chainlink/pull/14924) [`161d227`](https://github.com/smartcontractkit/chainlink/commit/161d227575bbeca0119e5eee8e5a54cf3b4df677) - New contract for deploying, CCIP-compatible token pools and configuring with the tokenAdminRegistry, and a new ERC20 with constructor compatible with the factory's deployment pattern. #internal

  PR issue: CCIP-3171

  Solidity Review issue: CCIP-3966

- [#15099](https://github.com/smartcontractkit/chainlink/pull/15099) [`9c79488`](https://github.com/smartcontractkit/chainlink/commit/9c79488e4259bd59aa4d25b2be9c2ffd9390333d) - #internal Add new event in setRateLimitAdmin for Atlas

  PR issue: CCIP-4099

  Solidity Review issue: CCIP-3966

- [#15737](https://github.com/smartcontractkit/chainlink/pull/15737) [`631cd8f`](https://github.com/smartcontractkit/chainlink/commit/631cd8fae34c73801e223f9ef96f23a032cf407f) - #internal Account for tokenTransferBytesOverhead in exec cost

  PR issue: CCIP-4646

  Solidity Review issue: CCIP-3966

- [#16298](https://github.com/smartcontractkit/chainlink/pull/16298) [`ab06bf0`](https://github.com/smartcontractkit/chainlink/commit/ab06bf0277b5641c24596d05351dd23df544a72c) - release CCIP 1.6, remove -dev suffix

  PR issue: CCIP-5181

  Solidity Review issue: CCIP-3966

- [#15123](https://github.com/smartcontractkit/chainlink/pull/15123) [`72da397`](https://github.com/smartcontractkit/chainlink/commit/72da397b9b1a8baa95129ffe635e5d60852e9ebc) - Add a new contract, BurnMintERC20, which is basically just our ERC677 implementation without the transferAndCall function. #internal

  PR issue: CCIP-4130

  Solidity Review issue: CCIP-3966

- [#14696](https://github.com/smartcontractkit/chainlink/pull/14696) [`072bfb6`](https://github.com/smartcontractkit/chainlink/commit/072bfb667a4e1f7cc0c874409ebfe6ef7f7b6cbe) - #internal add getNextDonId() and getNodes(bytes32[] calldata p2pIds) in CapabilitiesRegistry and define interface for node info

  PR issue: CCIP-3569

- release CCIP 1.6, remove -dev suffix

- [#14981](https://github.com/smartcontractkit/chainlink/pull/14981) [`a0309c9`](https://github.com/smartcontractkit/chainlink/commit/a0309c9a87ca7bcbe50db2fa272f9fab024bef13) - #internal skip stale price update from keystone instead of reverting

  PR issue: CCIP-3795

  Solidity Review issue: CCIP-3966

- [#15165](https://github.com/smartcontractkit/chainlink/pull/15165) [`03827b9`](https://github.com/smartcontractkit/chainlink/commit/03827b9d291cb11501162f9eafa1eba5619cabc9) - #internal CCIP test restructuring

  PR issue: CCIP-4116

  Solidity Review issue: CCIP-3966

- [#14817](https://github.com/smartcontractkit/chainlink/pull/14817) [`974def5`](https://github.com/smartcontractkit/chainlink/commit/974def52d97ee548b7568cf2facbc556dfa0e797) - change minSigners to f in RMNRemote/RMNHome

  PR issue: CCIP-3614

  Solidity Review issue: CCIP-3966

- [#15448](https://github.com/smartcontractkit/chainlink/pull/15448) [`c1ee7ab`](https://github.com/smartcontractkit/chainlink/commit/c1ee7ab715b524df6b580593e18f51637bd1500d) - #internal Add supportsInterface to FeeQuoter for Keystone

  PR issue: CCIP-4359

  Solidity Review issue: CCIP-3966

- [#14877](https://github.com/smartcontractkit/chainlink/pull/14877) [`317f930`](https://github.com/smartcontractkit/chainlink/commit/317f93014d9b5deb76e5b54685b020f44be9b46e) - #internal fix sender encoding and comments in CCIP Any2EVMMEssage and corrected comments

  PR issue: CCIP-3899

- [#15504](https://github.com/smartcontractkit/chainlink/pull/15504) [`437ef64`](https://github.com/smartcontractkit/chainlink/commit/437ef640db1d4455e7b4d90868c9e5c4d62054df) - #internal Fix gas estimation by adding a reverting clause

  PR issue: CCIP-4223

  Solidity Review issue: CCIP-3966

- [#14951](https://github.com/smartcontractkit/chainlink/pull/14951) [`2fab939`](https://github.com/smartcontractkit/chainlink/commit/2fab939ad3fbd49f79e96c177a9ffb11387f397e) - #internal Add a missing condition for the Execution plugin in the \_afterOCR3ConfigSet function. Now, the function correctly reverts if signature verification is enabled for the Execution plugin

  PR issue: CCIP-3799

  Solidity Review issue: CCIP-3966

- [#14918](https://github.com/smartcontractkit/chainlink/pull/14918) [`1c53ec2`](https://github.com/smartcontractkit/chainlink/commit/1c53ec25ed6c6bfee37e9052236fe595eeef8d89) - #internal applyRateLimiterConfigUpdates validation

  PR issue: CCIP-3797

  Solidity Review issue: CCIP-3966

- [#15983](https://github.com/smartcontractkit/chainlink/pull/15983) [`78d548d`](https://github.com/smartcontractkit/chainlink/commit/78d548d6101afd7fc82c2568a3665b4b57d5b1eb) - #internal Add EVM extraArgs encode & decode to MessageHasher

  PR issue: CCIP-4918

  Solidity Review issue: CCIP-3966

- [#15575](https://github.com/smartcontractkit/chainlink/pull/15575) [`ef0dd1c`](https://github.com/smartcontractkit/chainlink/commit/ef0dd1c1b9c9650a3dba7b5c3d5f2e81db90a555) - #internal make gas for call exact check immutable

  PR issue: CCIP-4477

  Solidity Review issue: CCIP-3966

### Patch Changes

- [#15422](https://github.com/smartcontractkit/chainlink/pull/15422) [`3cecd5f`](https://github.com/smartcontractkit/chainlink/commit/3cecd5f7dd5f8eaa624eafd8db701475facb8617) - add legacy fallback to RMN

  PR issue: CCIP-4261

  Solidity Review issue: CCIP-3966

- [#14798](https://github.com/smartcontractkit/chainlink/pull/14798) [`26e22eb`](https://github.com/smartcontractkit/chainlink/commit/26e22eb6cfc320d981c91b1bfeb87c3c645f10d6) - #internal fixing comments and data availability bytes length calculation

  PR issue : CCIP-3785

- [#15829](https://github.com/smartcontractkit/chainlink/pull/15829) [`6e65dee`](https://github.com/smartcontractkit/chainlink/commit/6e65deecae053ee1e885da7ce6d1d308364ced1d) - #internal Generate gethwrappers through Foundry instead of solc-select via python

  PR issue: CCIP-4737

  Solidity Review issue: CCIP-3966

- [#16201](https://github.com/smartcontractkit/chainlink/pull/16201) [`8ab88c2`](https://github.com/smartcontractkit/chainlink/commit/8ab88c2d9d5e5df4f8496763e74b0bb9c4c183a9) - #internal enable both blessed and unblessed roots in a single commit report

  PR issue: CCIP-5140

  Solidity Review issue: CCIP-3966

- [#15523](https://github.com/smartcontractkit/chainlink/pull/15523) [`baa225e`](https://github.com/smartcontractkit/chainlink/commit/baa225e7614f504a70cb51a8b0d0a98402971268) - remove legacy curse check from RMNRemote isCursed() method #bugfix

  PR issue: CCIP-4476

  Solidity Review issue: CCIP-3966

- [#16226](https://github.com/smartcontractkit/chainlink/pull/16226) [`b92a304`](https://github.com/smartcontractkit/chainlink/commit/b92a304b55570fdeb87a4f3e840819d1cf33f043) - #internal Add Support for SVM ATA to USDCTokenPool and send correct tokenReceiver to token pools

  PR issue: CCIP-5139

  Solidity Review issue: CCIP-3966

- [#16017](https://github.com/smartcontractkit/chainlink/pull/16017) [`17a9e2a`](https://github.com/smartcontractkit/chainlink/commit/17a9e2af202a313ea734a556af3f4460d3e8c795) - #internal decouple LiquidityManager tests with LockReleaseTokenPool + Test Rename

  PR issue: CCIP-4428

  Solidity Review issue: CCIP-3966

- [#15458](https://github.com/smartcontractkit/chainlink/pull/15458) [`5a3a99b`](https://github.com/smartcontractkit/chainlink/commit/5a3a99b7982dbf0f8aa57654dac356419736bd30) - Add token address to TokenHandlingError

  PR issue: CCIP-4174

  Solidity Review issue: CCIP-3966

- [#16141](https://github.com/smartcontractkit/chainlink/pull/16141) [`16a7985`](https://github.com/smartcontractkit/chainlink/commit/16a7985836d2055aca62d4cd331f2d374792d0d1) - #internal move multiple weth9 implementations to vendor

  PR issue: CCIP-5081

  Solidity Review issue: CCIP-3966

- [#14805](https://github.com/smartcontractkit/chainlink/pull/14805) [`b17c09d`](https://github.com/smartcontractkit/chainlink/commit/b17c09d252f75071f6ec54b3389257c1f27df9b2) - Gas optimizations and comment cleanup #internal

  PR issue: CCIP-3736

  Solidity Review issue: CCIP-3966

- [#15904](https://github.com/smartcontractkit/chainlink/pull/15904) [`5314b41`](https://github.com/smartcontractkit/chainlink/commit/5314b4127404aa2f402cd396c3088923136ac9d0) - #internal add EIP-7623 support

  PR issue: CCIP-4761

  Solidity Review issue: CCIP-3966

- [#15357](https://github.com/smartcontractkit/chainlink/pull/15357) [`18cb44e`](https://github.com/smartcontractkit/chainlink/commit/18cb44e891a00edff7486640ffc8e0c9275a04f8) - #added new function to CCIPReaderTester getLatestPriceSequenceNumber

  PR issue: CCIP-4239

  Solidity Review issue: CCIP-3966

- [#15067](https://github.com/smartcontractkit/chainlink/pull/15067) [`eeb58e2`](https://github.com/smartcontractkit/chainlink/commit/eeb58e2d3ae5d84826b31eaf805b30f722a8e87d) - #feature adds OZ AccessControl support to the registry module

  PR issue: CCIP-4105

  Solidity Review issue: CCIP-3966

- [#14972](https://github.com/smartcontractkit/chainlink/pull/14972) [`6db71d3`](https://github.com/smartcontractkit/chainlink/commit/6db71d32d756eba147ec69f385252aa51589d517) - #internal minor nits, allow Router updates even when the offRamp has been used. Remove getRouter from onRamp

  PR issue: CCIP-4010

  Solidity Review issue: CCIP-3966

- [#15293](https://github.com/smartcontractkit/chainlink/pull/15293) [`4665863`](https://github.com/smartcontractkit/chainlink/commit/466586309a8cbbfc1c793ff1021b7fcd3522dd3e) - allow multiple remote pools per chain selector

  PR issue: CCIP-4269

  Solidity Review issue: CCIP-3966

- [#14922](https://github.com/smartcontractkit/chainlink/pull/14922) [`42db9fd`](https://github.com/smartcontractkit/chainlink/commit/42db9fd17ca32d554f8bc9d0c6aab01e5c0e8c81) - Minor fixes to formatting, pragma, imports, etc. for Hybrid USDC Token Pools #bugfix

  PR issue: CCIP-3014

  Solidity Review issue: CCIP-3966

- [#14969](https://github.com/smartcontractkit/chainlink/pull/14969) [`ccd9956`](https://github.com/smartcontractkit/chainlink/commit/ccd9956ac6cec25770c93e57e283aa2a5ebb6737) - remove rawVs from RMNRemote

  PR issue: CCIP-4015

  Solidity Review issue: CCIP-3966

- [#14845](https://github.com/smartcontractkit/chainlink/pull/14845) [`3f955bf`](https://github.com/smartcontractkit/chainlink/commit/3f955bfde18bbad19f7195da1da179d04275a873) - #internal Minor fixes and changing the order of removes/adds in feeToken config

  CCIP-3730
  CCIP-3727
  CCIP-3725

- [#14809](https://github.com/smartcontractkit/chainlink/pull/14809) [`082e6fc`](https://github.com/smartcontractkit/chainlink/commit/082e6fc8f918dc69e9a5a2acbff6644dca74f2b1) - Modified TokenPriceFeedConfig to support tokens with zero decimals #bugfix

  PR issue: CCIP-3723

  Solidity Review issue: CCIP-3966

- [#15006](https://github.com/smartcontractkit/chainlink/pull/15006) [`33dba89`](https://github.com/smartcontractkit/chainlink/commit/33dba89a3beda1138fe7332a4947411e748fd99f) - minor gas optimizations and input sanity checks for CCIPHome #bugfix

  PR issue: CCIP-4075

  Solidity Review issue: CCIP-3966

- [#16102](https://github.com/smartcontractkit/chainlink/pull/16102) [`57ca0fb`](https://github.com/smartcontractkit/chainlink/commit/57ca0fb8f3c74fec461e2168d362b62374e26b63) - Comment and parameter validation fixes and remove outstandingTokens from BurnToAddressMintTokenPool #bugfix

  PR issue: CCIP-5061

  Solidity Review issue: CCIP-3966

- [#16090](https://github.com/smartcontractkit/chainlink/pull/16090) [`2efb46c`](https://github.com/smartcontractkit/chainlink/commit/2efb46c470867127671d39b214297f9419b24a48) - #internal Minor FeeQuoter audit fixes

  PR issue: CCIP-5046

  Solidity Review issue: CCIP-3966

- [#16175](https://github.com/smartcontractkit/chainlink/pull/16175) [`1c76b30`](https://github.com/smartcontractkit/chainlink/commit/1c76b30f78fdb542d9f5b7ee4c7238e6b8f408d2) - #internal cap max accounts in svm extra args

  PR issue: CCIP-5111

  Solidity Review issue: CCIP-3966

- [#15386](https://github.com/smartcontractkit/chainlink/pull/15386) [`62c2376`](https://github.com/smartcontractkit/chainlink/commit/62c23768cd483b179301625603a785dd773f2c78) - Modify TokenPool.sol function setChainRateLimiterConfig to now accept an array of configs and set sequentially. Requested by front-end. PR issue CCIP-4329 #bugfix

- [#15984](https://github.com/smartcontractkit/chainlink/pull/15984) [`86e9119`](https://github.com/smartcontractkit/chainlink/commit/86e9119d69b15c5e83ed38edc4debe1e6ce87674) - #internal fix missing case in gas estimation logic

  PR issue: CCIP-4919

  Solidity Review issue: CCIP-3966

- [#15605](https://github.com/smartcontractkit/chainlink/pull/15605) [`8c65527`](https://github.com/smartcontractkit/chainlink/commit/8c65527c82a20c74b2a4707221ef496802b21804) - replace f with fObserve in RMNHome and RMNRemote and update all tests CCIP-4058

- [#15405](https://github.com/smartcontractkit/chainlink/pull/15405) [`16f1529`](https://github.com/smartcontractkit/chainlink/commit/16f1529856a575c8d2091a16033a8d0371408d96) - enable via-ir in CCIP compilation

  PR issue: CCIP-4656

  Solidity Review issue: CCIP-3966

- [#14960](https://github.com/smartcontractkit/chainlink/pull/14960) [`1b1dc3b`](https://github.com/smartcontractkit/chainlink/commit/1b1dc3b8ee058901828963a0b9e59fc239444e93) - CCIP-3789 Add check on MultiAggregateRateLimiter:UpdateRateLimitTokens that remote token is not abi.encode(address(0)) #bugfix

- [#15020](https://github.com/smartcontractkit/chainlink/pull/15020) [`9ec788e`](https://github.com/smartcontractkit/chainlink/commit/9ec788e78b4fcf3266b5e0a9ed2e166cea51388f) - #internal more efficient ownership usage

  PR issue: CCIP-4083

  Solidity Review issue: CCIP-3966

- [#14734](https://github.com/smartcontractkit/chainlink/pull/14734) [`ca71878`](https://github.com/smartcontractkit/chainlink/commit/ca71878aa5a55fe239a456d7b564ffeba9bc84d7) - Make stalenessThreshold per dest chain and have 0 mean no staleness check.

  PR issue: CCIP-3414

- [#16299](https://github.com/smartcontractkit/chainlink/pull/16299) [`3b89c46`](https://github.com/smartcontractkit/chainlink/commit/3b89c464872609a87a172d93175f3314cb8558f6) - Additional security and parameter checks and comment fixes

  PR issue: CCIP-5183

  Solidity Review issue: CCIP-3966

- [#15743](https://github.com/smartcontractkit/chainlink/pull/15743) [`4a19318`](https://github.com/smartcontractkit/chainlink/commit/4a19318efa56c079da53777f82404a1fbb24479a) - Create a new version of the ERC165Checker library which checks for sufficient gas before making an external call to prevent message delivery issues. #bugfix

  PR issue: CCIP-4659

  Solidity Review issue: CCIP-3966

- [#15301](https://github.com/smartcontractkit/chainlink/pull/15301) [`6c4f1b9`](https://github.com/smartcontractkit/chainlink/commit/6c4f1b920c64b9c066b19119bb5990f0bb0714b0) - Refactor MockCCIPRouter to support EVMExtraArgsV2

  PR issue : CCIP-4288

- [#14863](https://github.com/smartcontractkit/chainlink/pull/14863) [`84bcbe0`](https://github.com/smartcontractkit/chainlink/commit/84bcbe03ebfa3ab7f58b97897eb0c55b45191859) - change else if to else..if..

  PR issue: CCIP-3726

- [#15570](https://github.com/smartcontractkit/chainlink/pull/15570) [`c1341a5`](https://github.com/smartcontractkit/chainlink/commit/c1341a5081d098bce04a7564a6525a91f2beeecf) - add getChainConfig to ccipHome

  PR issue: CCIP-4517

  Solidity Review issue: CCIP-3966

## 1.5.0

### Minor Changes

- [#14266](https://github.com/smartcontractkit/chainlink/pull/14266) [`c323e0d`](https://github.com/smartcontractkit/chainlink/commit/c323e0d600c659a4ea584dbae0a0db187afd51eb) Thanks [@asoliman92](https://github.com/asoliman92)! - #updated move latest ccip contracts code from ccip repo to chainlink repo
- [#13941](https://github.com/smartcontractkit/chainlink/pull/13941) [`9e74eee`](https://github.com/smartcontractkit/chainlink/commit/9e74eee9d415b386db33bdf2dd44facc82cd3551) Thanks [@RensR](https://github.com/RensR)! - add ccip contracts to the repo

### Patch Changes

- [#14345](https://github.com/smartcontractkit/chainlink/pull/14345) [`c83c687`](https://github.com/smartcontractkit/chainlink/commit/c83c68735bdee6bbd8510733b7415797cd08ecbd) Thanks [@makramkd](https://github.com/makramkd)! - #internal merge ccip contracts
- [#14516](https://github.com/smartcontractkit/chainlink/pull/14516) [`0e32c07`](https://github.com/smartcontractkit/chainlink/commit/0e32c07d22973343e722a228ff1c3b1e8f9bc04e) Thanks [@mateusz-sekara](https://github.com/mateusz-sekara)! - Adding USDCReaderTester contract for CCIP integration tests #internal
- [#14739](https://github.com/smartcontractkit/chainlink/pull/14739) [`4842271`](https://github.com/smartcontractkit/chainlink/commit/4842271b0f7054f5f1364c59d3d9da534c5d4f25) Thanks [@RensR](https://github.com/RensR)! - #internal remove CCIP 1.5
