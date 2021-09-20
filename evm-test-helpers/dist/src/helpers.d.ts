/// <reference types="node" />
import { ethers, utils } from 'ethers';
import { ContractReceipt } from 'ethers/contract';
import { EventDescription } from 'ethers/utils';
/**
 * Convert string to hex bytes
 * @param data string to onvert to hex bytes
 */
export declare function stringToBytes(data: string): string;
/**
 * Convert hex bytes to utf8 string
 * @param data bytes to convert to utf8 stirng
 */
export declare function bytesToString(data: string): string;
/**
 * Parse out an evm word (32 bytes) into an address (20 bytes) representation
 *
 * @param hex The evm word in hex string format to parse the address
 * out of.
 */
export declare function evmWordToAddress(hex?: string): string;
/**
 * Convert a number value to bytes32 format
 *
 * @param num The number value to convert to bytes32 format
 */
export declare function numToBytes32(num: Parameters<typeof utils.hexlify>[0]): string;
/**
 * Convert a UTF-8 string into a bytes32 hex string representation
 *
 * The inverse function of [[parseBytes32String]]
 *
 * @param args The UTF-8 string representation to convert to a bytes32 hex string representation
 */
export declare function toBytes32String(...args: Parameters<typeof utils.formatBytes32String>): ReturnType<typeof utils.formatBytes32String>;
/**
 * Convert a bytes32 formatted hex string into its UTF-8 representation
 *
 * The inverse function of [[toBytes32String]].
 *
 * @param args The bytes32 hex string representation to convert to an UTF-8 representation
 */
export declare function parseBytes32String(...args: Parameters<typeof utils.parseBytes32String>): ReturnType<typeof utils.parseBytes32String>;
/**
 * Converts any number, BigNumber, hex string or Arrayish to a hex string.
 *
 * @param args Value to convert to a hex string
 */
export declare function toHex(...args: Parameters<typeof utils.hexlify>): ReturnType<typeof utils.hexlify>;
/**
 * Create a buffer from a hex string
 *
 * @param hexstr The hex string to convert to a buffer
 */
export declare function hexToBuf(hexstr: string): Buffer;
/**
 * Convert an Ether value to a wei amount
 *
 * @param args Ether value to convert to an Ether amount
 */
export declare function toWei(...args: Parameters<typeof utils.parseEther>): ReturnType<typeof utils.parseEther>;
/**
 * Convert a value to an ethers BigNum
 *
 * @param num Value to convert to a BigNum
 */
export declare function bigNum(num: utils.BigNumberish): utils.BigNumber;
/**
 * Convert a UTF-8 string into a bytearray
 *
 * @param args The values needed to convert a string into a bytearray
 */
export declare function toUtf8Bytes(...args: Parameters<typeof utils.toUtf8Bytes>): ReturnType<typeof utils.toUtf8Bytes>;
/**
 * Turn a [x,y] coordinate into an ethereum address
 *
 * @param pubkey The x,y coordinate to turn into an ethereum address
 */
export declare function pubkeyToAddress(pubkey: utils.BigNumber[]): string;
/**
 * A wrapper function to make generated contracts compatible with truffle test suites.
 *
 * Note that the returned contract is an instance of ethers.Contract, not a @truffle/contract, so there are slight
 * api differences, though largely the same.
 *
 * @see https://docs.ethers.io/ethers.js/html/api-contract.html
 * @param contractFactory The ethers based contract factory to interop with
 * @param address The address to supply as the signer
 */
export declare function create<T extends new (...args: any[]) => any>(contractFactory: T, address: string): InstanceType<T>;
/**
 * Decodes a CBOR hex string, and adds opening and closing brackets to the CBOR if they are not present.
 *
 * @param hexstr The hex string to decode
 */
export declare function decodeDietCBOR(hexstr: string): any;
/**
 * Add a starting and closing map characters to a CBOR encoding if they are not already present.
 */
export declare function addCBORMapDelimiters(buffer: Buffer): Buffer;
/**
 * Add a hex prefix to a hex string
 *
 * @param hex The hex string to prepend the hex prefix to
 */
export declare function addHexPrefix(hex: string): string;
/**
 * Strip the leading 0x hex prefix from a hex string
 *
 * @param hex The hex string to strip the leading hex prefix out of
 */
export declare function stripHexPrefix(hex: string): string;
/**
 * Compute the keccak256 cryptographic hash of a value, returned as a hex string.
 * (Note: often Ethereum documentation refers to this, incorrectly, as SHA3)
 *
 * @param args The data to compute the keccak256 hash of
 */
export declare function keccak(...args: Parameters<typeof utils.keccak256>): ReturnType<typeof utils.keccak256>;
/**
 * Increase the current time within the evm to "n" seconds past the current time
 *
 * @param seconds The number of seconds to increase to the current time by
 * @param provider The ethers provider to send the time increase request to
 */
export declare function increaseTimeBy(seconds: number, provider: ethers.providers.JsonRpcProvider): Promise<void>;
/**
 * Instruct the provider to mine an additional block
 *
 * @param provider The ethers provider to instruct to mine an additional block
 */
export declare function mineBlock(provider: ethers.providers.JsonRpcProvider): Promise<void>;
/**
 * Increase the current time within the evm to 5 minutes past the current time
 *
 * @param provider The ethers provider to send the time increase request to
 */
export declare function increaseTime5Minutes(provider: ethers.providers.JsonRpcProvider): Promise<void>;
interface EventArgsArray extends Array<any> {
    [key: string]: any;
}
/**
 * Typecast an ethers event to its proper type, until
 * https://github.com/ethers-io/ethers.js/pull/698 is addressed
 *
 * @param event The event to typecast
 */
export declare function eventArgs(event?: ethers.Event): EventArgsArray;
/**
 * Find an event within a transaction receipt by its event description
 *
 * @param receipt The events array to search through
 * @param eventDescription The event description to pass to check its name by
 */
export declare function findEventIn(receipt: ContractReceipt, eventDescription: EventDescription): ethers.Event | undefined;
/**
 * Calculate six months from the current date in seconds
 */
export declare function sixMonthsFromNow(): utils.BigNumber;
/**
 * Extract array of logs from a transaction
 *
 * @param tx The transaction to wait for, then extract logs from
 */
export declare function getLogs(tx: ethers.ContractTransaction): Promise<ethers.providers.Log[]>;
/**
 * Retrieve single log from transaction
 *
 * @param tx The transaction to wait for, then extract logs from
 * @param index The index of the log to retrieve
 */
export declare function getLog(tx: ethers.ContractTransaction, index: number): Promise<ethers.providers.Log>;
export {};
//# sourceMappingURL=helpers.d.ts.map