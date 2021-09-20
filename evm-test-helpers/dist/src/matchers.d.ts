import { ethers } from 'ethers';
import { ContractReceipt } from 'ethers/contract';
import { BigNumberish, EventDescription } from 'ethers/utils';
/**
 * Check that two big numbers are the same value.
 *
 * @param expected The expected value to match against
 * @param actual The actual value to match against the expected value
 * @param failureMessage Failure message to display if the actual value does not match the expected value.
 */
export declare function bigNum(expected: BigNumberish, actual: BigNumberish, failureMessage?: string): void;
/**
 * Check that an evm operation reverts
 *
 * @param action The asynchronous action to execute, which should cause an evm revert.
 * @param msg The failure message to display if the action __does not__ throw
 */
export declare function evmRevert(action: (() => Promise<any>) | Promise<any>, msg?: string): Promise<void>;
/**
 * Check that a contract's abi exposes the expected interface.
 *
 * @param contract The contract with the actual abi to check the expected exposed methods and getters against.
 * @param expectedPublic The expected public exposed methods and getters to match against the actual abi.
 */
export declare function publicAbi(contract: ethers.Contract | ethers.ContractFactory, expectedPublic: string[]): void;
/**
 * Assert that an event exists
 *
 * @param receipt The contract receipt to find the event in
 * @param eventDescription A description of the event to search by
 */
export declare function eventExists(receipt: ContractReceipt, eventDescription: EventDescription): ethers.Event;
/**
 * Assert that an event doesnt exist
 *
 * @param receipt The contract receipt to find the event in
 * @param eventDescription A description of the event to search by
 */
export declare function eventDoesNotExist(receipt: ContractReceipt, eventDescription: EventDescription): void;
/**
 * Assert that an event doesnt exist
 *
 * @param max The maximum allowable gas difference
 * @param receipt1 The contract receipt to compare to
 * @param receipt2 The contract receipt with a gas difference
 */
export declare function gasDiffLessThan(max: number, receipt1: ContractReceipt, receipt2: ContractReceipt): void;
//# sourceMappingURL=matchers.d.ts.map