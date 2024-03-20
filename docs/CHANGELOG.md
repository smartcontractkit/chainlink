# chainlink

## 2.10.0

### Minor Changes

- [#12378](https://github.com/smartcontractkit/chainlink/pull/12378) [`18c7237181`](https://github.com/smartcontractkit/chainlink/commit/18c7237181e8f9134e2f4992ba16b64f3548725d) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - update AutomationBase interface to check for ready only address on polygon zkEVM

- [#12348](https://github.com/smartcontractkit/chainlink/pull/12348) [`efead72965`](https://github.com/smartcontractkit/chainlink/commit/efead72965fec7e822a16f4d50cc0e5a27dd4640) Thanks [@reductionista](https://github.com/reductionista)! - Update config for zkevm polygon chains

- [#12431](https://github.com/smartcontractkit/chainlink/pull/12431) [`5546698edc`](https://github.com/smartcontractkit/chainlink/commit/5546698edc0a8b7ab2959aabe9772ba0e5b52a63) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - remove registerUpkeep from auto v21 common

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

- [#12510](https://github.com/smartcontractkit/chainlink/pull/12510) [`d01d8418ef`](https://github.com/smartcontractkit/chainlink/commit/d01d8418ef56b34349c70f9f38424bff1eeb8d8a) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - Helper VRF CLI command

- [#12325](https://github.com/smartcontractkit/chainlink/pull/12325) [`a2db9de8a5`](https://github.com/smartcontractkit/chainlink/commit/a2db9de8a5d3da4a541e3807fe3139d1aa1b0375) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - VRF V2+ Coordinator msg.data len validation

- [#12510](https://github.com/smartcontractkit/chainlink/pull/12510) [`d01d8418ef`](https://github.com/smartcontractkit/chainlink/commit/d01d8418ef56b34349c70f9f38424bff1eeb8d8a) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - Helper VRF CLI command

- [#12360](https://github.com/smartcontractkit/chainlink/pull/12360) [`8241e811b2`](https://github.com/smartcontractkit/chainlink/commit/8241e811b2ed37ccd3bc11674735e0599c43429c) Thanks [@reductionista](https://github.com/reductionista)! - Add support for eth_getLogs & finality tags in simulated_backend_client.go

- [#12355](https://github.com/smartcontractkit/chainlink/pull/12355) [`2e08d9be68`](https://github.com/smartcontractkit/chainlink/commit/2e08d9be685f6a9d6acce9a656ed92a028539157) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - Validate if flat fee configs are configured correctly

### Patch Changes

- [#12371](https://github.com/smartcontractkit/chainlink/pull/12371) [`710c60c5ee`](https://github.com/smartcontractkit/chainlink/commit/710c60c5eeaf0043a88555038fecfee0621eb397) Thanks [@anirudhwarrier](https://github.com/anirudhwarrier)! - Update automation smoke test to use UpkeepCounter with time based counter

- [#12497](https://github.com/smartcontractkit/chainlink/pull/12497) [`3ca3494450`](https://github.com/smartcontractkit/chainlink/commit/3ca34944507b01b7d4511d8ea8aff402c0a7bb85) Thanks [@RyanRHall](https://github.com/RyanRHall)! - added logic C contract to automation 2.3

- [#12389](https://github.com/smartcontractkit/chainlink/pull/12389) [`9f44174dd6`](https://github.com/smartcontractkit/chainlink/commit/9f44174dd60ecb29839fc1ce517c31bbbe474835) Thanks [@bolekk](https://github.com/bolekk)! - External peering core service

- [#12405](https://github.com/smartcontractkit/chainlink/pull/12405) [`2bd210bfa8`](https://github.com/smartcontractkit/chainlink/commit/2bd210bfa8c4705b0981a315cba939b0281d7bf3) Thanks [@jinhoonbang](https://github.com/jinhoonbang)! - Soft delete consumer nonce in VRF coordinator v2.5

- [#12499](https://github.com/smartcontractkit/chainlink/pull/12499) [`1a36386481`](https://github.com/smartcontractkit/chainlink/commit/1a363864816a3e7821d5a5844f13be360f0ecb58) Thanks [@RyanRHall](https://github.com/RyanRHall)! - refactor foundry tests for auto 2.3

- [#12388](https://github.com/smartcontractkit/chainlink/pull/12388) [`30b73a804d`](https://github.com/smartcontractkit/chainlink/commit/30b73a804dfba394180abe354569dade80a71be5) Thanks [@justinkaseman](https://github.com/justinkaseman)! - Chainlink Functions contracts v1.3 audit findings

- [#12296](https://github.com/smartcontractkit/chainlink/pull/12296) [`19b048561d`](https://github.com/smartcontractkit/chainlink/commit/19b048561dcb2e565adbfff1f745da51fea94df4) Thanks [@bolekk](https://github.com/bolekk)! - Added a RageP2P wrapper

- [#12392](https://github.com/smartcontractkit/chainlink/pull/12392) [`8626f1b83d`](https://github.com/smartcontractkit/chainlink/commit/8626f1b83df0fc5725d46874fd6e973567ce8edd) Thanks [@ilija42](https://github.com/ilija42)! - Add kv store tied to jobs and use it for juels fee per coin cache to store persisted values for backup

- [#12413](https://github.com/smartcontractkit/chainlink/pull/12413) [`e6843e8d9b`](https://github.com/smartcontractkit/chainlink/commit/e6843e8d9b99bac8c8fa724768a497f43ee1fb9d) Thanks [@shileiwill](https://github.com/shileiwill)! - make reserveAmounts to be a map

- [#12332](https://github.com/smartcontractkit/chainlink/pull/12332) [`89abd726b6`](https://github.com/smartcontractkit/chainlink/commit/89abd726b6c3f29a84e0fc5d230a1324f622755b) Thanks [@Tofel](https://github.com/Tofel)! - Add new pipeline for testing EVM node compatibility on go-ethereum dependency bump

- [#12425](https://github.com/smartcontractkit/chainlink/pull/12425) [`e3f4a6c4b3`](https://github.com/smartcontractkit/chainlink/commit/e3f4a6c4b331d7a7f5c3be2ddaf0c118993ff84e) Thanks [@jinhoonbang](https://github.com/jinhoonbang)! - add pending request counter for vrf v2.5 coordinator

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

- [#12443](https://github.com/smartcontractkit/chainlink/pull/12443) [`e604a73d7b`](https://github.com/smartcontractkit/chainlink/commit/e604a73d7b21c5f053631d9c8afeb0eaf7203310) Thanks [@shileiwill](https://github.com/shileiwill)! - use common interface for v2.3

- [#12353](https://github.com/smartcontractkit/chainlink/pull/12353) [`07c9f6cadd`](https://github.com/smartcontractkit/chainlink/commit/07c9f6cadd449989b21977af461305ded8e5b2f0) Thanks [@amit-momin](https://github.com/amit-momin)! - Fixed a race condition bug around EVM nonce management, which could cause the Node to skip a nonce and get stuck.

- [#12344](https://github.com/smartcontractkit/chainlink/pull/12344) [`6fa1f5dddc`](https://github.com/smartcontractkit/chainlink/commit/6fa1f5dddc6e257c2223503f1592297ca69521bd) Thanks [@eutopian](https://github.com/eutopian)! - Add rebalancer support for feeds manager ocr2 plugins

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

- [#12375](https://github.com/smartcontractkit/chainlink/pull/12375) [`831aea819d`](https://github.com/smartcontractkit/chainlink/commit/831aea819dd6b3415770cc927c4857a1da4557b5) Thanks [@shileiwill](https://github.com/shileiwill)! - add liquidity pool for automation 2.3

- [#12314](https://github.com/smartcontractkit/chainlink/pull/12314) [`15103b8ced`](https://github.com/smartcontractkit/chainlink/commit/15103b8ced1d931244d915c912a506b165fefb84) Thanks [@ibrajer](https://github.com/ibrajer)! - Validation for premium limits added to VRFCoordinatorV2_5 contract

- [#12412](https://github.com/smartcontractkit/chainlink/pull/12412) [`83c8688a14`](https://github.com/smartcontractkit/chainlink/commit/83c8688a14ac04111f999d132673ebaf6a364b4a) Thanks [@poopoothegorilla](https://github.com/poopoothegorilla)! - bump grafana to 1.1.1

- [#12248](https://github.com/smartcontractkit/chainlink/pull/12248) [`e1950769ee`](https://github.com/smartcontractkit/chainlink/commit/e1950769ee3ff2a40ca5772b9634c45f8be241cc) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - add version support for automation registry 2.\*
