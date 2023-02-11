// Package coordinator implements the VRF request tracking and accounting
// CoordinatorInterface from the ocr2vrf types package. This is a
// plain-langugage outline of its logic as of commit xxxxxxxxxxxxxxxx. Be sure
// to diff the version you're working on against that, to check for changes
// which may invalidate the claims here.
//
// The ReportBlocks method is the primary service. It depends on two caches,
// toBeTransmittedBlocks and toBeTransmittedCallbacks. Both have hard-coded
// eviction deadlines of one minute. (Set in New function.) When ReportBlocks is
// called stale elements of the caches are evicted.
//
// ReportBlocks gathers onchain logs for
//
//   - randomness requests (requests to provide randomness for the next beacon
//     block, which the client will retrieve themselves),
//   - randomness fulfillments (requests that the VRF service call back into the
//     client with randomness for the next beacon block),
//   - random words fulfilled (fulfilled callbacks to clients
//   - outputs served (the beacon blocks for which randomness has been
//     provided.)
//
// It only gathers these logs from the last 1,000 blocks (hard-coded in New.)
//
//

package coordinator
