'use strict';

var assert = require('chai').assert;
var abi = require('../index.js');
var contracts = require('./contracts.json');

describe('test basic encoding and decoding functionality', function () {
  it('should encode and decode contract data nicely', function () {
    var BalanceClaimInterface = JSON.parse(contracts.BalanceClaim['interface']);
    var encodeBalanceClaimMethod1 = abi.encodeMethod(BalanceClaimInterface[0], []);
    assert.equal(encodeBalanceClaimMethod1, '0x30509bca');
    var interfaceABI = [{ 'constant': false, 'inputs': [{ 'name': '_value', 'type': 'uint256' }], 'name': 'set', 'outputs': [{ 'name': '', 'type': 'bool' }], 'payable': false, 'type': 'function' }, { 'constant': false, 'inputs': [], 'name': 'get', 'outputs': [{ 'name': 'storeValue', 'type': 'uint256' }], 'payable': false, 'type': 'function' }, { 'anonymous': false, 'inputs': [{ 'indexed': false, 'name': '_newValue', 'type': 'uint256' }, { 'indexed': false, 'name': '_sender', 'type': 'address' }], 'name': 'SetComplete', 'type': 'event' }]; // eslint-disable-line

    var setMethodInputBytecode = abi.encodeMethod(interfaceABI[0], [24000]);
    abi.decodeMethod(interfaceABI[0], '0x0000000000000000000000000000000000000000000000000000000000000001');

    abi.encodeMethod(interfaceABI[1], []);
    abi.decodeMethod(interfaceABI[1], '0x000000000000000000000000000000000000000000000000000000000000b26e');

    abi.encodeEvent(interfaceABI[2], [24000, '0xca35b7d915458ef540ade6068dfe2f44e8fa733c']);
    abi.decodeEvent(interfaceABI[2], '0x0000000000000000000000000000000000000000000000000000000000000d7d000000000000000000000000ca35b7d915458ef540ade6068dfe2f44e8fa733c');

    assert.equal(setMethodInputBytecode, '0x60fe47b10000000000000000000000000000000000000000000000000000000000005dc0');
  });
});