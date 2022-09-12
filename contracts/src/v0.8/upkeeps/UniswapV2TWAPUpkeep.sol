// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/proxy/Proxy.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@uniswap/v2-core/contracts/UniswapV2Pair.sol";
import "@uniswap/v2-core/contracts/UniswapV2Factory.sol";
import "@uniswap/v2-periphery/contracts/libraries/UniswapV2OracleLibrary.sol";
import "@uniswap/lib/contracts/libraries/FixedPoint.sol";
import "../ConfirmedOwner.sol";
import "../KeeperBase.sol";
import "../interfaces/KeeperCompatibleInterface.sol";

contract UniswapV2TWAPUpkeep is KeeperCompatibleInterface, KeeperBase, ConfirmedOwner, Pausable, Proxy {
    using EnumerableSet for EnumerableSet.UintSet;

    event JobExecuted(uint256 indexed jobID, uint256 timestamp);
    event JobCreated(uint256 indexed jobID);
    event JobUpdated(uint256 indexed jobID);
    event JobDeleted(uint256 indexed jobID);

    error TooSoonToPerform(uint256 jobID);
    error JobIDNotFound(uint256 jobID);
    error ExceedsMaxJobs();
    error BadJobSpec();

    struct Job {
        UniswapV2Factory[] unifactories;
        mapping(UniswapV2Factory => Pair) pairs;
        Observation[] averagedObservations;
        uint256 periodSize;
        uint256 lastObservationTimestamp;
    }

    struct Pair {
        UniswapV2Pair unipair;
        address token0;
        address token1;
    }

    struct Observation {
        uint256 price0Cumulative;
        uint256 price1Cumulative;
        uint32  timestamp;
    }

    mapping(uint256 => Job) private s_jobs;
    EnumerableSet.UintSet   private s_activeJobIDs;
    uint256                 private s_nextJobID = 1;

    constructor(address owner, uint256 maxJobs)
        ConfirmedOwner(owner)
    {
        s_maxJobs = maxJobs;
    }

    function query(uint256 jobID, address tokenIn, uint256 amountIn, uint256 startPeriod, uint256 endPeriod)
        public
        view
        returns (uint256 amountOut)
    {
        Job storage job = s_jobs[jobID];
        Observation memory observationStart = job.averagedObservations[startPeriod];
        Observation memory observationEnd   = job.averagedObservations[endPeriod];

        if (tokenIn == token0) {
            return amountIn * (observationEnd.price0Cumulative - observationStart.price0Cumulative) / (observationEnd.timestamp - observationStart.timestamp);
        } else {
            return amountIn * (observationEnd.price1Cumulative - observationStart.price1Cumulative) / (observationEnd.timestamp - observationStart.timestamp);
        }
    }

    function observationIndexOf(uint256 jobID, uint256 timestamp) public view returns (uint8 index) {
        Job storage job = s_jobs[jobID];
        uint256 epochPeriod = timestamp / job.periodSize;
        return uint8(epochPeriod % job.granularity);
    }

    function checkUpkeep(bytes calldata performData)
        external
        override
        whenNotPaused
        cannotExecute
        returns (bool, bytes memory)
    {
        (uint256 jobID) = abi.decode(performData, (uint256));
        if (s_jobs[jobID].lastObservationTimestamp + s_jobs[jobID].granularity > block.timestamp) {
            return (false, bytes(0));
        }
        return (true, performData);
    }

    function performUpkeep(bytes calldata performData)
        external
        override
        whenNotPaused
    {
        (uint256 jobID) = abi.decode(performData, (uint256));

        Job memory job = s_jobs[jobID];
        if (!s_activeJobIDs.contains(jobID)) {
            revert JobIDNotFound(jobID);
        } else if (job.lastObservationTimestamp + job.granularity > block.timestamp) {
            revert TooSoonToPerform(jobID);
        }

        uint256 summedPrice0Cumulative;
        uint256 summedPrice1Cumulative;
        uint256 jul = job.unifactories.length;
        for (uint256 i = 0; i < jul; ) {
            UniswapV2Factory unifactory = job.unifactories[i];
            (uint256 price0Cumulative, uint256 price1Cumulative, ) = UniswapV2OracleLibrary.currentCumulativePrices(address(job.pairs[unifactory].unipair));
            summedPrice0Cumulative += price0Cumulative;
            summedPrice1Cumulative += price1Cumulative;
            unchecked { ++i; }
        }

        s_jobs[jobID].averagedObservations[observationIndexOf(block.timestamp)] = Observation({
            price0Cumulative: summedPrice0Cumulative / jul,
            price1Cumulative: summedPrice1Cumulative / jul,
            timestamp: block.timestamp
        });

        job.lastObservationTimestamp = block.timestamp;

        emit JobExecuted(jobID, block.timestamp);
    }

    function createJob(UniswapV2Factory[] calldata unifactories, address token0, address token1, uint256 periodSize)
        external
        onlyOwner
    {
        if (s_activeJobIDs.length() >= s_maxJobs) {
            revert ExceedsMaxJobs();
        } else if (unifactories.length == 0 || token0 == address(0) || token1 == address(0) || granularity <= 1) {
            revert BadJobSpec();
        }

        uint256 jobID = s_nextJobID;
        s_nextJobID++;
        s_activeJobIDs.add(jobID);

        _setJob(jobID, unifactories, token0, token1, periodSize);

        // Populate the array with empty observations
        for (uint256 i = 0; i < s_jobs[jobID].granularity; ) {
            s_jobs[jobID].averagedObservations.push();
            unchecked { ++i; }
        }

        emit JobCreated(jobID);
    }

    function updateJob(uint256 jobID, UniswapV2Factory[] calldata unifactories, uint256 periodSize)
        external
        onlyOwner
    {
        if (!s_activeJobIDs.contains(jobID)) {
            revert JobIDNotFound(jobID);
        }
        _setJob(jobID, unifactories, token0, token1, periodSize);
        emit JobUpdated(id, newTarget, newHandler);
    }

    function _setJob(uint256 jobID, UniswapV2Factory[] calldata unifactories, address token0, address token1, uint256 periodSize)
        internal
    {
        s_jobs[jobID] = Job({
            unifactories: unifactories,
            periodSize: periodSize
        });

        Job storage job = s_jobs[jobID];
        uint256 fl = unifactories.length;
        for (uint256 i = 0; i < fl; ) {
            job.pairs[unifactories[i]] = Pair({
                unipair: unifactories[i].getPair(token0, token1),
                token0: ,
                token1: ,
            });
            unchecked { ++i; }
        }
    }

    /**
      * @notice Deletes the job matching the provided id. Reverts if
      * the id is not found.
      * @param id the id of the job to delete
      */
    function deleteJob(uint256 jobID) external onlyOwner {
        if (!s_activeJobIDs.contains(jobID)) {
            revert JobIDNotFound(jobID);
        }
        delete s_jobs[jobID];
        s_activeJobIDs.remove(jobID);
        emit JobDeleted(jobID);
    }

    /**
      * @notice gets a list of active job IDs
      * @return list of active job IDs
      */
    function getActiveJobIDs() external view returns (uint256[] memory) {
        uint256 length = s_activeJobIDs.length();
        uint256[] memory jobIDs = new uint256[](length);
        for (uint256 idx = 0; idx < length; idx++) {
            jobIDs[idx] = s_activeJobIDs.at(idx);
        }
        return jobIDs;
    }

    /**
      * @notice Pauses the contract, which prevents executing performUpkeep
      */
    function pause() external onlyOwner {
        _pause();
    }

    /**
      * @notice Unpauses the contract
      */
    function unpause() external onlyOwner {
        _unpause();
    }

}
