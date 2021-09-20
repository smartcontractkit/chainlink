"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.callable = void 0;
const tslib_1 = require("tslib");
/**
 * @packageDocumentation
 *
 * This file deals with contract helpers to deal with ethers.js contract abstractions
 */
const ethers_1 = require("ethers");
tslib_1.__exportStar(require("./generated/factories/LinkToken__factory"), exports);
function callable(oldContract, methods) {
    var _a;
    const oldAbi = oldContract.interface.abi;
    const newAbi = oldAbi.map((fragment) => {
        var _a, _b;
        if (!methods.includes((_a = fragment.name) !== null && _a !== void 0 ? _a : '')) {
            return fragment;
        }
        if (((_b = fragment) === null || _b === void 0 ? void 0 : _b.constant) === false) {
            return {
                ...fragment,
                stateMutability: 'view',
                constant: true,
            };
        }
        return {
            ...fragment,
            stateMutability: 'view',
        };
    });
    const contract = new ethers_1.ethers.Contract(oldContract.address, newAbi, (_a = oldContract.signer) !== null && _a !== void 0 ? _a : oldContract.provider);
    return contract;
}
exports.callable = callable;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29udHJhY3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9zcmMvY29udHJhY3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7OztBQUFBOzs7O0dBSUc7QUFDSCxtQ0FBNEQ7QUFHNUQsbUZBQXdEO0FBZ0N4RCxTQUFnQixRQUFRLENBQUMsV0FBNEIsRUFBRSxPQUFpQjs7SUFDdEUsTUFBTSxNQUFNLEdBQUcsV0FBVyxDQUFDLFNBQVMsQ0FBQyxHQUFHLENBQUE7SUFDeEMsTUFBTSxNQUFNLEdBQUcsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLFFBQVEsRUFBRSxFQUFFOztRQUNyQyxJQUFJLENBQUMsT0FBTyxDQUFDLFFBQVEsQ0FBQyxNQUFBLFFBQVEsQ0FBQyxJQUFJLG1DQUFJLEVBQUUsQ0FBQyxFQUFFO1lBQzFDLE9BQU8sUUFBUSxDQUFBO1NBQ2hCO1FBRUQsSUFBSSxDQUFBLE1BQUMsUUFBNkIsMENBQUUsUUFBUSxNQUFLLEtBQUssRUFBRTtZQUN0RCxPQUFPO2dCQUNMLEdBQUcsUUFBUTtnQkFDWCxlQUFlLEVBQUUsTUFBTTtnQkFDdkIsUUFBUSxFQUFFLElBQUk7YUFDZixDQUFBO1NBQ0Y7UUFDRCxPQUFPO1lBQ0wsR0FBRyxRQUFRO1lBQ1gsZUFBZSxFQUFFLE1BQU07U0FDeEIsQ0FBQTtJQUNILENBQUMsQ0FBQyxDQUFBO0lBQ0YsTUFBTSxRQUFRLEdBQUcsSUFBSSxlQUFNLENBQUMsUUFBUSxDQUNsQyxXQUFXLENBQUMsT0FBTyxFQUNuQixNQUFNLEVBQ04sTUFBQSxXQUFXLENBQUMsTUFBTSxtQ0FBSSxXQUFXLENBQUMsUUFBUSxDQUMzQyxDQUFBO0lBRUQsT0FBTyxRQUFRLENBQUE7QUFDakIsQ0FBQztBQTFCRCw0QkEwQkMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBwYWNrYWdlRG9jdW1lbnRhdGlvblxuICpcbiAqIFRoaXMgZmlsZSBkZWFscyB3aXRoIGNvbnRyYWN0IGhlbHBlcnMgdG8gZGVhbCB3aXRoIGV0aGVycy5qcyBjb250cmFjdCBhYnN0cmFjdGlvbnNcbiAqL1xuaW1wb3J0IHsgZXRoZXJzLCBTaWduZXIsIENvbnRyYWN0VHJhbnNhY3Rpb24gfSBmcm9tICdldGhlcnMnXG5pbXBvcnQgeyBQcm92aWRlciB9IGZyb20gJ2V0aGVycy9wcm92aWRlcnMnXG5pbXBvcnQgeyBGdW5jdGlvbkZyYWdtZW50IH0gZnJvbSAnZXRoZXJzL3V0aWxzJ1xuZXhwb3J0ICogZnJvbSAnLi9nZW5lcmF0ZWQvZmFjdG9yaWVzL0xpbmtUb2tlbl9fZmFjdG9yeSdcblxuLyoqXG4gKiBUaGUgdHlwZSBvZiBhbnkgZnVuY3Rpb24gdGhhdCBpcyBkZXBsb3lhYmxlXG4gKi9cbnR5cGUgRGVwbG95YWJsZSA9IHtcbiAgZGVwbG95OiAoLi4uZGVwbG95QXJnczogYW55W10pID0+IFByb21pc2U8YW55PlxufVxuXG4vKipcbiAqIEdldCB0aGUgcmV0dXJuIHR5cGUgb2YgYSBmdW5jdGlvbiwgYW5kIHVuYm94IGFueSBwcm9taXNlc1xuICovXG5leHBvcnQgdHlwZSBJbnN0YW5jZTxUIGV4dGVuZHMgRGVwbG95YWJsZT4gPSBUIGV4dGVuZHMge1xuICBkZXBsb3k6ICguLi5kZXBsb3lBcmdzOiBhbnlbXSkgPT4gUHJvbWlzZTxpbmZlciBVPlxufVxuICA/IFVcbiAgOiBuZXZlclxuXG50eXBlIE92ZXJyaWRlPFQ+ID0ge1xuICBbSyBpbiBrZXlvZiBUXTogVFtLXSBleHRlbmRzICguLi5hcmdzOiBhbnlbXSkgPT4gUHJvbWlzZTxDb250cmFjdFRyYW5zYWN0aW9uPlxuICAgID8gKC4uLmFyZ3M6IGFueVtdKSA9PiBQcm9taXNlPGFueT5cbiAgICA6IFRbS11cbn1cblxuZXhwb3J0IHR5cGUgQ2FsbGFibGVPdmVycmlkZUluc3RhbmNlPFQgZXh0ZW5kcyBEZXBsb3lhYmxlPiA9IFQgZXh0ZW5kcyB7XG4gIGRlcGxveTogKC4uLmRlcGxveUFyZ3M6IGFueVtdKSA9PiBQcm9taXNlPGluZmVyIENvbnRyYWN0SW50ZXJmYWNlPlxufVxuICA/IE9taXQ8T3ZlcnJpZGU8Q29udHJhY3RJbnRlcmZhY2U+LCAnY29ubmVjdCc+ICYge1xuICAgICAgY29ubmVjdChzaWduZXI6IHN0cmluZyB8IFNpZ25lciB8IFByb3ZpZGVyKTogQ2FsbGFibGVPdmVycmlkZUluc3RhbmNlPFQ+XG4gICAgfVxuICA6IG5ldmVyXG5cbmV4cG9ydCBmdW5jdGlvbiBjYWxsYWJsZShvbGRDb250cmFjdDogZXRoZXJzLkNvbnRyYWN0LCBtZXRob2RzOiBzdHJpbmdbXSk6IGFueSB7XG4gIGNvbnN0IG9sZEFiaSA9IG9sZENvbnRyYWN0LmludGVyZmFjZS5hYmlcbiAgY29uc3QgbmV3QWJpID0gb2xkQWJpLm1hcCgoZnJhZ21lbnQpID0+IHtcbiAgICBpZiAoIW1ldGhvZHMuaW5jbHVkZXMoZnJhZ21lbnQubmFtZSA/PyAnJykpIHtcbiAgICAgIHJldHVybiBmcmFnbWVudFxuICAgIH1cblxuICAgIGlmICgoZnJhZ21lbnQgYXMgRnVuY3Rpb25GcmFnbWVudCk/LmNvbnN0YW50ID09PSBmYWxzZSkge1xuICAgICAgcmV0dXJuIHtcbiAgICAgICAgLi4uZnJhZ21lbnQsXG4gICAgICAgIHN0YXRlTXV0YWJpbGl0eTogJ3ZpZXcnLFxuICAgICAgICBjb25zdGFudDogdHJ1ZSxcbiAgICAgIH1cbiAgICB9XG4gICAgcmV0dXJuIHtcbiAgICAgIC4uLmZyYWdtZW50LFxuICAgICAgc3RhdGVNdXRhYmlsaXR5OiAndmlldycsXG4gICAgfVxuICB9KVxuICBjb25zdCBjb250cmFjdCA9IG5ldyBldGhlcnMuQ29udHJhY3QoXG4gICAgb2xkQ29udHJhY3QuYWRkcmVzcyxcbiAgICBuZXdBYmksXG4gICAgb2xkQ29udHJhY3Quc2lnbmVyID8/IG9sZENvbnRyYWN0LnByb3ZpZGVyLFxuICApXG5cbiAgcmV0dXJuIGNvbnRyYWN0XG59XG4iXX0=