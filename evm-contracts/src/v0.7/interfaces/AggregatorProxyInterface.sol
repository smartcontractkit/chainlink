pragma solidity >=0.7.0;

import "./AggregatorV2V3Interface.sol";

interface AggregatorProxyInterface is AggregatorV2V3Interface {
	// phaseAggregators
	function phaseAggregators(uint16 phaseId) external view returns (address);
	// phaseId
	function phaseId() external view returns (uint16);
	// proposedAggregator
	function proposedAggregator() external view returns (address);
	// proposedGetRoundData
	function proposedGetRoundData(uint80 roundId) external view returns (uint80 id, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound);
	// proposedLatestRoundData
	function proposedLatestRoundData() external view returns (uint80 id, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound);
}