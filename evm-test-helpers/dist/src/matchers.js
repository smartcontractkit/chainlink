"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.gasDiffLessThan = exports.eventDoesNotExist = exports.eventExists = exports.publicAbi = exports.evmRevert = exports.bigNum = void 0;
/**
 * @packageDocumentation
 *
 * This file contains a number of matcher functions meant to perform common assertions for
 * ethereum based tests. Specific assertion functions targeting chainlink smart contracts live in
 * their respective contracts/<contract>.ts file.
 */
const chai_1 = require("chai");
const ethers_1 = require("ethers");
const debug_1 = require("./debug");
const helpers_1 = require("./helpers");
const debug = debug_1.makeDebug('helpers');
/**
 * Check that two big numbers are the same value.
 *
 * @param expected The expected value to match against
 * @param actual The actual value to match against the expected value
 * @param failureMessage Failure message to display if the actual value does not match the expected value.
 */
function bigNum(expected, actual, failureMessage) {
    const msg = failureMessage ? ': ' + failureMessage : '';
    chai_1.assert(ethers_1.ethers.utils.bigNumberify(expected).eq(ethers_1.ethers.utils.bigNumberify(actual)), `BigNum (expected)${expected} is not (actual)${actual} ${msg}`);
}
exports.bigNum = bigNum;
/**
 * Check that an evm operation reverts
 *
 * @param action The asynchronous action to execute, which should cause an evm revert.
 * @param msg The failure message to display if the action __does not__ throw
 */
async function evmRevert(action, msg) {
    const d = debug.extend('assertActionThrows');
    let e = undefined;
    try {
        if (typeof action === 'function') {
            await action();
        }
        else {
            await action;
        }
    }
    catch (error) {
        e = error;
    }
    d(e);
    if (!e) {
        chai_1.assert.exists(e, 'Expected an error to be raised');
        return;
    }
    chai_1.assert(e.message, 'Expected an error to contain a message');
    const ERROR_MESSAGES = ['invalid opcode', 'revert'];
    const hasErrored = ERROR_MESSAGES.some((msg) => { var _a; return (_a = e === null || e === void 0 ? void 0 : e.message) === null || _a === void 0 ? void 0 : _a.includes(msg); });
    if (msg) {
        expect(e.message).toMatch(msg);
    }
    chai_1.assert(hasErrored, `expected following error message to include ${ERROR_MESSAGES.join(' or ')}. Got: "${e.message}"`);
}
exports.evmRevert = evmRevert;
/**
 * Check that a contract's abi exposes the expected interface.
 *
 * @param contract The contract with the actual abi to check the expected exposed methods and getters against.
 * @param expectedPublic The expected public exposed methods and getters to match against the actual abi.
 */
function publicAbi(contract, expectedPublic) {
    const actualPublic = [];
    for (const method of contract.interface.abi) {
        if (method.type === 'function') {
            actualPublic.push(method.name);
        }
    }
    for (const method of actualPublic) {
        const index = expectedPublic.indexOf(method);
        chai_1.assert.isAtLeast(index, 0, `#${method} is NOT expected to be public`);
    }
    for (const method of expectedPublic) {
        const index = actualPublic.indexOf(method);
        chai_1.assert.isAtLeast(index, 0, `#${method} is expected to be public`);
    }
}
exports.publicAbi = publicAbi;
/**
 * Assert that an event exists
 *
 * @param receipt The contract receipt to find the event in
 * @param eventDescription A description of the event to search by
 */
function eventExists(receipt, eventDescription) {
    const event = helpers_1.findEventIn(receipt, eventDescription);
    if (!event) {
        throw Error(`Unable to find ${eventDescription.name} in transaction receipt`);
    }
    return event;
}
exports.eventExists = eventExists;
/**
 * Assert that an event doesnt exist
 *
 * @param receipt The contract receipt to find the event in
 * @param eventDescription A description of the event to search by
 */
