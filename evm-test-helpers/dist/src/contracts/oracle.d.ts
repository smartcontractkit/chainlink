/// <reference types="node" />
/**
 * @packageDocumentation
 *
 * This file provides convenience functions to interact with existing solidity contract abstraction libraries, such as
 * @truffle/contract and ethers.js specifically for our `Oracle.sol` solidity smart contract.
 */
import { ethers } from 'ethers';
import { BigNumberish } from 'ethers/utils';
/**
 * Transaction options such as gasLimit, gasPrice, data, ...
 */
declare type TxOptions = Omit<ethers.providers.TransactionRequest, 'to' | 'from'>;
/**
 * A run request is an event emitted by `Oracle.sol` which triggers a job run
 * on a receiving chainlink node watching for RunRequests coming from that
 * specId + optionally requester.
 */
export interface RunRequest {
    /**
     * The ID of the job spec this request is targeting
     *
     * @solformat bytes32
     */
    specId: string;
    /**
     * The requester of the run
     *
     * @solformat address
     */
    requester: string;
    /**
     * The ID of the request, check Oracle.sol#oracleRequest to see how its computed
     *
     * @solformat bytes32
     */
    requestId: string;
    /**
     * The amount of LINK used for payment
     *
     * @solformat uint256
     */
    payment: string;
    /**
     * The address of the contract instance to callback with the fulfillment result
     *
     * @solformat address
     */
    callbackAddr: string;
    /**
     * The function selector of the method that the oracle should call after fulfillment
     *
     * @solformat bytes4
     */
    callbackFunc: string;
    /**
     * The expiration that the node should respond by before the requester can cancel
     *
     * @solformat uint256
     */
    expiration: string;
    /**
     * The specified data version
     *
     * @solformat uint256
     */
    dataVersion: number;
    /**
     * The CBOR encoded payload of the request
     *
     * @solformat bytes
     */
    data: Buffer;
    /**
     * The hash of the signature of the OracleRequest event.
     * ```solidity
     *  event OracleRequest(
     *    bytes32 indexed specId,
     *    address requester,
     *    bytes32 requestId,
     *    uint256 payment,
     *    address callbackAddr,
     *    bytes4 callbackFunctionId,
     *    uint256 cancelExpiration,
     *    uint256 dataVersion,
     *    bytes data
     *  );
     * ```
     * Note: this is a property used for testing purposes only.
     * It is not part of the actual run request.
     *
     * @solformat bytes32
     */
    topic: string;
}
/**
 * Convert the javascript format of the parameters needed to call the
 * ```solidity
 *  function fulfillOracleRequest(
 *    bytes32 _requestId,
 *    uint256 _payment,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _expiration,
 *    bytes32 _data
 *  )
 * ```
 * method on an Oracle.sol contract.
 *
 * @param runRequest The run request to flatten into the correct order to perform the `fulfillOracleRequest` function
 * @param response The response to fulfill the run request with, if it is an ascii string, it is converted to bytes32 string
 * @param txOpts Additional ethereum tx options
 */
export declare function convertFufillParams(runRequest: RunRequest, response: string, txOpts?: TxOptions): [string, string, string, string, string, string, TxOptions];
/**
 * Convert the javascript format of the parameters needed to call the
 * ```solidity
 *  function fulfillOracleRequest2(
 *    bytes32 _requestId,
 *    uint256 _payment,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _expiration,
 *    bytes memory _data
 *  )
 * ```
 * method on an Oracle.sol contract.
 *
 * @param runRequest The run request to flatten into the correct order to perform the `fulfillOracleRequest` function
 * @param response The response to fulfill the run request with, if it is an ascii string, it is converted to bytes32 string
 * @param txOpts Additional ethereum tx options
 */
export declare function convertFulfill2Params(runRequest: RunRequest, responseTypes: string[], responseValues: string[], txOpts?: TxOptions): [string, string, string, string, string, string, TxOptions];
/**
 * Convert the javascript format of the parameters needed to call the
 * ```solidity
 *  function cancelOracleRequest(
 *    bytes32 _requestId,
 *    uint256 _payment,
 *    bytes4 _callbackFunc,
 *    uint256 _expiration
 *  )
 * ```
 * method on an Oracle.sol contract.
 *
 * @param runRequest The run request to flatten into the correct order to perform the `cancelOracleRequest` function
 * @param txOpts Additional ethereum tx options
 */
export declare function convertCancelParams(runRequest: RunRequest, txOpts?: TxOptions): [string, string, string, string, TxOptions];
/**
 * Abi encode parameters to call the `oracleRequest` method on the Oracle.sol contract.
 * ```solidity
 *  function oracleRequest(
 *    address _sender,
 *    uint256 _payment,
 *    bytes32 _specId,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _nonce,
 *    uint256 _dataVersion,
 *    bytes _data
 *  )
 * ```
 *
 * @param specId The Job Specification ID
 * @param callbackAddr The callback contract address for the response
 * @param callbackFunctionId The callback function id for the response
 * @param nonce The nonce sent by the requester
 * @param data The CBOR payload of the request
 */
export declare function encodeOracleRequest(specId: string, callbackAddr: string, callbackFunctionId: string, nonce: number, data: BigNumberish, dataVersion?: BigNumberish): string;
/**
 * Abi encode parameters to call the `requestOracleData` method on the Operator.sol contract.
 * ```solidity
 *  function requestOracleData(
 *    address _sender,
 *    uint256 _payment,
 *    bytes32 _specId,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _nonce,
 *    uint256 _dataVersion,
 *    bytes _data
 *  )
 * ```
 *
 * @param specId The Job Specification ID
 * @param callbackAddr The callback contract address for the response
 * @param callbackFunctionId The callback function id for the response
 * @param nonce The nonce sent by the requester
 * @param data The CBOR payload of the request
 */
export declare function encodeRequestOracleData(specId: string, callbackAddr: string, callbackFunctionId: string, nonce: number, data: BigNumberish, dataVersion?: BigNumberish): string;
/**
 * Extract a javascript representation of a run request from the data
 * contained within a EVM log.
 * ```solidity
 *  event OracleRequest(
 *    bytes32 indexed specId,
 *    address requester,
 *    bytes32 requestId,
 *    uint256 payment,
 *    address callbackAddr,
 *    bytes4 callbackFunctionId,
 *    uint256 cancelExpiration,
 *    uint256 dataVersion,
 *    bytes data
 *  );
 * ```
 *
 * @param log The log to extract the run request from
 */
export declare function decodeRunRequest(log?: ethers.providers.Log): RunRequest;
/**
 * Extract a javascript representation of a ConcreteChainlinked#Request event
 * from an EVM log.
 * ```solidity
 *  event Request(
 *    bytes32 id,
 *    address callbackAddress,
 *    bytes4 callbackfunctionSelector,
 *    bytes data
 *  );
 * ```
 * The request event is emitted from the `ConcreteChainlinked.sol` testing contract.
 *
 * @param log The log to decode
 */
export declare function decodeCCRequest(log: ethers.providers.Log): [string, string, string, string];
export {};
//# sourceMappingURL=oracle.d.ts.map