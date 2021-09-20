/// <reference types="node" />
import { ethers, utils } from 'ethers';
import { BigNumberish } from 'ethers/utils';
export interface ServiceAgreement {
    /**
     * Price in LINK to request a report based on this agreement
     *
     * @solformat uint256
     */
    payment: ethers.utils.BigNumberish;
    /**
     * Expiration is the amount of time an oracle has to answer a request
     *
     * @solformat uint256
     */
    expiration: ethers.utils.BigNumberish;
    /**
     * The service agreement is valid until this time
     *
     * @solformat uint256
     */
    endAt: ethers.utils.BigNumberish;
    /**
     * An array of oracle addresses to use within the process of aggregation
     *
     * @solformat address[]
     */
    oracles: (string | ethers.Wallet)[];
    /**
     * This effectively functions as an ID tag for the off-chain job of the
     * service agreement. It is calculated as the keccak256 hash of the
     * normalized JSON request to create the ServiceAgreement, but that identity
     * is unused, and its value is essentially arbitrary.
     *
     * @solformat bytes32
     */
    requestDigest: string;
    /**
     *  Specification of aggregator interface. See ../../../evm/contracts/tests/MeanAggregator.sol
     *  for example.
     */
    /**
     * Address of where the aggregator instance is held
     *
     * @solformat address
     */
    aggregator: string;
    /**
     * Selectors for the interface methods must be specified, because their
     * arguments can vary from aggregator to aggregator.
     *
     * Function selector for aggregator initiateJob method
     *
     * @solformat bytes4
     */
    aggInitiateJobSelector: string;
    /**
     * Function selector for aggregator fulfill method
     *
     * @solformat bytes4
     */
    aggFulfillSelector: string;
}
/**
 * A collection of multiple oracle signatures stored via parallel arrays
 */
export interface OracleSignatures {
    /**
     * The recovery parameters normalized for Solidity, either 27 or 28
     *
     * @solformat uint8[]
     */
    vs: ethers.utils.BigNumberish[];
    /**
     * the r coordinate within (r, s) public point of a signature
     *
     * @solformat bytes32[]
     */
    rs: string[];
    /**
     * the s coordinate within (r, s) public point of a signature
     *
     * @solformat  bytes32[]
     */
    ss: string[];
}
/**
 * Create a service agreement with sane testing defaults
 *
 * @param overrides Values to override service agreement defaults
 */
export declare function serviceAgreement(overrides: Partial<ServiceAgreement>): ServiceAgreement;
/**
 * Check that all values for the struct at this SAID have default values.
 *
 * For example, when an invalid service agreement initialization request is made to a `Coordinator`, we want to make sure that
 * it did not initialize its service agreement struct to any value, hence checking for it being empty.
 *
 * @param coordinator The coordinator contract
 * @param serviceAgreementID The service agreement ID
 *
 * @throws when any of payment, expiration, endAt, requestDigest are non-empty
 */
export declare function assertServiceAgreementEmpty(sa: Omit<ServiceAgreement, 'oracles'>): void;
/**
 * Create parameters needed for the
 * ```solidity
 *   function initiateServiceAgreement(
 *    bytes memory _serviceAgreementData,
 *    bytes memory _oracleSignaturesData
 *  )
 * ```
 * method of the `Coordinator.sol` contract
 *
 * @param overrides Values to override the defaults for creating a service agreement
 */
export declare function initiateSAParams(overrides: Partial<ServiceAgreement>): Promise<[string, string]>;
/**
 * ABI encode a service agreement object
 *
 * @param sa The service agreement to encode
 */
export declare function encodeServiceAgreement(sa: ServiceAgreement): string;
/**
 * Generate the unique identifier of a service agreement by computing its
 * digest.
 *
 * @param sa The service agreement to compute the digest of
 */
export declare function generateSAID(sa: ServiceAgreement): ReturnType<typeof ethers.utils.keccak256>;
/**
 * ABI encode the javascript representation of OracleSignatures
 *```solidity
 *  struct OracleSignatures {
 *    uint8[] vs;
 *    bytes32[] rs;
 *    bytes32[] ss;
 *  }
 * ```
 *
 * @param os The oracle signatures to ABI encode
 */
export declare function encodeOracleSignatures(os: OracleSignatures): string;
/**
 * Abi encode the oracleRequest() method for `Coordinator.sol`
 * ```solidity
 *  function oracleRequest(
 *    address _sender,
 *    uint256 _amount,
 *    bytes32 _sAId,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _nonce,
 *    uint256 _dataVersion,
 *    bytes calldata _data
 *  )
 * ```
 *
 * @param sAID The service agreement ID
 * @param callbackAddr The callback contract address for the response
 * @param callbackFunctionId The callback function id for the response
 * @param nonce The nonce sent by the requester
 * @param data The CBOR payload of the request
 */
export declare function encodeOracleRequest(specId: string, to: string, fHash: string, nonce: BigNumberish, dataBytes: string): string;
/**
 * Generates the oracle signatures on a ServiceAgreement
 *
 * @param serviceAgreement The service agreement to sign
 * @param signers The list oracles that will sign the service agreement
 */
export declare function generateOracleSignatures(serviceAgreement: ServiceAgreement): Promise<OracleSignatures>;
/**
 * Signs a message according to ethereum specs by first appending
 * "\x19Ethereum Signed Message:\n' + <message.length>" to the message
 *
 * @param message The message to sign - either a Buffer or a hex string
 * @param wallet The wallet of the signer
 */
export declare function personalSign(message: Buffer | string, wallet: ethers.Wallet): Promise<Required<utils.Signature>>;
/**
 * Recovers the address of the signer of a message
 *
 * @param message The message that was signed
 * @param signature The signature on the message
 */
export declare function recoverAddressFromSignature(message: string | Buffer, signature: Required<utils.Signature>): string;
/**
 * Combine v, r, and s params of multiple signatures into format expected by contracts
 *
 * @param signatures The list of signatures to combine
 */
export declare function combineOracleSignatures(signatures: Required<utils.Signature>[]): OracleSignatures;
//# sourceMappingURL=coordinator.d.ts.map