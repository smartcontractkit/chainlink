"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var ethers_1 = require("ethers");
var chai_1 = require("chai");
exports.constants = {
    ZERO_ADDRESS: "0x0000000000000000000000000000000000000000",
    ZERO_BYTES32: "0x0000000000000000000000000000000000000000000000000000000000000000",
    MAX_UINT256: ethers_1.BigNumber.from("2").pow(ethers_1.BigNumber.from("256")).sub(ethers_1.BigNumber.from("1")),
    MAX_INT256: ethers_1.BigNumber.from("2").pow(ethers_1.BigNumber.from("255")).sub(ethers_1.BigNumber.from("1")),
    MIN_INT256: ethers_1.BigNumber.from("2").pow(ethers_1.BigNumber.from("255")).mul(ethers_1.BigNumber.from("-1")),
};
/**
 * Check that a contract's abi exposes the expected interface.
 *
 * @param contract The contract with the actual abi to check the expected exposed methods and getters against.
 * @param expectedPublic The expected public exposed methods and getters to match against the actual abi.
 */
function publicAbi(contract, expectedPublic) {
    var actualPublic = [];
    for (var m in contract.functions) {
        if (!m.includes("(")) {
            actualPublic.push(m);
        }
    }
    for (var _i = 0, actualPublic_1 = actualPublic; _i < actualPublic_1.length; _i++) {
        var method = actualPublic_1[_i];
        var index = expectedPublic.indexOf(method);
        chai_1.assert.isAtLeast(index, 0, "#" + method + " is NOT expected to be public");
    }
    for (var _a = 0, expectedPublic_1 = expectedPublic; _a < expectedPublic_1.length; _a++) {
        var method = expectedPublic_1[_a];
        var index = actualPublic.indexOf(method);
        chai_1.assert.isAtLeast(index, 0, "#" + method + " is expected to be public");
    }
}
exports.publicAbi = publicAbi;
//# sourceMappingURL=helpers.js.map