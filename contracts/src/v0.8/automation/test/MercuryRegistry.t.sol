pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";
import "../dev/MercuryRegistry.sol";
import "../dev/MercuryRegistryBatchUpkeep.sol";
import "../interfaces/StreamsLookupCompatibleInterface.sol";

contract MercuryRegistryTest is Test {
  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;
  int192 internal constant DEVIATION_THRESHOLD = 10_000; // 1%
  uint32 internal constant STALENESS_SECONDS = 3600; // 1 hour

  address s_verifier = 0x60448B880c9f3B501af3f343DA9284148BD7D77C;

  string[] feedIds;
  string s_BTCUSDFeedId = "0x6962e629c3a0f5b7e3e9294b0c283c9b20f94f1c89c8ba8c1ee4650738f20fb2";
  string s_ETHUSDFeedId = "0xf753e1201d54ac94dfd9334c542562ff7e42993419a661261d010af0cbfd4e34";
  MercuryRegistry s_testRegistry;

  // Feed: BTC/USD
  // Date: Tuesday, August 22, 2023 7:29:28 PM
  // Price: $25,857.11126720
  bytes s_august22BTCUSDMercuryReport =
    hex"0006a2f7f9b6c10385739c687064aa1e457812927f59446cccddf7740cc025ad00000000000000000000000000000000000000000000000000000000014cb94e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001206962e629c3a0f5b7e3e9294b0c283c9b20f94f1c89c8ba8c1ee4650738f20fb20000000000000000000000000000000000000000000000000000000064e50c980000000000000000000000000000000000000000000000000000025a0864a8c00000000000000000000000000000000000000000000000000000025a063481720000000000000000000000000000000000000000000000000000025a0a94d00f000000000000000000000000000000000000000000000000000000000226181f4733a6d98892d1821771c041d5d69298210fdca9d643ad74477423b6a3045647000000000000000000000000000000000000000000000000000000000226181f0000000000000000000000000000000000000000000000000000000064e50c9700000000000000000000000000000000000000000000000000000000000000027f3056b1b71dd516037afd2e636f8afb39853f5cb3ccaa4b02d6f9a2a64622534e94aa1f794f6a72478deb7e0eb2942864b7fac76d6e120bd809530b1b74a32b00000000000000000000000000000000000000000000000000000000000000027bd3b385c0812dfcad2652d225410a014a0b836cd9635a6e7fb404f65f7a912f0b193db57e5c4f38ce71f29170f7eadfa94d972338858bacd59ab224245206db";

  // Feed: BTC/USD
  // Date: Wednesday, August 23, 2023 7:55:02 PM
  // Price: $26,720.37346975
  bytes s_august23BTCUSDMercuryReport =
    hex"0006a2f7f9b6c10385739c687064aa1e457812927f59446cccddf7740cc025ad000000000000000000000000000000000000000000000000000000000159a630000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001206962e629c3a0f5b7e3e9294b0c283c9b20f94f1c89c8ba8c1ee4650738f20fb20000000000000000000000000000000000000000000000000000000064e664160000000000000000000000000000000000000000000000000000026e21d63e9f0000000000000000000000000000000000000000000000000000026e2147576a0000000000000000000000000000000000000000000000000000026e226525d30000000000000000000000000000000000000000000000000000000002286ce7c44fa27f67f6dd0a8bb40c12f0f050231845789f022a82aa5f4b3fe5bf2068fb0000000000000000000000000000000000000000000000000000000002286ce70000000000000000000000000000000000000000000000000000000064e664150000000000000000000000000000000000000000000000000000000000000002e9c5857631172082a47a20aa2fd9f580c1c48275d030c17a2dff77da04f88708ce776ef74c04b9ef6ba87c56d8f8c57e80ddd5298b477d60dd49fb8120f1b9ce000000000000000000000000000000000000000000000000000000000000000248624e0e2341cdaf989098f8b3dee2660b792b24e5251d6e48e3abe0a879c0683163a3a199969010e15353a99926d113f6d4cbab9d82ae90a159af9f74f8c157";

  // Feed: BTC/USD
  // Date: Wednesday, August 23, 2023 8:13:28 PM
  // Price: $26,559.67100000
  bytes s_august23BTCUSDMercuryReport_2 =
    hex"0006a2f7f9b6c10385739c687064aa1e457812927f59446cccddf7740cc025ad000000000000000000000000000000000000000000000000000000000159d009000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001206962e629c3a0f5b7e3e9294b0c283c9b20f94f1c89c8ba8c1ee4650738f20fb20000000000000000000000000000000000000000000000000000000064e668690000000000000000000000000000000000000000000000000000026a63f9bc600000000000000000000000000000000000000000000000000000026a635984c00000000000000000000000000000000000000000000000000000026a67bb929d00000000000000000000000000000000000000000000000000000000022873e999d3ff9b644bba530af933dfaa6c59e31c3e232fcaa1e5f7304e2e79d939da1900000000000000000000000000000000000000000000000000000000022873e80000000000000000000000000000000000000000000000000000000064e66868000000000000000000000000000000000000000000000000000000000000000247c21657a6c2795986e95081876bf8b5f24bf72abd2dc4c601e7c96d654bcf543b5bb730e3d4736a308095e4531e7c03f581ac364f0889922ba3ae24b7cf968000000000000000000000000000000000000000000000000000000000000000020d3037d9f55256a001a2aa79ea746526c7cb36747e1deb4c804311394b4027667e5b711bcecfe60632e86cf8e83c28d1465e2d8d90bc0638dad8347f55488e8e";

  // Feed: ETH/USD
  // Date: Wednesday, August 23, 2023 7:55:01 PM
  // Price: $1,690.76482169
  bytes s_august23ETHUSDMercuryReport =
    hex"0006c41ec94138ae62cce3f1a2b852e42fe70359502fa7b6bdbf81207970d88e00000000000000000000000000000000000000000000000000000000016d874d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000028000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000120f753e1201d54ac94dfd9334c542562ff7e42993419a661261d010af0cbfd4e340000000000000000000000000000000000000000000000000000000064e66415000000000000000000000000000000000000000000000000000000275dbe6079000000000000000000000000000000000000000000000000000000275c905eba000000000000000000000000000000000000000000000000000000275e5693080000000000000000000000000000000000000000000000000000000002286ce7c44fa27f67f6dd0a8bb40c12f0f050231845789f022a82aa5f4b3fe5bf2068fb0000000000000000000000000000000000000000000000000000000002286ce70000000000000000000000000000000000000000000000000000000064e664150000000000000000000000000000000000000000000000000000000000000002a2b01f7741563cfe305efaec43e56cd85731e3a8e2396f7c625bd16adca7b39c97805b6170adc84d065f9d68c87104c3509aeefef42c0d1711e028ace633888000000000000000000000000000000000000000000000000000000000000000025d984ad476bda9547cf0f90d32732dc5a0d84b0e2fe9795149b786fb05332d4c092e278b4dddeef45c070b818c6e221db2633b573d616ef923c755a145ea099c";

  // Feed: USDC/USD
  // Date: Wednesday, August 30, 2023 5:05:01 PM
  // Price: $1.00035464
  bytes s_august30USDCUSDMercuryReport =
    hex"0006970c13551e2a390246f5eccb62b9be26848e72026830f4688f49201b5a050000000000000000000000000000000000000000000000000000000001c89843000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000028000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000120a5b07943b89e2c278fc8a2754e2854316e03cb959f6d323c2d5da218fb6b0ff80000000000000000000000000000000000000000000000000000000064ef69fa0000000000000000000000000000000000000000000000000000000005f5da000000000000000000000000000000000000000000000000000000000005f5b0f80000000000000000000000000000000000000000000000000000000005f5f8b0000000000000000000000000000000000000000000000000000000000240057307d0a0421d25328cb6dcfc5d0e211ff0580baaaf104e9877fc52cf2e8ec0aa7d00000000000000000000000000000000000000000000000000000000024005730000000000000000000000000000000000000000000000000000000064ef69fa0000000000000000000000000000000000000000000000000000000000000002b9e7fb46f1e9d22a1156024dc2bbf2bc6d337e0a2d78aaa3fb6e43b880217e5897732b516e39074ef4dcda488733bfee80c0a10714b94621cd93df6842373cf5000000000000000000000000000000000000000000000000000000000000000205ca5f8da9d6ae01ec6d85c681e536043323405b3b8a15e4d2a288e02dac32f10b2294593e270a4bbf53b0c4978b725293e85e49685f1d3ce915ff670ab6612f";

  function setUp() public virtual {
    // Set owner, and fork Arbitrum Goerli Testnet (chain ID 421613).
    // The fork is only used with the `FORK_TEST` flag enabeld, as to not disrupt CI. For CI, a mock verifier is used instead.
    vm.startPrank(OWNER);
    try vm.envBool("FORK_TEST") returns (bool /* fork testing enabled */) {
      vm.selectFork(vm.createFork("https://goerli-rollup.arbitrum.io/rpc"));
    } catch {
      s_verifier = address(new MockVerifierProxy());
    }
    vm.chainId(31337); // restore chain Id

    // Use a BTC feed and ETH feed.
    feedIds = new string[](2);
    feedIds[0] = s_BTCUSDFeedId;
    feedIds[1] = s_ETHUSDFeedId;

    // Deviation threshold and staleness are the same for all feeds.
    int192[] memory thresholds = new int192[](1);
    thresholds[0] = DEVIATION_THRESHOLD;
    uint32[] memory stalenessSeconds = new uint32[](1);
    stalenessSeconds[0] = STALENESS_SECONDS;

    // Initialize with BTC feed.
    string[] memory initialFeedIds = new string[](1);
    initialFeedIds[0] = feedIds[0];
    string[] memory initialFeedNames = new string[](1);
    initialFeedNames[0] = "BTC/USD";
    s_testRegistry = new MercuryRegistry(
      initialFeedIds,
      initialFeedNames,
      thresholds,
      stalenessSeconds,
      address(0) // verifier unset
    );
    s_testRegistry.setVerifier(s_verifier); // set verifier

    // Add ETH feed.
    string[] memory addedFeedIds = new string[](1);
    addedFeedIds[0] = feedIds[1];
    string[] memory addedFeedNames = new string[](1);
    addedFeedNames[0] = "ETH/USD";
    s_testRegistry.addFeeds(addedFeedIds, addedFeedNames, thresholds, stalenessSeconds);
  }

  function testMercuryRegistry() public {
    // Check upkeep, receive Mercury revert.
    uint256 blockNumber = block.number;
    vm.expectRevert(
      abi.encodeWithSelector(
        StreamsLookupCompatibleInterface.StreamsLookup.selector,
        "feedIdHex", // feedParamKey
        feedIds, // feed Ids
        "blockNumber", // timeParamKey
        blockNumber, // block number on which request is occuring
        "" // extra data
      )
    );
    s_testRegistry.checkUpkeep("");

    // Obtain mercury report off-chain (for August 22 BTC/USD price)
    bytes[] memory values = new bytes[](1);
    values[0] = s_august22BTCUSDMercuryReport;

    // Pass the obtained mercury report into checkCallback, to assert that an update is warranted.
    (bool shouldPerformUpkeep, bytes memory performData) = s_testRegistry.checkCallback(values, bytes(""));
    assertEq(shouldPerformUpkeep, true);

    // Perform upkeep to update on-chain state.
    s_testRegistry.performUpkeep(performData);

    // Check state of BTC/USD feed to ensure update was propagated.
    bytes memory oldPerformData;
    uint32 oldObservationsTimestamp;
    {
      // scoped to prevent stack-too-deep error
      (
        uint32 observationsTimestamp,
        int192 price,
        int192 bid,
        int192 ask,
        string memory feedName,
        string memory localFeedId,
        bool active,
        int192 deviationPercentagePPM,
        uint32 stalenessSeconds
      ) = s_testRegistry.s_feedMapping(s_BTCUSDFeedId);
      assertEq(observationsTimestamp, 1692732568); // Tuesday, August 22, 2023 7:29:28 PM
      assertEq(bid, 2585674416498); //   $25,856.74416498
      assertEq(price, 2585711126720); // $25,857.11126720
      assertEq(ask, 2585747836943); //   $25,857.47836943
      assertEq(feedName, "BTC/USD");
      assertEq(localFeedId, s_BTCUSDFeedId);
      assertEq(active, true);
      assertEq(deviationPercentagePPM, DEVIATION_THRESHOLD);
      assertEq(stalenessSeconds, STALENESS_SECONDS);

      // Save this for later in the test.
      oldPerformData = performData;
      oldObservationsTimestamp = observationsTimestamp;
    }
    // Obtain mercury report off-chain (for August 23 BTC/USD price & ETH/USD price)
    values = new bytes[](2);
    values[0] = s_august23BTCUSDMercuryReport;
    values[1] = s_august23ETHUSDMercuryReport;

    // Pass the obtained mercury report into checkCallback, to assert that an update is warranted.
    (shouldPerformUpkeep, performData) = s_testRegistry.checkCallback(values, bytes(""));
    assertEq(shouldPerformUpkeep, true);

    // Perform upkeep to update on-chain state.
    s_testRegistry.performUpkeep(performData);

    // Make a batch request for both the BTC/USD feed data and the ETH/USD feed data.
    MercuryRegistry.Feed[] memory feeds = s_testRegistry.getLatestFeedData(feedIds);

    // Check state of BTC/USD feed to ensure update was propagated.
    assertEq(feeds[0].observationsTimestamp, 1692820502); // Wednesday, August 23, 2023 7:55:02 PM
    assertEq(feeds[0].bid, 2672027981674); //   $26,720.27981674
    assertEq(feeds[0].price, 2672037346975); // $26,720.37346975
    assertEq(feeds[0].ask, 2672046712275); //   $26,720.46712275
    assertEq(feeds[0].feedName, "BTC/USD");
    assertEq(feeds[0].feedId, s_BTCUSDFeedId);

    // Check state of ETH/USD feed to ensure update was propagated.
    assertEq(feeds[1].observationsTimestamp, 1692820501); // Wednesday, August 23, 2023 7:55:01 PM
    assertEq(feeds[1].bid, 169056689850); //   $1,690.56689850
    assertEq(feeds[1].price, 169076482169); // $1,690.76482169
    assertEq(feeds[1].ask, 169086456584); //   $16,90.86456584
    assertEq(feeds[1].feedName, "ETH/USD");
    assertEq(feeds[1].feedId, s_ETHUSDFeedId);
    assertEq(feeds[1].active, true);
    assertEq(feeds[1].deviationPercentagePPM, DEVIATION_THRESHOLD);
    assertEq(feeds[1].stalenessSeconds, STALENESS_SECONDS);

    // Obtain mercury report off-chain for August 23 BTC/USD price (second report of the day).
    // The price of this incoming report will not deviate enough from the on-chain value to trigger an update,
    // nor is the on-chain data stale enough.
    values = new bytes[](1);
    values[0] = s_august23BTCUSDMercuryReport_2;

    // Pass the obtained mercury report into checkCallback, to assert that an update is not warranted.
    (shouldPerformUpkeep, performData) = s_testRegistry.checkCallback(values, bytes(""));
    assertEq(shouldPerformUpkeep, false);

    // Ensure stale reports cannot be included.
    vm.expectRevert(
      abi.encodeWithSelector(
        MercuryRegistry.StaleReport.selector,
        feedIds[0],
        feeds[0].observationsTimestamp,
        oldObservationsTimestamp
      )
    );
    s_testRegistry.performUpkeep(oldPerformData);

    // Ensure reports for inactive feeds cannot be included.
    bytes[] memory inactiveFeedReports = new bytes[](1);
    inactiveFeedReports[0] = s_august30USDCUSDMercuryReport;
    bytes memory lookupData = "";
    vm.expectRevert(
      abi.encodeWithSelector(
        MercuryRegistry.FeedNotActive.selector,
        "0xa5b07943b89e2c278fc8a2754e2854316e03cb959f6d323c2d5da218fb6b0ff8" // USDC/USD feed id
      )
    );
    s_testRegistry.performUpkeep(abi.encode(inactiveFeedReports, lookupData));
  }

  // Below are the same tests as `testMercuryRegistry`, except done via a batching Mercury registry that
  // consumes the test registry. This is to assert that batching can be accomplished by multiple different
  // upkeep jobs, which can populate the same
  function testMercuryRegistryBatchUpkeep() public {
    MercuryRegistryBatchUpkeep batchedRegistry = new MercuryRegistryBatchUpkeep(
      address(s_testRegistry), // use the test registry as master registry
      0, // start batch at index 0.
      50 // end batch beyond length of feed Ids (take responsibility for all feeds)
    );
    // Check upkeep, receive Mercury revert.
    uint256 blockNumber = block.number;
    vm.expectRevert(
      abi.encodeWithSelector(
        StreamsLookupCompatibleInterface.StreamsLookup.selector,
        "feedIdHex", // feedParamKey
        feedIds, // feed Ids
        "blockNumber", // timeParamKey
        blockNumber, // block number on which request is occuring
        "" // extra data
      )
    );
    batchedRegistry.checkUpkeep("");

    // Obtain mercury report off-chain (for August 22 BTC/USD price)
    bytes[] memory values = new bytes[](1);
    values[0] = s_august22BTCUSDMercuryReport;

    // Pass the obtained mercury report into checkCallback, to assert that an update is warranted.
    (bool shouldPerformUpkeep, bytes memory performData) = batchedRegistry.checkCallback(values, bytes(""));
    assertEq(shouldPerformUpkeep, true);

    // Perform upkeep to update on-chain state.
    batchedRegistry.performUpkeep(performData);

    // Check state of BTC/USD feed to ensure update was propagated.
    (
      uint32 observationsTimestamp,
      int192 price,
      int192 bid,
      int192 ask,
      string memory feedName,
      string memory localFeedId,
      bool active,
      int192 deviationPercentagePPM,
      uint32 stalenessSeconds
    ) = s_testRegistry.s_feedMapping(s_BTCUSDFeedId);
    assertEq(observationsTimestamp, 1692732568); // Tuesday, August 22, 2023 7:29:28 PM
    assertEq(bid, 2585674416498); //   $25,856.74416498
    assertEq(price, 2585711126720); // $25,857.11126720
    assertEq(ask, 2585747836943); //   $25,857.47836943
    assertEq(feedName, "BTC/USD");
    assertEq(localFeedId, s_BTCUSDFeedId);
    assertEq(active, true);
    assertEq(deviationPercentagePPM, DEVIATION_THRESHOLD);
    assertEq(stalenessSeconds, STALENESS_SECONDS);

    // Obtain mercury report off-chain (for August 23 BTC/USD price & ETH/USD price)
    values = new bytes[](2);
    values[0] = s_august23BTCUSDMercuryReport;
    values[1] = s_august23ETHUSDMercuryReport;

    // Pass the obtained mercury report into checkCallback, to assert that an update is warranted.
    (shouldPerformUpkeep, performData) = batchedRegistry.checkCallback(values, bytes(""));
    assertEq(shouldPerformUpkeep, true);

    // Perform upkeep to update on-chain state, but with not enough gas to update both feeds.
    batchedRegistry.performUpkeep{gas: 250_000}(performData);

    // Make a batch request for both the BTC/USD feed data and the ETH/USD feed data.
    MercuryRegistry.Feed[] memory feeds = s_testRegistry.getLatestFeedData(feedIds);

    // Check state of BTC/USD feed to ensure update was propagated.
    assertEq(feeds[0].observationsTimestamp, 1692820502); // Wednesday, August 23, 2023 7:55:02 PM
    assertEq(feeds[0].bid, 2672027981674); //   $26,720.27981674
    assertEq(feeds[0].price, 2672037346975); // $26,720.37346975
    assertEq(feeds[0].ask, 2672046712275); //   $26,720.46712275
    assertEq(feeds[0].feedName, "BTC/USD");
    assertEq(feeds[0].feedId, s_BTCUSDFeedId);

    // Check state of ETH/USD feed to observe that the update was not propagated.
    assertEq(feeds[1].observationsTimestamp, 0);
    assertEq(feeds[1].bid, 0);
    assertEq(feeds[1].price, 0);
    assertEq(feeds[1].ask, 0);
    assertEq(feeds[1].feedName, "ETH/USD");
    assertEq(feeds[1].feedId, s_ETHUSDFeedId);
    assertEq(feeds[1].active, true);
    assertEq(feeds[1].deviationPercentagePPM, DEVIATION_THRESHOLD);
    assertEq(feeds[1].stalenessSeconds, STALENESS_SECONDS);

    // Try again, with sufficient gas to update both feeds.
    batchedRegistry.performUpkeep{gas: 2_500_000}(performData);
    feeds = s_testRegistry.getLatestFeedData(feedIds);

    // Check state of ETH/USD feed to ensure update was propagated.
    assertEq(feeds[1].observationsTimestamp, 1692820501); // Wednesday, August 23, 2023 7:55:01 PM
    assertEq(feeds[1].bid, 169056689850); //   $1,690.56689850
    assertEq(feeds[1].price, 169076482169); // $1,690.76482169
    assertEq(feeds[1].ask, 169086456584); //   $16,90.86456584
    assertEq(feeds[1].feedName, "ETH/USD");
    assertEq(feeds[1].feedId, s_ETHUSDFeedId);

    // Obtain mercury report off-chain for August 23 BTC/USD price (second report of the day).
    // The price of this incoming report will not deviate enough from the on-chain value to trigger an update.
    values = new bytes[](1);
    values[0] = s_august23BTCUSDMercuryReport_2;

    // Pass the obtained mercury report into checkCallback, to assert that an update is not warranted.
    (shouldPerformUpkeep, performData) = batchedRegistry.checkCallback(values, bytes(""));
    assertEq(shouldPerformUpkeep, false);
  }
}

contract MockVerifierProxy is IVerifierProxy {
  function verify(bytes calldata payload) external payable override returns (bytes memory) {
    (, bytes memory reportData, , , ) = abi.decode(payload, (bytes32[3], bytes, bytes32[], bytes32[], bytes32));
    return reportData;
  }
}