function eventDoesNotExist(receipt, eventDescription) {
    const event = helpers_1.findEventIn(receipt, eventDescription);
    if (event) {
        throw Error(`Found ${eventDescription.name} in transaction receipt, when expecting no instances`);
    }
}
exports.eventDoesNotExist = eventDoesNotExist;
/**
 * Assert that an event doesnt exist
 *
 * @param max The maximum allowable gas difference
 * @param receipt1 The contract receipt to compare to
 * @param receipt2 The contract receipt with a gas difference
 */
function gasDiffLessThan(max, receipt1, receipt2) {
    var _a;
    chai_1.assert(receipt1, 'receipt1 is not present for gas comparison');
    chai_1.assert(receipt2, 'receipt2 is not present for gas comparison');
    const diff = (_a = receipt2.gasUsed) === null || _a === void 0 ? void 0 : _a.sub(receipt1.gasUsed || 0);
    chai_1.assert.isAbove(max, (diff === null || diff === void 0 ? void 0 : diff.toNumber()) || Infinity);
}
exports.gasDiffLessThan = gasDiffLessThan;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibWF0Y2hlcnMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9zcmMvbWF0Y2hlcnMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUE7Ozs7OztHQU1HO0FBQ0gsK0JBQTZCO0FBQzdCLG1DQUErQjtBQUcvQixtQ0FBbUM7QUFDbkMsdUNBQXVDO0FBQ3ZDLE1BQU0sS0FBSyxHQUFHLGlCQUFTLENBQUMsU0FBUyxDQUFDLENBQUE7QUFFbEM7Ozs7OztHQU1HO0FBQ0gsU0FBZ0IsTUFBTSxDQUNwQixRQUFzQixFQUN0QixNQUFvQixFQUNwQixjQUF1QjtJQUV2QixNQUFNLEdBQUcsR0FBRyxjQUFjLENBQUMsQ0FBQyxDQUFDLElBQUksR0FBRyxjQUFjLENBQUMsQ0FBQyxDQUFDLEVBQUUsQ0FBQTtJQUN2RCxhQUFNLENBQ0osZUFBTSxDQUFDLEtBQUssQ0FBQyxZQUFZLENBQUMsUUFBUSxDQUFDLENBQUMsRUFBRSxDQUFDLGVBQU0sQ0FBQyxLQUFLLENBQUMsWUFBWSxDQUFDLE1BQU0sQ0FBQyxDQUFDLEVBQ3pFLG9CQUFvQixRQUFRLG1CQUFtQixNQUFNLElBQUksR0FBRyxFQUFFLENBQy9ELENBQUE7QUFDSCxDQUFDO0FBVkQsd0JBVUM7QUFFRDs7Ozs7R0FLRztBQUNJLEtBQUssVUFBVSxTQUFTLENBQzdCLE1BQTJDLEVBQzNDLEdBQVk7SUFFWixNQUFNLENBQUMsR0FBRyxLQUFLLENBQUMsTUFBTSxDQUFDLG9CQUFvQixDQUFDLENBQUE7SUFDNUMsSUFBSSxDQUFDLEdBQXNCLFNBQVMsQ0FBQTtJQUVwQyxJQUFJO1FBQ0YsSUFBSSxPQUFPLE1BQU0sS0FBSyxVQUFVLEVBQUU7WUFDaEMsTUFBTSxNQUFNLEVBQUUsQ0FBQTtTQUNmO2FBQU07WUFDTCxNQUFNLE1BQU0sQ0FBQTtTQUNiO0tBQ0Y7SUFBQyxPQUFPLEtBQUssRUFBRTtRQUNkLENBQUMsR0FBRyxLQUFLLENBQUE7S0FDVjtJQUNELENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtJQUNKLElBQUksQ0FBQyxDQUFDLEVBQUU7UUFDTixhQUFNLENBQUMsTUFBTSxDQUFDLENBQUMsRUFBRSxnQ0FBZ0MsQ0FBQyxDQUFBO1FBQ2xELE9BQU07S0FDUDtJQUVELGFBQU0sQ0FBQyxDQUFDLENBQUMsT0FBTyxFQUFFLHdDQUF3QyxDQUFDLENBQUE7SUFFM0QsTUFBTSxjQUFjLEdBQUcsQ0FBQyxnQkFBZ0IsRUFBRSxRQUFRLENBQUMsQ0FBQTtJQUNuRCxNQUFNLFVBQVUsR0FBRyxjQUFjLENBQUMsSUFBSSxDQUFDLENBQUMsR0FBRyxFQUFFLEVBQUUsV0FBQyxPQUFBLE1BQUEsQ0FBQyxhQUFELENBQUMsdUJBQUQsQ0FBQyxDQUFFLE9BQU8sMENBQUUsUUFBUSxDQUFDLEdBQUcsQ0FBQyxDQUFBLEVBQUEsQ0FBQyxDQUFBO0lBRTFFLElBQUksR0FBRyxFQUFFO1FBQ1AsTUFBTSxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUMsQ0FBQyxPQUFPLENBQUMsR0FBRyxDQUFDLENBQUE7S0FDL0I7SUFFRCxhQUFNLENBQ0osVUFBVSxFQUNWLCtDQUErQyxjQUFjLENBQUMsSUFBSSxDQUNoRSxNQUFNLENBQ1AsV0FBVyxDQUFDLENBQUMsT0FBTyxHQUFHLENBQ3pCLENBQUE7QUFDSCxDQUFDO0FBckNELDhCQXFDQztBQUVEOzs7OztHQUtHO0FBQ0gsU0FBZ0IsU0FBUyxDQUN2QixRQUFrRCxFQUNsRCxjQUF3QjtJQUV4QixNQUFNLFlBQVksR0FBRyxFQUFFLENBQUE7SUFDdkIsS0FBSyxNQUFNLE1BQU0sSUFBSSxRQUFRLENBQUMsU0FBUyxDQUFDLEdBQUcsRUFBRTtRQUMzQyxJQUFJLE1BQU0sQ0FBQyxJQUFJLEtBQUssVUFBVSxFQUFFO1lBQzlCLFlBQVksQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFBO1NBQy9CO0tBQ0Y7SUFFRCxLQUFLLE1BQU0sTUFBTSxJQUFJLFlBQVksRUFBRTtRQUNqQyxNQUFNLEtBQUssR0FBRyxjQUFjLENBQUMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxDQUFBO1FBQzVDLGFBQU0sQ0FBQyxTQUFTLENBQUMsS0FBSyxFQUFFLENBQUMsRUFBRSxJQUFJLE1BQU0sK0JBQStCLENBQUMsQ0FBQTtLQUN0RTtJQUVELEtBQUssTUFBTSxNQUFNLElBQUksY0FBYyxFQUFFO1FBQ25DLE1BQU0sS0FBSyxHQUFHLFlBQVksQ0FBQyxPQUFPLENBQUMsTUFBTSxDQUFDLENBQUE7UUFDMUMsYUFBTSxDQUFDLFNBQVMsQ0FBQyxLQUFLLEVBQUUsQ0FBQyxFQUFFLElBQUksTUFBTSwyQkFBMkIsQ0FBQyxDQUFBO0tBQ2xFO0FBQ0gsQ0FBQztBQXBCRCw4QkFvQkM7QUFFRDs7Ozs7R0FLRztBQUNILFNBQWdCLFdBQVcsQ0FDekIsT0FBd0IsRUFDeEIsZ0JBQWtDO0lBRWxDLE1BQU0sS0FBSyxHQUFHLHFCQUFXLENBQUMsT0FBTyxFQUFFLGdCQUFnQixDQUFDLENBQUE7SUFDcEQsSUFBSSxDQUFDLEtBQUssRUFBRTtRQUNWLE1BQU0sS0FBSyxDQUNULGtCQUFrQixnQkFBZ0IsQ0FBQyxJQUFJLHlCQUF5QixDQUNqRSxDQUFBO0tBQ0Y7SUFFRCxPQUFPLEtBQUssQ0FBQTtBQUNkLENBQUM7QUFaRCxrQ0FZQztBQUVEOzs7OztHQUtHO0FBQ0gsU0FBZ0IsaUJBQWlCLENBQy9CLE9BQXdCLEVBQ3hCLGdCQUFrQztJQUVsQyxNQUFNLEtBQUssR0FBRyxxQkFBVyxDQUFDLE9BQU8sRUFBRSxnQkFBZ0IsQ0FBQyxDQUFBO0lBQ3BELElBQUksS0FBSyxFQUFFO1FBQ1QsTUFBTSxLQUFLLENBQ1QsU0FBUyxnQkFBZ0IsQ0FBQyxJQUFJLHNEQUFzRCxDQUNyRixDQUFBO0tBQ0Y7QUFDSCxDQUFDO0FBVkQsOENBVUM7QUFFRDs7Ozs7O0dBTUc7QUFDSCxTQUFnQixlQUFlLENBQzdCLEdBQVcsRUFDWCxRQUF5QixFQUN6QixRQUF5Qjs7SUFFekIsYUFBTSxDQUFDLFFBQVEsRUFBRSw0Q0FBNEMsQ0FBQyxDQUFBO0lBQzlELGFBQU0sQ0FBQyxRQUFRLEVBQUUsNENBQTRDLENBQUMsQ0FBQTtJQUM5RCxNQUFNLElBQUksR0FBRyxNQUFBLFFBQVEsQ0FBQyxPQUFPLDBDQUFFLEdBQUcsQ0FBQyxRQUFRLENBQUMsT0FBTyxJQUFJLENBQUMsQ0FBQyxDQUFBO0lBQ3pELGFBQU0sQ0FBQyxPQUFPLENBQUMsR0FBRyxFQUFFLENBQUEsSUFBSSxhQUFKLElBQUksdUJBQUosSUFBSSxDQUFFLFFBQVEsRUFBRSxLQUFJLFFBQVEsQ0FBQyxDQUFBO0FBQ25ELENBQUM7QUFURCwwQ0FTQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKlxuICogVGhpcyBmaWxlIGNvbnRhaW5zIGEgbnVtYmVyIG9mIG1hdGNoZXIgZnVuY3Rpb25zIG1lYW50IHRvIHBlcmZvcm0gY29tbW9uIGFzc2VydGlvbnMgZm9yXG4gKiBldGhlcmV1bSBiYXNlZCB0ZXN0cy4gU3BlY2lmaWMgYXNzZXJ0aW9uIGZ1bmN0aW9ucyB0YXJnZXRpbmcgY2hhaW5saW5rIHNtYXJ0IGNvbnRyYWN0cyBsaXZlIGluXG4gKiB0aGVpciByZXNwZWN0aXZlIGNvbnRyYWN0cy88Y29udHJhY3Q+LnRzIGZpbGUuXG4gKi9cbmltcG9ydCB7IGFzc2VydCB9IGZyb20gJ2NoYWknXG5pbXBvcnQgeyBldGhlcnMgfSBmcm9tICdldGhlcnMnXG5pbXBvcnQgeyBDb250cmFjdFJlY2VpcHQgfSBmcm9tICdldGhlcnMvY29udHJhY3QnXG5pbXBvcnQgeyBCaWdOdW1iZXJpc2gsIEV2ZW50RGVzY3JpcHRpb24gfSBmcm9tICdldGhlcnMvdXRpbHMnXG5pbXBvcnQgeyBtYWtlRGVidWcgfSBmcm9tICcuL2RlYnVnJ1xuaW1wb3J0IHsgZmluZEV2ZW50SW4gfSBmcm9tICcuL2hlbHBlcnMnXG5jb25zdCBkZWJ1ZyA9IG1ha2VEZWJ1ZygnaGVscGVycycpXG5cbi8qKlxuICogQ2hlY2sgdGhhdCB0d28gYmlnIG51bWJlcnMgYXJlIHRoZSBzYW1lIHZhbHVlLlxuICpcbiAqIEBwYXJhbSBleHBlY3RlZCBUaGUgZXhwZWN0ZWQgdmFsdWUgdG8gbWF0Y2ggYWdhaW5zdFxuICogQHBhcmFtIGFjdHVhbCBUaGUgYWN0dWFsIHZhbHVlIHRvIG1hdGNoIGFnYWluc3QgdGhlIGV4cGVjdGVkIHZhbHVlXG4gKiBAcGFyYW0gZmFpbHVyZU1lc3NhZ2UgRmFpbHVyZSBtZXNzYWdlIHRvIGRpc3BsYXkgaWYgdGhlIGFjdHVhbCB2YWx1ZSBkb2VzIG5vdCBtYXRjaCB0aGUgZXhwZWN0ZWQgdmFsdWUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBiaWdOdW0oXG4gIGV4cGVjdGVkOiBCaWdOdW1iZXJpc2gsXG4gIGFjdHVhbDogQmlnTnVtYmVyaXNoLFxuICBmYWlsdXJlTWVzc2FnZT86IHN0cmluZyxcbik6IHZvaWQge1xuICBjb25zdCBtc2cgPSBmYWlsdXJlTWVzc2FnZSA/ICc6ICcgKyBmYWlsdXJlTWVzc2FnZSA6ICcnXG4gIGFzc2VydChcbiAgICBldGhlcnMudXRpbHMuYmlnTnVtYmVyaWZ5KGV4cGVjdGVkKS5lcShldGhlcnMudXRpbHMuYmlnTnVtYmVyaWZ5KGFjdHVhbCkpLFxuICAgIGBCaWdOdW0gKGV4cGVjdGVkKSR7ZXhwZWN0ZWR9IGlzIG5vdCAoYWN0dWFsKSR7YWN0dWFsfSAke21zZ31gLFxuICApXG59XG5cbi8qKlxuICogQ2hlY2sgdGhhdCBhbiBldm0gb3BlcmF0aW9uIHJldmVydHNcbiAqXG4gKiBAcGFyYW0gYWN0aW9uIFRoZSBhc3luY2hyb25vdXMgYWN0aW9uIHRvIGV4ZWN1dGUsIHdoaWNoIHNob3VsZCBjYXVzZSBhbiBldm0gcmV2ZXJ0LlxuICogQHBhcmFtIG1zZyBUaGUgZmFpbHVyZSBtZXNzYWdlIHRvIGRpc3BsYXkgaWYgdGhlIGFjdGlvbiBfX2RvZXMgbm90X18gdGhyb3dcbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIGV2bVJldmVydChcbiAgYWN0aW9uOiAoKCkgPT4gUHJvbWlzZTxhbnk+KSB8IFByb21pc2U8YW55PixcbiAgbXNnPzogc3RyaW5nLFxuKSB7XG4gIGNvbnN0IGQgPSBkZWJ1Zy5leHRlbmQoJ2Fzc2VydEFjdGlvblRocm93cycpXG4gIGxldCBlOiBFcnJvciB8IHVuZGVmaW5lZCA9IHVuZGVmaW5lZFxuXG4gIHRyeSB7XG4gICAgaWYgKHR5cGVvZiBhY3Rpb24gPT09ICdmdW5jdGlvbicpIHtcbiAgICAgIGF3YWl0IGFjdGlvbigpXG4gICAgfSBlbHNlIHtcbiAgICAgIGF3YWl0IGFjdGlvblxuICAgIH1cbiAgfSBjYXRjaCAoZXJyb3IpIHtcbiAgICBlID0gZXJyb3JcbiAgfVxuICBkKGUpXG4gIGlmICghZSkge1xuICAgIGFzc2VydC5leGlzdHMoZSwgJ0V4cGVjdGVkIGFuIGVycm9yIHRvIGJlIHJhaXNlZCcpXG4gICAgcmV0dXJuXG4gIH1cblxuICBhc3NlcnQoZS5tZXNzYWdlLCAnRXhwZWN0ZWQgYW4gZXJyb3IgdG8gY29udGFpbiBhIG1lc3NhZ2UnKVxuXG4gIGNvbnN0IEVSUk9SX01FU1NBR0VTID0gWydpbnZhbGlkIG9wY29kZScsICdyZXZlcnQnXVxuICBjb25zdCBoYXNFcnJvcmVkID0gRVJST1JfTUVTU0FHRVMuc29tZSgobXNnKSA9PiBlPy5tZXNzYWdlPy5pbmNsdWRlcyhtc2cpKVxuXG4gIGlmIChtc2cpIHtcbiAgICBleHBlY3QoZS5tZXNzYWdlKS50b01hdGNoKG1zZylcbiAgfVxuXG4gIGFzc2VydChcbiAgICBoYXNFcnJvcmVkLFxuICAgIGBleHBlY3RlZCBmb2xsb3dpbmcgZXJyb3IgbWVzc2FnZSB0byBpbmNsdWRlICR7RVJST1JfTUVTU0FHRVMuam9pbihcbiAgICAgICcgb3IgJyxcbiAgICApfS4gR290OiBcIiR7ZS5tZXNzYWdlfVwiYCxcbiAgKVxufVxuXG4vKipcbiAqIENoZWNrIHRoYXQgYSBjb250cmFjdCdzIGFiaSBleHBvc2VzIHRoZSBleHBlY3RlZCBpbnRlcmZhY2UuXG4gKlxuICogQHBhcmFtIGNvbnRyYWN0IFRoZSBjb250cmFjdCB3aXRoIHRoZSBhY3R1YWwgYWJpIHRvIGNoZWNrIHRoZSBleHBlY3RlZCBleHBvc2VkIG1ldGhvZHMgYW5kIGdldHRlcnMgYWdhaW5zdC5cbiAqIEBwYXJhbSBleHBlY3RlZFB1YmxpYyBUaGUgZXhwZWN0ZWQgcHVibGljIGV4cG9zZWQgbWV0aG9kcyBhbmQgZ2V0dGVycyB0byBtYXRjaCBhZ2FpbnN0IHRoZSBhY3R1YWwgYWJpLlxuICovXG5leHBvcnQgZnVuY3Rpb24gcHVibGljQWJpKFxuICBjb250cmFjdDogZXRoZXJzLkNvbnRyYWN0IHwgZXRoZXJzLkNvbnRyYWN0RmFjdG9yeSxcbiAgZXhwZWN0ZWRQdWJsaWM6IHN0cmluZ1tdLFxuKSB7XG4gIGNvbnN0IGFjdHVhbFB1YmxpYyA9IFtdXG4gIGZvciAoY29uc3QgbWV0aG9kIG9mIGNvbnRyYWN0LmludGVyZmFjZS5hYmkpIHtcbiAgICBpZiAobWV0aG9kLnR5cGUgPT09ICdmdW5jdGlvbicpIHtcbiAgICAgIGFjdHVhbFB1YmxpYy5wdXNoKG1ldGhvZC5uYW1lKVxuICAgIH1cbiAgfVxuXG4gIGZvciAoY29uc3QgbWV0aG9kIG9mIGFjdHVhbFB1YmxpYykge1xuICAgIGNvbnN0IGluZGV4ID0gZXhwZWN0ZWRQdWJsaWMuaW5kZXhPZihtZXRob2QpXG4gICAgYXNzZXJ0LmlzQXRMZWFzdChpbmRleCwgMCwgYCMke21ldGhvZH0gaXMgTk9UIGV4cGVjdGVkIHRvIGJlIHB1YmxpY2ApXG4gIH1cblxuICBmb3IgKGNvbnN0IG1ldGhvZCBvZiBleHBlY3RlZFB1YmxpYykge1xuICAgIGNvbnN0IGluZGV4ID0gYWN0dWFsUHVibGljLmluZGV4T2YobWV0aG9kKVxuICAgIGFzc2VydC5pc0F0TGVhc3QoaW5kZXgsIDAsIGAjJHttZXRob2R9IGlzIGV4cGVjdGVkIHRvIGJlIHB1YmxpY2ApXG4gIH1cbn1cblxuLyoqXG4gKiBBc3NlcnQgdGhhdCBhbiBldmVudCBleGlzdHNcbiAqXG4gKiBAcGFyYW0gcmVjZWlwdCBUaGUgY29udHJhY3QgcmVjZWlwdCB0byBmaW5kIHRoZSBldmVudCBpblxuICogQHBhcmFtIGV2ZW50RGVzY3JpcHRpb24gQSBkZXNjcmlwdGlvbiBvZiB0aGUgZXZlbnQgdG8gc2VhcmNoIGJ5XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBldmVudEV4aXN0cyhcbiAgcmVjZWlwdDogQ29udHJhY3RSZWNlaXB0LFxuICBldmVudERlc2NyaXB0aW9uOiBFdmVudERlc2NyaXB0aW9uLFxuKTogZXRoZXJzLkV2ZW50IHtcbiAgY29uc3QgZXZlbnQgPSBmaW5kRXZlbnRJbihyZWNlaXB0LCBldmVudERlc2NyaXB0aW9uKVxuICBpZiAoIWV2ZW50KSB7XG4gICAgdGhyb3cgRXJyb3IoXG4gICAgICBgVW5hYmxlIHRvIGZpbmQgJHtldmVudERlc2NyaXB0aW9uLm5hbWV9IGluIHRyYW5zYWN0aW9uIHJlY2VpcHRgLFxuICAgIClcbiAgfVxuXG4gIHJldHVybiBldmVudFxufVxuXG4vKipcbiAqIEFzc2VydCB0aGF0IGFuIGV2ZW50IGRvZXNudCBleGlzdFxuICpcbiAqIEBwYXJhbSByZWNlaXB0IFRoZSBjb250cmFjdCByZWNlaXB0IHRvIGZpbmQgdGhlIGV2ZW50IGluXG4gKiBAcGFyYW0gZXZlbnREZXNjcmlwdGlvbiBBIGRlc2NyaXB0aW9uIG9mIHRoZSBldmVudCB0byBzZWFyY2ggYnlcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGV2ZW50RG9lc05vdEV4aXN0KFxuICByZWNlaXB0OiBDb250cmFjdFJlY2VpcHQsXG4gIGV2ZW50RGVzY3JpcHRpb246IEV2ZW50RGVzY3JpcHRpb24sXG4pIHtcbiAgY29uc3QgZXZlbnQgPSBmaW5kRXZlbnRJbihyZWNlaXB0LCBldmVudERlc2NyaXB0aW9uKVxuICBpZiAoZXZlbnQpIHtcbiAgICB0aHJvdyBFcnJvcihcbiAgICAgIGBGb3VuZCAke2V2ZW50RGVzY3JpcHRpb24ubmFtZX0gaW4gdHJhbnNhY3Rpb24gcmVjZWlwdCwgd2hlbiBleHBlY3Rpbmcgbm8gaW5zdGFuY2VzYCxcbiAgICApXG4gIH1cbn1cblxuLyoqXG4gKiBBc3NlcnQgdGhhdCBhbiBldmVudCBkb2VzbnQgZXhpc3RcbiAqXG4gKiBAcGFyYW0gbWF4IFRoZSBtYXhpbXVtIGFsbG93YWJsZSBnYXMgZGlmZmVyZW5jZVxuICogQHBhcmFtIHJlY2VpcHQxIFRoZSBjb250cmFjdCByZWNlaXB0IHRvIGNvbXBhcmUgdG9cbiAqIEBwYXJhbSByZWNlaXB0MiBUaGUgY29udHJhY3QgcmVjZWlwdCB3aXRoIGEgZ2FzIGRpZmZlcmVuY2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGdhc0RpZmZMZXNzVGhhbihcbiAgbWF4OiBudW1iZXIsXG4gIHJlY2VpcHQxOiBDb250cmFjdFJlY2VpcHQsXG4gIHJlY2VpcHQyOiBDb250cmFjdFJlY2VpcHQsXG4pIHtcbiAgYXNzZXJ0KHJlY2VpcHQxLCAncmVjZWlwdDEgaXMgbm90IHByZXNlbnQgZm9yIGdhcyBjb21wYXJpc29uJylcbiAgYXNzZXJ0KHJlY2VpcHQyLCAncmVjZWlwdDIgaXMgbm90IHByZXNlbnQgZm9yIGdhcyBjb21wYXJpc29uJylcbiAgY29uc3QgZGlmZiA9IHJlY2VpcHQyLmdhc1VzZWQ/LnN1YihyZWNlaXB0MS5nYXNVc2VkIHx8IDApXG4gIGFzc2VydC5pc0Fib3ZlKG1heCwgZGlmZj8udG9OdW1iZXIoKSB8fCBJbmZpbml0eSlcbn1cbiJdfQ==