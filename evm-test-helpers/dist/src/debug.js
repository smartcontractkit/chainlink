"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.makeDebug = void 0;
const tslib_1 = require("tslib");
/**
 * @packageDocumentation
 *
 * This file contains functionality for debugging tests, like creating loggers.
 */
const debug_1 = tslib_1.__importDefault(require("debug"));
/**
 * This creates a debug logger instance to be used within our internal code.
 *
 * @see https://www.npmjs.com/package/debug to see how to use the logger at runtime
 * @see wallet.ts makes extensive use of this function.
 * @param name The root namespace to assign to the log messages
 */
function makeDebug(name) {
    return debug_1.default(name);
}
exports.makeDebug = makeDebug;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZGVidWcuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9zcmMvZGVidWcudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7OztBQUFBOzs7O0dBSUc7QUFDSCwwREFBeUI7QUFFekI7Ozs7OztHQU1HO0FBQ0gsU0FBZ0IsU0FBUyxDQUFDLElBQVk7SUFDcEMsT0FBTyxlQUFLLENBQUMsSUFBSSxDQUFDLENBQUE7QUFDcEIsQ0FBQztBQUZELDhCQUVDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqXG4gKiBUaGlzIGZpbGUgY29udGFpbnMgZnVuY3Rpb25hbGl0eSBmb3IgZGVidWdnaW5nIHRlc3RzLCBsaWtlIGNyZWF0aW5nIGxvZ2dlcnMuXG4gKi9cbmltcG9ydCBkZWJ1ZyBmcm9tICdkZWJ1ZydcblxuLyoqXG4gKiBUaGlzIGNyZWF0ZXMgYSBkZWJ1ZyBsb2dnZXIgaW5zdGFuY2UgdG8gYmUgdXNlZCB3aXRoaW4gb3VyIGludGVybmFsIGNvZGUuXG4gKlxuICogQHNlZSBodHRwczovL3d3dy5ucG1qcy5jb20vcGFja2FnZS9kZWJ1ZyB0byBzZWUgaG93IHRvIHVzZSB0aGUgbG9nZ2VyIGF0IHJ1bnRpbWVcbiAqIEBzZWUgd2FsbGV0LnRzIG1ha2VzIGV4dGVuc2l2ZSB1c2Ugb2YgdGhpcyBmdW5jdGlvbi5cbiAqIEBwYXJhbSBuYW1lIFRoZSByb290IG5hbWVzcGFjZSB0byBhc3NpZ24gdG8gdGhlIGxvZyBtZXNzYWdlc1xuICovXG5leHBvcnQgZnVuY3Rpb24gbWFrZURlYnVnKG5hbWU6IHN0cmluZyk6IGRlYnVnLkRlYnVnZ2VyIHtcbiAgcmV0dXJuIGRlYnVnKG5hbWUpXG59XG4iXX0=