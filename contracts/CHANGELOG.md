# @chainlink/contracts

## 1.2.0 - 2024-07-18

### Minor Changes

- [#13428](https://github.com/smartcontractkit/chainlink/pull/13428) [`a02bb0a`](https://github.com/smartcontractkit/chainlink/commit/a02bb0a48cb78c815ec3a2f3bd351f2efc3869d7) Thanks [@elatoskinas](https://github.com/elatoskinas)! - #added EnumerableMapAddresses shared lib for AddressToAddress and AddressToBytes32 maps

- [#13774](https://github.com/smartcontractkit/chainlink/pull/13774) [`b5c0ea9`](https://github.com/smartcontractkit/chainlink/commit/b5c0ea9f106ced6daf8e1f35e88bc39d7b82de28) Thanks [@0xsuryansh](https://github.com/0xsuryansh)! - We have multiple validation use-cases on-chain which requires the inputs to be a set, sorted-set or we need to do subset checks.Adding a library for these validations

- [#13514](https://github.com/smartcontractkit/chainlink/pull/13514) [`f84a3f2`](https://github.com/smartcontractkit/chainlink/commit/f84a3f2f276847d26c94bf67215e2a3600951c9c) Thanks [@ilija42](https://github.com/ilija42)! - #internal Modify ChainReader tester contract.

### Patch Changes

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - capability registry internal review

- [#13546](https://github.com/smartcontractkit/chainlink/pull/13546) [`10ddafa`](https://github.com/smartcontractkit/chainlink/commit/10ddafaebedb94ad5a59968d19256b8c4592857f) Thanks [@cds95](https://github.com/cds95)! - upgrade keystone contracts to 0.8.24

- [#13376](https://github.com/smartcontractkit/chainlink/pull/13376) [`bb40d51`](https://github.com/smartcontractkit/chainlink/commit/bb40d51502487b010a1d6621db42458356dbbdc0) Thanks [@shileiwill](https://github.com/shileiwill)! - add events, add getter and add comments #bugfix

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - add capability registry comment explaining why we do not validate node operator name

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - capability registry informational findings

- [#13644](https://github.com/smartcontractkit/chainlink/pull/13644) [`2ed4478`](https://github.com/smartcontractkit/chainlink/commit/2ed4478ecc7fa6bada1138d204c12236fe56a810) Thanks [@shileiwill](https://github.com/shileiwill)! - add test for v23 #added

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - allow updating capabilities and to add/remove multiple capabilities at once from the capability registry

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - remove tracking deprecated arrays

- [#13183](https://github.com/smartcontractkit/chainlink/pull/13183) [`9630475`](https://github.com/smartcontractkit/chainlink/commit/96304756a77cdb2acf251d21d59b6aa8b55bf61a) Thanks [@cds95](https://github.com/cds95)! - track config count in DON struct

- [#13326](https://github.com/smartcontractkit/chainlink/pull/13326) [`732eea0`](https://github.com/smartcontractkit/chainlink/commit/732eea0f834e1f84c3a7520e131698d38762c78b) Thanks [@jhweintraub](https://github.com/jhweintraub)! - Add Additional tests to the callWithExactGasLibrary to ensure proper gas usage

- [#13183](https://github.com/smartcontractkit/chainlink/pull/13183) [`9630475`](https://github.com/smartcontractkit/chainlink/commit/96304756a77cdb2acf251d21d59b6aa8b55bf61a) Thanks [@cds95](https://github.com/cds95)! - implement remove DONs in capability registry

- [#13425](https://github.com/smartcontractkit/chainlink/pull/13425) [`eeb363f`](https://github.com/smartcontractkit/chainlink/commit/eeb363f1230415dde573607a095b177c612d3bef) Thanks [@DeividasK](https://github.com/DeividasK)! - #internal

- [#13453](https://github.com/smartcontractkit/chainlink/pull/13453) [`8c98c80`](https://github.com/smartcontractkit/chainlink/commit/8c98c80376c3b6d72bffeab62ee45a74449b6ef5) Thanks [@cds95](https://github.com/cds95)! - update error message when node does not exist

- [#13366](https://github.com/smartcontractkit/chainlink/pull/13366) [`d53d6d0`](https://github.com/smartcontractkit/chainlink/commit/d53d6d08dac5d3ee27ae89012669c6c2455295c8) Thanks [@bolekk](https://github.com/bolekk)! - #internal stub of keystone feed consumer contract

- [#13287](https://github.com/smartcontractkit/chainlink/pull/13287) [`8491b24`](https://github.com/smartcontractkit/chainlink/commit/8491b247cd5ff240d281ecfbabc0fe0fa19134f1) Thanks [@shileiwill](https://github.com/shileiwill)! - roundup #bugfix

- [#13426](https://github.com/smartcontractkit/chainlink/pull/13426) [`592b2bb`](https://github.com/smartcontractkit/chainlink/commit/592b2bb5a84a0e8858f77c5faa99e881f911878c) Thanks [@archseer](https://github.com/archseer)! - #internal

- [#13600](https://github.com/smartcontractkit/chainlink/pull/13600) [`4406364`](https://github.com/smartcontractkit/chainlink/commit/4406364702dd5af0b34c0397c00f8489510ecb11) Thanks [@jinhoonbang](https://github.com/jinhoonbang)! - relax pragma solidity version in VRF SubscriptionAPI.sol

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - #internal internal-review-for-capability-registry

- [#13528](https://github.com/smartcontractkit/chainlink/pull/13528) [`95502ad`](https://github.com/smartcontractkit/chainlink/commit/95502ad2699d63891662594f70e82e76682f2ed8) Thanks [@DeividasK](https://github.com/DeividasK)! - #internal

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - validate that a node is not part of a DON when removing

- [#13504](https://github.com/smartcontractkit/chainlink/pull/13504) [`815c5ea`](https://github.com/smartcontractkit/chainlink/commit/815c5ea8715462e00f6ea10cdc0b93ec3e1ba505) Thanks [@shileiwill](https://github.com/shileiwill)! - move v23 contracts #bugfix

- [#13335](https://github.com/smartcontractkit/chainlink/pull/13335) [`697e469`](https://github.com/smartcontractkit/chainlink/commit/697e469e41e640c8c71214461426174340527b4b) Thanks [@ibrajer](https://github.com/ibrajer)! - VRFV2Plus coordinator and wrapper split contracts between L1 and L2 chains #updated

- [#13487](https://github.com/smartcontractkit/chainlink/pull/13487) [`5e27da9`](https://github.com/smartcontractkit/chainlink/commit/5e27da95f09f21272e93f086bc2de5a9bc2ae399) Thanks [@bolekk](https://github.com/bolekk)! - #internal Use audited version of OCR2Base.sol in OCR3Capability.sol

- [#13183](https://github.com/smartcontractkit/chainlink/pull/13183) [`9630475`](https://github.com/smartcontractkit/chainlink/commit/96304756a77cdb2acf251d21d59b6aa8b55bf61a) Thanks [@cds95](https://github.com/cds95)! - Add update DON function to capability registry

- [#13566](https://github.com/smartcontractkit/chainlink/pull/13566) [`f1d478d`](https://github.com/smartcontractkit/chainlink/commit/f1d478d8bbba92f9753f7a1d6733ba4efd1f616a) Thanks [@DeividasK](https://github.com/DeividasK)! - #internal

- [#13389](https://github.com/smartcontractkit/chainlink/pull/13389) [`3959091`](https://github.com/smartcontractkit/chainlink/commit/3959091d4f3925b64cb6b0b55b7f7c72a4f924b9) Thanks [@bolekk](https://github.com/bolekk)! - #internal Keystone Forwarder and Feeds Consumer

- [#13183](https://github.com/smartcontractkit/chainlink/pull/13183) [`9630475`](https://github.com/smartcontractkit/chainlink/commit/96304756a77cdb2acf251d21d59b6aa8b55bf61a) Thanks [@cds95](https://github.com/cds95)! - add getters in capability registry

- [#13815](https://github.com/smartcontractkit/chainlink/pull/13815) [`fb177f4`](https://github.com/smartcontractkit/chainlink/commit/fb177f4ee77898dd12e20499e421a4d591fb92ef) Thanks [@KuphJr](https://github.com/KuphJr)! - Implemented improved L1 fee calculation for L2 chains in Functions contracts

- [#13237](https://github.com/smartcontractkit/chainlink/pull/13237) [`53312f0`](https://github.com/smartcontractkit/chainlink/commit/53312f01f3459408402ed97a9e3935ec63ccee7a) Thanks [@erikburt](https://github.com/erikburt)! - increase solhint max-warnings to 2 (from 0) to allow workflow to pass

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - remove update capabilities from capability registry

- [#13183](https://github.com/smartcontractkit/chainlink/pull/13183) [`9630475`](https://github.com/smartcontractkit/chainlink/commit/96304756a77cdb2acf251d21d59b6aa8b55bf61a) Thanks [@cds95](https://github.com/cds95)! - update ICapabilityConfiguration interface

- [#13216](https://github.com/smartcontractkit/chainlink/pull/13216) [`6099abb`](https://github.com/smartcontractkit/chainlink/commit/6099abbdbfb3ad396ca1ed5138ecd7a13159de19) Thanks [@ibrajer](https://github.com/ibrajer)! - Added Base Sepolia to ChainUtils #changed

- [#13352](https://github.com/smartcontractkit/chainlink/pull/13352) [`33a9cdf`](https://github.com/smartcontractkit/chainlink/commit/33a9cdf8d7c42cab4682121b50d5fa12d9b5ff27) Thanks [@shileiwill](https://github.com/shileiwill)! - check data size #bugfix

- [#13177](https://github.com/smartcontractkit/chainlink/pull/13177) [`0d58a8d`](https://github.com/smartcontractkit/chainlink/commit/0d58a8d5db24f42720226e73328e501637ba59c5) Thanks [@shileiwill](https://github.com/shileiwill)! - link transfer status check #bugfix

- [#13183](https://github.com/smartcontractkit/chainlink/pull/13183) [`9630475`](https://github.com/smartcontractkit/chainlink/commit/96304756a77cdb2acf251d21d59b6aa8b55bf61a) Thanks [@cds95](https://github.com/cds95)! - update node signer type

- [#13580](https://github.com/smartcontractkit/chainlink/pull/13580) [`0d4a3b2`](https://github.com/smartcontractkit/chainlink/commit/0d4a3b2cd8ff938ba018d982ef514c754a7df345) Thanks [@DeividasK](https://github.com/DeividasK)! - #internal

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - verify valid node operator when adding nodes to capability registry

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - update uint256 to uint32 for donId declaration in capability config interface

- [#13453](https://github.com/smartcontractkit/chainlink/pull/13453) [`8c98c80`](https://github.com/smartcontractkit/chainlink/commit/8c98c80376c3b6d72bffeab62ee45a74449b6ef5) Thanks [@cds95](https://github.com/cds95)! - return hashed capability ids

- [#13676](https://github.com/smartcontractkit/chainlink/pull/13676) [`ed6b9ad`](https://github.com/smartcontractkit/chainlink/commit/ed6b9ad7909e09d0ff9850b7a1e34a0137762642) Thanks [@RyanRHall](https://github.com/RyanRHall)! - remove stale automation HH tests

- [#13364](https://github.com/smartcontractkit/chainlink/pull/13364) [`fc007a9`](https://github.com/smartcontractkit/chainlink/commit/fc007a94846c178bc9d5203dbff6b6b8c7546a71) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - #bugfix
  fix a funding bug in LinkAvailableBalanceMonitor

- [#13569](https://github.com/smartcontractkit/chainlink/pull/13569) [`f5a70eb`](https://github.com/smartcontractkit/chainlink/commit/f5a70eb09abc9a4d859442c9bd062a74a7ec9c54) Thanks [@DeividasK](https://github.com/DeividasK)! - #internal

- [#13436](https://github.com/smartcontractkit/chainlink/pull/13436) [`f37afb9`](https://github.com/smartcontractkit/chainlink/commit/f37afb9ebaeda10f8b3873b069b8a824e60a81c3) Thanks [@bolekk](https://github.com/bolekk)! - #internal KeystoneFeedsConsumer bytes10 decoding bugfix

- [#13368](https://github.com/smartcontractkit/chainlink/pull/13368) [`000f2cb`](https://github.com/smartcontractkit/chainlink/commit/000f2cb36b7d9b6d046d383c85996ae1ae7a606e) Thanks [@cds95](https://github.com/cds95)! - return don capability config contract config from capability registry

## 1.1.1 - 2024-05-23

### Patch Changes

- [#13031](https://github.com/smartcontractkit/chainlink/pull/13031) [`04b42f1dd7`](https://github.com/smartcontractkit/chainlink/commit/04b42f1dd7db449e5267e2491a9ba6971f41b1bf) Thanks [@DeividasK](https://github.com/DeividasK)! - #internal

- [#13058](https://github.com/smartcontractkit/chainlink/pull/13058) [`a34a17ae9d`](https://github.com/smartcontractkit/chainlink/commit/a34a17ae9d62679a1ff15a7703f5cbcf6dfd1d0f) Thanks [@shileiwill](https://github.com/shileiwill)! - withdraw in offchain mode #bugfix

- [#13096](https://github.com/smartcontractkit/chainlink/pull/13096) [`2c08c8c1a5`](https://github.com/smartcontractkit/chainlink/commit/2c08c8c1a58ea4b7c09b0d5a5ca3b8a677beb9f4) Thanks [@shileiwill](https://github.com/shileiwill)! - add upkeepCharged event #bugfix

- [#13185](https://github.com/smartcontractkit/chainlink/pull/13185) [`286e6af423`](https://github.com/smartcontractkit/chainlink/commit/286e6af42357dc08acc23659ac6d486dff458478) Thanks [@RensR](https://github.com/RensR)! - fix solhint issues

- [#13129](https://github.com/smartcontractkit/chainlink/pull/13129) [`811fe603ae`](https://github.com/smartcontractkit/chainlink/commit/811fe603ae67ec13b7f6ca1fb034840d33ad5303) Thanks [@RensR](https://github.com/RensR)! - bump versions and fix solhint

- [#13008](https://github.com/smartcontractkit/chainlink/pull/13008) [`841fe61daa`](https://github.com/smartcontractkit/chainlink/commit/841fe61daa90b980f1e1622d2f7bd8f850b55462) Thanks [@HenryNguyen5](https://github.com/HenryNguyen5)! - #internal Keystone - rename type to id

- [#13036](https://github.com/smartcontractkit/chainlink/pull/13036) [`77efb61e56`](https://github.com/smartcontractkit/chainlink/commit/77efb61e562fb456ac5bef3b6a528904cda000ec) Thanks [@DeividasK](https://github.com/DeividasK)! - #internal

- [#13088](https://github.com/smartcontractkit/chainlink/pull/13088) [`29b16360fb`](https://github.com/smartcontractkit/chainlink/commit/29b16360fb41e4372f72fe744aaf3ee8234a9b67) Thanks [@shileiwill](https://github.com/shileiwill)! - get available erc20s for payment #bugfix

- [#12962](https://github.com/smartcontractkit/chainlink/pull/12962) [`62d31d0a6e`](https://github.com/smartcontractkit/chainlink/commit/62d31d0a6ea724e7ac00d5860d60cc93d88a3f35) Thanks [@archseer](https://github.com/archseer)! - #internal

- [#13103](https://github.com/smartcontractkit/chainlink/pull/13103) [`54f7c9c8f5`](https://github.com/smartcontractkit/chainlink/commit/54f7c9c8f5508d0d0a063eb435404b4164723300) Thanks [@DeividasK](https://github.com/DeividasK)! - implement add don for capability regsitry

- [#12983](https://github.com/smartcontractkit/chainlink/pull/12983) [`644f5f271d`](https://github.com/smartcontractkit/chainlink/commit/644f5f271d9ed47e999e1d9aa4b99e5de0fd8b89) Thanks [@austinborn](https://github.com/austinborn)! - Update operatorforwarder tests and pull out of dev/

- [#13102](https://github.com/smartcontractkit/chainlink/pull/13102) [`700a827194`](https://github.com/smartcontractkit/chainlink/commit/700a82719451611381ab5dbb94fe00547660440b) Thanks [@cds95](https://github.com/cds95)! - implement remove nodes on capability registry

- [#13080](https://github.com/smartcontractkit/chainlink/pull/13080) [`36cc95f625`](https://github.com/smartcontractkit/chainlink/commit/36cc95f6256b5ba418a916de2c9dc9597508147a) Thanks [@cds95](https://github.com/cds95)! - Add function to update nodes in capability registry

- [#13022](https://github.com/smartcontractkit/chainlink/pull/13022) [`2805fa6c9b`](https://github.com/smartcontractkit/chainlink/commit/2805fa6c9b469d535edcd3d66c08e1d22bbaa2d0) Thanks [@cds95](https://github.com/cds95)! - #internal

- [#12812](https://github.com/smartcontractkit/chainlink/pull/12812) [`5b33a3296f`](https://github.com/smartcontractkit/chainlink/commit/5b33a3296f895cec8a23ba2e235989868f398ddb) Thanks [@shileiwill](https://github.com/shileiwill)! - Support decimals #added

- [#12979](https://github.com/smartcontractkit/chainlink/pull/12979) [`0c4c24ad8c`](https://github.com/smartcontractkit/chainlink/commit/0c4c24ad8c95e505cd2a29be711cc40e612658b0) Thanks [@cds95](https://github.com/cds95)! - Add function to remove node operators from capability registry

- [#12991](https://github.com/smartcontractkit/chainlink/pull/12991) [`929312681f`](https://github.com/smartcontractkit/chainlink/commit/929312681fb27529915912e8bd6e4000559ea77f) Thanks [@cds95](https://github.com/cds95)! - Add function to update node operator'

- [#12669](https://github.com/smartcontractkit/chainlink/pull/12669) [`3134ce8868`](https://github.com/smartcontractkit/chainlink/commit/3134ce8868ccc22bd4ae670c8b0bfda5fa78a332) Thanks [@leeyikjiun](https://github.com/leeyikjiun)! - vrfv2plus - account for num words in coordinator gas overhead in v2plus wrapper

- [#12906](https://github.com/smartcontractkit/chainlink/pull/12906) [`365c38be8b`](https://github.com/smartcontractkit/chainlink/commit/365c38be8b589d5ffa0b21755dcb40e2e4205652) Thanks [@cds95](https://github.com/cds95)! - Implement function to add node operators to the capability registry'

## 1.1.0 - 2024-04-23

### Minor Changes

- [#12481](https://github.com/smartcontractkit/chainlink/pull/12481) [`daa90db289`](https://github.com/smartcontractkit/chainlink/commit/daa90db289f84829a607b41792f7d231871a5462) Thanks [@justinkaseman](https://github.com/justinkaseman)! - Chainlink Functions contracts v1.3.0

- [#12489](https://github.com/smartcontractkit/chainlink/pull/12489) [`3a49094db2`](https://github.com/smartcontractkit/chainlink/commit/3a49094db25036e1948818e4030fca11be748914) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - - Misc VRF V2+ contract changes

  - Reuse struct RequestCommitmentV2Plus from VRFTypes
  - Fix interface name IVRFCoordinatorV2PlusFulfill in BatchVRFCoordinatorV2Plus to avoid confusion with IVRFCoordinatorV2Plus.sol
  - Remove unused errors
  - Rename variables for readability
  - Fix comments
  - Minor gas optimisation (++i)

- [#12522](https://github.com/smartcontractkit/chainlink/pull/12522) [`886201638e`](https://github.com/smartcontractkit/chainlink/commit/886201638e14dc478ae7104b4a5aed9ac8af5bba) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - improve foundry tests

- [#12581](https://github.com/smartcontractkit/chainlink/pull/12581) [`6fcc73983e`](https://github.com/smartcontractkit/chainlink/commit/6fcc73983e5b782bb4ac577cb33093bf80e3a582) Thanks [@kidambisrinivas](https://github.com/kidambisrinivas)! - VRFV2PlusWrapper config refactor

- [#12547](https://github.com/smartcontractkit/chainlink/pull/12547) [`8162f7b101`](https://github.com/smartcontractkit/chainlink/commit/8162f7b1012dd669e51bbb4038a6d5df29906267) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - pay deactivated transmitters in offchain settlement

- [#12466](https://github.com/smartcontractkit/chainlink/pull/12466) [`f9d02e3192`](https://github.com/smartcontractkit/chainlink/commit/f9d02e3192f1a35fda05ca69a50f986c9149748f) Thanks [@vreff](https://github.com/vreff)! - Update type and version name for VRFv2+ Wrapper

- [#12469](https://github.com/smartcontractkit/chainlink/pull/12469) [`1370133b72`](https://github.com/smartcontractkit/chainlink/commit/1370133b722ab97650b15c6aeab72bb494790b63) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - implement offchain settlement for NOPs payment

- [#12578](https://github.com/smartcontractkit/chainlink/pull/12578) [`ffd492295f`](https://github.com/smartcontractkit/chainlink/commit/ffd492295f03de8c3b946a003dacbded731d7899) Thanks [@RensR](https://github.com/RensR)! - Removed 0.6 and 0.7 Solidity source code

- [#12418](https://github.com/smartcontractkit/chainlink/pull/12418) [`22114fb20a`](https://github.com/smartcontractkit/chainlink/commit/22114fb20a67e2263ffb6d445530559f02423809) Thanks [@RyanRHall](https://github.com/RyanRHall)! - introduce native billing support to automation registry v2.3

- [#12583](https://github.com/smartcontractkit/chainlink/pull/12583) [`50724c3bb1`](https://github.com/smartcontractkit/chainlink/commit/50724c3bb1fb959f85d361bc0615f58cc16e4fc9) Thanks [@FelixFan1992](https://github.com/FelixFan1992)! - bug fixes in s_reserveAmount accounting

- [#12569](https://github.com/smartcontractkit/chainlink/pull/12569) [`98ef65add8`](https://github.com/smartcontractkit/chainlink/commit/98ef65add85dc4c22333bf413fa7b593c501212d) Thanks [@RensR](https://github.com/RensR)! - removed 0.4 and 0.5 contracts

### Patch Changes

- [#12612](https://github.com/smartcontractkit/chainlink/pull/12612) [`d44abe3769`](https://github.com/smartcontractkit/chainlink/commit/d44abe37693d6995377fa1329e433e7fba26885d) Thanks [@RensR](https://github.com/RensR)! - upgraded transmission to 0.8.19

- [#12582](https://github.com/smartcontractkit/chainlink/pull/12582) [`684afa4e1f`](https://github.com/smartcontractkit/chainlink/commit/684afa4e1fcb2cad292cbc3b97ebeda3e3ef7bc8) Thanks [@RyanRHall](https://github.com/RyanRHall)! - fix bug in auto2.3 withdrawERC20Fees

- [#12591](https://github.com/smartcontractkit/chainlink/pull/12591) [`b3086d0ec2`](https://github.com/smartcontractkit/chainlink/commit/b3086d0ec2565badaafdbb9c26e30caeb6fb41c9) Thanks [@RyanRHall](https://github.com/RyanRHall)! - fix withdraw LINK bug in auto 2.3

- [#12497](https://github.com/smartcontractkit/chainlink/pull/12497) [`3ca3494450`](https://github.com/smartcontractkit/chainlink/commit/3ca34944507b01b7d4511d8ea8aff402c0a7bb85) Thanks [@RyanRHall](https://github.com/RyanRHall)! - added logic C contract to automation 2.3

- [#12479](https://github.com/smartcontractkit/chainlink/pull/12479) [`93762ccbd8`](https://github.com/smartcontractkit/chainlink/commit/93762ccbd868b9e227abf3220afb9ad22ba41b92) Thanks [@jinhoonbang](https://github.com/jinhoonbang)! - upgrade solc version to 0.8.19 for vrf v2.5 contracts

- [#12619](https://github.com/smartcontractkit/chainlink/pull/12619) [`6531e34865`](https://github.com/smartcontractkit/chainlink/commit/6531e348659e2b7048b25183eadddbcb10426741) Thanks [@RyanRHall](https://github.com/RyanRHall)! - enable gas tests for auto 2.3

- [#12536](https://github.com/smartcontractkit/chainlink/pull/12536) [`87b0d8f309`](https://github.com/smartcontractkit/chainlink/commit/87b0d8f3091e3276cd049d3a852ab63e4d6bda5b) Thanks [@shileiwill](https://github.com/shileiwill)! - billing overrides

- [#12614](https://github.com/smartcontractkit/chainlink/pull/12614) [`93ff878b2d`](https://github.com/smartcontractkit/chainlink/commit/93ff878b2d88f6e928cdb6a8a830fb8ee100bddd) Thanks [@RensR](https://github.com/RensR)! - rm hh coverage

- [#12613](https://github.com/smartcontractkit/chainlink/pull/12613) [`dd333e977f`](https://github.com/smartcontractkit/chainlink/commit/dd333e977f0c39509250a2bd40295da279726496) Thanks [@RensR](https://github.com/RensR)! - mv vrf foundry tests

- [#12529](https://github.com/smartcontractkit/chainlink/pull/12529) [`84913bfcfc`](https://github.com/smartcontractkit/chainlink/commit/84913bfcfcfcf6f93fb359814208a32e3e659d23) Thanks [@ibrajer](https://github.com/ibrajer)! - VRFV2PlusWrapper contract: subID param added to the constructor, removed migrate() method

- [#12555](https://github.com/smartcontractkit/chainlink/pull/12555) [`cda84cb1b7`](https://github.com/smartcontractkit/chainlink/commit/cda84cb1b7582379ac140b3a88da6179275dbefb) Thanks [@shileiwill](https://github.com/shileiwill)! - safeTransfer and cleanups

- [#12640](https://github.com/smartcontractkit/chainlink/pull/12640) [`76e507e849`](https://github.com/smartcontractkit/chainlink/commit/76e507e8490933318e5b36cf103d6157f9fa2f34) Thanks [@RensR](https://github.com/RensR)! - fix solhint issues in automation folder

- [#12553](https://github.com/smartcontractkit/chainlink/pull/12553) [`4892376917`](https://github.com/smartcontractkit/chainlink/commit/4892376917a16253165dc761f8efec41da32ec9c) Thanks [@RyanRHall](https://github.com/RyanRHall)! - address TODOs and docs for 2.3

- [#12499](https://github.com/smartcontractkit/chainlink/pull/12499) [`1a36386481`](https://github.com/smartcontractkit/chainlink/commit/1a363864816a3e7821d5a5844f13be360f0ecb58) Thanks [@RyanRHall](https://github.com/RyanRHall)! - auto 2.3 foundry test refactor

- [#12661](https://github.com/smartcontractkit/chainlink/pull/12661) [`3b02047754`](https://github.com/smartcontractkit/chainlink/commit/3b020477548c17ed786036494ccc733107ca4152) Thanks [@RyanRHall](https://github.com/RyanRHall)! - more auto 2.3 tests

- [#12598](https://github.com/smartcontractkit/chainlink/pull/12598) [`e753637e01`](https://github.com/smartcontractkit/chainlink/commit/e753637e01fabb8ea3760eb14204124c8d3b88e1) Thanks [@RyanRHall](https://github.com/RyanRHall)! - small gas fix

- [#12540](https://github.com/smartcontractkit/chainlink/pull/12540) [`17c037678d`](https://github.com/smartcontractkit/chainlink/commit/17c037678d05c88f28a28a3ac760c742f549d5ec) Thanks [@RyanRHall](https://github.com/RyanRHall)! - change auto 2.3 flat fees from link to USD

- [#12592](https://github.com/smartcontractkit/chainlink/pull/12592) [`b512ef5a7d`](https://github.com/smartcontractkit/chainlink/commit/b512ef5a7d1bc87d0cbd5357c5c47cc0dcb75e0b) Thanks [@ibrajer](https://github.com/ibrajer)! - Set LINK native feed in VRFV2PlusWrapper to immutable

- [#12577](https://github.com/smartcontractkit/chainlink/pull/12577) [`cf00183f62`](https://github.com/smartcontractkit/chainlink/commit/cf00183f6295fe95979b460f89bcc65f22237fd4) Thanks [@shileiwill](https://github.com/shileiwill)! - add billing override test

- [#12443](https://github.com/smartcontractkit/chainlink/pull/12443) [`e604a73d7b`](https://github.com/smartcontractkit/chainlink/commit/e604a73d7b21c5f053631d9c8afeb0eaf7203310) Thanks [@shileiwill](https://github.com/shileiwill)! - use common interface for v2.3

- [#12565](https://github.com/smartcontractkit/chainlink/pull/12565) [`b673505a91`](https://github.com/smartcontractkit/chainlink/commit/b673505a91719d42ff1a60623f1cfea26d186e56) Thanks [@RensR](https://github.com/RensR)! - bump solhint and address issues, remove unused imports

- [#12477](https://github.com/smartcontractkit/chainlink/pull/12477) [`b2576475fc`](https://github.com/smartcontractkit/chainlink/commit/b2576475fc5c8ac037fc569fddc56e9d515ae7ca) Thanks [@shileiwill](https://github.com/shileiwill)! - support native payment

- [#12531](https://github.com/smartcontractkit/chainlink/pull/12531) [`88e010d604`](https://github.com/smartcontractkit/chainlink/commit/88e010d604682c54c4f99e0a0916f94c0d13ece6) Thanks [@jinhoonbang](https://github.com/jinhoonbang)! - increase num optimizations to 500 for v2.5 coordinator

## 1.0.0 - 2024-03-25

- Moved `VRFCoordinatorV2Mock.sol` to src/v0.8/vrf/mocks
- Moved `VRFCoordinatorMock.sol` to src/v0.8/vrf/mocks
- Move Functions v1.0.0 contracts out of dev. New dev folder for v1.X (#10941)
- Release Functions v1.1.0 contracts. Move v1.1.0 out of dev (#11431)
  - Add minimumEstimateGasPriceWei to Functions Coordinator config (#10916)
  - Remove redundant Functions Coordinator commitment & request id checks (#10975)
  - Add L2 fee contract for Arbitrum, Optimism, and Base (#11102 & #11275)
  - Functions Request IDs are now globally unique (#10891)
  - Add an event for broken down billing costs (#11185)
  - Add custom errors to OCR2Base contract (#11249)
- Updated AutomationBase interface to check for ready only address on polygon

### Removed

- Removed all code related to versions prior to Solidity 0.8.0 (#10931)

## 0.8.0 - 2023-10-04

### Changed

- Add a re-entrancy guard to VRFCoordinatorV2Mock to mimic VRFCoordinatorV2's behavior (#10585)
- Enhanced support for destination configs in Data Streams verifiers (#10472)
- Update Data Streams proxy and billing interfaces for better UX (#10603)
- Allow new reward recipients to be added to pools in Data Streams reward management (#10658)
- Reorganize Data Streams contracts (llo-feeds/) (#10727)
- Release automation 2.1 contracts (#10587)
  - Note: consumers should only use IKeeperRegistryMaster when interacting with the registry contract
- Fix Functions v1 OracleWithdrawAll to correctly use transmitters (#10392)
- Clean up unused Functions v1 code: FunctionsBilling.sol maxCallbackGasLimit & FunctionsRequest.sol requestSignature (#10509)
- Fix Functions v1 FunctionsBilling.sol gas price naming to reflect that it is in wei, not gwei (#10509)
- Use Natspec comment lines in Functions v1 contracts (#10567)
- Functions v1 Subscriptions now require a minimum number of requests to release a deposit amount (#10513)
- Fix Functions v1 Subscriptions add consumer checks for when maximum consumers changes in contract configuration (#10511)
- Functions v1 Router no longer reverts during fulfillment on an invalid client (#10511)
- Functions v1 Coordinator oracleWithdrawAll checks for 0 balances (#10511)

## 0.7.1 - 2023-09-20

### Changed

- Add Chainlink Functions v1.0.0 (#9365)
- Change Functions Client variables to internal for use when integrating Automation (#8429)
- Make Functions Billing Registry and Functions Oracle upgradable using the transparent proxy pattern (#8371)
- Update dependency hardhat from version 2.10.1 to 2.12.7 (#8464)
- Fix Functions cost estimation not correctly using registry fee (#8502)
- Fix Functions transmitter NOP fee payment (#8557)
- Functions library uses solidty-cborutils CBOR v2.0.0 and ENS Buffer v0.1.0(#8485)
- Gas optimization to AuthorizedOriginReceiverUpgradable by using EnumberableSet .values()
- Remove support for inline secrets in Functions requests (#8847)
- Moved versioned directories to use v prefix

## 0.6.1 - 2023-02-06

### Added

- Support for off-chain secrets in Functions Oracle contract

### Changed

- Modified FunctionsClientExample.sol to use constant amount of gas regardless of response size
- Fixed comments in FunctionsBillingRegistry.sol
- Make Functions billing registry's timeoutRequest pausable (#8299)
- Remove user specified gas price from Functions Oracle sendRequest
  (#8320)

## 0.6.0 - 2023-01-11

### Added

- Added a Solidity style guide.

### Changed

- Migrated and improved `AuthorizedReceiverInterface` and `AuthorizedReceiver` from 0.7.0
- Added `Chainlink Functions` interfaces and contracts (initial version for PoC)

## 0.5.1 - 2022-09-27

- Rename `KeeperBase` -> `AutomationBase` and add alias for backwards compatibility
- Rename `KeeperCompatible` -> `AutomationCompatible` and add alias for backwards compatibility
- Rename `KeeperCompatibleInterface` -> `AutomationCompatibleInterface` and add alias for backwards compatibility
- Rename `KeeperRegistryInterface1_2` -> `AutomationRegistryInterface1_2` and add alias for backwards compatibility

## 0.5.0 - 2022-09-26

### Changed

- Fix EIP-150 Bug in VRFV2Wrapper.sol (b9d8261eaa05838b9b609ea02005ecca3b6adca3)
- Added a new UpkeepFormat version `V2` in `UpkeepFormat`
- Renamed `KeeperRegistry` to `KeeperRegistry1_2` and `KeeperRegistryInterface` to `KeeperRegistryInterface1_2`
- Updated `UpkeepTranscoder` to only do a pass-through for upkeep bytes

## 0.4.2 - 2022-07-20

### Changed

- Downgrade 0.8.13 contracts to 0.8.6 due to [this solc bug](https://medium.com/certora/overly-optimistic-optimizer-certora-bug-disclosure-2101e3f7994d).
- Reintroduce v0.6 `EACAggregatorProxy` after removing it in [this commit](https://github.com/smartcontractkit/chainlink/commit/558f42f5122779cb2e05dc8c2b84d1ae78cc0d71)
- Ignore status update in `ArbitrumSequencerUptimeFeed` if incoming update has stale timestamp
- Revert to using current Arbitrum seq status flag in `ArbitrumSequencerUptimeFeed`
- Moved `VRFV2Wrapper`, `VRFV2WrapperConsumerBase` and `interfaces/VRFV2WrapperInterface` out of `dev` folder.

## 0.4.1 - 2022-05-09

### Changed

- VRFv2 contract pragma versions changed from `^0.8.0` to `^0.8.4`.

## 0.4.0 - 2022-02-07

### Added

- `ArbitrumSequencerUptimeFeedInterface` and `ArbitrumSequencerUptimeFeed` added in v0.8.

### Changed

- Changed `ArbitrumValidator#validate` target to `ArbitrumSequencerUptimeFeed` instead of
  Flags contract.
- Moved `VRFConsumerBaseV2` out of dev

## 0.3.1 - 2022-01-05

### Changed:

- Fixed install issue with npm.

## 0.3.0 - 2021-12-09

### Added

- Prettier Solidity formatting applied to v0.7 and above.
- ERC677ReceiverInterface added in v0.8.
- `KeeperBase.sol` and `KeeperCompatible.sol` in Solidity v0.6 and v0.8

### Changed:

- Operator Contract and Chainlink Client are officially supported. This enables
  multiword requests/response are available through the ChainlinkClient by using
  the newly enabled `buildOperatorRequest` along with `sendOperatorRequest` or
  `sendOperatorRequestTo`.
- `ChainlinkClient` functions `requestOracleData` and `requestOracleDataFrom` have been changed to `sendChainlinkRequest` and
  `sendChainlinkRequestTo` respectively.
- Updated function comments in `v0.6/interfaces/KeeperCompatibleInterface.sol` and `v0.8/interfaces/KeeperCompatibleInterface.sol` to match the latest in v0.7.
- Add `DelegateForwarderInterface` interface and `CrossDomainDelegateForwarder` base contract which implements a new `forwardDelegate()` function to forward delegatecalls from L1 to L2.

## 0.2.2 - 2021-09-21

### Added:

- v0.8 Access Controlled contracts (`SimpleWriteAccessController` and `SimpleReadAccessController`).
- v0.8 Flags contracts (`Flags`).
- v0.8 Contracts for the V2 VRF. `VRFCoordinatorV2.sol`, `VRF.sol`,
  `VRFConsumerBaseV2.sol`, `VRFCoordinatorV2Interface.sol`. Along
  with related test contract `VRFConsumerV2.sol` and example contracts
  `VRFSingleConsumerExample.sol` and `VRFConsumerExternalSubOwnerExampl.sol`.
- v0.6 `MockV3Aggregator` in src/v0.6/tests/.
- v0.7 Added keeper-related smart contracts from the keeper repo. Added tests for `KeeperRegistry` and `UpkeepRegistrationRequests` in `test/v0.7/`.

### Changed:

- Move `Operator` and associated contracts (`AuthorizedForwarder`, `AuthorizedReceiver`, `LinkTokenReceiver`, `OperatorFactory`) from `./src/v0.7/dev/` to `./src/v0.7/`.
- Updated `Denominations` in `./src/<version>` to include additional fiat currencies.
- Updated `./src/v0.8/vender/BufferChainlink.sol` with latest unchecked math version.

## 0.2.1 - 2021-07-13

### Changed:

- Bump hardhat from 2.3.3 to 2.4.1
- Move Solidity version 0.8.x contracts `ChainlinkClient.sol`, `Chainlink.sol`, `VRFConsumerBase.sol` and `VRFRequestIDBase.sol` from `./src/v0.8/dev/` to `./src/v0.8/`.
- Updated `FeedRegistryInterface` to use `base` and `quote` parameter names.
- Move `Denominations` from `./src/<version>/dev/` to `./src/<version>`

## 0.2.0 - 2021-07-01

### Added:

- `@chainlink/contracts` package changelog.
- `KeeperCompatibleInterface` contracts.
- Feeds Registry contracts: `FeedRegistryInterface` and `Denominations`.
- v0.8 Consumable contracts (`ChainlinkClient`, `VRFConsumerBase` and aggregator interfaces).
- Multi-word response handling in v0.7 and v0.8 `ChainlinkClient` contracts.

### Changed:

- Added missing licensees to `KeeperComptibleInterface`'s
- Upgrade solidity v8 compiler version from 0.8.4 to 0.8.6
- Tests converted to Hardhat.
- Ethers upgraded from v4 to v5.
- Contract artifacts in `abi/` are now raw abi .json files, and do not include bytecode or other supplimentary data.

### Removed:

- Removed dependencies: `@chainlink/belt`, `@chainlink/test-helpers` and `@truffle`.
- Ethers and Truffle contract artifacts are no longer published.
