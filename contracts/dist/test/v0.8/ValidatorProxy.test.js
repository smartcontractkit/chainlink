"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var hardhat_1 = require("hardhat");
var helpers_1 = require("../helpers");
var chai_1 = require("chai");
describe("ValidatorProxy", function () {
    var accounts;
    var owner;
    var ownerAddress;
    var aggregator;
    var aggregatorAddress;
    var validator;
    var validatorAddress;
    var validatorProxy;
    beforeEach(function () { return __awaiter(void 0, void 0, void 0, function () {
        var vpf;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, hardhat_1.ethers.getSigners()];
                case 1:
                    accounts = _a.sent();
                    owner = accounts[0];
                    aggregator = accounts[1];
                    validator = accounts[2];
                    return [4 /*yield*/, owner.getAddress()];
                case 2:
                    ownerAddress = _a.sent();
                    return [4 /*yield*/, aggregator.getAddress()];
                case 3:
                    aggregatorAddress = _a.sent();
                    return [4 /*yield*/, validator.getAddress()];
                case 4:
                    validatorAddress = _a.sent();
                    return [4 /*yield*/, hardhat_1.ethers.getContractFactory("ValidatorProxy", owner)];
                case 5:
                    vpf = _a.sent();
                    return [4 /*yield*/, vpf.deploy(aggregatorAddress, validatorAddress)];
                case 6:
                    validatorProxy = _a.sent();
                    return [4 /*yield*/, validatorProxy.deployed()];
                case 7:
                    validatorProxy = _a.sent();
                    return [2 /*return*/];
            }
        });
    }); });
    it("has a limited public interface", function () { return __awaiter(void 0, void 0, void 0, function () {
        return __generator(this, function (_a) {
            helpers_1.publicAbi(validatorProxy, [
                // ConfirmedOwner functions
                "acceptOwnership",
                "owner",
                "transferOwnership",
                // ValidatorProxy functions
                "validate",
                "proposeNewAggregator",
                "upgradeAggregator",
                "getAggregators",
                "proposeNewValidator",
                "upgradeValidator",
                "getValidators",
                "typeAndVersion",
            ]);
            return [2 /*return*/];
        });
    }); });
    describe("#constructor", function () {
        it("should set the aggregator addresses correctly", function () { return __awaiter(void 0, void 0, void 0, function () {
            var response;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, validatorProxy.getAggregators()];
                    case 1:
                        response = _a.sent();
                        chai_1.assert.equal(response.current, aggregatorAddress);
                        chai_1.assert.equal(response.hasProposal, false);
                        chai_1.assert.equal(response.proposed, helpers_1.constants.ZERO_ADDRESS);
                        return [2 /*return*/];
                }
            });
        }); });
        it("should set the validator addresses conrrectly", function () { return __awaiter(void 0, void 0, void 0, function () {
            var response;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, validatorProxy.getValidators()];
                    case 1:
                        response = _a.sent();
                        chai_1.assert.equal(response.current, validatorAddress);
                        chai_1.assert.equal(response.hasProposal, false);
                        chai_1.assert.equal(response.proposed, helpers_1.constants.ZERO_ADDRESS);
                        return [2 /*return*/];
                }
            });
        }); });
        it("should set the owner correctly", function () { return __awaiter(void 0, void 0, void 0, function () {
            var response;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, validatorProxy.owner()];
                    case 1:
                        response = _a.sent();
                        chai_1.assert.equal(response, ownerAddress);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    describe("#proposeNewAggregator", function () {
        var newAggregator;
        var newAggregatorAddress;
        beforeEach(function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        newAggregator = accounts[3];
                        return [4 /*yield*/, newAggregator.getAddress()];
                    case 1:
                        newAggregatorAddress = _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
        describe("failure", function () {
            it("should only be called by the owner", function () { return __awaiter(void 0, void 0, void 0, function () {
                var stranger;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            stranger = accounts[4];
                            return [4 /*yield*/, chai_1.expect(validatorProxy.connect(stranger).proposeNewAggregator(newAggregatorAddress)).to.be.revertedWith("Only callable by owner")];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should revert if no change in proposal", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, validatorProxy.proposeNewAggregator(newAggregatorAddress)];
                        case 1:
                            _a.sent();
                            return [4 /*yield*/, chai_1.expect(validatorProxy.proposeNewAggregator(newAggregatorAddress)).to.be.revertedWith("Invalid proposal")];
                        case 2:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should revert if the proposal is the same as the current", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.proposeNewAggregator(aggregatorAddress)).to.be.revertedWith("Invalid proposal")];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
        });
        describe("success", function () {
            it("should emit an event", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.proposeNewAggregator(newAggregatorAddress))
                                .to.emit(validatorProxy, "AggregatorProposed")
                                .withArgs(newAggregatorAddress)];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should set the correct address and hasProposal is true", function () { return __awaiter(void 0, void 0, void 0, function () {
                var response;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, validatorProxy.proposeNewAggregator(newAggregatorAddress)];
                        case 1:
                            _a.sent();
                            return [4 /*yield*/, validatorProxy.getAggregators()];
                        case 2:
                            response = _a.sent();
                            chai_1.assert.equal(response.current, aggregatorAddress);
                            chai_1.assert.equal(response.hasProposal, true);
                            chai_1.assert.equal(response.proposed, newAggregatorAddress);
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should set a zero address and hasProposal is false", function () { return __awaiter(void 0, void 0, void 0, function () {
                var response;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, validatorProxy.proposeNewAggregator(newAggregatorAddress)];
                        case 1:
                            _a.sent();
                            return [4 /*yield*/, validatorProxy.proposeNewAggregator(helpers_1.constants.ZERO_ADDRESS)];
                        case 2:
                            _a.sent();
                            return [4 /*yield*/, validatorProxy.getAggregators()];
                        case 3:
                            response = _a.sent();
                            chai_1.assert.equal(response.current, aggregatorAddress);
                            chai_1.assert.equal(response.hasProposal, false);
                            chai_1.assert.equal(response.proposed, helpers_1.constants.ZERO_ADDRESS);
                            return [2 /*return*/];
                    }
                });
            }); });
        });
    });
    describe("#upgradeAggregator", function () {
        describe("failure", function () {
            it("should only be called by the owner", function () { return __awaiter(void 0, void 0, void 0, function () {
                var stranger;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            stranger = accounts[4];
                            return [4 /*yield*/, chai_1.expect(validatorProxy.connect(stranger).upgradeAggregator()).to.be.revertedWith("Only callable by owner")];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should revert if there is no proposal", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.upgradeAggregator()).to.be.revertedWith("No proposal")];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
        });
        describe("success", function () {
            var newAggregator;
            var newAggregatorAddress;
            beforeEach(function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            newAggregator = accounts[3];
                            return [4 /*yield*/, newAggregator.getAddress()];
                        case 1:
                            newAggregatorAddress = _a.sent();
                            return [4 /*yield*/, validatorProxy.proposeNewAggregator(newAggregatorAddress)];
                        case 2:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should emit an event", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.upgradeAggregator())
                                .to.emit(validatorProxy, "AggregatorUpgraded")
                                .withArgs(aggregatorAddress, newAggregatorAddress)];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should upgrade the addresses", function () { return __awaiter(void 0, void 0, void 0, function () {
                var response;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, validatorProxy.upgradeAggregator()];
                        case 1:
                            _a.sent();
                            return [4 /*yield*/, validatorProxy.getAggregators()];
                        case 2:
                            response = _a.sent();
                            chai_1.assert.equal(response.current, newAggregatorAddress);
                            chai_1.assert.equal(response.hasProposal, false);
                            chai_1.assert.equal(response.proposed, helpers_1.constants.ZERO_ADDRESS);
                            return [2 /*return*/];
                    }
                });
            }); });
        });
    });
    describe("#proposeNewValidator", function () {
        var newValidator;
        var newValidatorAddress;
        beforeEach(function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        newValidator = accounts[3];
                        return [4 /*yield*/, newValidator.getAddress()];
                    case 1:
                        newValidatorAddress = _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
        describe("failure", function () {
            it("should only be called by the owner", function () { return __awaiter(void 0, void 0, void 0, function () {
                var stranger;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            stranger = accounts[4];
                            return [4 /*yield*/, chai_1.expect(validatorProxy.connect(stranger).proposeNewAggregator(newValidatorAddress)).to.be.revertedWith("Only callable by owner")];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should revert if no change in proposal", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, validatorProxy.proposeNewValidator(newValidatorAddress)];
                        case 1:
                            _a.sent();
                            return [4 /*yield*/, chai_1.expect(validatorProxy.proposeNewValidator(newValidatorAddress)).to.be.revertedWith("Invalid proposal")];
                        case 2:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should revert if the proposal is the same as the current", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.proposeNewValidator(validatorAddress)).to.be.revertedWith("Invalid proposal")];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
        });
        describe("success", function () {
            it("should emit an event", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.proposeNewValidator(newValidatorAddress))
                                .to.emit(validatorProxy, "ValidatorProposed")
                                .withArgs(newValidatorAddress)];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should set the correct address and hasProposal is true", function () { return __awaiter(void 0, void 0, void 0, function () {
                var response;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, validatorProxy.proposeNewValidator(newValidatorAddress)];
                        case 1:
                            _a.sent();
                            return [4 /*yield*/, validatorProxy.getValidators()];
                        case 2:
                            response = _a.sent();
                            chai_1.assert.equal(response.current, validatorAddress);
                            chai_1.assert.equal(response.hasProposal, true);
                            chai_1.assert.equal(response.proposed, newValidatorAddress);
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should set a zero address and hasProposal is false", function () { return __awaiter(void 0, void 0, void 0, function () {
                var response;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, validatorProxy.proposeNewValidator(newValidatorAddress)];
                        case 1:
                            _a.sent();
                            return [4 /*yield*/, validatorProxy.proposeNewValidator(helpers_1.constants.ZERO_ADDRESS)];
                        case 2:
                            _a.sent();
                            return [4 /*yield*/, validatorProxy.getValidators()];
                        case 3:
                            response = _a.sent();
                            chai_1.assert.equal(response.current, validatorAddress);
                            chai_1.assert.equal(response.hasProposal, false);
                            chai_1.assert.equal(response.proposed, helpers_1.constants.ZERO_ADDRESS);
                            return [2 /*return*/];
                    }
                });
            }); });
        });
    });
    describe("#upgradeValidator", function () {
        describe("failure", function () {
            it("should only be called by the owner", function () { return __awaiter(void 0, void 0, void 0, function () {
                var stranger;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            stranger = accounts[4];
                            return [4 /*yield*/, chai_1.expect(validatorProxy.connect(stranger).upgradeValidator()).to.be.revertedWith("Only callable by owner")];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should revert if there is no proposal", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.upgradeValidator()).to.be.revertedWith("No proposal")];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
        });
        describe("success", function () {
            var newValidator;
            var newValidatorAddress;
            beforeEach(function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            newValidator = accounts[3];
                            return [4 /*yield*/, newValidator.getAddress()];
                        case 1:
                            newValidatorAddress = _a.sent();
                            return [4 /*yield*/, validatorProxy.proposeNewValidator(newValidatorAddress)];
                        case 2:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should emit an event", function () { return __awaiter(void 0, void 0, void 0, function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.upgradeValidator())
                                .to.emit(validatorProxy, "ValidatorUpgraded")
                                .withArgs(validatorAddress, newValidatorAddress)];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("should upgrade the addresses", function () { return __awaiter(void 0, void 0, void 0, function () {
                var response;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, validatorProxy.upgradeValidator()];
                        case 1:
                            _a.sent();
                            return [4 /*yield*/, validatorProxy.getValidators()];
                        case 2:
                            response = _a.sent();
                            chai_1.assert.equal(response.current, newValidatorAddress);
                            chai_1.assert.equal(response.hasProposal, false);
                            chai_1.assert.equal(response.proposed, helpers_1.constants.ZERO_ADDRESS);
                            return [2 /*return*/];
                    }
                });
            }); });
        });
    });
    describe("#validate", function () {
        describe("failure", function () {
            it("reverts when not called by aggregator or proposed aggregator", function () { return __awaiter(void 0, void 0, void 0, function () {
                var stranger;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            stranger = accounts[9];
                            return [4 /*yield*/, chai_1.expect(validatorProxy.connect(stranger).validate(99, 88, 77, 66)).to.be.revertedWith("Not a configured aggregator")];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
            it("reverts when there is no validator set", function () { return __awaiter(void 0, void 0, void 0, function () {
                var vpf;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4 /*yield*/, hardhat_1.ethers.getContractFactory("ValidatorProxy", owner)];
                        case 1:
                            vpf = _a.sent();
                            return [4 /*yield*/, vpf.deploy(aggregatorAddress, helpers_1.constants.ZERO_ADDRESS)];
                        case 2:
                            validatorProxy = _a.sent();
                            return [4 /*yield*/, validatorProxy.deployed()];
                        case 3:
                            _a.sent();
                            return [4 /*yield*/, chai_1.expect(validatorProxy.connect(aggregator).validate(99, 88, 77, 66)).to.be.revertedWith("No validator set")];
                        case 4:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            }); });
        });
        describe("success", function () {
            describe("from the aggregator", function () {
                var mockValidator1;
                beforeEach(function () { return __awaiter(void 0, void 0, void 0, function () {
                    var mvf, vpf;
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0: return [4 /*yield*/, hardhat_1.ethers.getContractFactory("MockAggregatorValidator", owner)];
                            case 1:
                                mvf = _a.sent();
                                return [4 /*yield*/, mvf.deploy(1)];
                            case 2:
                                mockValidator1 = _a.sent();
                                return [4 /*yield*/, mockValidator1.deployed()];
                            case 3:
                                mockValidator1 = _a.sent();
                                return [4 /*yield*/, hardhat_1.ethers.getContractFactory("ValidatorProxy", owner)];
                            case 4:
                                vpf = _a.sent();
                                return [4 /*yield*/, vpf.deploy(aggregatorAddress, mockValidator1.address)];
                            case 5:
                                validatorProxy = _a.sent();
                                return [4 /*yield*/, validatorProxy.deployed()];
                            case 6:
                                validatorProxy = _a.sent();
                                return [2 /*return*/];
                        }
                    });
                }); });
                describe("for a single validator", function () {
                    it("calls validate on the validator", function () { return __awaiter(void 0, void 0, void 0, function () {
                        return __generator(this, function (_a) {
                            switch (_a.label) {
                                case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.connect(aggregator).validate(200, 300, 400, 500))
                                        .to.emit(mockValidator1, "ValidateCalled")
                                        .withArgs(1, 200, 300, 400, 500)];
                                case 1:
                                    _a.sent();
                                    return [2 /*return*/];
                            }
                        });
                    }); });
                    it("uses a specific amount of gas", function () { return __awaiter(void 0, void 0, void 0, function () {
                        var resp, receipt;
                        return __generator(this, function (_a) {
                            switch (_a.label) {
                                case 0: return [4 /*yield*/, validatorProxy.connect(aggregator).validate(200, 300, 400, 500)];
                                case 1:
                                    resp = _a.sent();
                                    return [4 /*yield*/, resp.wait()];
                                case 2:
                                    receipt = _a.sent();
                                    chai_1.assert.equal(receipt.gasUsed.toString(), "35371");
                                    return [2 /*return*/];
                            }
                        });
                    }); });
                });
                describe("for a validator and a proposed validator", function () {
                    var mockValidator2;
                    beforeEach(function () { return __awaiter(void 0, void 0, void 0, function () {
                        var mvf;
                        return __generator(this, function (_a) {
                            switch (_a.label) {
                                case 0: return [4 /*yield*/, hardhat_1.ethers.getContractFactory("MockAggregatorValidator", owner)];
                                case 1:
                                    mvf = _a.sent();
                                    return [4 /*yield*/, mvf.deploy(2)];
                                case 2:
                                    mockValidator2 = _a.sent();
                                    return [4 /*yield*/, mockValidator2.deployed()];
                                case 3:
                                    mockValidator2 = _a.sent();
                                    return [4 /*yield*/, validatorProxy.proposeNewValidator(mockValidator2.address)];
                                case 4:
                                    _a.sent();
                                    return [2 /*return*/];
                            }
                        });
                    }); });
                    it("calls validate on the validator", function () { return __awaiter(void 0, void 0, void 0, function () {
                        return __generator(this, function (_a) {
                            switch (_a.label) {
                                case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.connect(aggregator).validate(2000, 3000, 4000, 5000))
                                        .to.emit(mockValidator1, "ValidateCalled")
                                        .withArgs(1, 2000, 3000, 4000, 5000)];
                                case 1:
                                    _a.sent();
                                    return [2 /*return*/];
                            }
                        });
                    }); });
                    it("also calls validate on the proposed validator", function () { return __awaiter(void 0, void 0, void 0, function () {
                        return __generator(this, function (_a) {
                            switch (_a.label) {
                                case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.connect(aggregator).validate(2000, 3000, 4000, 5000))
                                        .to.emit(mockValidator2, "ValidateCalled")
                                        .withArgs(2, 2000, 3000, 4000, 5000)];
                                case 1:
                                    _a.sent();
                                    return [2 /*return*/];
                            }
                        });
                    }); });
                    it("uses a specific amount of gas", function () { return __awaiter(void 0, void 0, void 0, function () {
                        var resp, receipt;
                        return __generator(this, function (_a) {
                            switch (_a.label) {
                                case 0: return [4 /*yield*/, validatorProxy.connect(aggregator).validate(2000, 3000, 4000, 5000)];
                                case 1:
                                    resp = _a.sent();
                                    return [4 /*yield*/, resp.wait()];
                                case 2:
                                    receipt = _a.sent();
                                    chai_1.assert.equal(receipt.gasUsed.toString(), "45318");
                                    return [2 /*return*/];
                            }
                        });
                    }); });
                });
            });
            describe("from the proposed aggregator", function () {
                var newAggregator;
                var newAggregatorAddress;
                beforeEach(function () { return __awaiter(void 0, void 0, void 0, function () {
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0:
                                newAggregator = accounts[3];
                                return [4 /*yield*/, newAggregator.getAddress()];
                            case 1:
                                newAggregatorAddress = _a.sent();
                                return [4 /*yield*/, validatorProxy.connect(owner).proposeNewAggregator(newAggregatorAddress)];
                            case 2:
                                _a.sent();
                                return [2 /*return*/];
                        }
                    });
                }); });
                it("emits an event", function () { return __awaiter(void 0, void 0, void 0, function () {
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0: return [4 /*yield*/, chai_1.expect(validatorProxy.connect(newAggregator).validate(555, 666, 777, 888))
                                    .to.emit(validatorProxy, "ProposedAggregatorValidateCall")
                                    .withArgs(newAggregatorAddress, 555, 666, 777, 888)];
                            case 1:
                                _a.sent();
                                return [2 /*return*/];
                        }
                    });
                }); });
            });
        });
    });
});
//# sourceMappingURL=ValidatorProxy.test.js.map