# chainlink

## 3.0.0

### Major Changes

- [#12629](https://github.com/smartcontractkit/chainlink/pull/12629) [`3ec8cc914b`](https://github.com/smartcontractkit/chainlink/commit/3ec8cc914b6f8b3c889592ddb54a5801b5c0d5c6) Thanks [@ilija42](https://github.com/ilija42)! - Fix kv_store migration fk cascade deletion

### Minor Changes

- [#12586](https://github.com/smartcontractkit/chainlink/pull/12586) [`7987045897`](https://github.com/smartcontractkit/chainlink/commit/7987045897b4549760811de3aa661520672d361e) Thanks [@ilija42](https://github.com/ilija42)! - Fix error log formatting for in memory data source cache for juels fee per coin

- [#12378](https://github.com/smartcontractkit/chainlink/pull/12378) [`18c7237181`](https://github.com/smartcontractkit/chainlink/commit/18c7237181e8f9134e2f4992ba16b64f3548725d) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - update AutomationBase interface to check for ready only address on polygon zkEVM

- [#12348](https://github.com/smartcontractkit/chainlink/pull/12348) [`efead72965`](https://github.com/smartcontractkit/chainlink/commit/efead72965fec7e822a16f4d50cc0e5a27dd4640) Thanks [@reductionista](https://github.com/reductionista)! - Update config for zkevm polygon chains

- [#12647](https://github.com/smartcontractkit/chainlink/pull/12647) [`bc4fbbdb61`](https://github.com/smartcontractkit/chainlink/commit/bc4fbbdb616b050a4b7861f5c10c5d3ee0ddad75) Thanks [@ilija42](https://github.com/ilija42)! - fix jfpc cache cleanup

- [#12431](https://github.com/smartcontractkit/chainlink/pull/12431) [`5546698edc`](https://github.com/smartcontractkit/chainlink/commit/5546698edc0a8b7ab2959aabe9772ba0e5b52a63) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - remove registerUpkeep from auto v21 common

- [#12469](https://github.com/smartcontractkit/chainlink/pull/12469) [`1370133b72`](https://github.com/smartcontractkit/chainlink/commit/1370133b722ab97650b15c6aeab72bb494790b63) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - implement offchain settlement for NOPs payment

- [#12082](https://github.com/smartcontractkit/chainlink/pull/12082) [`608ea0a467`](https://github.com/smartcontractkit/chainlink/commit/608ea0a467ee36e15fdc654a88494ae579d778a6) Thanks [@dhaidashenko](https://github.com/dhaidashenko)! - HeadTracker now respects the `FinalityTagEnabled` config option. If the flag is enabled, HeadTracker backfills blocks up to the latest finalized block provided by the corresponding RPC call. To address potential misconfigurations, `HistoryDepth` is now calculated from the latest finalized block instead of the head. NOTE: Consumers (e.g. TXM and LogPoller) do not fully utilize Finality Tag yet.

- [#12489](https://github.com/smartcontractkit/chainlink/pull/12489) [`3a49094db2`](https://github.com/smartcontractkit/chainlink/commit/3a49094db25036e1948818e4030fca11be748914) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - - Misc VRF V2+ contract changes

  - Reuse struct RequestCommitmentV2Plus from VRFTypes
  - Fix interface name IVRFCoordinatorV2PlusFulfill in BatchVRFCoordinatorV2Plus to avoid confusion with IVRFCoordinatorV2Plus.sol
  - Remove unused errors
  - Rename variables for readability
  - Fix comments
  - Minor gas optimisation (++i)
  - Fix integration tests

- [#12093](https://github.com/smartcontractkit/chainlink/pull/12093) [`3f6d901fe6`](https://github.com/smartcontractkit/chainlink/commit/3f6d901fe676698769cb6713250152e322747145) Thanks [@friedemannf](https://github.com/friedemannf)! - The `xdai` `ChainType` has been renamed to `gnosis` to match the chain's new name. The old value is still supported but has been deprecated and will be removed in v2.13.0.

- [#12503](https://github.com/smartcontractkit/chainlink/pull/12503) [`dc224a2924`](https://github.com/smartcontractkit/chainlink/commit/dc224a29249c83c74a38d9ca9d16fb00e192a4e2) Thanks [@amit-momin](https://github.com/amit-momin)! - Added a tx simulation feature to the chain client to enable testing for zk out-of-counter (OOC) errors

- [#12510](https://github.com/smartcontractkit/chainlink/pull/12510) [`d01d8418ef`](https://github.com/smartcontractkit/chainlink/commit/d01d8418ef56b34349c70f9f38424bff1eeb8d8a) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - Helper VRF CLI command

- [#12634](https://github.com/smartcontractkit/chainlink/pull/12634) [`e9e903bf4b`](https://github.com/smartcontractkit/chainlink/commit/e9e903bf4b34099f8b274eb1e0f013b4ab326bb4) Thanks [@ettec](https://github.com/ettec)! - Update keyvalue store to be compatible with the interface required in chainlink common

- [#12496](https://github.com/smartcontractkit/chainlink/pull/12496) [`31350477ae`](https://github.com/smartcontractkit/chainlink/commit/31350477ae51f00e035b1b8c50775e5955258ac1) Thanks [@silaslenihan](https://github.com/silaslenihan)! - Change LimitTransfer gasLimit type from uint32 to uint64

- [#12522](https://github.com/smartcontractkit/chainlink/pull/12522) [`886201638e`](https://github.com/smartcontractkit/chainlink/commit/886201638e14dc478ae7104b4a5aed9ac8af5bba) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - improve foundry tests and fix nits

- [#12339](https://github.com/smartcontractkit/chainlink/pull/12339) [`96d2fe13b8`](https://github.com/smartcontractkit/chainlink/commit/96d2fe13b8510631bbc92ffd20a4d923b93002e6) Thanks [@dhaidashenko](https://github.com/dhaidashenko)! - Add the `pool_rpc_node_highest_finalized_block` metric that tracks the highest finalized block seen per RPC. If `FinalityTagEnabled = true`, a positive `NodePool.FinalizedBlockPollInterval` is needed to collect the metric. If the finality tag is not enabled, the metric is populated with a calculated latest finalized block based on the latest head and finality depth.

- [#12325](https://github.com/smartcontractkit/chainlink/pull/12325) [`a2db9de8a5`](https://github.com/smartcontractkit/chainlink/commit/a2db9de8a5d3da4a541e3807fe3139d1aa1b0375) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - VRF V2+ Coordinator msg.data len validation

- [#12586](https://github.com/smartcontractkit/chainlink/pull/12586) [`7987045897`](https://github.com/smartcontractkit/chainlink/commit/7987045897b4549760811de3aa661520672d361e) Thanks [@ilija42](https://github.com/ilija42)! - Add error log if juels fee per coin cache is over 24h old and lower other logs severity in cache to warn

- [#12510](https://github.com/smartcontractkit/chainlink/pull/12510) [`d01d8418ef`](https://github.com/smartcontractkit/chainlink/commit/d01d8418ef56b34349c70f9f38424bff1eeb8d8a) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - Helper VRF CLI command

- [#12132](https://github.com/smartcontractkit/chainlink/pull/12132) [`478f73b9cf`](https://github.com/smartcontractkit/chainlink/commit/478f73b9cfe8013316546df7a057a784e2e7bf01) Thanks [@vreff](https://github.com/vreff)! - Remove noisy log poller warning in VRFv2 & VRFv2+ listener loops

- [#12583](https://github.com/smartcontractkit/chainlink/pull/12583) [`50724c3bb1`](https://github.com/smartcontractkit/chainlink/commit/50724c3bb1fb959f85d361bc0615f58cc16e4fc9) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - bug fixes in s_reserveAmount accounting

- [#12360](https://github.com/smartcontractkit/chainlink/pull/12360) [`8241e811b2`](https://github.com/smartcontractkit/chainlink/commit/8241e811b2ed37ccd3bc11674735e0599c43429c) Thanks [@reductionista](https://github.com/reductionista)! - Add support for eth_getLogs & finality tags in simulated_backend_client.go

- [#12473](https://github.com/smartcontractkit/chainlink/pull/12473) [`f1d1f249eb`](https://github.com/smartcontractkit/chainlink/commit/f1d1f249ebecb37da7eacbc4cc12e1eb0205f29a) Thanks [@justinkaseman](https://github.com/justinkaseman)! - Copy common transmitter methods into FunctionsContractTransmitter to enable product specific modification

- [#12581](https://github.com/smartcontractkit/chainlink/pull/12581) [`6fcc73983e`](https://github.com/smartcontractkit/chainlink/commit/6fcc73983e5b782bb4ac577cb33093bf80e3a582) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - VRFV2PlusWrapper config refactor

- [#12534](https://github.com/smartcontractkit/chainlink/pull/12534) [`bd532b5e2a`](https://github.com/smartcontractkit/chainlink/commit/bd532b5e2a1bebe8c9fe689d059464c43365ced1) Thanks [@silaslenihan](https://github.com/silaslenihan)! - Extracted Gas Limit Multiplier from gas estimators to WrappedEvmEstimator.

- [#12355](https://github.com/smartcontractkit/chainlink/pull/12355) [`2e08d9be68`](https://github.com/smartcontractkit/chainlink/commit/2e08d9be685f6a9d6acce9a656ed92a028539157) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - Validate if flat fee configs are configured correctly

- [#12578](https://github.com/smartcontractkit/chainlink/pull/12578) [`ffd492295f`](https://github.com/smartcontractkit/chainlink/commit/ffd492295f03de8c3b946a003dacbded731d7899) Thanks [@RensR](https://github.com/RensR)! - Remove 0.6 and 0.7 Solidity source code

- [#12547](https://github.com/smartcontractkit/chainlink/pull/12547) [`8162f7b101`](https://github.com/smartcontractkit/chainlink/commit/8162f7b1012dd669e51bbb4038a6d5df29906267) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - pay deactivated transmitters in offchain settlement

### Patch Changes

- [#12570](https://github.com/smartcontractkit/chainlink/pull/12570) [`2d33524a35`](https://github.com/smartcontractkit/chainlink/commit/2d33524a3539e32ac32a84c4600e6cdfb8e01cf3) Thanks [@samsondav](https://github.com/samsondav)! - VerboseLogging is now turned on by default.

  You may disable if this results in excessive log volume. Disable like so:

  ```
  [Pipeline]
  VerboseLogging = false
  ```

- [#12458](https://github.com/smartcontractkit/chainlink/pull/12458) [`51b134700a`](https://github.com/smartcontractkit/chainlink/commit/51b134700afe6daa1a10692e6365fdbbaf3b1396) Thanks [@HenryNguyen5](https://github.com/HenryNguyen5)! - Add json schema support to workflows

- [#12502](https://github.com/smartcontractkit/chainlink/pull/12502) [`ca14ccd3c6`](https://github.com/smartcontractkit/chainlink/commit/ca14ccd3c64bea128e12a0d37d399f400ff62584) Thanks [@bolekk](https://github.com/bolekk)! - Dispatcher service for external peering

- [#12598](https://github.com/smartcontractkit/chainlink/pull/12598) [`e753637e01`](https://github.com/smartcontractkit/chainlink/commit/e753637e01fabb8ea3760eb14204124c8d3b88e1) Thanks [@RyanRHall](https://github.com/RyanRHall)! - small gas fix

- [#12371](https://github.com/smartcontractkit/chainlink/pull/12371) [`710c60c5ee`](https://github.com/smartcontractkit/chainlink/commit/710c60c5eeaf0043a88555038fecfee0621eb397) Thanks [@anirudhwarrier](https://github.com/anirudhwarrier)! - Update automation smoke test to use UpkeepCounter with time based counter

- [#12540](https://github.com/smartcontractkit/chainlink/pull/12540) [`17c037678d`](https://github.com/smartcontractkit/chainlink/commit/17c037678d05c88f28a28a3ac760c742f549d5ec) Thanks [@RyanRHall](https://github.com/RyanRHall)! - change auto 2.3 flat fees from link to USD

- [#12555](https://github.com/smartcontractkit/chainlink/pull/12555) [`cda84cb1b7`](https://github.com/smartcontractkit/chainlink/commit/cda84cb1b7582379ac140b3a88da6179275dbefb) Thanks [@shileiwill](https://github.com/shileiwill)! - safeTransfer and cleanups

- [#12497](https://github.com/smartcontractkit/chainlink/pull/12497) [`3ca3494450`](https://github.com/smartcontractkit/chainlink/commit/3ca34944507b01b7d4511d8ea8aff402c0a7bb85) Thanks [@RyanRHall](https://github.com/RyanRHall)! - added logic C contract to automation 2.3

- [#12389](https://github.com/smartcontractkit/chainlink/pull/12389) [`9f44174dd6`](https://github.com/smartcontractkit/chainlink/commit/9f44174dd60ecb29839fc1ce517c31bbbe474835) Thanks [@bolekk](https://github.com/bolekk)! - External peering core service

- [#12405](https://github.com/smartcontractkit/chainlink/pull/12405) [`2bd210bfa8`](https://github.com/smartcontractkit/chainlink/commit/2bd210bfa8c4705b0981a315cba939b0281d7bf3) Thanks [@jinhoonbang](https://github.com/jinhoonbang)! - Soft delete consumer nonce in VRF coordinator v2.5

- [#12499](https://github.com/smartcontractkit/chainlink/pull/12499) [`1a36386481`](https://github.com/smartcontractkit/chainlink/commit/1a363864816a3e7821d5a5844f13be360f0ecb58) Thanks [@RyanRHall](https://github.com/RyanRHall)! - refactor foundry tests for auto 2.3

- [#12529](https://github.com/smartcontractkit/chainlink/pull/12529) [`84913bfcfc`](https://github.com/smartcontractkit/chainlink/commit/84913bfcfcfcf6f93fb359814208a32e3e659d23) Thanks [@ibrajer](https://github.com/ibrajer)! - VRFV2PlusWrapper contract: subID param added to the constructor, removed migrate() method

- [#12387](https://github.com/smartcontractkit/chainlink/pull/12387) [`42e72d2d26`](https://github.com/smartcontractkit/chainlink/commit/42e72d2d2610d2481c5a9469fc9b49c167d37f79) Thanks [@ogtownsend](https://github.com/ogtownsend)! - Adds prometheus metrics for automation streams error handling

- [#12388](https://github.com/smartcontractkit/chainlink/pull/12388) [`30b73a804d`](https://github.com/smartcontractkit/chainlink/commit/30b73a804dfba394180abe354569dade80a71be5) Thanks [@justinkaseman](https://github.com/justinkaseman)! - Chainlink Functions contracts v1.3 audit findings

- [#12296](https://github.com/smartcontractkit/chainlink/pull/12296) [`19b048561d`](https://github.com/smartcontractkit/chainlink/commit/19b048561dcb2e565adbfff1f745da51fea94df4) Thanks [@bolekk](https://github.com/bolekk)! - Added a RageP2P wrapper

- [#12392](https://github.com/smartcontractkit/chainlink/pull/12392) [`8626f1b83d`](https://github.com/smartcontractkit/chainlink/commit/8626f1b83df0fc5725d46874fd6e973567ce8edd) Thanks [@ilija42](https://github.com/ilija42)! - Add kv store tied to jobs and use it for juels fee per coin cache to store persisted values for backup

- [#12413](https://github.com/smartcontractkit/chainlink/pull/12413) [`e6843e8d9b`](https://github.com/smartcontractkit/chainlink/commit/e6843e8d9b99bac8c8fa724768a497f43ee1fb9d) Thanks [@shileiwill](https://github.com/shileiwill)! - make reserveAmounts to be a map

- [#12536](https://github.com/smartcontractkit/chainlink/pull/12536) [`87b0d8f309`](https://github.com/smartcontractkit/chainlink/commit/87b0d8f3091e3276cd049d3a852ab63e4d6bda5b) Thanks [@shileiwill](https://github.com/shileiwill)! - billing overrides

- [#12332](https://github.com/smartcontractkit/chainlink/pull/12332) [`89abd726b6`](https://github.com/smartcontractkit/chainlink/commit/89abd726b6c3f29a84e0fc5d230a1324f622755b) Thanks [@Tofel](https://github.com/Tofel)! - Add new pipeline for testing EVM node compatibility on go-ethereum dependency bump

- [#12425](https://github.com/smartcontractkit/chainlink/pull/12425) [`e3f4a6c4b3`](https://github.com/smartcontractkit/chainlink/commit/e3f4a6c4b331d7a7f5c3be2ddaf0c118993ff84e) Thanks [@jinhoonbang](https://github.com/jinhoonbang)! - add pending request counter for vrf v2.5 coordinator

- [#12621](https://github.com/smartcontractkit/chainlink/pull/12621) [`9c2764adbf`](https://github.com/smartcontractkit/chainlink/commit/9c2764adbf3969654795ed2c35c5fb56eaf70785) Thanks [@KuphJr](https://github.com/KuphJr)! - Add GetFilters function to the log_poller

- [#12592](https://github.com/smartcontractkit/chainlink/pull/12592) [`b512ef5a7d`](https://github.com/smartcontractkit/chainlink/commit/b512ef5a7d1bc87d0cbd5357c5c47cc0dcb75e0b) Thanks [@ibrajer](https://github.com/ibrajer)! - Set LINK native feed in VRFV2PlusWrapper to immutable

- [#12553](https://github.com/smartcontractkit/chainlink/pull/12553) [`4892376917`](https://github.com/smartcontractkit/chainlink/commit/4892376917a16253165dc761f8efec41da32ec9c) Thanks [@RyanRHall](https://github.com/RyanRHall)! - address TODOs and docs for 2.3

- [#12498](https://github.com/smartcontractkit/chainlink/pull/12498) [`1c576d0e34`](https://github.com/smartcontractkit/chainlink/commit/1c576d0e34d93a6298ddcb662ee89fd04eeda53e) Thanks [@samsondav](https://github.com/samsondav)! - Add new config option Pipeline.VerboseLogging

  VerboseLogging enables detailed logging of pipeline execution steps. This is
  disabled by default because it increases log volume for pipeline runs, but can
  be useful for debugging failed runs without relying on the UI or database.
  Consider enabling this if you disabled run saving by setting MaxSuccessfulRuns
  to zero.

  Set it like the following example:

  ```
  [Pipeline]
  VerboseLogging = true
  ```

- [#12582](https://github.com/smartcontractkit/chainlink/pull/12582) [`684afa4e1f`](https://github.com/smartcontractkit/chainlink/commit/684afa4e1fcb2cad292cbc3b97ebeda3e3ef7bc8) Thanks [@RyanRHall](https://github.com/RyanRHall)! - fix bug in auto2.3 withdrawERC20Fees

- [#12661](https://github.com/smartcontractkit/chainlink/pull/12661) [`3b02047754`](https://github.com/smartcontractkit/chainlink/commit/3b020477548c17ed786036494ccc733107ca4152) Thanks [@RyanRHall](https://github.com/RyanRHall)! - more auto 2.3 tests

- [#12152](https://github.com/smartcontractkit/chainlink/pull/12152) [`a6a2acfe20`](https://github.com/smartcontractkit/chainlink/commit/a6a2acfe2017dc766d401d55627f0c5016c824b9) Thanks [@ferglor](https://github.com/ferglor)! - Calculate blockRate and logLimit defaults in the log provider based on chain ID

- [#12584](https://github.com/smartcontractkit/chainlink/pull/12584) [`c7cacd0710`](https://github.com/smartcontractkit/chainlink/commit/c7cacd0710f5040a46532e6dae7eac1b9eafe645) Thanks [@matYang](https://github.com/matYang)! - L1Oracle handles OP Stack Ecotone encoded l1 gas price

- [#12189](https://github.com/smartcontractkit/chainlink/pull/12189) [`79db1206af`](https://github.com/smartcontractkit/chainlink/commit/79db1206aff6b49238a36561d1f8dc86210736e6) Thanks [@dimriou](https://github.com/dimriou)! - Refactor Log and TxStore ORMs

- [#12443](https://github.com/smartcontractkit/chainlink/pull/12443) [`e604a73d7b`](https://github.com/smartcontractkit/chainlink/commit/e604a73d7b21c5f053631d9c8afeb0eaf7203310) Thanks [@shileiwill](https://github.com/shileiwill)! - use common interface for v2.3

- [#12564](https://github.com/smartcontractkit/chainlink/pull/12564) [`246762ceeb`](https://github.com/smartcontractkit/chainlink/commit/246762ceebba7923641ec00e66ae1aaf59bbcdc2) Thanks [@mateusz-sekara](https://github.com/mateusz-sekara)! - Exposing information about LogPoller finality violation via Healthy method. It's raised whenever LogPoller sees reorg deeper than the finality

- [#12575](https://github.com/smartcontractkit/chainlink/pull/12575) [`23254c4bf5`](https://github.com/smartcontractkit/chainlink/commit/23254c4bf577e84b71bda1d9a8b2c11e7b548267) Thanks [@augustbleeds](https://github.com/augustbleeds)! - update starknet relayer to fix nonce issue. introduces optional api-key for starknet toml config.

- [#12353](https://github.com/smartcontractkit/chainlink/pull/12353) [`07c9f6cadd`](https://github.com/smartcontractkit/chainlink/commit/07c9f6cadd449989b21977af461305ded8e5b2f0) Thanks [@amit-momin](https://github.com/amit-momin)! - Fixed a race condition bug around EVM nonce management, which could cause the Node to skip a nonce and get stuck.

- [#12344](https://github.com/smartcontractkit/chainlink/pull/12344) [`6fa1f5dddc`](https://github.com/smartcontractkit/chainlink/commit/6fa1f5dddc6e257c2223503f1592297ca69521bd) Thanks [@eutopian](https://github.com/eutopian)! - Add rebalancer support for feeds manager ocr2 plugins

- [#12484](https://github.com/smartcontractkit/chainlink/pull/12484) [`590cad6126`](https://github.com/smartcontractkit/chainlink/commit/590cad61269c75a6b22be1f6a73c74adfd1baa40) Thanks [@mateusz-sekara](https://github.com/mateusz-sekara)! - Making LogPoller's replay more robust by backfilling up to finalized block and processing rest in the main loop

- [#12477](https://github.com/smartcontractkit/chainlink/pull/12477) [`b2576475fc`](https://github.com/smartcontractkit/chainlink/commit/b2576475fc5c8ac037fc569fddc56e9d515ae7ca) Thanks [@shileiwill](https://github.com/shileiwill)! - native support

- [#12612](https://github.com/smartcontractkit/chainlink/pull/12612) [`d44abe3769`](https://github.com/smartcontractkit/chainlink/commit/d44abe37693d6995377fa1329e433e7fba26885d) Thanks [@RensR](https://github.com/RensR)! - upgraded transmission to 0.8.19

- [#12518](https://github.com/smartcontractkit/chainlink/pull/12518) [`e74aeab286`](https://github.com/smartcontractkit/chainlink/commit/e74aeab286f642bdc5b168d8e6f716d92bfcc8ea) Thanks [@erikburt](https://github.com/erikburt)! - docs: remove repeated words in documentation and comments

- [#12444](https://github.com/smartcontractkit/chainlink/pull/12444) [`dde7fdff33`](https://github.com/smartcontractkit/chainlink/commit/dde7fdff33cfc0690844cf0a88295bef57e2a269) Thanks [@ogtownsend](https://github.com/ogtownsend)! - Updating prometheus metrics for Automation log triggers

- [#12479](https://github.com/smartcontractkit/chainlink/pull/12479) [`93762ccbd8`](https://github.com/smartcontractkit/chainlink/commit/93762ccbd868b9e227abf3220afb9ad22ba41b92) Thanks [@jinhoonbang](https://github.com/jinhoonbang)! - update solc version for vrf v2.5 coordinators

- [#12606](https://github.com/smartcontractkit/chainlink/pull/12606) [`90ea24bc33`](https://github.com/smartcontractkit/chainlink/commit/90ea24bc33f2663c414a57e1eb71a2bc98b5ddfa) Thanks [@shileiwill](https://github.com/shileiwill)! - remove trailing slash

- [#12337](https://github.com/smartcontractkit/chainlink/pull/12337) [`195b504a93`](https://github.com/smartcontractkit/chainlink/commit/195b504a93b1a241c1981ec21726e4b722d40b2b) Thanks [@samsondav](https://github.com/samsondav)! - Mercury jobs can now broadcast to multiple mercury servers.

  Previously, a single mercury server would be specified in a job spec as so:

  ```toml
  [pluginConfig]
  serverURL = "example.com/foo"
  serverPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
  ```

  You may now specify multiple mercury servers, as so:

  ```toml
  [pluginConfig]
  servers = { "example.com/foo" = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", "mercury2.example:1234/bar" = "524ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93" }
  ```

- [#11899](https://github.com/smartcontractkit/chainlink/pull/11899) [`67560b9f1d`](https://github.com/smartcontractkit/chainlink/commit/67560b9f1dc052712a76eeb245fba12f2daf8e8d) Thanks [@DylanTinianov](https://github.com/DylanTinianov)! - Refactor EVM ORMs to remove pg dependency

- [#12531](https://github.com/smartcontractkit/chainlink/pull/12531) [`88e010d604`](https://github.com/smartcontractkit/chainlink/commit/88e010d604682c54c4f99e0a0916f94c0d13ece6) Thanks [@jinhoonbang](https://github.com/jinhoonbang)! - increase num optimizations to 500 for vrf v2.5 coordinator

- [#12375](https://github.com/smartcontractkit/chainlink/pull/12375) [`831aea819d`](https://github.com/smartcontractkit/chainlink/commit/831aea819dd6b3415770cc927c4857a1da4557b5) Thanks [@shileiwill](https://github.com/shileiwill)! - add liquidity pool for automation 2.3

- [#12314](https://github.com/smartcontractkit/chainlink/pull/12314) [`15103b8ced`](https://github.com/smartcontractkit/chainlink/commit/15103b8ced1d931244d915c912a506b165fefb84) Thanks [@ibrajer](https://github.com/ibrajer)! - Validation for premium limits added to VRFCoordinatorV2_5 contract

- [#12521](https://github.com/smartcontractkit/chainlink/pull/12521) [`a27d466135`](https://github.com/smartcontractkit/chainlink/commit/a27d4661359d95712bf25d127f35732ca491d4dc) Thanks [@DylanTinianov](https://github.com/DylanTinianov)! - Remove pg from evm tests

- [#12591](https://github.com/smartcontractkit/chainlink/pull/12591) [`b3086d0ec2`](https://github.com/smartcontractkit/chainlink/commit/b3086d0ec2565badaafdbb9c26e30caeb6fb41c9) Thanks [@RyanRHall](https://github.com/RyanRHall)! - fix withdraw LINK bug in auto 2.3

- [#12412](https://github.com/smartcontractkit/chainlink/pull/12412) [`83c8688a14`](https://github.com/smartcontractkit/chainlink/commit/83c8688a14ac04111f999d132673ebaf6a364b4a) Thanks [@poopoothegorilla](https://github.com/poopoothegorilla)! - bump grafana to 1.1.1

- [#12338](https://github.com/smartcontractkit/chainlink/pull/12338) [`1853ccaaf9`](https://github.com/smartcontractkit/chainlink/commit/1853ccaaf9887b22556fb43b81bd9a65b8ed414f) Thanks [@friedemannf](https://github.com/friedemannf)! - Handle zkSync specific known transaction error

- [#12577](https://github.com/smartcontractkit/chainlink/pull/12577) [`cf00183f62`](https://github.com/smartcontractkit/chainlink/commit/cf00183f6295fe95979b460f89bcc65f22237fd4) Thanks [@shileiwill](https://github.com/shileiwill)! - add test for billing override

- [#12619](https://github.com/smartcontractkit/chainlink/pull/12619) [`6531e34865`](https://github.com/smartcontractkit/chainlink/commit/6531e348659e2b7048b25183eadddbcb10426741) Thanks [@RyanRHall](https://github.com/RyanRHall)! - enable gas tests for auto 2.3

- [#12248](https://github.com/smartcontractkit/chainlink/pull/12248) [`e1950769ee`](https://github.com/smartcontractkit/chainlink/commit/e1950769ee3ff2a40ca5772b9634c45f8be241cc) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - add version support for automation registry 2.\*
