// Package coordinator implements the VRF request tracking and accounting
// [ocr2vrftypes.CoordinatorInterface]. This is a plain-langugage outline of its
// logic as of commit xxxxxxxxxxxxxxxx. Be sure to diff the version you're
// working on against that, to check for changes which may invalidate the claims
// here. These notes are highly implementation-specific, and expected to be
// ephemeral (they are primarily intended as preparation for an audit) so it may
// make sense to remove these notes.
//
// Methods on the private coordinator struct, which implements
// [ocr2vrftypes.CoordinatorInterface], are referred to here without
// qualification.
//
// # Recommendations
//
//   - Update eviction times on the blocks/callbacks caches in
//     SetOffChainConfig, based on
//     [ocr2vrftypes.CoordinatorConfig.CacheEvictionWindowSeconds]
//   - Consider ways that an attacker could force this system to request
//     approximately [ocr2vrftypes.CoordinatorConfig.LookbackBlocks] blocks from
//     the [logpoller] on every observation (via getBlockhashesMapping), and
//     whether that would be a significant cost.
//   - There is a potential race condition with the logpoller, where the current
//     height, onchain logs and blockhashes are retrieved in separate calls to
//     the logpoller ([logpoller.LogPoller.LatestBlock],
//     [logpoller.LogPoller.GetBlocksRange] and
//     [logpoller.LogPoller.LogsWithSigs].) Without some kind of check on the
//     consistency of these values, we might end up with a report with responses
//     to requests from one chain, attested with a recentBlockHash from a
//     different chain, due to an intervening re-org. We should probably add an
//     extra method ValidateReport to [ocr2vrftypes.CoordinatorInterface], which
//     checks that all hashes mentioned in a report come from the same chain.
//
// # ReportBlocks
//
// The ReportBlocks method is the primary service. It depends on two caches,
// toBeTransmittedBlocks and toBeTransmittedCallbacks. Both have hard-coded
// initial eviction deadlines of one minute. (Set in [New].) When
// ReportBlocks is called stale elements of the caches are evicted. There is
// also a [ocr2vrftypes.CoordinatorConfig.CacheEvictionWindowSeconds] field,
// which can be set via SetOffChainConfig. However, that does not
// seem to be referenced anywhere in the [coordinator] package, and
// SetOffChainConfig does not update the eviction times on the
// block/callback caches.
//
// ReportBlocks gathers onchain logs for
//
//   - RandomnessRequested (requests to provide randomness for the next beacon
//     block, which the client will retrieve themselves),
//   - RandomnessFulfillmentRequested (requests that the VRF service call back
//     into the client with randomness for the next beacon block),
//   - RandomWordsFulfilled (fulfilled callbacks to clients),
//   - OutputsServed (the beacon blocks for which randomness has been provided.)
//
// It only gathers these logs from the last 1,000 blocks (initially hard-coded
// in [New] as [ocr2vrftypes.CoordinatorConfig.LookbackBlocks]. This setting is
// correctly updated by SetOffChainConfig, and respected by ReportBlocks in the
// call to LogsWithSigs.)
//
// Based on the RandomnessRequested and RandomnessFulfillmentRequested logs,
// ReportBlocks aggregates (via getBlockhashesMappingFromRequests)the set of
// blockhashes for
//  1. beacon blocks where a randomness output is required, and are old enough
//     to be eligible for service, according to the isBlockEligible function.
//  2. the recentBlockHeight for any such beacon block, which has recently been
//     included in an onchain report.
//
// It does this via getBlockhashesMapping, via [logpoller.GetBlocksRange]. The
// implementation of this in [logpoller] requests all blocks between the minimum
// and maximum requested block heights. We need to consider ways that an
// attacker could force this system to request approximately LookbackBlocks
// blocks on every observation, and whether that would be a significant cost.
//
// Next, ReportBlocks, via filterEligibleRandomnessRequests, scans through the
// RandomnessRequested logs for those which are
//  1. old enough to be eligible for service, and
//  2. either not included in the toBeTransmittedBlocks cache, or, if found
//     there, included in a report intended for a different chain than the
//     canonical one in view of the local logpoller (based on the report having
//     a different recentBlockHash for its specified recentBlockHeight from the
//     one returned by the logpoller.)
//
// filterEligibleRandomnessRequests returns a list of such blocks, which are
// then deduplicated into a set by ReportBlocks, blocksRequested.
//
// Next, ReportBlocks goes through the same process for callbacks, via
// filterEligibleCallbacks. The unfulfilled beacon blocks from this are merged
// into blocksRequested, and a list of unfulfilled callbacks is kept in in
// callbacksRequested.
//
// Next, the unfulfilled blocks are compared to the blocks reported in
// OutputsServed logs, via getFulfilledBlocks. Any blockheight/confdelay pairs
// returned from getFulfilledBlocks are deleted from blocksRequested.
//
// The blocks remaining in blocksRequested are successively stored in the return
// value, blocks, until the gas budget for the prospective report is fully used.
// This budget is configured in
// [ocr2vrftypes.CoordinatorConfig.CoordinatorOverhead].
//
// Next, filterUnfulfilledCallbacks produces a list of callbacks which
//  1. Do not have a requestID appearing in any RandomWordsFulfilled logs
//  2. Are sorted according to increasing height at which they'll be eligible,
//     and, at a given eligibility height, their gas allowances.
//  3. Are successively taken from that sorted list as long as the estimated
//     total gas usage is
package coordinator
