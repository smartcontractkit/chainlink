// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

/* External Imports */
import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title OVM_GasPriceOracle
 * @dev This contract exposes the current l2 gas price, a measure of how congested the network
 * currently is. This measure is used by the Sequencer to determine what fee to charge for
 * transactions. When the system is more congested, the l2 gas price will increase and fees
 * will also increase as a result.
 *
 * All public variables are set while generating the initial L2 state. The
 * constructor doesn't run in practice as the L2 state generation script uses
 * the deployed bytecode instead of running the initcode.
 */
contract OVM_GasPriceOracle is Ownable {
    /*************
     * Variables *
     *************/

    // Current L2 gas price
    uint256 public gasPrice;
    // Current L1 base fee
    uint256 public l1BaseFee;
    // Amortized cost of batch submission per transaction
    uint256 public overhead;
    // Value to scale the fee up by
    uint256 public scalar;
    // Number of decimals of the scalar
    uint256 public decimals;

    /***************
     * Constructor *
     ***************/

    /**
     * @param _owner Address that will initially own this contract.
     */
    constructor(address _owner) Ownable() {
        transferOwnership(_owner);
    }

    /**********
     * Events *
     **********/

    event GasPriceUpdated(uint256);
    event L1BaseFeeUpdated(uint256);
    event OverheadUpdated(uint256);
    event ScalarUpdated(uint256);
    event DecimalsUpdated(uint256);

    /********************
     * Public Functions *
     ********************/

    /**
     * Allows the owner to modify the l2 gas price.
     * @param _gasPrice New l2 gas price.
     */
    // slither-disable-next-line external-function
    function setGasPrice(uint256 _gasPrice) public onlyOwner {
        gasPrice = _gasPrice;
        emit GasPriceUpdated(_gasPrice);
    }

    /**
     * Allows the owner to modify the l1 base fee.
     * @param _baseFee New l1 base fee
     */
    // slither-disable-next-line external-function
    function setL1BaseFee(uint256 _baseFee) public onlyOwner {
        l1BaseFee = _baseFee;
        emit L1BaseFeeUpdated(_baseFee);
    }

    /**
     * Allows the owner to modify the overhead.
     * @param _overhead New overhead
     */
    // slither-disable-next-line external-function
    function setOverhead(uint256 _overhead) public onlyOwner {
        overhead = _overhead;
        emit OverheadUpdated(_overhead);
    }

    /**
     * Allows the owner to modify the scalar.
     * @param _scalar New scalar
     */
    // slither-disable-next-line external-function
    function setScalar(uint256 _scalar) public onlyOwner {
        scalar = _scalar;
        emit ScalarUpdated(_scalar);
    }

    /**
     * Allows the owner to modify the decimals.
     * @param _decimals New decimals
     */
    // slither-disable-next-line external-function
    function setDecimals(uint256 _decimals) public onlyOwner {
        decimals = _decimals;
        emit DecimalsUpdated(_decimals);
    }

    /**
     * Computes the L1 portion of the fee
     * based on the size of the RLP encoded tx
     * and the current l1BaseFee
     * @param _data Unsigned RLP encoded tx, 6 elements
     * @return L1 fee that should be paid for the tx
     */
    // slither-disable-next-line external-function
    function getL1Fee(bytes memory _data) public view returns (uint256) {
        uint256 l1GasUsed = getL1GasUsed(_data);
        uint256 l1Fee = l1GasUsed * l1BaseFee;
        uint256 divisor = 10**decimals;
        uint256 unscaled = l1Fee * scalar;
        uint256 scaled = unscaled / divisor;
        return scaled;
    }

    // solhint-disable max-line-length
    /**
     * Computes the amount of L1 gas used for a transaction
     * The overhead represents the per batch gas overhead of
     * posting both transaction and state roots to L1 given larger
     * batch sizes.
     * 4 gas for 0 byte
     * https://github.com/ethereum/go-ethereum/blob/9ada4a2e2c415e6b0b51c50e901336872e028872/params/protocol_params.go#L33
     * 16 gas for non zero byte
     * https://github.com/ethereum/go-ethereum/blob/9ada4a2e2c415e6b0b51c50e901336872e028872/params/protocol_params.go#L87
     * This will need to be updated if calldata gas prices change
     * Account for the transaction being unsigned
     * Padding is added to account for lack of signature on transaction
     * 1 byte for RLP V prefix
     * 1 byte for V
     * 1 byte for RLP R prefix
     * 32 bytes for R
     * 1 byte for RLP S prefix
     * 32 bytes for S
     * Total: 68 bytes of padding
     * @param _data Unsigned RLP encoded tx, 6 elements
     * @return Amount of L1 gas used for a transaction
     */
    // solhint-enable max-line-length
    function getL1GasUsed(bytes memory _data) public view returns (uint256) {
        uint256 total = 0;
        for (uint256 i = 0; i < _data.length; i++) {
            if (_data[i] == 0) {
                total += 4;
            } else {
                total += 16;
            }
        }
        uint256 unsigned = total + overhead;
        return unsigned + (68 * 16);
    }
}